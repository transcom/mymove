package ghcapi

import (
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	orderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	orderservice "github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/swagger/nullable"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trace"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestGetOrderHandlerIntegration() {
	officeUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
		factory.GetTraitOfficeUserStubbed,
	})

	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	params := orderop.GetOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
	}
	handlerConfig := suite.HandlerConfig()
	handler := GetOrdersHandler{
		handlerConfig,
		orderservice.NewOrderFetcher(),
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	suite.IsType(&orderop.GetOrderOK{}, response)
	orderOK := response.(*orderop.GetOrderOK)
	ordersPayload := orderOK.Payload

	// Validate outgoing payload
	suite.NoError(ordersPayload.Validate(strfmt.Default))

	suite.Equal(order.ID.String(), ordersPayload.ID.String())
	suite.Equal(move.Locator, ordersPayload.MoveCode)
	suite.Equal(order.ServiceMemberID.String(), ordersPayload.Customer.ID.String())
	suite.Equal(order.NewDutyLocationID.String(), ordersPayload.DestinationDutyLocation.ID.String())
	suite.NotNil(order.NewDutyLocation)
	payloadEntitlement := ordersPayload.Entitlement
	suite.Equal((*order.EntitlementID).String(), payloadEntitlement.ID.String())
	orderEntitlement := order.Entitlement
	suite.NotNil(orderEntitlement)
	suite.EqualValues(orderEntitlement.ProGearWeight, payloadEntitlement.ProGearWeight)
	suite.EqualValues(orderEntitlement.ProGearWeightSpouse, payloadEntitlement.ProGearWeightSpouse)
	suite.EqualValues(orderEntitlement.RequiredMedicalEquipmentWeight, payloadEntitlement.RequiredMedicalEquipmentWeight)
	suite.EqualValues(orderEntitlement.OrganizationalClothingAndIndividualEquipment, payloadEntitlement.OrganizationalClothingAndIndividualEquipment)
	suite.Equal(order.OriginDutyLocation.ID.String(), ordersPayload.OriginDutyLocation.ID.String())
	suite.NotZero(order.OriginDutyLocation)
	suite.NotZero(ordersPayload.DateIssued)
}

func (suite *HandlerSuite) TestWeightAllowances() {
	suite.Run("With E-1 rank and no dependents", func() {
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Stub: true,
			Order: models.Order{
				ID:            uuid.Must(uuid.NewV4()),
				HasDependents: *swag.Bool(false),
			},
			Entitlement: models.Entitlement{
				ID:                   uuid.Must(uuid.NewV4()),
				DependentsAuthorized: swag.Bool(false),
				ProGearWeight:        2000,
				ProGearWeightSpouse:  500,
			},
		})
		request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
		params := orderop.GetOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
		}
		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("FetchOrder", mock.AnythingOfType("*appcontext.appContext"),
			order.ID).Return(&order, nil)

		handlerConfig := suite.HandlerConfig()
		handler := GetOrdersHandler{
			handlerConfig,
			&orderFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.GetOrderOK{}, response)
		orderOK := response.(*orderop.GetOrderOK)
		orderPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(orderPayload.Validate(strfmt.Default))

		payloadEntitlement := orderPayload.Entitlement
		orderEntitlement := order.Entitlement
		expectedAllowance := int64(orderEntitlement.WeightAllotment().TotalWeightSelf)

		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeight), payloadEntitlement.ProGearWeight)
		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
		suite.Equal(expectedAllowance, payloadEntitlement.TotalWeight)
		suite.Equal(int64(*orderEntitlement.AuthorizedWeight()), *payloadEntitlement.AuthorizedWeight)
	})

	suite.Run("With E-1 rank and dependents", func() {
		order := testdatagen.MakeOrder(suite.DB(), testdatagen.Assertions{
			Stub: true,
			Order: models.Order{
				ID:            uuid.Must(uuid.NewV4()),
				HasDependents: *swag.Bool(true),
			},
		})

		request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
		params := orderop.GetOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
		}

		orderFetcher := mocks.OrderFetcher{}
		orderFetcher.On("FetchOrder", mock.AnythingOfType("*appcontext.appContext"),
			order.ID).Return(&order, nil)

		handlerConfig := suite.HandlerConfig()
		handler := GetOrdersHandler{
			handlerConfig,
			&orderFetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.GetOrderOK{}, response)
		orderOK := response.(*orderop.GetOrderOK)
		orderPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(orderPayload.Validate(strfmt.Default))

		payloadEntitlement := orderPayload.Entitlement
		orderEntitlement := order.Entitlement
		expectedAllowance := int64(orderEntitlement.WeightAllotment().TotalWeightSelfPlusDependents)

		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeight), payloadEntitlement.ProGearWeight)
		suite.Equal(int64(orderEntitlement.WeightAllotment().ProGearWeightSpouse), payloadEntitlement.ProGearWeightSpouse)
		suite.Equal(expectedAllowance, payloadEntitlement.TotalWeight)
		suite.Equal(int64(*orderEntitlement.AuthorizedWeight()), *payloadEntitlement.AuthorizedWeight)
	})
}

type updateOrderHandlerAmendedUploadSubtestData struct {
	handlerConfig           handlers.HandlerConfig
	userUploader            *uploader.UserUploader
	amendedOrder            models.Order
	approvalsRequestedMove  models.Move
	originDutyLocation      models.DutyLocation
	destinationDutyLocation models.DutyLocation
}

func (suite *HandlerSuite) makeUpdateOrderHandlerAmendedUploadSubtestData() (subtestData *updateOrderHandlerAmendedUploadSubtestData) {
	subtestData = &updateOrderHandlerAmendedUploadSubtestData{}
	subtestData.handlerConfig = suite.createS3HandlerConfig()

	var err error
	subtestData.userUploader, err = uploader.NewUserUploader(subtestData.handlerConfig.FileStorer(), 100*uploader.MB)
	assert.NoError(suite.T(), err, "failed to create user uploader for amended orders")
	amendedDocument := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
	amendedUpload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
		UserUpload: models.UserUpload{
			DocumentID: &amendedDocument.ID,
			Document:   amendedDocument,
			UploaderID: amendedDocument.ServiceMember.UserID,
		},
		UserUploader: subtestData.userUploader,
	})

	amendedDocument.UserUploads = append(amendedDocument.UserUploads, amendedUpload)
	subtestData.approvalsRequestedMove = testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
		Order: models.Order{
			UploadedAmendedOrders:   &amendedDocument,
			UploadedAmendedOrdersID: &amendedDocument.ID,
			ServiceMember:           amendedDocument.ServiceMember,
			ServiceMemberID:         amendedDocument.ServiceMemberID,
		},
	})

	subtestData.amendedOrder = subtestData.approvalsRequestedMove.Orders

	subtestData.originDutyLocation = testdatagen.MakeDefaultDutyLocation(suite.DB())
	subtestData.destinationDutyLocation = testdatagen.MakeDefaultDutyLocation(suite.DB())

	return subtestData
}

func (suite *HandlerSuite) TestUpdateOrderHandlerWithAmendedUploads() {

	queryBuilder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
		queryBuilder,
		mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter),
		moveRouter,
	)

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	ordersAcknowledgement := true

	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)

	suite.Run("Returns 200 when acknowledging orders", func() {
		subtestData := suite.makeUpdateOrderHandlerAmendedUploadSubtestData()
		handlerConfig := subtestData.handlerConfig
		userUploader := subtestData.userUploader
		destinationDutyLocation := subtestData.destinationDutyLocation
		originDutyLocation := subtestData.originDutyLocation

		document := testdatagen.MakeDocument(suite.DB(), testdatagen.Assertions{})
		upload := testdatagen.MakeUserUpload(suite.DB(), testdatagen.Assertions{
			UserUpload: models.UserUpload{
				DocumentID: &document.ID,
				Document:   document,
				UploaderID: document.ServiceMember.UserID,
			},
			UserUploader: userUploader,
		})

		document.UserUploads = append(document.UserUploads, upload)
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				UploadedAmendedOrders:   &document,
				UploadedAmendedOrdersID: &document.ID,
				ServiceMember:           document.ServiceMember,
				ServiceMemberID:         document.ServiceMemberID,
			},
		})

		order := move.Orders

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		body := &ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:   &deptIndicator,
			IssueDate:             handlers.FmtDatePtr(&issueDate),
			ReportByDate:          handlers.FmtDatePtr(&reportByDate),
			OrdersType:            ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
			OrdersTypeDetail:      &ordersTypeDetail,
			OrdersNumber:          handlers.FmtString("ORDER100"),
			NewDutyLocationID:     handlers.FmtUUID(destinationDutyLocation.ID),
			OriginDutyLocationID:  handlers.FmtUUID(originDutyLocation.ID),
			Tac:                   handlers.FmtString("E19A"),
			Sac:                   nullable.NewString("987654321"),
			NtsTac:                nullable.NewString("E19A"),
			NtsSac:                nullable.NewString("987654321"),
			OrdersAcknowledgement: &ordersAcknowledgement,
		}

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		handler := UpdateOrderHandler{
			handlerConfig,
			orderservice.NewOrderUpdater(moveRouter),
			moveTaskOrderUpdater,
		}

		suite.Nil(order.AmendedOrdersAcknowledgedAt)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		suite.IsType(&orderop.UpdateOrderOK{}, response)
		orderOK := response.(*orderop.UpdateOrderOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyLocationID.String(), ordersPayload.DestinationDutyLocation.ID.String())
		suite.Equal(body.OriginDutyLocationID.String(), ordersPayload.OriginDutyLocation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(*body.OrdersType, ordersPayload.OrderType)
		suite.Equal(body.OrdersTypeDetail, ordersPayload.OrderTypeDetail)
		suite.Equal(body.OrdersNumber, ordersPayload.OrderNumber)
		suite.Equal(body.DepartmentIndicator, ordersPayload.DepartmentIndicator)
		suite.Equal(body.Tac, ordersPayload.Tac)
		suite.Equal(body.Sac.Value, ordersPayload.Sac)
		suite.Equal(body.NtsTac.Value, ordersPayload.NtsTac)
		suite.Equal(body.NtsSac.Value, ordersPayload.NtsSac)
		suite.NotNil(ordersPayload.AmendedOrdersAcknowledgedAt)

		reloadErr := suite.DB().Reload(&move)
		suite.NoError(reloadErr, "error reloading move of amended orders")

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
	})

	suite.Run("Does not update move status if orders are not acknowledged", func() {
		subtestData := suite.makeUpdateOrderHandlerAmendedUploadSubtestData()
		handlerConfig := subtestData.handlerConfig
		destinationDutyLocation := subtestData.destinationDutyLocation
		originDutyLocation := subtestData.originDutyLocation
		amendedOrder := subtestData.amendedOrder
		approvalsRequestedMove := subtestData.approvalsRequestedMove

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		unacknowledgedOrders := false
		body := &ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:   &deptIndicator,
			IssueDate:             handlers.FmtDatePtr(&issueDate),
			ReportByDate:          handlers.FmtDatePtr(&reportByDate),
			OrdersType:            ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
			OrdersTypeDetail:      &ordersTypeDetail,
			OrdersNumber:          handlers.FmtString("ORDER100"),
			NewDutyLocationID:     handlers.FmtUUID(destinationDutyLocation.ID),
			OriginDutyLocationID:  handlers.FmtUUID(originDutyLocation.ID),
			Tac:                   handlers.FmtString("E19A"),
			Sac:                   nullable.NewString("987654321"),
			NtsTac:                nullable.NewString("E19A"),
			NtsSac:                nullable.NewString("987654321"),
			OrdersAcknowledgement: &unacknowledgedOrders,
		}

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(amendedOrder.ID.String()),
			IfMatch:     etag.GenerateEtag(amendedOrder.UpdatedAt),
			Body:        body,
		}

		orderUpdater := orderservice.NewOrderUpdater(moveRouter)
		handler := UpdateOrderHandler{
			handlerConfig,
			orderUpdater,
			moveTaskOrderUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.UpdateOrderOK{}, response)
		payload := response.(*orderop.UpdateOrderOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		var moveInDB models.Move
		err := suite.DB().Find(&moveInDB, approvalsRequestedMove.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, moveInDB.Status)
	})

	suite.Run("Does not update move status if move status is not APPROVALS_REQUESTED", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		move := subtestData.move
		order := subtestData.move.Orders
		destinationDutyLocation := order.NewDutyLocation
		originDutyStation := order.OriginDutyLocation

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		unacknowledgedOrders := false
		body := &ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:   &deptIndicator,
			IssueDate:             handlers.FmtDatePtr(&issueDate),
			ReportByDate:          handlers.FmtDatePtr(&reportByDate),
			OrdersType:            ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
			OrdersTypeDetail:      &ordersTypeDetail,
			OrdersNumber:          handlers.FmtString("ORDER100"),
			NewDutyLocationID:     handlers.FmtUUID(destinationDutyLocation.ID),
			OriginDutyLocationID:  handlers.FmtUUID(originDutyStation.ID),
			Tac:                   handlers.FmtString("E19A"),
			Sac:                   nullable.NewString("987654321"),
			NtsTac:                nullable.NewString("E19A"),
			NtsSac:                nullable.NewString("987654321"),
			OrdersAcknowledgement: &unacknowledgedOrders,
		}

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		orderUpdater := orderservice.NewOrderUpdater(moveRouter)
		handler := UpdateOrderHandler{
			handlerConfig,
			orderUpdater,
			moveTaskOrderUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.UpdateOrderOK{}, response)
		payload := response.(*orderop.UpdateOrderOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		var moveInDB models.Move
		err := suite.DB().Find(&moveInDB, move.ID)
		suite.NoError(err)
		suite.Equal(move.Status, moveInDB.Status)
	})
}

type updateOrderHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateOrderPayload
}

func (suite *HandlerSuite) makeUpdateOrderHandlerSubtestData() (subtestData *updateOrderHandlerSubtestData) {
	subtestData = &updateOrderHandlerSubtestData{}

	subtestData.move = testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{})
	subtestData.order = subtestData.move.Orders

	originDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
	destinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	subtestData.body = &ghcmessages.UpdateOrderPayload{
		DepartmentIndicator:  &deptIndicator,
		IssueDate:            handlers.FmtDatePtr(&issueDate),
		ReportByDate:         handlers.FmtDatePtr(&reportByDate),
		OrdersType:           ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
		OrdersTypeDetail:     &ordersTypeDetail,
		OrdersNumber:         handlers.FmtString("ORDER100"),
		NewDutyLocationID:    handlers.FmtUUID(destinationDutyLocation.ID),
		OriginDutyLocationID: handlers.FmtUUID(originDutyLocation.ID),
		Tac:                  handlers.FmtString("E19A"),
		Sac:                  nullable.NewString("987654321"),
		NtsTac:               nullable.NewString("E19A"),
		NtsSac:               nullable.NewString("987654321"),
	}

	return subtestData
}

func (suite *HandlerSuite) TestUpdateOrderHandler() {
	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		moveTaskOrderUpdater := mocks.MoveTaskOrderUpdater{}
		moveRouter := moverouter.NewMoveRouter()
		handler := UpdateOrderHandler{
			handlerConfig,
			orderservice.NewOrderUpdater(moveRouter),
			&moveTaskOrderUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.UpdateOrderOK{}, response)
		orderOK := response.(*orderop.UpdateOrderOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyLocationID.String(), ordersPayload.DestinationDutyLocation.ID.String())
		suite.Equal(body.OriginDutyLocationID.String(), ordersPayload.OriginDutyLocation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(*body.OrdersType, ordersPayload.OrderType)
		suite.Equal(body.OrdersTypeDetail, ordersPayload.OrderTypeDetail)
		suite.Equal(body.OrdersNumber, ordersPayload.OrderNumber)
		suite.Equal(body.DepartmentIndicator, ordersPayload.DepartmentIndicator)
		suite.Equal(body.Tac, ordersPayload.Tac)
		suite.Equal(body.Sac.Value, ordersPayload.Sac)
		suite.Equal(body.NtsTac.Value, ordersPayload.NtsTac)
		suite.Equal(body.NtsSac.Value, ordersPayload.NtsSac)
	})

	// We need to confirm whether a user who only has the TIO role should indeed
	// be authorized to update orders. If not, we also need to prevent them from
	// clicking the Edit Orders button in the frontend.
	suite.Run("Allows a TIO to update orders", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		move := subtestData.move
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTIO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			handlerConfig,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderOK{}, response)
		payload := response.(*orderop.UpdateOrderOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			handlerConfig,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.NotFoundError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderNotFound{}, response)
		payload := response.(*orderop.UpdateOrderNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			handlerConfig,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.PreconditionFailedError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderPreconditionFailed{}, response)
		payload := response.(*orderop.UpdateOrderPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			handlerConfig,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.InvalidInputError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderUnprocessableEntity{}, response)
		payload := response.(*orderop.UpdateOrderUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateOrderEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders

	body := &ghcmessages.UpdateOrderPayload{}

	requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
		factory.GetTraitOfficeUserStubbed,
	})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	request = request.WithContext(trace.NewContext(request.Context(), traceID))

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	updater := &mocks.OrderUpdater{}
	updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
		order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

	handlerConfig := suite.HandlerConfig()
	handler := UpdateOrderHandler{
		handlerConfig,
		updater,
		&mocks.MoveTaskOrderUpdater{},
	}

	// Validate incoming payload: not needed since we're mocking UpdateOrderAsTOO

	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)
	suite.IsType(&orderop.UpdateOrderOK{}, response)
	orderOK := response.(*orderop.UpdateOrderOK)
	ordersPayload := orderOK.Payload

	// Validate outgoing payload
	suite.NoError(ordersPayload.Validate(strfmt.Default))

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

type counselingUpdateOrderHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.CounselingUpdateOrderPayload
}

func (suite *HandlerSuite) makeCounselingUpdateOrderHandlerSubtestData() (subtestData *counselingUpdateOrderHandlerSubtestData) {
	subtestData = &counselingUpdateOrderHandlerSubtestData{}

	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	subtestData.move = testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
	subtestData.order = subtestData.move.Orders
	originDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())
	destinationDutyLocation := testdatagen.MakeDefaultDutyLocation(suite.DB())

	subtestData.body = &ghcmessages.CounselingUpdateOrderPayload{
		IssueDate:            handlers.FmtDatePtr(&issueDate),
		ReportByDate:         handlers.FmtDatePtr(&reportByDate),
		OrdersType:           ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
		NewDutyLocationID:    handlers.FmtUUID(destinationDutyLocation.ID),
		OriginDutyLocationID: handlers.FmtUUID(originDutyLocation.ID),
		Tac:                  handlers.FmtString("E19A"),
		Sac:                  nullable.NewString("987654321"),
		NtsTac:               nullable.NewString("E19A"),
		NtsSac:               nullable.NewString("987654321"),
	}

	return subtestData
}

func (suite *HandlerSuite) TestCounselingUpdateOrderHandler() {
	request := httptest.NewRequest("PATCH", "/counseling/orders/{orderID}", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		moveRouter := moverouter.NewMoveRouter()
		handler := CounselingUpdateOrderHandler{
			handlerConfig,
			orderservice.NewOrderUpdater(moveRouter),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.CounselingUpdateOrderOK{}, response)
		orderOK := response.(*orderop.CounselingUpdateOrderOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyLocationID.String(), ordersPayload.DestinationDutyLocation.ID.String())
		suite.Equal(body.OriginDutyLocationID.String(), ordersPayload.OriginDutyLocation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(*body.OrdersType, ordersPayload.OrderType)
		suite.Equal(body.Tac, ordersPayload.Tac)
		suite.Equal(body.Sac.Value, ordersPayload.Sac)
		suite.Equal(body.NtsTac.Value, ordersPayload.NtsTac)
		suite.Equal(body.NtsSac.Value, ordersPayload.NtsSac)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor, factory.GetTraitOfficeUserStubbed})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.NotFoundError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderNotFound{}, response)
		payload := response.(*orderop.CounselingUpdateOrderNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor, factory.GetTraitOfficeUserStubbed})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.PreconditionFailedError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderPreconditionFailed{}, response)
		payload := response.(*orderop.CounselingUpdateOrderPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor, factory.GetTraitOfficeUserStubbed})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.InvalidInputError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderUnprocessableEntity{}, response)
		payload := response.(*orderop.CounselingUpdateOrderUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

type updateAllowanceHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateAllowancePayload
}

type updateBillableWeightHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateBillableWeightPayload
}

type updateMaxBillableWeightAsTIOHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateMaxBillableWeightAsTIOPayload
}

func (suite *HandlerSuite) makeUpdateAllowanceHandlerSubtestData() (subtestData *updateAllowanceHandlerSubtestData) {
	subtestData = &updateAllowanceHandlerSubtestData{}

	subtestData.move = testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{})
	subtestData.order = subtestData.move.Orders

	newAuthorizedWeight := int64(10000)
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.AffiliationAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	subtestData.body = &ghcmessages.UpdateAllowancePayload{
		Agency:               &affiliation,
		AuthorizedWeight:     &newAuthorizedWeight,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		OrganizationalClothingAndIndividualEquipment: &ocie,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		RequiredMedicalEquipmentWeight: rmeWeight,
		StorageInTransit:               swag.Int64(60),
	}
	return subtestData
}

func (suite *HandlerSuite) makeUpdateMaxBillableWeightAsTIOHandlerSubtestData() (subtestData *updateMaxBillableWeightAsTIOHandlerSubtestData) {
	subtestData = &updateMaxBillableWeightAsTIOHandlerSubtestData{}
	now := time.Now()
	subtestData.move = testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{ExcessWeightQualifiedAt: &now},
	})
	subtestData.order = subtestData.move.Orders

	newAuthorizedWeight := int64(10000)
	newRemarks := "TIO remarks"

	subtestData.body = &ghcmessages.UpdateMaxBillableWeightAsTIOPayload{
		AuthorizedWeight: &newAuthorizedWeight,
		TioRemarks:       &newRemarks,
	}
	return subtestData
}

func (suite *HandlerSuite) makeUpdateBillableWeightHandlerSubtestData() (subtestData *updateBillableWeightHandlerSubtestData) {
	subtestData = &updateBillableWeightHandlerSubtestData{}
	now := time.Now()
	subtestData.move = testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{ExcessWeightQualifiedAt: &now},
	})
	subtestData.order = subtestData.move.Orders

	newAuthorizedWeight := int64(10000)

	subtestData.body = &ghcmessages.UpdateBillableWeightPayload{
		AuthorizedWeight: &newAuthorizedWeight,
	}
	return subtestData
}

func (suite *HandlerSuite) TestUpdateAllowanceHandler() {
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		moveRouter := moverouter.NewMoveRouter()
		handler := UpdateAllowanceHandler{
			handlerConfig,
			orderservice.NewOrderUpdater(moveRouter),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.UpdateAllowanceOK{}, response)
		orderOK := response.(*orderop.UpdateAllowanceOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*body.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*body.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(*body.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(*body.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
		suite.Equal(*body.StorageInTransit, *ordersPayload.Entitlement.StorageInTransit)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.NotFoundError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceNotFound{}, response)
		payload := response.(*orderop.UpdateAllowanceNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.PreconditionFailedError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowancePreconditionFailed{}, response)
		payload := response.(*orderop.UpdateAllowancePreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.InvalidInputError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceUnprocessableEntity{}, response)
		payload := response.(*orderop.UpdateAllowanceUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateAllowanceEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders

	body := &ghcmessages.UpdateAllowancePayload{}

	requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
		factory.GetTraitOfficeUserStubbed,
	})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	request = request.WithContext(trace.NewContext(request.Context(), traceID))

	params := orderop.UpdateAllowanceParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	updater := &mocks.OrderUpdater{}
	updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
		order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

	handlerConfig := suite.HandlerConfig()
	handler := UpdateAllowanceHandler{
		handlerConfig,
		updater,
	}

	// Validate incoming payload
	suite.NoError(params.Body.Validate(strfmt.Default))

	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)
	suite.IsType(&orderop.UpdateAllowanceOK{}, response)
	orderOK := response.(*orderop.UpdateAllowanceOK)
	ordersPayload := orderOK.Payload

	// Validate outgoing payload
	suite.NoError(ordersPayload.Validate(strfmt.Default))

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

func (suite *HandlerSuite) TestCounselingUpdateAllowanceHandler() {
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.AffiliationAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	body := &ghcmessages.CounselingUpdateAllowancePayload{
		Agency:               &affiliation,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		OrganizationalClothingAndIndividualEquipment: &ocie,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		RequiredMedicalEquipmentWeight: rmeWeight,
		StorageInTransit:               swag.Int64(80),
	}

	request := httptest.NewRequest("PATCH", "/counseling/orders/{orderID}/allowances", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		moveRouter := moverouter.NewMoveRouter()
		handler := CounselingUpdateAllowanceHandler{
			handlerConfig,
			orderservice.NewOrderUpdater(moveRouter),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.CounselingUpdateAllowanceOK{}, response)
		orderOK := response.(*orderop.CounselingUpdateAllowanceOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*body.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*body.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(*body.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(*body.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
		suite.Equal(*body.StorageInTransit, *ordersPayload.Entitlement.StorageInTransit)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor, factory.GetTraitOfficeUserStubbed})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.NotFoundError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceNotFound{}, response)
		payload := response.(*orderop.CounselingUpdateAllowanceNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor, factory.GetTraitOfficeUserStubbed})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.PreconditionFailedError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowancePreconditionFailed{}, response)
		payload := response.(*orderop.CounselingUpdateAllowancePreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{factory.GetTraitOfficeUserServicesCounselor, factory.GetTraitOfficeUserStubbed})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			handlerConfig,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, apperror.InvalidInputError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceUnprocessableEntity{}, response)
		payload := response.(*orderop.CounselingUpdateAllowanceUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestUpdateMaxBillableWeightAsTIOHandler() {
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/update-max-billable-weight/tio", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateMaxBillableWeightAsTIOHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateMaxBillableWeightAsTIOParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		router := moverouter.NewMoveRouter()
		handler := UpdateMaxBillableWeightAsTIOHandler{
			handlerConfig,
			orderservice.NewExcessWeightRiskManager(router),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.UpdateMaxBillableWeightAsTIOOK{}, response)
		orderOK := response.(*orderop.UpdateMaxBillableWeightAsTIOOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateMaxBillableWeightAsTIOHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTIO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateMaxBillableWeightAsTIOParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := UpdateMaxBillableWeightAsTIOHandler{
			handlerConfig,
			updater,
		}
		dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))
		tioRemarks := params.Body.TioRemarks

		updater.On("UpdateMaxBillableWeightAsTIO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, dbAuthorizedWeight, tioRemarks, params.IfMatch).Return(nil, nil, apperror.NotFoundError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateMaxBillableWeightAsTIONotFound{}, response)
		payload := response.(*orderop.UpdateMaxBillableWeightAsTIONotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateMaxBillableWeightAsTIOHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTIO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateMaxBillableWeightAsTIOParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := UpdateMaxBillableWeightAsTIOHandler{
			handlerConfig,
			updater,
		}
		dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))
		tioRemarks := params.Body.TioRemarks

		updater.On("UpdateMaxBillableWeightAsTIO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, dbAuthorizedWeight, tioRemarks, params.IfMatch).Return(nil, nil, apperror.PreconditionFailedError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateMaxBillableWeightAsTIOPreconditionFailed{}, response)
		payload := response.(*orderop.UpdateMaxBillableWeightAsTIOPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateMaxBillableWeightAsTIOHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTIO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateMaxBillableWeightAsTIOParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := UpdateMaxBillableWeightAsTIOHandler{
			handlerConfig,
			updater,
		}
		dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))
		tioRemarks := params.Body.TioRemarks

		verrs := validate.NewErrors()
		verrs.Add("some key", "some validation error")
		invalidInputError := apperror.NewInvalidInputError(order.ID, nil, verrs, "")
		updater.On("UpdateMaxBillableWeightAsTIO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, dbAuthorizedWeight, tioRemarks, params.IfMatch).Return(nil, nil, invalidInputError)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateMaxBillableWeightAsTIOUnprocessableEntity{}, response)
		payload := response.(*orderop.UpdateMaxBillableWeightAsTIOUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestUpdateBillableWeightHandler() {
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/update-billable-weight", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateBillableWeightHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateBillableWeightParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		router := moverouter.NewMoveRouter()
		handler := UpdateBillableWeightHandler{
			handlerConfig,
			orderservice.NewExcessWeightRiskManager(router),
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.UpdateBillableWeightOK{}, response)
		orderOK := response.(*orderop.UpdateBillableWeightOK)
		ordersPayload := orderOK.Payload

		// Validate outgoing payload
		suite.NoError(ordersPayload.Validate(strfmt.Default))

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateBillableWeightHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateBillableWeightParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := UpdateBillableWeightHandler{
			handlerConfig,
			updater,
		}
		dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))

		updater.On("UpdateBillableWeightAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, dbAuthorizedWeight, params.IfMatch).Return(nil, nil, apperror.NotFoundError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateBillableWeightNotFound{}, response)
		payload := response.(*orderop.UpdateBillableWeightNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateBillableWeightHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateBillableWeightParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := UpdateBillableWeightHandler{
			handlerConfig,
			updater,
		}
		dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))

		updater.On("UpdateBillableWeightAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, dbAuthorizedWeight, params.IfMatch).Return(nil, nil, apperror.PreconditionFailedError{})

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateBillableWeightPreconditionFailed{}, response)
		payload := response.(*orderop.UpdateBillableWeightPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		subtestData := suite.makeUpdateBillableWeightHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateBillableWeightParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := UpdateBillableWeightHandler{
			handlerConfig,
			updater,
		}
		dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))

		verrs := validate.NewErrors()
		verrs.Add("some key", "some validation error")
		invalidInputError := apperror.NewInvalidInputError(order.ID, nil, verrs, "")
		updater.On("UpdateBillableWeightAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, dbAuthorizedWeight, params.IfMatch).Return(nil, nil, invalidInputError)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateBillableWeightUnprocessableEntity{}, response)
		payload := response.(*orderop.UpdateBillableWeightUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

// Test that an order notification got stored successfully
func (suite *HandlerSuite) TestUpdateBillableWeightEventTrigger() {
	subtestData := suite.makeUpdateBillableWeightHandlerSubtestData()
	order := subtestData.order
	body := subtestData.body
	move := subtestData.move

	requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
		factory.GetTraitOfficeUserStubbed,
	})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/update-billable-weight", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	request = request.WithContext(trace.NewContext(request.Context(), traceID))

	params := orderop.UpdateBillableWeightParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}
	dbAuthorizedWeight := swag.Int(int(*params.Body.AuthorizedWeight))

	updater := &mocks.ExcessWeightRiskManager{}
	updater.On("UpdateBillableWeightAsTOO", mock.AnythingOfType("*appcontext.appContext"),
		order.ID, dbAuthorizedWeight, params.IfMatch).Return(&order, move.ID, nil)

	handlerConfig := suite.HandlerConfig()
	handler := UpdateBillableWeightHandler{
		handlerConfig,
		updater,
	}

	// Validate incoming payload
	suite.NoError(params.Body.Validate(strfmt.Default))

	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)
	suite.IsType(&orderop.UpdateBillableWeightOK{}, response)
	orderOK := response.(*orderop.UpdateBillableWeightOK)
	ordersPayload := orderOK.Payload

	// Validate outgoing payload
	suite.NoError(ordersPayload.Validate(strfmt.Default))

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

func (suite *HandlerSuite) TestAcknowledgeExcessWeightRiskHandler() {
	request := httptest.NewRequest("POST", "/orders/{orderID}/acknowledge-excess-weight-risk", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		handlerConfig := suite.HandlerConfig()
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
		order := move.Orders

		requestUser := factory.BuildOfficeUserWithRoles(nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.AcknowledgeExcessWeightRiskParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(move.UpdatedAt),
		}

		router := moverouter.NewMoveRouter()
		handler := AcknowledgeExcessWeightRiskHandler{
			handlerConfig,
			orderservice.NewExcessWeightRiskManager(router),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		suite.IsType(&orderop.AcknowledgeExcessWeightRiskOK{}, response)
		moveOK := response.(*orderop.AcknowledgeExcessWeightRiskOK)
		movePayload := moveOK.Payload

		// Validate outgoing payload
		suite.NoError(movePayload.Validate(strfmt.Default))

		suite.Equal(move.ID.String(), movePayload.ID.String())
		suite.NotNil(movePayload.ExcessWeightAcknowledgedAt)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		handlerConfig := suite.HandlerConfig()
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
		order := move.Orders

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.AcknowledgeExcessWeightRiskParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := AcknowledgeExcessWeightRiskHandler{
			handlerConfig,
			updater,
		}

		updater.On("AcknowledgeExcessWeightRisk", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, params.IfMatch).Return(nil, apperror.NotFoundError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&orderop.AcknowledgeExcessWeightRiskNotFound{}, response)
		payload := response.(*orderop.AcknowledgeExcessWeightRiskNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		handlerConfig := suite.HandlerConfig()
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
		order := move.Orders

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.AcknowledgeExcessWeightRiskParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := AcknowledgeExcessWeightRiskHandler{
			handlerConfig,
			updater,
		}

		updater.On("AcknowledgeExcessWeightRisk", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, params.IfMatch).Return(nil, apperror.PreconditionFailedError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&orderop.AcknowledgeExcessWeightRiskPreconditionFailed{}, response)
		payload := response.(*orderop.AcknowledgeExcessWeightRiskPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		handlerConfig := suite.HandlerConfig()
		now := time.Now()
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
			Move: models.Move{ExcessWeightQualifiedAt: &now},
		})
		order := move.Orders

		requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
			factory.GetTraitOfficeUserTOO,
			factory.GetTraitOfficeUserStubbed,
		})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.AcknowledgeExcessWeightRiskParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
		}

		updater := &mocks.ExcessWeightRiskManager{}
		handler := AcknowledgeExcessWeightRiskHandler{
			handlerConfig,
			updater,
		}

		verrs := validate.NewErrors()
		verrs.Add("some key", "some validation error")
		invalidInputError := apperror.NewInvalidInputError(order.ID, nil, verrs, "")
		updater.On("AcknowledgeExcessWeightRisk", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, params.IfMatch).Return(nil, invalidInputError)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)

		suite.IsType(&orderop.AcknowledgeExcessWeightRiskUnprocessableEntity{}, response)
		payload := response.(*orderop.AcknowledgeExcessWeightRiskUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

// Test that an order notification got stored successfully
func (suite *HandlerSuite) TestAcknowledgeExcessWeightRiskEventTrigger() {
	now := time.Now()
	move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{
		Move: models.Move{ExcessWeightQualifiedAt: &now},
	})
	order := move.Orders

	requestUser := factory.BuildOfficeUser(nil, nil, []factory.Trait{
		factory.GetTraitOfficeUserTOO,
		factory.GetTraitOfficeUserStubbed,
	})
	request := httptest.NewRequest("POST", "/orders/{orderID}/acknowledge-excess-weight-risk", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	request = request.WithContext(trace.NewContext(request.Context(), traceID))

	params := orderop.AcknowledgeExcessWeightRiskParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
	}

	updater := &mocks.ExcessWeightRiskManager{}
	updater.On("AcknowledgeExcessWeightRisk", mock.AnythingOfType("*appcontext.appContext"),
		order.ID, params.IfMatch).Return(&move, nil)

	handlerConfig := suite.HandlerConfig()
	handler := AcknowledgeExcessWeightRiskHandler{
		handlerConfig,
		updater,
	}

	// Validate incoming payload: no body to validate

	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)
	suite.IsType(&orderop.AcknowledgeExcessWeightRiskOK{}, response)
	moveOK := response.(*orderop.AcknowledgeExcessWeightRiskOK)
	movePayload := moveOK.Payload

	// Validate outgoing payload
	suite.NoError(movePayload.Validate(strfmt.Default))

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.Equal(movePayload.ID, strfmt.UUID(move.ID.String()))
	suite.HasWebhookNotification(move.ID, traceID)
}
