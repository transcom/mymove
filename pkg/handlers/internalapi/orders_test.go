package internalapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	ordersop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/orders"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/move"
	orderservice "github.com/transcom/mymove/pkg/services/order"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestCreateOrder() {
	sm := factory.BuildExtendedServiceMember(suite.DB(), nil, nil)
	dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())
	factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), dutyLocation.Address.PostalCode, "KKFA")
	factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

	req := httptest.NewRequest("POST", "/orders", nil)
	req = suite.AuthenticateRequest(req, sm)

	hasDependents := true
	spouseHasProGear := true
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	deptIndicator := internalmessages.DeptIndicatorAIRFORCE
	payload := &internalmessages.CreateUpdateOrders{
		HasDependents:       handlers.FmtBool(hasDependents),
		SpouseHasProGear:    handlers.FmtBool(spouseHasProGear),
		IssueDate:           handlers.FmtDate(issueDate),
		ReportByDate:        handlers.FmtDate(reportByDate),
		OrdersType:          internalmessages.NewOrdersType(ordersType),
		NewDutyLocationID:   handlers.FmtUUID(dutyLocation.ID),
		ServiceMemberID:     handlers.FmtUUID(sm.ID),
		OrdersNumber:        handlers.FmtString("123456"),
		Tac:                 handlers.FmtString("E19A"),
		Sac:                 handlers.FmtString("SacNumber"),
		DepartmentIndicator: internalmessages.NewDeptIndicator(deptIndicator),
	}

	params := ordersop.CreateOrdersParams{
		HTTPRequest:  req,
		CreateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	createHandler := CreateOrdersHandler{handlerConfig}

	response := createHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.CreateOrdersCreated{}, response)
	okResponse := response.(*ordersop.CreateOrdersCreated)
	orderID := okResponse.Payload.ID.String()
	createdOrder, _ := models.FetchOrder(suite.DB(), uuid.FromStringOrNil(orderID))

	suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Len(okResponse.Payload.Moves, 1)
	suite.Assertions.Equal(ordersType, *okResponse.Payload.OrdersType)
	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(handlers.FmtString("E19A"), okResponse.Payload.Tac)
	suite.Assertions.Equal(handlers.FmtString("SacNumber"), okResponse.Payload.Sac)
	suite.Assertions.Equal(&deptIndicator, okResponse.Payload.DepartmentIndicator)
	suite.Equal(sm.DutyLocationID, createdOrder.OriginDutyLocationID)
	suite.Equal((*string)(sm.Rank), createdOrder.Grade)
	suite.Assertions.Equal(*models.Int64Pointer(8000), *okResponse.Payload.AuthorizedWeight)
	suite.NotNil(&createdOrder.Entitlement)
	suite.NotEmpty(createdOrder.SupplyAndServicesCostEstimate)
	suite.NotEmpty(createdOrder.PackingAndShippingInstructions)
	suite.NotEmpty(createdOrder.MethodOfPayment)
	suite.NotEmpty(createdOrder.NAICS)
}

func (suite *HandlerSuite) TestShowOrder() {
	dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
			LinkOnly: true,
		},
	}, nil)
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    dutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)
	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("GET", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.ShowOrdersParams{
		HTTPRequest: req,
		OrdersID:    *handlers.FmtUUID(order.ID),
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)
	showHandler := ShowOrdersHandler{handlerConfig}

	response := showHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.ShowOrdersOK{}, response)
	okResponse := response.(*ordersop.ShowOrdersOK)

	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String())
	suite.Assertions.Equal(order.OrdersType, *okResponse.Payload.OrdersType)
	suite.Assertions.Equal(order.OrdersTypeDetail, okResponse.Payload.OrdersTypeDetail)
	suite.Assertions.Equal(*order.Grade, *okResponse.Payload.Grade)
	suite.Assertions.Equal(*order.TAC, *okResponse.Payload.Tac)
	suite.Assertions.Equal(*order.DepartmentIndicator, string(*okResponse.Payload.DepartmentIndicator))
	//suite.Assertions.Equal(order.IssueDate.String(), okResponse.Payload.IssueDate.String()) // TODO: get date formats aligned
	//suite.Assertions.Equal(order.ReportByDate.String(), okResponse.Payload.ReportByDate.String())
	suite.Assertions.Equal(order.HasDependents, *okResponse.Payload.HasDependents)
	suite.Assertions.Equal(order.SpouseHasProGear, *okResponse.Payload.SpouseHasProGear)
}

func setUpMockOrders() models.Order {
	orders := factory.BuildOrderWithoutDefaults(nil, nil, nil)

	orders.ID = uuid.Must(uuid.NewV4())

	orders.ServiceMemberID = uuid.Must(uuid.NewV4())
	orders.ServiceMember.ID = orders.ServiceMemberID
	orders.ServiceMember.UserID = uuid.Must(uuid.NewV4())
	orders.ServiceMember.User.ID = orders.ServiceMember.UserID

	return orders
}

func (suite *HandlerSuite) TestUploadAmendedOrdersHandlerUnit() {

	setUpRequestAndParams := func(orders models.Order) *ordersop.UploadAmendedOrdersParams {
		endpoint := fmt.Sprintf("/orders/%v/upload_amended_orders", orders.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		req = suite.AuthenticateRequest(req, orders.ServiceMember)

		params := ordersop.UploadAmendedOrdersParams{
			HTTPRequest: req,
			File:        suite.Fixture("filled-out-orders.pdf"),
			OrdersID:    *handlers.FmtUUID(orders.ID),
		}

		return &params
	}

	setUpOrOrderUpdater := func(returnValues ...interface{}) services.OrderUpdater {
		mockOrderUpdater := &mocks.OrderUpdater{}

		mockOrderUpdater.On(
			"UploadAmendedOrdersAsCustomer",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("uuid.UUID"),
			mock.AnythingOfType("*os.File"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("*test.FakeS3Storage"),
		).Return(returnValues...)

		return mockOrderUpdater
	}

	setUpHandler := func(orderUpdater services.OrderUpdater) UploadAmendedOrdersHandler {
		return UploadAmendedOrdersHandler{
			suite.createS3HandlerConfig(),
			orderUpdater,
		}
	}

	suite.Run("Returns a server error if there is an issue with the file type", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		params.File = factory.FixtureOpen("test.pdf")

		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, nil)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersInternalServerError{}, response)
	})

	suite.Run("Returns an error if the Orders ID in the URL is invalid", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		badUUID := "badUUID"
		params.HTTPRequest.URL.Path = fmt.Sprintf("/orders/%s/upload_amended_orders", badUUID)
		params.OrdersID = strfmt.UUID(badUUID)

		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, nil)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		if suite.IsType(&handlers.ErrResponse{}, response) {
			errResponse := response.(*handlers.ErrResponse)

			suite.Equal(http.StatusInternalServerError, errResponse.Code)
			suite.Contains(errResponse.Err.Error(), "incorrect UUID")
		}
	})

	suite.Run("Returns a 413 - Content Too Large if the file is too large", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		fakeErr := uploader.ErrTooLarge{
			FileSize:      uploader.MaxCustomerUserUploadFileSizeLimit + 1,
			FileSizeLimit: uploader.MaxCustomerUserUploadFileSizeLimit,
		}
		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, fakeErr)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersRequestEntityTooLarge{}, response)
	})

	suite.Run("Returns a server error if there is an error with the file", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		fakeErr := uploader.ErrFile{}

		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, fakeErr)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersInternalServerError{}, response)
	})

	suite.Run("Returns a server error if there is an error initializing the uploader", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		fakeErr := uploader.ErrFailedToInitUploader{}

		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, fakeErr)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersInternalServerError{}, response)
	})

	suite.Run("Returns a 404 if the order updater returns a NotFoundError", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		fakeErr := apperror.NotFoundError{}

		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, fakeErr)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersNotFound{}, response)
	})

	suite.Run("Returns a 500 if the order updater returns an unexpected error", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		fakeErr := apperror.NewBadDataError("Bad data")

		mockOrderUpdater := setUpOrOrderUpdater(models.Upload{}, "", nil, fakeErr)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		if suite.IsType(&handlers.ErrResponse{}, response) {
			errResponse := response.(*handlers.ErrResponse)

			suite.Equal(http.StatusInternalServerError, errResponse.Code)
			suite.Equal(fakeErr.Error(), errResponse.Err.Error())
		}
	})

	suite.Run("Returns a 201 if the amended orders are uploaded successfully", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		upload := factory.BuildUpload(suite.DB(), nil, nil)

		fakeURL := "https://fake.s3.url"
		mockOrderUpdater := setUpOrOrderUpdater(upload, fakeURL, nil, nil)

		handler := setUpHandler(mockOrderUpdater)

		response := handler.Handle(*params)

		if suite.IsType(&ordersop.UploadAmendedOrdersCreated{}, response) {
			payload := response.(*ordersop.UploadAmendedOrdersCreated).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.Equal(upload.ID.String(), payload.ID.String())
			suite.Equal(upload.ContentType, payload.ContentType)
			suite.Equal(upload.Filename, payload.Filename)
			suite.Equal(fakeURL, string(payload.URL))
		}
	})
}

func (suite *HandlerSuite) TestUploadAmendedOrdersHandlerIntegration() {
	orderUpdater := orderservice.NewOrderUpdater(move.NewMoveRouter())

	setUpRequestAndParams := func(orders models.Order) *ordersop.UploadAmendedOrdersParams {
		endpoint := fmt.Sprintf("/orders/%v/upload_amended_orders", orders.ID.String())

		req := httptest.NewRequest("PATCH", endpoint, nil)

		req = suite.AuthenticateRequest(req, orders.ServiceMember)

		params := ordersop.UploadAmendedOrdersParams{
			HTTPRequest: req,
			File:        suite.Fixture("filled-out-orders.pdf"),
			OrdersID:    *handlers.FmtUUID(orders.ID),
		}

		return &params
	}

	setUpHandler := func() UploadAmendedOrdersHandler {
		return UploadAmendedOrdersHandler{
			suite.createS3HandlerConfig(),
			orderUpdater,
		}
	}

	suite.Run("Returns a 404 if the orders aren't found", func() {
		orders := setUpMockOrders()

		params := setUpRequestAndParams(orders)

		handler := setUpHandler()

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersNotFound{}, response)
	})

	suite.Run("Returns a 400 - Bad Request if there is an issue with the file being uploaded", func() {
		orders := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

		params := setUpRequestAndParams(orders)
		params.File = suite.Fixture("empty.pdf")

		handler := setUpHandler()

		response := handler.Handle(*params)

		if suite.IsType(&handlers.ErrResponse{}, response) {
			errResponse := response.(*handlers.ErrResponse)

			suite.Equal(http.StatusBadRequest, errResponse.Code)
			suite.Equal(uploader.ErrZeroLengthFile.Error(), errResponse.Err.Error())
		}
	})

	suite.Run("Returns a 201 if the amended orders are uploaded successfully", func() {
		orders := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

		params := setUpRequestAndParams(orders)

		handler := setUpHandler()

		response := handler.Handle(*params)

		if suite.IsType(&ordersop.UploadAmendedOrdersCreated{}, response) {
			payload := response.(*ordersop.UploadAmendedOrdersCreated).Payload

			suite.NoError(payload.Validate(strfmt.Default))

			suite.NotEqual("", string(payload.ID))
			suite.Equal("filled-out-orders.pdf", payload.Filename)
			suite.Equal(uploader.FileTypePDF, payload.ContentType)
			suite.NotEqual("", string(payload.URL))
		}
	})
}

// TODO: Fix now that we capture transaction error. May be a data setup problem
/*
func (suite *HandlerSuite) TestUpdateOrder() {
	order := factory.BuildOrder(suite.DB(), nil, nil)

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("PUT", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	newOrdersTypeDetail := internalmessages.OrdersTypeDetailHHGPERMITTED
	departmentIndicator := internalmessages.DeptIndicatorAIRFORCE
	otherServiceMemberUUID := uuid.Must(uuid.NewV4())

	payload := &internalmessages.CreateUpdateOrders{
		OrdersNumber:        handlers.FmtString("123456"),
		HasDependents:       handlers.FmtBool(order.HasDependents),
		SpouseHasProGear:    handlers.FmtBool(order.SpouseHasProGear),
		IssueDate:           handlers.FmtDate(order.IssueDate),
		ReportByDate:        handlers.FmtDate(order.ReportByDate),
		OrdersType:          newOrdersType,
		OrdersTypeDetail:    &newOrdersTypeDetail,
		NewDutyLocationID:    handlers.FmtUUID(order.NewDutyLocationID),
		Tac:                 order.TAC,
		Sac:                 handlers.FmtString("N3TEST"),
		DepartmentIndicator: &departmentIndicator,
		// Attempt to assign to another service member
		ServiceMemberID: handlers.FmtUUID(otherServiceMemberUUID),
	}

	params := ordersop.UpdateOrdersParams{
		HTTPRequest:  req,
		OrdersID:     *handlers.FmtUUID(order.ID),
		UpdateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := handlers.NewHandlerCOnfig(suite.DB(), suite.TestLogger())
	handlerConfig.SetFileStorer(fakeS3)
	updateHandler := UpdateOrdersHandler{context}

	response := updateHandler.Handle(params)

	suite.Assertions.IsType(&ordersop.UpdateOrdersOK{}, response)
	okResponse := response.(*ordersop.UpdateOrdersOK)

	suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
	suite.Assertions.Equal(order.ServiceMember.ID.String(), okResponse.Payload.ServiceMemberID.String(), "service member id should not change")
	suite.Assertions.Equal(newOrdersType, okResponse.Payload.OrdersType)
	suite.Assertions.Equal(newOrdersTypeDetail, *okResponse.Payload.OrdersTypeDetail)
	suite.Assertions.Equal(handlers.FmtString("N3TEST"), okResponse.Payload.Sac)
}
*/
