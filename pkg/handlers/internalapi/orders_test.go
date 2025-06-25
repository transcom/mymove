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
	"github.com/transcom/mymove/pkg/services/entitlements"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/services/move"
	orderservice "github.com/transcom/mymove/pkg/services/order"
	transportationoffice "github.com/transcom/mymove/pkg/services/transportation_office"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

func (suite *HandlerSuite) TestCreateOrder() {
	customAffiliation := models.AffiliationARMY
	sm := factory.BuildExtendedServiceMember(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				Affiliation: &customAffiliation,
			},
		},
	}, nil)

	rank := factory.FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), customAffiliation.String())
	parameterName := "maxGunSafeAllowance"
	parameterValue := "500"
	param := models.ApplicationParameters{
		ParameterName:  &parameterName,
		ParameterValue: &parameterValue,
	}
	suite.MustSave(&param)

	suite.Run("can create conus orders", func() {
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(false),
				},
			},
		}, nil)

		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: factory.MakeRandomString(8),
				},
			},
			{
				Model:    address,
				LinkOnly: true,
			},
		}, nil)

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
		deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE

		payload := &internalmessages.CreateUpdateOrders{
			HasDependents:        handlers.FmtBool(hasDependents),
			SpouseHasProGear:     handlers.FmtBool(spouseHasProGear),
			IssueDate:            handlers.FmtDate(issueDate),
			ReportByDate:         handlers.FmtDate(reportByDate),
			OrdersType:           internalmessages.NewOrdersType(ordersType),
			OriginDutyLocationID: *handlers.FmtUUIDPtr(&originDutyLocation.ID),
			NewDutyLocationID:    handlers.FmtUUID(dutyLocation.ID),
			ServiceMemberID:      handlers.FmtUUID(sm.ID),
			OrdersNumber:         handlers.FmtString("123456"),
			Tac:                  handlers.FmtString("E19A"),
			Sac:                  handlers.FmtString("SacNumber"),
			DepartmentIndicator:  internalmessages.NewDeptIndicator(deptIndicator),
			Grade:                models.ServiceMemberGradeE1.Pointer(),
			Rank:                 strfmt.UUID(rank.ID.String()),
		}

		params := ordersop.CreateOrdersParams{
			HTTPRequest:  req,
			CreateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		createHandler := CreateOrdersHandler{handlerConfig}

		response := createHandler.Handle(params)

		suite.Assertions.IsType(&ordersop.CreateOrdersCreated{}, response)
		okResponse := response.(*ordersop.CreateOrdersCreated)
		orderID := okResponse.Payload.ID.String()
		createdOrder, _ := models.FetchOrder(suite.DB(), uuid.FromStringOrNil(orderID))
		var createdEntitlement models.Entitlement
		err := suite.DB().Find(&createdEntitlement, createdOrder.EntitlementID)
		suite.NoError(err)
		suite.NotEmpty(createdEntitlement)
		suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
		suite.Assertions.Len(okResponse.Payload.Moves, 1)
		suite.Assertions.Equal(ordersType, *okResponse.Payload.OrdersType)
		suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
		suite.Assertions.Equal(handlers.FmtString("E19A"), okResponse.Payload.Tac)
		suite.Assertions.Equal(handlers.FmtString("SacNumber"), okResponse.Payload.Sac)
		suite.Assertions.Equal(&deptIndicator, okResponse.Payload.DepartmentIndicator)
		suite.Assertions.Equal(*models.Int64Pointer(8000), *okResponse.Payload.AuthorizedWeight)
		suite.NotNil(&createdOrder.Entitlement)
		suite.NotEmpty(createdOrder.SupplyAndServicesCostEstimate)
		suite.NotEmpty(createdOrder.PackingAndShippingInstructions)
		suite.NotEmpty(createdOrder.MethodOfPayment)
		suite.NotEmpty(createdOrder.NAICS)
		suite.Nil(createdEntitlement.AccompaniedTour)
		suite.Nil(createdEntitlement.DependentsTwelveAndOver)
		suite.Nil(createdEntitlement.DependentsUnderTwelve)
		suite.Assertions.Equal(rank.RankAbbv, okResponse.Payload.Rank.RankAbbv)
		suite.Equal(createdEntitlement.GunSafeWeight, 500)
	})

	suite.Run("can create oconus orders", func() {
		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99801")
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &usprc.ID,
				},
			},
		}, nil)

		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name:      factory.MakeRandomString(8),
					AddressID: address.ID,
				},
			},
		}, nil)

		dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				ContractID: contract.ID,
				IsOconus:   true,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			},
		})
		suite.NotNil(rateArea)
		suite.Nil(err)

		us_country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		oconusRateArea, err := models.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
		suite.NotNil(oconusRateArea)
		suite.Nil(err)

		jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MAPK")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)
		req := httptest.NewRequest("POST", "/orders", nil)
		req = suite.AuthenticateRequest(req, sm)

		hasDependents := true
		spouseHasProGear := true
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE
		payload := &internalmessages.CreateUpdateOrders{
			HasDependents:           handlers.FmtBool(hasDependents),
			SpouseHasProGear:        handlers.FmtBool(spouseHasProGear),
			IssueDate:               handlers.FmtDate(issueDate),
			ReportByDate:            handlers.FmtDate(reportByDate),
			OrdersType:              internalmessages.NewOrdersType(ordersType),
			OriginDutyLocationID:    *handlers.FmtUUIDPtr(&originDutyLocation.ID),
			NewDutyLocationID:       handlers.FmtUUID(dutyLocation.ID),
			ServiceMemberID:         handlers.FmtUUID(sm.ID),
			OrdersNumber:            handlers.FmtString("123456"),
			Tac:                     handlers.FmtString("E19A"),
			Sac:                     handlers.FmtString("SacNumber"),
			DepartmentIndicator:     internalmessages.NewDeptIndicator(deptIndicator),
			Grade:                   models.ServiceMemberGradeE1.Pointer(),
			AccompaniedTour:         models.BoolPointer(true),
			DependentsTwelveAndOver: models.Int64Pointer(5),
			DependentsUnderTwelve:   models.Int64Pointer(5),
			CivilianTdyUbAllowance:  models.Int64Pointer(350),
			Rank:                    strfmt.UUID(rank.ID.String()),
		}

		params := ordersop.CreateOrdersParams{
			HTTPRequest:  req,
			CreateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		createHandler := CreateOrdersHandler{handlerConfig}
		response := createHandler.Handle(params)

		suite.Assertions.IsType(&ordersop.CreateOrdersCreated{}, response)
		okResponse := response.(*ordersop.CreateOrdersCreated)

		orderID := okResponse.Payload.ID.String()
		createdOrder, _ := models.FetchOrder(suite.DB(), uuid.FromStringOrNil(orderID))
		var createdEntitlement models.Entitlement
		err = suite.DB().Find(&createdEntitlement, createdOrder.EntitlementID)
		suite.NoError(err)
		suite.NotEmpty(createdEntitlement)
		suite.Assertions.Equal(sm.ID.String(), okResponse.Payload.ServiceMemberID.String())
		suite.Assertions.Len(okResponse.Payload.Moves, 1)
		suite.Assertions.Equal(ordersType, *okResponse.Payload.OrdersType)
		suite.Assertions.Equal(handlers.FmtString("123456"), okResponse.Payload.OrdersNumber)
		suite.Assertions.Equal(handlers.FmtString("MAPK"), okResponse.Payload.OriginDutyLocationGbloc)
		suite.Assertions.Equal(handlers.FmtString("E19A"), okResponse.Payload.Tac)
		suite.Assertions.Equal(handlers.FmtString("SacNumber"), okResponse.Payload.Sac)
		suite.Assertions.Equal(&deptIndicator, okResponse.Payload.DepartmentIndicator)
		suite.Assertions.Equal(*models.Int64Pointer(8000), *okResponse.Payload.AuthorizedWeight)
		suite.NotNil(&createdOrder.Entitlement)
		suite.NotEmpty(createdOrder.SupplyAndServicesCostEstimate)
		suite.NotEmpty(createdOrder.PackingAndShippingInstructions)
		suite.NotEmpty(createdOrder.MethodOfPayment)
		suite.NotEmpty(createdOrder.NAICS)
		suite.NotNil(createdEntitlement.AccompaniedTour)
		suite.NotNil(createdEntitlement.DependentsTwelveAndOver)
		suite.NotNil(createdEntitlement.DependentsUnderTwelve)
		suite.NotNil(createdEntitlement.UBAllowance)
	})

	suite.Run("properly handles entitlement validation", func() {
		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99506")
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &usprc.ID,
				},
			},
		}, nil)

		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name:      factory.MakeRandomString(8),
					AddressID: address.ID,
				},
			},
		}, nil)

		dutyLocation := factory.FetchOrBuildCurrentDutyLocation(suite.DB())

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				ContractID: contract.ID,
				IsOconus:   true,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			},
		})
		suite.NotNil(rateArea)

		us_country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		oconusRateArea, err := models.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
		suite.NotNil(oconusRateArea)
		suite.Nil(err)

		jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MBFL")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		factory.FetchOrBuildDefaultContractor(suite.DB(), nil, nil)

		req := httptest.NewRequest("POST", "/orders", nil)
		req = suite.AuthenticateRequest(req, sm)

		hasDependents := true
		spouseHasProGear := true
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		ordersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE
		payload := &internalmessages.CreateUpdateOrders{
			HasDependents:           handlers.FmtBool(hasDependents),
			SpouseHasProGear:        handlers.FmtBool(spouseHasProGear),
			IssueDate:               handlers.FmtDate(issueDate),
			ReportByDate:            handlers.FmtDate(reportByDate),
			OrdersType:              internalmessages.NewOrdersType(ordersType),
			OriginDutyLocationID:    *handlers.FmtUUIDPtr(&originDutyLocation.ID),
			NewDutyLocationID:       handlers.FmtUUID(dutyLocation.ID),
			ServiceMemberID:         handlers.FmtUUID(sm.ID),
			OrdersNumber:            handlers.FmtString("123456"),
			Tac:                     handlers.FmtString("E19A"),
			Sac:                     handlers.FmtString("SacNumber"),
			DepartmentIndicator:     internalmessages.NewDeptIndicator(deptIndicator),
			Grade:                   models.ServiceMemberGradeE1.Pointer(),
			DependentsTwelveAndOver: models.Int64Pointer(-2),
		}

		params := ordersop.CreateOrdersParams{
			HTTPRequest:  req,
			CreateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		createHandler := CreateOrdersHandler{handlerConfig}

		response := createHandler.Handle(params)
		suite.IsType(&handlers.ValidationErrorsResponse{}, response)
		verrsResponse, ok := response.(*handlers.ValidationErrorsResponse)
		suite.True(ok)
		suite.Contains(verrsResponse.Errors, "dependents_twelve_and_over")
	})
}

func (suite *HandlerSuite) TestShowOrder() {
	dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
			LinkOnly: true,
		},
	}, nil)
	rank := factory.FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), models.AffiliationAIRFORCE.String())
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    dutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				RankID: models.UUIDPointer(rank.ID),
			},
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
	handlerConfig := suite.NewHandlerConfig()
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
	suite.Assertions.Equal(order.HasDependents, *okResponse.Payload.HasDependents)
	suite.Assertions.Equal(order.SpouseHasProGear, *okResponse.Payload.SpouseHasProGear)
	suite.Assertions.Equal(rank.RankAbbv, okResponse.Payload.Rank.RankAbbv)
}

func (suite *HandlerSuite) TestPayloadForOrdersModel() {
	dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model:    factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
			LinkOnly: true,
		},
	}, nil)
	rank := factory.FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), models.AffiliationAIRFORCE.String())
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    dutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.Order{
				RankID: &rank.ID,
			},
		},
	}, nil)

	fakeS3 := storageTest.NewFakeS3Storage(true)

	payload, err := payloadForOrdersModel(fakeS3, order)
	suite.NoError(err)
	suite.NotNil(payload)
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
	orderUpdater := orderservice.NewOrderUpdater(move.NewMoveRouter(transportationoffice.NewTransportationOfficesFetcher()))

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

	suite.Run("Returns a 404 if the service member attempting to upload the orders is not the service member associated with the orders", func() {
		orders := factory.BuildOrderWithoutDefaults(suite.DB(), nil, nil)

		otherServiceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		// temporarily set the orders to be associated with a different service member so that the request session
		// has the info for the wrong service member
		orders.ServiceMemberID = otherServiceMember.ID
		orders.ServiceMember = otherServiceMember

		params := setUpRequestAndParams(orders)

		handler := setUpHandler()

		response := handler.Handle(*params)

		suite.IsType(&ordersop.UploadAmendedOrdersNotFound{}, response)
	})

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

func (suite *HandlerSuite) TestUpdateOrdersHandler() {
	waf := entitlements.NewWeightAllotmentFetcher()
	rank := factory.FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), models.AffiliationAIRFORCE.String())
	suite.Run("Can update CONUS orders", func() {
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus: models.BoolPointer(false),
				},
			},
		}, nil)

		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: factory.MakeRandomString(8),
				},
			},
			{
				Model:    address,
				LinkOnly: true,
			},
		}, nil)
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    originDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
			{
				Model: models.Order{
					RankID: models.UUIDPointer(rank.ID),
				},
			},
		}, nil)
		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			}}, nil)

		newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)

		newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		newOrdersNumber := "123456"
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE

		payload := &internalmessages.CreateUpdateOrders{
			OrdersNumber:         handlers.FmtString(newOrdersNumber),
			OrdersType:           &newOrdersType,
			NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
			OriginDutyLocationID: *handlers.FmtUUID(*order.OriginDutyLocationID),
			IssueDate:            handlers.FmtDate(issueDate),
			ReportByDate:         handlers.FmtDate(reportByDate),
			DepartmentIndicator:  &deptIndicator,
			HasDependents:        handlers.FmtBool(false),
			SpouseHasProGear:     handlers.FmtBool(false),
			Grade:                models.ServiceMemberGradeE4.Pointer(),
			MoveID:               *handlers.FmtUUID(move.ID),
			CounselingOfficeID:   handlers.FmtUUID(*newDutyLocation.TransportationOfficeID),
			ServiceMemberID:      handlers.FmtUUID(order.ServiceMemberID),
			Rank:                 strfmt.UUID(rank.ID.String()),
		}

		path := fmt.Sprintf("/orders/%v", order.ID.String())
		req := httptest.NewRequest("PUT", path, nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		params := ordersop.UpdateOrdersParams{
			HTTPRequest:  req,
			OrdersID:     *handlers.FmtUUID(order.ID),
			UpdateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := UpdateOrdersHandler{handlerConfig}

		response := handler.Handle(params)

		suite.IsType(&ordersop.UpdateOrdersOK{}, response)
		okResponse := response.(*ordersop.UpdateOrdersOK)

		suite.NoError(okResponse.Payload.Validate(strfmt.Default))
		suite.Equal(string(newOrdersType), string(*okResponse.Payload.OrdersType))
		suite.Equal(newOrdersNumber, *okResponse.Payload.OrdersNumber)

		updatedOrder, err := models.FetchOrder(suite.DB(), order.ID)
		suite.NoError(err)
		suite.Equal(payload.Grade, updatedOrder.Grade)
		suite.Equal(*okResponse.Payload.AuthorizedWeight, int64(7000)) // E4 authorized weight is 7000, make sure we return that in the response
		expectedUpdatedOrderWeightAllotment, err := waf.GetWeightAllotment(suite.AppContextForTest(), string(*updatedOrder.Grade), updatedOrder.OrdersType)
		suite.NoError(err)
		expectedUpdatedOrderAuthorizedWeight := expectedUpdatedOrderWeightAllotment.TotalWeightSelf
		if *payload.HasDependents {
			expectedUpdatedOrderAuthorizedWeight = expectedUpdatedOrderWeightAllotment.TotalWeightSelfPlusDependents
		}

		expectedOriginalOrderWeightAllotment, err := waf.GetWeightAllotment(suite.AppContextForTest(), string(*order.Grade), updatedOrder.OrdersType)
		suite.NoError(err)
		expectedOriginalOrderAuthorizedWeight := expectedOriginalOrderWeightAllotment.TotalWeightSelf
		if *payload.HasDependents {
			expectedUpdatedOrderAuthorizedWeight = expectedOriginalOrderWeightAllotment.TotalWeightSelfPlusDependents
		}

		suite.Equal(expectedUpdatedOrderAuthorizedWeight, 7000)  // Ensure that when GetWeightAllotment is recalculated that it also returns 7000. This ensures that the database stored the correct information
		suite.Equal(expectedOriginalOrderAuthorizedWeight, 5000) // The order was created as an E1. Ensure that the E1 authorized weight is 5000.
		suite.Equal(string(newOrdersType), string(updatedOrder.OrdersType))
		// Check updated entitlement
		var updatedEntitlement models.Entitlement
		err = suite.DB().Find(&updatedEntitlement, updatedOrder.EntitlementID)
		suite.NoError(err)
		suite.NotEmpty(updatedEntitlement)

		suite.Nil(updatedEntitlement.AccompaniedTour)
		suite.Nil(updatedEntitlement.DependentsTwelveAndOver)
		suite.Nil(updatedEntitlement.DependentsUnderTwelve)
	})

	suite.Run("Can update OCONUS orders", func() {
		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "99801")
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					IsOconus:           models.BoolPointer(true),
					UsPostRegionCityID: &usprc.ID,
				},
			},
		}, nil)

		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name:      factory.MakeRandomString(8),
					AddressID: address.ID,
				},
			},
		}, nil)

		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model:    originDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
			{
				Model: models.Order{
					RankID: models.UUIDPointer(rank.ID),
				},
			},
		}, nil)

		move := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			}}, nil)

		newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				ContractID: contract.ID,
				IsOconus:   true,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			},
		})
		suite.NotNil(rateArea)
		suite.Nil(err)

		us_country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		oconusRateArea, err := models.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
		suite.NotNil(oconusRateArea)
		suite.Nil(err)

		jppsoRegion, err := models.FetchJppsoRegionByCode(suite.DB(), "MAPK")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := models.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, models.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		newOrdersNumber := "123456"
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE

		payload := &internalmessages.CreateUpdateOrders{
			OrdersNumber:         handlers.FmtString(newOrdersNumber),
			OrdersType:           &newOrdersType,
			NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
			OriginDutyLocationID: *handlers.FmtUUID(*order.OriginDutyLocationID),
			IssueDate:            handlers.FmtDate(issueDate),
			ReportByDate:         handlers.FmtDate(reportByDate),
			DepartmentIndicator:  &deptIndicator,
			HasDependents:        handlers.FmtBool(false),
			SpouseHasProGear:     handlers.FmtBool(false),
			Grade:                models.ServiceMemberGradeE4.Pointer(),
			MoveID:               *handlers.FmtUUID(move.ID),
			CounselingOfficeID:   handlers.FmtUUID(*newDutyLocation.TransportationOfficeID),
			ServiceMemberID:      handlers.FmtUUID(order.ServiceMemberID),
			Rank:                 strfmt.UUID(rank.ID.String()),
		}

		payload.AccompaniedTour = models.BoolPointer(true)
		payload.DependentsTwelveAndOver = models.Int64Pointer(5)
		payload.DependentsUnderTwelve = models.Int64Pointer(5)

		path := fmt.Sprintf("/orders/%v", order.ID.String())
		req := httptest.NewRequest("PUT", path, nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		params := ordersop.UpdateOrdersParams{
			HTTPRequest:  req,
			OrdersID:     *handlers.FmtUUID(order.ID),
			UpdateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := UpdateOrdersHandler{handlerConfig}

		response := handler.Handle(params)

		suite.IsType(&ordersop.UpdateOrdersOK{}, response)
		okResponse := response.(*ordersop.UpdateOrdersOK)

		suite.NoError(okResponse.Payload.Validate(strfmt.Default))
		suite.Equal(string(newOrdersType), string(*okResponse.Payload.OrdersType))
		suite.Equal(newOrdersNumber, *okResponse.Payload.OrdersNumber)

		updatedOrder, err := models.FetchOrder(suite.DB(), order.ID)
		suite.NoError(err)
		suite.Equal(payload.Grade, updatedOrder.Grade)
		suite.Equal(*okResponse.Payload.AuthorizedWeight, int64(7000)) // E4 authorized weight is 7000, make sure we return that in the response
		expectedUpdatedOrderWeightAllotment, err := waf.GetWeightAllotment(suite.AppContextForTest(), string(*updatedOrder.Grade), updatedOrder.OrdersType)
		suite.NoError(err)
		expectedUpdatedOrderAuthorizedWeight := expectedUpdatedOrderWeightAllotment.TotalWeightSelf
		if *payload.HasDependents {
			expectedUpdatedOrderAuthorizedWeight = expectedUpdatedOrderWeightAllotment.TotalWeightSelfPlusDependents
		}

		expectedOriginalOrderWeightAllotment, err := waf.GetWeightAllotment(suite.AppContextForTest(), string(*order.Grade), updatedOrder.OrdersType)
		suite.NoError(err)
		expectedOriginalOrderAuthorizedWeight := expectedOriginalOrderWeightAllotment.TotalWeightSelf
		if *payload.HasDependents {
			expectedUpdatedOrderAuthorizedWeight = expectedOriginalOrderWeightAllotment.TotalWeightSelfPlusDependents
		}

		suite.Equal(expectedUpdatedOrderAuthorizedWeight, 7000)  // Ensure that when GetWeightAllotment is recalculated that it also returns 7000. This ensures that the database stored the correct information
		suite.Equal(expectedOriginalOrderAuthorizedWeight, 5000) // The order was created as an E1. Ensure that the E1 authorized weight is 5000.
		suite.Equal(string(newOrdersType), string(updatedOrder.OrdersType))
		// Check updated entitlement
		var updatedEntitlement models.Entitlement
		err = suite.DB().Find(&updatedEntitlement, updatedOrder.EntitlementID)
		suite.NoError(err)
		suite.NotEmpty(updatedEntitlement)

		suite.NotNil(updatedEntitlement.AccompaniedTour)
		suite.NotNil(updatedEntitlement.DependentsTwelveAndOver)
		suite.NotNil(updatedEntitlement.DependentsUnderTwelve)
	})

	suite.Run("Updating order grade to civilian changes PPM type to ACTUAL_EXPENSE", func() {
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					Grade:  models.ServiceMemberGradeE7.Pointer(),
					RankID: models.UUIDPointer(rank.ID),
				},
			},
		}, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			},
			{
				Model: models.PPMShipment{
					PPMType: models.PPMTypeIncentiveBased,
					Status:  models.PPMShipmentStatusDraft,
				},
			},
		}, nil)

		newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		newOrdersNumber := "123456"
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		deptIndicator := internalmessages.DeptIndicatorARMY
		payload := &internalmessages.CreateUpdateOrders{
			OrdersNumber:         handlers.FmtString(newOrdersNumber),
			OrdersType:           &newOrdersType,
			NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
			OriginDutyLocationID: *handlers.FmtUUID(*order.OriginDutyLocationID),
			IssueDate:            handlers.FmtDate(issueDate),
			ReportByDate:         handlers.FmtDate(reportByDate),
			DepartmentIndicator:  &deptIndicator,
			HasDependents:        handlers.FmtBool(false),
			SpouseHasProGear:     handlers.FmtBool(false),
			Grade:                models.ServiceMemberGradeCIVILIANEMPLOYEE.Pointer(),
			MoveID:               *handlers.FmtUUID(ppmShipment.Shipment.MoveTaskOrderID),
			CounselingOfficeID:   handlers.FmtUUID(*newDutyLocation.TransportationOfficeID),
			ServiceMemberID:      handlers.FmtUUID(order.ServiceMemberID),
			Rank:                 strfmt.UUID(rank.ID.String()),
		}

		path := fmt.Sprintf("/orders/%v", order.ID.String())
		req := httptest.NewRequest("PUT", path, nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		params := ordersop.UpdateOrdersParams{
			HTTPRequest:  req,
			OrdersID:     *handlers.FmtUUID(order.ID),
			UpdateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := UpdateOrdersHandler{handlerConfig}

		response := handler.Handle(params)

		suite.IsType(&ordersop.UpdateOrdersOK{}, response)
		okResponse := response.(*ordersop.UpdateOrdersOK)
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		updatedPPM, err := models.FetchPPMShipmentByPPMShipmentID(suite.DB(), ppmShipment.ID)
		suite.NoError(err)
		suite.Equal(updatedPPM.PPMType, models.PPMTypeActualExpense)
		suite.True(*updatedPPM.IsActualExpenseReimbursement)
	})

	suite.Run("Updating order grade FROM civilian to non-civilian changes PPM type to INCENTIVE_BASED", func() {
		order := factory.BuildOrder(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					Grade: models.ServiceMemberGradeCIVILIANEMPLOYEE.Pointer(),
				},
			},
		}, nil)

		ppmShipment := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    order,
				LinkOnly: true,
			},
			{
				Model: models.PPMShipment{
					PPMType: models.PPMTypeActualExpense,
					Status:  models.PPMShipmentStatusDraft,
				},
			},
		}, nil)

		newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)
		newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
		newOrdersNumber := "123456"
		issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
		reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
		deptIndicator := internalmessages.DeptIndicatorARMY
		payload := &internalmessages.CreateUpdateOrders{
			OrdersNumber:         handlers.FmtString(newOrdersNumber),
			OrdersType:           &newOrdersType,
			NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
			OriginDutyLocationID: *handlers.FmtUUID(*order.OriginDutyLocationID),
			IssueDate:            handlers.FmtDate(issueDate),
			ReportByDate:         handlers.FmtDate(reportByDate),
			DepartmentIndicator:  &deptIndicator,
			HasDependents:        handlers.FmtBool(false),
			SpouseHasProGear:     handlers.FmtBool(false),
			Grade:                models.ServiceMemberGradeE4.Pointer(),
			MoveID:               *handlers.FmtUUID(ppmShipment.Shipment.MoveTaskOrderID),
			CounselingOfficeID:   handlers.FmtUUID(*newDutyLocation.TransportationOfficeID),
			ServiceMemberID:      handlers.FmtUUID(order.ServiceMemberID),
			Rank:                 strfmt.UUID(rank.ID.String()),
		}

		path := fmt.Sprintf("/orders/%v", order.ID.String())
		req := httptest.NewRequest("PUT", path, nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		params := ordersop.UpdateOrdersParams{
			HTTPRequest:  req,
			OrdersID:     *handlers.FmtUUID(order.ID),
			UpdateOrders: payload,
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)

		handler := UpdateOrdersHandler{handlerConfig}

		response := handler.Handle(params)

		suite.IsType(&ordersop.UpdateOrdersOK{}, response)
		okResponse := response.(*ordersop.UpdateOrdersOK)
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		updatedPPM, err := models.FetchPPMShipmentByPPMShipmentID(suite.DB(), ppmShipment.ID)
		suite.NoError(err)
		suite.Equal(updatedPPM.PPMType, models.PPMTypeIncentiveBased)
		suite.False(*updatedPPM.IsActualExpenseReimbursement)
	})

}

func (suite *HandlerSuite) TestUpdateOrdersHandlerOriginPostalCodeAndGBLOC() {

	firstAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				PostalCode: "90210",
			},
		},
	}, nil)
	updatedAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
		{
			Model: models.Address{
				PostalCode: "35023",
			},
		},
	}, nil)
	dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				AddressID: firstAddress.ID,
			},
		},
	}, nil)
	updatedDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				AddressID: updatedAddress.ID,
			},
		},
	}, nil)
	newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)

	rank := factory.FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), models.AffiliationAIRFORCE.String())
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model: models.Order{
				OriginDutyLocationID: &dutyLocation.ID,
				NewDutyLocationID:    newDutyLocation.ID,
			},
		},
	}, nil)

	factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    order,
			LinkOnly: true,
		}}, nil)

	fetchedOrder, err := models.FetchOrder(suite.DB(), order.ID)
	suite.NoError(err)

	var fetchedPostalCode, fetchedGBLOC string
	fetchedPostalCode, err = fetchedOrder.GetOriginPostalCode(suite.DB())
	suite.NoError(err)
	fetchedGBLOC, err = fetchedOrder.GetOriginGBLOC(suite.DB())
	suite.NoError(err)

	suite.Equal("90210", fetchedPostalCode)
	suite.Equal("KKFA", fetchedGBLOC)

	newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE

	payload := &internalmessages.CreateUpdateOrders{
		OrdersType:           &order.OrdersType,
		NewDutyLocationID:    handlers.FmtUUID(order.NewDutyLocationID),
		OriginDutyLocationID: *handlers.FmtUUID(updatedDutyLocation.ID),
		IssueDate:            handlers.FmtDate(issueDate),
		ReportByDate:         handlers.FmtDate(reportByDate),
		DepartmentIndicator:  &deptIndicator,
		HasDependents:        handlers.FmtBool(false),
		SpouseHasProGear:     handlers.FmtBool(false),
		ServiceMemberID:      handlers.FmtUUID(order.ServiceMemberID),
		Grade:                models.ServiceMemberGradeE4.Pointer(),
		Rank:                 strfmt.UUID(rank.ID.String()),
	}

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("PUT", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.UpdateOrdersParams{
		HTTPRequest:  req,
		OrdersID:     *handlers.FmtUUID(order.ID),
		UpdateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.NewHandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)

	handler := UpdateOrdersHandler{handlerConfig}

	response := handler.Handle(params)

	suite.IsType(&ordersop.UpdateOrdersOK{}, response)

	okResponse := response.(*ordersop.UpdateOrdersOK)

	suite.NoError(okResponse.Payload.Validate(strfmt.Default))
	suite.Equal(string(newOrdersType), string(*okResponse.Payload.OrdersType))

	fetchedOrder, err = models.FetchOrder(suite.DB(), order.ID)
	suite.NoError(err)

	fetchedPostalCode, err = fetchedOrder.GetOriginPostalCode(suite.DB())
	suite.NoError(err)
	fetchedGBLOC = *fetchedOrder.OriginDutyLocationGBLOC
	suite.NoError(err)

	suite.Equal("35023", fetchedPostalCode)
	suite.Equal("CNNQ", fetchedGBLOC)
}

func (suite *HandlerSuite) TestEntitlementHelperFunc() {
	orderGrade := internalmessages.OrderPayGrade("O-3")
	int64Dependents := int64(2)
	intDependents := int(int64Dependents)
	suite.Run("Can fully cover the hasEntitlementChangedFunc", func() {
		testCases := []struct {
			order                          models.Order
			payloadPayGrade                *internalmessages.OrderPayGrade
			payloadDependentsUnderTwelve   *int64
			payloadDependentsTwelveAndOver *int64
			payloadAccompaniedTour         *bool
			shouldReturnFalse              *bool
			payloadOrdersType              *internalmessages.OrdersType
		}{
			{
				order: models.Order{
					Grade: &orderGrade,
				},
			},
			{
				order: models.Order{
					Entitlement: &models.Entitlement{
						DependentsUnderTwelve: &intDependents,
					},
				},
			},
			{
				order: models.Order{
					Entitlement: &models.Entitlement{
						DependentsTwelveAndOver: &intDependents,
					},
				},
			},
			{
				order: models.Order{
					Entitlement: &models.Entitlement{
						AccompaniedTour: models.BoolPointer(true),
					},
				},
			},
			{
				order:             models.Order{},
				shouldReturnFalse: models.BoolPointer(true),
			},
		}
		for _, tc := range testCases {
			if tc.shouldReturnFalse != nil && *tc.shouldReturnFalse {
				// Test should return false
				suite.False(hasEntitlementChanged(tc.order, tc.payloadOrdersType, tc.payloadPayGrade, tc.payloadDependentsUnderTwelve, tc.payloadDependentsTwelveAndOver, tc.payloadAccompaniedTour))
			} else {
				// Test defaults to returning true
				suite.True(hasEntitlementChanged(tc.order, tc.payloadOrdersType, tc.payloadPayGrade, tc.payloadDependentsUnderTwelve, tc.payloadDependentsTwelveAndOver, tc.payloadAccompaniedTour))
			}

		}
	})
}

func (suite *HandlerSuite) TestUpdateOrdersHandlerWithCounselingOffice() {
	originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
		{
			Model: models.DutyLocation{
				ProvidesServicesCounseling: true,
			},
		},
	}, nil)
	rank := factory.FetchOrBuildRankByPayGradeAndAffiliation(suite.DB(), string(models.ServiceMemberGradeE4), models.AffiliationAIRFORCE.String())
	order := factory.BuildOrder(suite.DB(), []factory.Customization{
		{
			Model:    originDutyLocation,
			LinkOnly: true,
			Type:     &factory.DutyLocations.OriginDutyLocation,
		},
	}, nil)

	newDutyLocation := factory.BuildDutyLocation(suite.DB(), nil, nil)

	newOrdersType := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	newOrdersNumber := "123456"
	issueDate := time.Date(2018, time.March, 10, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(2018, time.August, 1, 0, 0, 0, 0, time.UTC)
	deptIndicator := internalmessages.DeptIndicatorAIRANDSPACEFORCE
	move := factory.BuildMove(suite.DB(), []factory.Customization{
		{
			Model:    order,
			LinkOnly: true,
		}}, nil)
	payload := &internalmessages.CreateUpdateOrders{
		OrdersNumber:         handlers.FmtString(newOrdersNumber),
		OrdersType:           &newOrdersType,
		NewDutyLocationID:    handlers.FmtUUID(newDutyLocation.ID),
		OriginDutyLocationID: *handlers.FmtUUID(*order.OriginDutyLocationID),
		IssueDate:            handlers.FmtDate(issueDate),
		ReportByDate:         handlers.FmtDate(reportByDate),
		DepartmentIndicator:  &deptIndicator,
		HasDependents:        handlers.FmtBool(false),
		SpouseHasProGear:     handlers.FmtBool(false),
		Grade:                models.ServiceMemberGradeE4.Pointer(),
		MoveID:               *handlers.FmtUUID(move.ID),
		CounselingOfficeID:   handlers.FmtUUID(*newDutyLocation.TransportationOfficeID),
		ServiceMemberID:      handlers.FmtUUID(order.ServiceMemberID),
		Rank:                 strfmt.UUID(rank.ID.String()),
	}

	path := fmt.Sprintf("/orders/%v", order.ID.String())
	req := httptest.NewRequest("PUT", path, nil)
	req = suite.AuthenticateRequest(req, order.ServiceMember)

	params := ordersop.UpdateOrdersParams{
		HTTPRequest:  req,
		OrdersID:     *handlers.FmtUUID(order.ID),
		UpdateOrders: payload,
	}

	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig := suite.NewHandlerConfig()
	handlerConfig.SetFileStorer(fakeS3)

	handler := UpdateOrdersHandler{handlerConfig}

	response := handler.Handle(params)

	suite.IsType(&ordersop.UpdateOrdersOK{}, response)
	okResponse := response.(*ordersop.UpdateOrdersOK)

	suite.NoError(okResponse.Payload.Validate(strfmt.Default))
	suite.Equal(string(newOrdersType), string(*okResponse.Payload.OrdersType))
	suite.Equal(newOrdersNumber, *okResponse.Payload.OrdersNumber)

}

func (suite *HandlerSuite) TestGetRanksHandler() {
	suite.Run("happy path", func() {
		order := factory.BuildOrder(suite.DB(), nil, nil)
		affiliation := internalmessages.AffiliationAIRFORCE
		grade := models.ServiceMemberGradeE2
		path := fmt.Sprintf("/ranks/%v&%v", affiliation, grade)
		req := httptest.NewRequest("GET", path, nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		params := ordersop.GetRanksParams{
			HTTPRequest: req,
			Affiliation: string(affiliation),
			Grade:       string(grade),
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		showHandler := GetRanksHandler{handlerConfig}

		response := showHandler.Handle(params)

		suite.Assertions.IsType(&ordersop.GetRanksOK{}, response)
		okResponse := response.(*ordersop.GetRanksOK)

		suite.Assertions.Equal(1, len(okResponse.Payload))
	})

	suite.Run("test a bad affiliation", func() {
		order := factory.BuildOrder(suite.DB(), nil, nil)
		grade := models.ServiceMemberGradeE2
		path := fmt.Sprintf("/ranks/%v&%v", "FAKE", grade)
		req := httptest.NewRequest("GET", path, nil)
		req = suite.AuthenticateRequest(req, order.ServiceMember)

		params := ordersop.GetRanksParams{
			HTTPRequest: req,
			Affiliation: "FAKE",
		}

		fakeS3 := storageTest.NewFakeS3Storage(true)
		handlerConfig := suite.NewHandlerConfig()
		handlerConfig.SetFileStorer(fakeS3)
		showHandler := GetRanksHandler{handlerConfig}

		response := showHandler.Handle(params)

		suite.Assertions.IsType(&ordersop.GetRanksNotFound{}, response)
	})
}
