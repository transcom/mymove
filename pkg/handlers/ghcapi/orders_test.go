package ghcapi

import (
	"net/http/httptest"
	"time"

	moverouter "github.com/transcom/mymove/pkg/services/move"

	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/uploader"

	"github.com/transcom/mymove/pkg/services"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	orderop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/order"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services/mocks"
	orderservice "github.com/transcom/mymove/pkg/services/order"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestGetOrderHandlerIntegration() {
	officeUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})

	move := testdatagen.MakeDefaultMove(suite.DB())
	order := move.Orders
	request := httptest.NewRequest("GET", "/orders/{orderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, officeUser)

	params := orderop.GetOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
	}
	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := GetOrdersHandler{
		context,
		orderservice.NewOrderFetcher(),
	}

	response := handler.Handle(params)
	suite.IsNotErrResponse(response)
	orderOK := response.(*orderop.GetOrderOK)
	ordersPayload := orderOK.Payload

	suite.Assertions.IsType(&orderop.GetOrderOK{}, response)
	suite.Equal(order.ID.String(), ordersPayload.ID.String())
	suite.Equal(move.Locator, ordersPayload.MoveCode)
	suite.Equal(order.ServiceMemberID.String(), ordersPayload.Customer.ID.String())
	suite.Equal(order.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
	suite.NotNil(order.NewDutyStation)
	payloadEntitlement := ordersPayload.Entitlement
	suite.Equal((*order.EntitlementID).String(), payloadEntitlement.ID.String())
	orderEntitlement := order.Entitlement
	suite.NotNil(orderEntitlement)
	suite.EqualValues(orderEntitlement.ProGearWeight, payloadEntitlement.ProGearWeight)
	suite.EqualValues(orderEntitlement.ProGearWeightSpouse, payloadEntitlement.ProGearWeightSpouse)
	suite.EqualValues(orderEntitlement.RequiredMedicalEquipmentWeight, payloadEntitlement.RequiredMedicalEquipmentWeight)
	suite.EqualValues(orderEntitlement.OrganizationalClothingAndIndividualEquipment, payloadEntitlement.OrganizationalClothingAndIndividualEquipment)
	suite.Equal(order.OriginDutyStation.ID.String(), ordersPayload.OriginDutyStation.ID.String())
	suite.NotZero(order.OriginDutyStation)
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

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetOrdersHandler{
			context,
			&orderFetcher,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*orderop.GetOrderOK)
		orderPayload := orderOK.Payload
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

		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handler := GetOrdersHandler{
			context,
			&orderFetcher,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*orderop.GetOrderOK)
		orderPayload := orderOK.Payload
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
	handlerContext         handlers.HandlerContext
	userUploader           *uploader.UserUploader
	amendedOrder           models.Order
	approvalsRequestedMove models.Move
	originDutyStation      models.DutyStation
	destinationDutyStation models.DutyStation
}

func (suite *HandlerSuite) makeUpdateOrderHandlerAmendedUploadSubtestData() (subtestData *updateOrderHandlerAmendedUploadSubtestData) {
	subtestData = &updateOrderHandlerAmendedUploadSubtestData{}
	subtestData.handlerContext = suite.createHandlerContext()

	var err error
	subtestData.userUploader, err = uploader.NewUserUploader(subtestData.handlerContext.FileStorer(), 100*uploader.MB)
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

	subtestData.originDutyStation = testdatagen.MakeDefaultDutyStation(suite.DB())
	subtestData.destinationDutyStation = testdatagen.MakeDefaultDutyStation(suite.DB())

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
		context := subtestData.handlerContext
		userUploader := subtestData.userUploader
		destinationDutyStation := subtestData.destinationDutyStation
		originDutyStation := subtestData.originDutyStation

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

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		body := &ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:   &deptIndicator,
			IssueDate:             handlers.FmtDatePtr(&issueDate),
			ReportByDate:          handlers.FmtDatePtr(&reportByDate),
			OrdersType:            ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
			OrdersTypeDetail:      &ordersTypeDetail,
			OrdersNumber:          handlers.FmtString("ORDER100"),
			NewDutyStationID:      handlers.FmtUUID(destinationDutyStation.ID),
			OriginDutyStationID:   handlers.FmtUUID(originDutyStation.ID),
			Tac:                   handlers.FmtString("E19A"),
			Sac:                   handlers.FmtString("987654321"),
			OrdersAcknowledgement: &ordersAcknowledgement,
		}

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := UpdateOrderHandler{
			context,
			orderservice.NewOrderUpdater(),
			moveTaskOrderUpdater,
		}

		suite.Nil(order.AmendedOrdersAcknowledgedAt)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, move.Status)

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&orderop.UpdateOrderOK{}, response)
		orderOK := response.(*orderop.UpdateOrderOK)
		ordersPayload := orderOK.Payload

		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
		suite.Equal(body.OriginDutyStationID.String(), ordersPayload.OriginDutyStation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(*body.OrdersType, ordersPayload.OrderType)
		suite.Equal(body.OrdersTypeDetail, ordersPayload.OrderTypeDetail)
		suite.Equal(body.OrdersNumber, ordersPayload.OrderNumber)
		suite.Equal(body.DepartmentIndicator, ordersPayload.DepartmentIndicator)
		suite.Equal(body.Tac, ordersPayload.Tac)
		suite.Equal(body.Sac, ordersPayload.Sac)
		suite.NotNil(ordersPayload.AmendedOrdersAcknowledgedAt)

		reloadErr := suite.DB().Reload(&move)
		suite.NoError(reloadErr, "error reloading move of amended orders")

		suite.Equal(models.MoveStatusAPPROVED, move.Status)
	})

	suite.Run("Does not update move status if orders are not acknowledged", func() {
		subtestData := suite.makeUpdateOrderHandlerAmendedUploadSubtestData()
		context := subtestData.handlerContext
		destinationDutyStation := subtestData.destinationDutyStation
		originDutyStation := subtestData.originDutyStation
		amendedOrder := subtestData.amendedOrder
		approvalsRequestedMove := subtestData.approvalsRequestedMove

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		unacknowledgedOrders := false
		body := &ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:   &deptIndicator,
			IssueDate:             handlers.FmtDatePtr(&issueDate),
			ReportByDate:          handlers.FmtDatePtr(&reportByDate),
			OrdersType:            ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
			OrdersTypeDetail:      &ordersTypeDetail,
			OrdersNumber:          handlers.FmtString("ORDER100"),
			NewDutyStationID:      handlers.FmtUUID(destinationDutyStation.ID),
			OriginDutyStationID:   handlers.FmtUUID(originDutyStation.ID),
			Tac:                   handlers.FmtString("E19A"),
			Sac:                   handlers.FmtString("987654321"),
			OrdersAcknowledgement: &unacknowledgedOrders,
		}

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(amendedOrder.ID.String()),
			IfMatch:     etag.GenerateEtag(amendedOrder.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		orderUpdater := mocks.OrderUpdater{}
		// This is not modified but we're relying on the check of the params short circuiting
		orderUpdater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			amendedOrder.ID, *body, params.IfMatch).Return(&amendedOrder, approvalsRequestedMove.ID, nil)

		moveUpdater := mocks.MoveTaskOrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			&orderUpdater,
			&moveUpdater,
		}

		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		suite.Assertions.IsType(&orderop.UpdateOrderOK{}, response)

		suite.True(moveUpdater.AssertNotCalled(suite.T(), "UpdateApprovedAmendedOrders"))
	})

	suite.Run("Returns a 409 conflict error if move status is in invalid state", func() {
		subtestData := suite.makeUpdateOrderHandlerAmendedUploadSubtestData()
		userUploader := subtestData.userUploader
		context := subtestData.handlerContext
		destinationDutyStation := subtestData.destinationDutyStation
		originDutyStation := subtestData.originDutyStation

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

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
		move := testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{
			Order: models.Order{
				UploadedAmendedOrders:   &document,
				UploadedAmendedOrdersID: &document.ID,
				ServiceMember:           document.ServiceMember,
				ServiceMemberID:         document.ServiceMemberID,
			},
		})

		order := move.Orders

		body := &ghcmessages.UpdateOrderPayload{
			DepartmentIndicator:   &deptIndicator,
			IssueDate:             handlers.FmtDatePtr(&issueDate),
			ReportByDate:          handlers.FmtDatePtr(&reportByDate),
			OrdersType:            ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
			OrdersTypeDetail:      &ordersTypeDetail,
			OrdersNumber:          handlers.FmtString("ORDER100"),
			NewDutyStationID:      handlers.FmtUUID(destinationDutyStation.ID),
			OriginDutyStationID:   handlers.FmtUUID(originDutyStation.ID),
			Tac:                   handlers.FmtString("E19A"),
			Sac:                   handlers.FmtString("987654321"),
			OrdersAcknowledgement: &ordersAcknowledgement,
		}

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		orderUpdater := mocks.OrderUpdater{}
		acknowledgedAt := time.Now()
		order.AmendedOrdersAcknowledgedAt = &acknowledgedAt
		orderUpdater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *body, params.IfMatch).Return(&order, move.ID, nil)

		moveUpdater := mocks.MoveTaskOrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			&orderUpdater,
			&moveUpdater,
		}

		response := handler.Handle(params)

		suite.Assertions.IsType(&orderop.UpdateOrderConflict{}, response)
		conflictErr := response.(*orderop.UpdateOrderConflict)

		suite.Contains(*conflictErr.Payload.Message, "Cannot approve move with amended orders because the move status is not APPROVALS REQUESTED")
	})
}

type updateOrderHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateOrderPayload
}

func (suite *HandlerSuite) makeUpdateOrderHandlerSubtestData() (subtestData *updateOrderHandlerSubtestData) {
	subtestData = &updateOrderHandlerSubtestData{}

	subtestData.move = testdatagen.MakeDefaultMove(suite.DB())
	subtestData.order = subtestData.move.Orders

	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	issueDate, _ := time.Parse("2006-01-02", "2020-08-01")
	reportByDate, _ := time.Parse("2006-01-02", "2020-10-31")
	deptIndicator := ghcmessages.DeptIndicatorCOASTGUARD
	ordersTypeDetail := ghcmessages.OrdersTypeDetail("INSTRUCTION_20_WEEKS")
	subtestData.body = &ghcmessages.UpdateOrderPayload{
		DepartmentIndicator: &deptIndicator,
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
		OrdersTypeDetail:    &ordersTypeDetail,
		OrdersNumber:        handlers.FmtString("ORDER100"),
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
		Tac:                 handlers.FmtString("E19A"),
		Sac:                 handlers.FmtString("987654321"),
	}

	return subtestData
}

func (suite *HandlerSuite) TestUpdateOrderHandler() {
	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		moveTaskOrderUpdater := mocks.MoveTaskOrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			orderservice.NewOrderUpdater(),
			&moveTaskOrderUpdater,
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)

		orderOK := response.(*orderop.UpdateOrderOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.UpdateOrderOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
		suite.Equal(body.OriginDutyStationID.String(), ordersPayload.OriginDutyStation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(*body.OrdersType, ordersPayload.OrderType)
		suite.Equal(body.OrdersTypeDetail, ordersPayload.OrderTypeDetail)
		suite.Equal(body.OrdersNumber, ordersPayload.OrderNumber)
		suite.Equal(body.DepartmentIndicator, ordersPayload.DepartmentIndicator)
		suite.Equal(body.Tac, ordersPayload.Tac)
		suite.Equal(body.Sac, ordersPayload.Sac)
	})

	suite.Run("Returns a 403 when the user does not have TXO role", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderForbidden{}, response)
	})

	// We need to confirm whether a user who only has the TIO role should indeed
	// be authorized to update orders. If not, we also need to prevent them from
	// clicking the Edit Orders button in the frontend.
	suite.Run("Allows a TIO to update orders", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		move := subtestData.move
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTIOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)
		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderOK{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateOrderHandler{
			context,
			updater,
			&mocks.MoveTaskOrderUpdater{},
		}

		updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateOrderUnprocessableEntity{}, response)
	})
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateOrderEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders

	body := &ghcmessages.UpdateOrderPayload{}

	requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	params := orderop.UpdateOrderParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	updater := &mocks.OrderUpdater{}
	updater.On("UpdateOrderAsTOO", mock.AnythingOfType("*appcontext.appContext"),
		order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateOrderHandler{
		context,
		updater,
		&mocks.MoveTaskOrderUpdater{},
	}

	traceID, err := uuid.NewV4()
	handler.SetTraceID(traceID)        // traceID is inserted into handler
	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)

	orderOK := response.(*orderop.UpdateOrderOK)
	ordersPayload := orderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.IsType(&orderop.UpdateOrderOK{}, response)
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
	originDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())
	destinationDutyStation := testdatagen.MakeDefaultDutyStation(suite.DB())

	subtestData.body = &ghcmessages.CounselingUpdateOrderPayload{
		IssueDate:           handlers.FmtDatePtr(&issueDate),
		ReportByDate:        handlers.FmtDatePtr(&reportByDate),
		OrdersType:          ghcmessages.NewOrdersType(ghcmessages.OrdersTypeRETIREMENT),
		NewDutyStationID:    handlers.FmtUUID(destinationDutyStation.ID),
		OriginDutyStationID: handlers.FmtUUID(originDutyStation.ID),
	}

	return subtestData
}

func (suite *HandlerSuite) TestCounselingUpdateOrderHandler() {
	request := httptest.NewRequest("PATCH", "/counseling/orders/{orderID}", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := CounselingUpdateOrderHandler{
			context,
			orderservice.NewOrderUpdater(),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.CounselingUpdateOrderOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.CounselingUpdateOrderOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.NewDutyStationID.String(), ordersPayload.DestinationDutyStation.ID.String())
		suite.Equal(body.OriginDutyStationID.String(), ordersPayload.OriginDutyStation.ID.String())
		suite.Equal(*body.IssueDate, ordersPayload.DateIssued)
		suite.Equal(*body.ReportByDate, ordersPayload.ReportByDate)
		suite.Equal(*body.OrdersType, ordersPayload.OrderType)
	})

	suite.Run("Returns a 403 when the user does not have Counselor role", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateOrderAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderForbidden{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderPreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeCounselingUpdateOrderHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateOrderParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateOrderHandler{
			context,
			updater,
		}

		updater.On("UpdateOrderAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateOrderUnprocessableEntity{}, response)
	})
}

type updateAllowanceHandlerSubtestData struct {
	move  models.Move
	order models.Order
	body  *ghcmessages.UpdateAllowancePayload
}

func (suite *HandlerSuite) makeUpdateAllowanceHandlerSubtestData() (subtestData *updateAllowanceHandlerSubtestData) {
	subtestData = &updateAllowanceHandlerSubtestData{}

	subtestData.move = testdatagen.MakeServiceCounselingCompletedMove(suite.DB(), testdatagen.Assertions{})
	subtestData.order = subtestData.move.Orders

	newAuthorizedWeight := int64(10000)
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.BranchAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	subtestData.body = &ghcmessages.UpdateAllowancePayload{
		Agency:               affiliation,
		AuthorizedWeight:     &newAuthorizedWeight,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		OrganizationalClothingAndIndividualEquipment: &ocie,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		RequiredMedicalEquipmentWeight: rmeWeight,
	}
	return subtestData
}

func (suite *HandlerSuite) TestUpdateAllowanceHandler() {
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := UpdateAllowanceHandler{
			context,
			orderservice.NewOrderUpdater(),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.UpdateAllowanceOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.UpdateAllowanceOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.AuthorizedWeight, ordersPayload.Entitlement.AuthorizedWeight)
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*body.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*body.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(*body.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(*body.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
	})

	suite.Run("Returns a 403 when the user does not have TOO role", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceForbidden{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowancePreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		subtestData := suite.makeUpdateAllowanceHandlerSubtestData()
		order := subtestData.order
		body := subtestData.body

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.UpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := UpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.UpdateAllowanceUnprocessableEntity{}, response)
	})
}

// Test that an order notification got stored Successfully
func (suite *HandlerSuite) TestUpdateAllowanceEventTrigger() {
	move := testdatagen.MakeAvailableMove(suite.DB())
	order := move.Orders

	body := &ghcmessages.UpdateAllowancePayload{}

	requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
	request := httptest.NewRequest("PATCH", "/orders/{orderID}/allowances", nil)
	request = suite.AuthenticateOfficeRequest(request, requestUser)

	params := orderop.UpdateAllowanceParams{
		HTTPRequest: request,
		OrderID:     strfmt.UUID(order.ID.String()),
		IfMatch:     etag.GenerateEtag(order.UpdatedAt), // This is broken if you get a preconditioned failed error
		Body:        body,
	}

	updater := &mocks.OrderUpdater{}
	updater.On("UpdateAllowanceAsTOO", mock.AnythingOfType("*appcontext.appContext"),
		order.ID, *params.Body, params.IfMatch).Return(&order, move.ID, nil)

	context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
	handler := UpdateAllowanceHandler{
		context,
		updater,
	}

	traceID, err := uuid.NewV4()
	handler.SetTraceID(traceID)        // traceID is inserted into handler
	response := handler.Handle(params) // This step also saves traceID into DB

	suite.IsNotErrResponse(response)

	orderOK := response.(*orderop.UpdateAllowanceOK)
	ordersPayload := orderOK.Payload

	suite.FatalNoError(err, "Error creating a new trace ID.")
	suite.IsType(&orderop.UpdateAllowanceOK{}, response)
	suite.Equal(ordersPayload.ID, strfmt.UUID(order.ID.String()))
	suite.HasWebhookNotification(order.ID, traceID)
}

func (suite *HandlerSuite) TestCounselingUpdateAllowanceHandler() {
	grade := ghcmessages.GradeO5
	affiliation := ghcmessages.BranchAIRFORCE
	ocie := false
	proGearWeight := swag.Int64(100)
	proGearWeightSpouse := swag.Int64(10)
	rmeWeight := swag.Int64(10000)

	body := &ghcmessages.CounselingUpdateAllowancePayload{
		Agency:               affiliation,
		DependentsAuthorized: swag.Bool(true),
		Grade:                &grade,
		OrganizationalClothingAndIndividualEquipment: &ocie,
		ProGearWeight:                  proGearWeight,
		ProGearWeightSpouse:            proGearWeightSpouse,
		RequiredMedicalEquipmentWeight: rmeWeight,
	}

	request := httptest.NewRequest("PATCH", "/counseling/orders/{orderID}/allowances", nil)

	suite.Run("Returns 200 when all validations pass", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		handler := CounselingUpdateAllowanceHandler{
			context,
			orderservice.NewOrderUpdater(),
		}
		response := handler.Handle(params)

		suite.IsNotErrResponse(response)
		orderOK := response.(*orderop.CounselingUpdateAllowanceOK)
		ordersPayload := orderOK.Payload

		suite.Assertions.IsType(&orderop.CounselingUpdateAllowanceOK{}, response)
		suite.Equal(order.ID.String(), ordersPayload.ID.String())
		suite.Equal(body.Grade, ordersPayload.Grade)
		suite.Equal(body.Agency, ordersPayload.Agency)
		suite.Equal(body.DependentsAuthorized, ordersPayload.Entitlement.DependentsAuthorized)
		suite.Equal(*body.OrganizationalClothingAndIndividualEquipment, ordersPayload.Entitlement.OrganizationalClothingAndIndividualEquipment)
		suite.Equal(*body.ProGearWeight, ordersPayload.Entitlement.ProGearWeight)
		suite.Equal(*body.ProGearWeightSpouse, ordersPayload.Entitlement.ProGearWeightSpouse)
		suite.Equal(*body.RequiredMedicalEquipmentWeight, ordersPayload.Entitlement.RequiredMedicalEquipmentWeight)
	})

	suite.Run("Returns a 403 when the user does not have Counselor role", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeTOOOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsTOO", 0)
		updater.AssertNumberOfCalls(suite.T(), "UpdateAllowanceAsCounselor", 0)

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceForbidden{}, response)
	})

	suite.Run("Returns 404 when updater returns NotFoundError", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.NotFoundError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceNotFound{}, response)
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     "",
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.PreconditionFailedError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowancePreconditionFailed{}, response)
	})

	suite.Run("Returns 422 when updater service returns validation errors", func() {
		context := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		move := testdatagen.MakeNeedsServiceCounselingMove(suite.DB())
		order := move.Orders

		requestUser := testdatagen.MakeServicesCounselorOfficeUser(suite.DB(), testdatagen.Assertions{Stub: true})
		request = suite.AuthenticateOfficeRequest(request, requestUser)

		params := orderop.CounselingUpdateAllowanceParams{
			HTTPRequest: request,
			OrderID:     strfmt.UUID(order.ID.String()),
			IfMatch:     etag.GenerateEtag(order.UpdatedAt),
			Body:        body,
		}

		updater := &mocks.OrderUpdater{}
		handler := CounselingUpdateAllowanceHandler{
			context,
			updater,
		}

		updater.On("UpdateAllowanceAsCounselor", mock.AnythingOfType("*appcontext.appContext"),
			order.ID, *params.Body, params.IfMatch).Return(nil, nil, services.InvalidInputError{})

		response := handler.Handle(params)

		suite.IsType(&orderop.CounselingUpdateAllowanceUnprocessableEntity{}, response)
	})
}
