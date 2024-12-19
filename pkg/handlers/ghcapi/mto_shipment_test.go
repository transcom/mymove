package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	shipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mobilehomeshipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	sitextension "github.com/transcom/mymove/pkg/services/sit_extension"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	"github.com/transcom/mymove/pkg/swagger/nullable"
	"github.com/transcom/mymove/pkg/trace"
	"github.com/transcom/mymove/pkg/unit"
)

type listMTOShipmentsSubtestData struct {
	mtoAgent       models.MTOAgent
	mtoServiceItem models.MTOServiceItem
	shipments      models.MTOShipments
	params         mtoshipmentops.ListMTOShipmentsParams
	sitExtension   models.SITDurationUpdate
	sit            models.MTOServiceItem
}

func (suite *HandlerSuite) makeListMTOShipmentsSubtestData() (subtestData *listMTOShipmentsSubtestData) {
	subtestData = &listMTOShipmentsSubtestData{}

	mto := factory.BuildMove(suite.DB(), nil, nil)

	storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

	sitAllowance := int(90)
	mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:           models.MTOShipmentStatusApproved,
				ShipmentType:     models.MTOShipmentTypeHHGIntoNTSDom,
				CounselorRemarks: handlers.FmtString("counselor remark"),
				SITDaysAllowance: &sitAllowance,
			},
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
	}, nil)

	secondShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	// third shipment with destination address and type
	destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationType := models.DestinationTypeHomeOfRecord
	thirdShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				Status:               models.MTOShipmentStatusSubmitted,
				DestinationAddressID: &destinationAddress.ID,
				DestinationType:      &destinationType,
			},
		},
	}, nil)

	subtestData.mtoAgent = factory.BuildMTOAgent(suite.DB(), []factory.Customization{
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
	}, nil)
	subtestData.mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				MTOShipmentID: &mtoShipment.ID,
			},
		},
	}, nil)

	ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// testdatagen.MakeDOFSITReService(suite.DB(), testdatagen.Assertions{})

	year, month, day := time.Now().Date()
	lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
	lastMonthDeparture := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
	factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				SITEntryDate:     &lastMonthEntry,
				SITDepartureDate: &lastMonthDeparture,
				Status:           models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
	}, nil)

	aWeekAgo := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
	departureDate := aWeekAgo.Add(time.Hour * 24 * 30)
	subtestData.sit = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				SITEntryDate:     &aWeekAgo,
				SITDepartureDate: &departureDate,
				Status:           models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
	}, nil)

	subtestData.sitExtension = factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
		{
			Model:    mtoShipment,
			LinkOnly: true,
		},
	}, []factory.Trait{factory.GetTraitApprovedSITDurationUpdate})
	subtestData.shipments = models.MTOShipments{mtoShipment, secondShipment, thirdShipment, ppm.Shipment}
	requestUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_shipments", mto.ID.String()), nil)
	req = suite.AuthenticateOfficeRequest(req, requestUser)

	subtestData.params = mtoshipmentops.ListMTOShipmentsParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
	}

	return subtestData
}

func (suite *HandlerSuite) TestListMTOShipmentsHandler() {
	suite.Run("Successful list fetch - Integration Test", func() {
		subtestData := suite.makeListMTOShipmentsSubtestData()
		params := subtestData.params
		shipments := subtestData.shipments
		mtoAgent := subtestData.mtoAgent
		mtoServiceItem := subtestData.mtoServiceItem
		sitExtension := subtestData.sitExtension

		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mtoshipment.NewMTOShipmentFetcher(),
			sitstatus.NewShipmentSITStatus(),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsOK{}, response)
		okResponse := response.(*mtoshipmentops.ListMTOShipmentsOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Len(okResponse.Payload, 4)

		payloadShipment := okResponse.Payload[0]
		suite.Equal(shipments[0].ID.String(), payloadShipment.ID.String())
		suite.Equal(*shipments[0].CounselorRemarks, *payloadShipment.CounselorRemarks)
		suite.Equal(mtoAgent.ID.String(), payloadShipment.MtoAgents[0].ID.String())
		suite.Equal(mtoServiceItem.ID.String(), payloadShipment.MtoServiceItems[0].ID.String())
		suite.Equal(sitExtension.ID.String(), payloadShipment.SitExtensions[0].ID.String())
		suite.Equal(shipments[0].StorageFacility.ID.String(), payloadShipment.StorageFacility.ID.String())
		suite.Equal(shipments[0].StorageFacility.Address.ID.String(), payloadShipment.StorageFacility.Address.ID.String())

		payloadShipment2 := okResponse.Payload[1]
		suite.Equal(shipments[1].ID.String(), payloadShipment2.ID.String())
		suite.Nil(payloadShipment2.StorageFacility)

		suite.Equal(int64(190), *payloadShipment.SitDaysAllowance)
		suite.Equal(sitstatus.OriginSITLocation, payloadShipment.SitStatus.CurrentSIT.Location)
		suite.Equal(int64(8), *payloadShipment.SitStatus.CurrentSIT.DaysInSIT)
		suite.Equal(int64(174), *payloadShipment.SitStatus.TotalDaysRemaining)
		suite.Equal(int64(16), *payloadShipment.SitStatus.TotalSITDaysUsed) // 7 from the previous SIT and 7 from the current (+2 for including last days)
		suite.Equal(int64(16), *payloadShipment.SitStatus.CalculatedTotalDaysInSIT)
		suite.Equal(subtestData.sit.SITEntryDate.Format("2006-01-02"), payloadShipment.SitStatus.CurrentSIT.SitEntryDate.String())
		suite.Equal(subtestData.sit.SITDepartureDate.Format("2006-01-02"), payloadShipment.SitStatus.CurrentSIT.SitDepartureDate.String())

		suite.Len(payloadShipment.SitStatus.PastSITServiceItemGroupings, 1)
		year, month, day := time.Now().Date()
		lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
		suite.Equal(lastMonthEntry.Format(strfmt.MarshalFormat), payloadShipment.SitStatus.PastSITServiceItemGroupings[0].Summary.SitEntryDate.String())

		// This one has a destination shipment type
		payloadShipment3 := okResponse.Payload[2]
		suite.Equal(string(models.DestinationTypeHomeOfRecord), string(*payloadShipment3.DestinationType))

		payloadShipment4 := okResponse.Payload[3]
		suite.NotNil(payloadShipment4.PpmShipment)
		suite.Equal(shipments[3].ID.String(), payloadShipment4.PpmShipment.ShipmentID.String())
		suite.Equal(shipments[3].PPMShipment.ID.String(), payloadShipment4.PpmShipment.ID.String())
	})

	suite.Run("Failure list fetch - Internal Server Error", func() {
		subtestData := suite.makeListMTOShipmentsSubtestData()
		params := subtestData.params
		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}

		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mockMTOShipmentFetcher,
			sitstatus.NewShipmentSITStatus(),
		}

		mockMTOShipmentFetcher.On("ListMTOShipments", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID")).Return(nil, apperror.NewQueryError("MTOShipment", errors.New("query error"), ""))

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsInternalServerError{}, response)
		payload := response.(*mtoshipmentops.ListMTOShipmentsInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Failure list fetch - 404 Not Found - Move Task Order ID", func() {
		subtestData := suite.makeListMTOShipmentsSubtestData()
		params := subtestData.params

		mockMTOShipmentFetcher := &mocks.MTOShipmentFetcher{}

		handler := ListMTOShipmentsHandler{
			suite.HandlerConfig(),
			mockMTOShipmentFetcher,
			sitstatus.NewShipmentSITStatus(),
		}

		mockMTOShipmentFetcher.On("ListMTOShipments", mock.AnythingOfType("*appcontext.appContext"), mock.AnythingOfType("uuid.UUID")).Return(nil, apperror.NewNotFoundError(uuid.FromStringOrNil(params.MoveTaskOrderID.String()), "move not found"))

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsNotFound{}, response)
		payload := response.(*mtoshipmentops.ListMTOShipmentsNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestDeleteShipmentHandler() {
	suite.Run("Returns a 403 when user is not a service counselor or TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeQae})
		uuid := uuid.Must(uuid.NewV4())
		deleter := &mocks.ShipmentDeleter{}

		deleter.AssertNumberOfCalls(suite.T(), "DeleteShipment", 0)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentForbidden{}, response)
		payload := response.(*shipmentops.DeleteShipmentForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 204 when all validations pass", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(shipment.MoveTaskOrderID, nil)

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentNoContent{}, response)

		// Validate outgoing payload: no payload
	})

	suite.Run("Returns 404 when deleter returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.NotFoundError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentNotFound{}, response)
		payload := response.(*shipmentops.DeleteShipmentNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 403 when deleter returns ForbiddenError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ForbiddenError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentForbidden{}, response)
		payload := response.(*shipmentops.DeleteShipmentForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 422 - Unprocessable Enitity error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.UnprocessableEntityError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentUnprocessableEntity{}, response)
		payload := response.(*shipmentops.DeleteShipmentUnprocessableEntity).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 - Conflict error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		deleter := &mocks.ShipmentDeleter{}

		deleter.On("DeleteShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID).Return(uuid.Nil, apperror.ConflictError{})

		req := httptest.NewRequest("DELETE", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DeleteShipmentHandler{
			handlerConfig,
			deleter,
		}
		deletionParams := shipmentops.DeleteShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(deletionParams)
		suite.IsType(&shipmentops.DeleteShipmentConflict{}, response)
		payload := response.(*shipmentops.DeleteShipmentConflict).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestGetShipmentHandler() {
	// Success integration test
	suite.Run("Successful fetch (integration) test", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		fetcher := mtoshipment.NewMTOShipmentFetcher()
		request := httptest.NewRequest("GET", fmt.Sprintf("/shipments/%s", shipment.ID.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := mtoshipmentops.GetShipmentParams{
			HTTPRequest: request,
			ShipmentID:  strfmt.UUID(shipment.ID.String()),
		}

		handler := GetMTOShipmentHandler{
			HandlerConfig:      handlerConfig,
			mtoShipmentFetcher: fetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.GetShipmentOK{}, response)
		payload := response.(*mtoshipmentops.GetShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// 404 response
	suite.Run("404 response when the service returns not found", func() {
		uuidForShipment, _ := uuid.NewV4()
		officeUser := factory.BuildOfficeUser(nil, nil, nil)
		handlerConfig := suite.HandlerConfig()
		fetcher := mtoshipment.NewMTOShipmentFetcher()
		request := httptest.NewRequest("GET", fmt.Sprintf("/shipments/%s", uuidForShipment.String()), nil)
		request = suite.AuthenticateOfficeRequest(request, officeUser)

		params := mtoshipmentops.GetShipmentParams{
			HTTPRequest: request,
			ShipmentID:  strfmt.UUID(uuidForShipment.String()),
		}

		handler := GetMTOShipmentHandler{
			HandlerConfig:      handlerConfig,
			mtoShipmentFetcher: fetcher,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.GetShipmentNotFound{}, response)
		payload := response.(*mtoshipmentops.GetShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestApproveShipmentHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		builder := query.NewQueryBuilder()
		moveRouter := moveservices.NewMoveRouter()
		planner := &routemocks.Planner{}
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		approver := mtoshipment.NewShipmentApprover(
			mtoshipment.NewShipmentRouter(),
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			&routemocks.Planner{},
			moveWeights,
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentOK{}, response)
		payload := response.(*shipmentops.ApproveShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		uuid := uuid.Must(uuid.NewV4())
		approver := &mocks.ShipmentApprover{}

		approver.AssertNumberOfCalls(suite.T(), "ApproveShipment", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentForbidden{}, response)
		payload := response.(*shipmentops.ApproveShipmentForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 404 when approver returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentNotFound{}, response)
		payload := response.(*shipmentops.ApproveShipmentNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 when approver returns Conflict Error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.ConflictError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentConflict{}, response)
		payload := response.(*shipmentops.ApproveShipmentConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(time.Now())
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentPreconditionFailed{}, response)
		payload := response.(*shipmentops.ApproveShipmentPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when approver returns validation errors", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentUnprocessableEntity{}, response)
		payload := response.(*shipmentops.ApproveShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when approver returns unexpected error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentApprover{}

		approver.On("ApproveShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentInternalServerError{}, response)
		payload := response.(*shipmentops.ApproveShipmentInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestRequestShipmentDiversionHandler() {
	diversionReason := "Test Reason"

	suite.Run("Returns 200 when all validations pass", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		requester := mtoshipment.NewShipmentDiversionRequester(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}

		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionOK{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		uuid := uuid.Must(uuid.NewV4())
		requester := &mocks.ShipmentDiversionRequester{}

		requester.AssertNumberOfCalls(suite.T(), "RequestShipmentDiversion", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}
		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionForbidden{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 404 when requester returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &diversionReason).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}
		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionNotFound{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 when requester returns Conflict Error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &diversionReason).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}
		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionConflict{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(time.Now())
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &diversionReason).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}
		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionPreconditionFailed{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when requester returns validation errors", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &diversionReason).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}
		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionUnprocessableEntity{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when requester returns unexpected error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		requester := &mocks.ShipmentDiversionRequester{}

		requester.On("RequestShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &diversionReason).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentDiversionHandler{
			handlerConfig,
			requester,
			sitstatus.NewShipmentSITStatus(),
		}
		requestParams := shipmentops.RequestShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RequestDiversion{
				DiversionReason: &diversionReason,
			},
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(requestParams)
		suite.IsType(&shipmentops.RequestShipmentDiversionInternalServerError{}, response)
		payload := response.(*shipmentops.RequestShipmentDiversionInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestApproveShipmentDiversionHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:    models.MTOShipmentStatusSubmitted,
					Diversion: true,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := mtoshipment.NewShipmentDiversionApprover(
			mtoshipment.NewShipmentRouter(),
			moveservices.NewMoveRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionOK{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		uuid := uuid.Must(uuid.NewV4())
		approver := &mocks.ShipmentDiversionApprover{}

		approver.AssertNumberOfCalls(suite.T(), "ApproveShipmentDiversion", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionForbidden{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 404 when approver returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionNotFound{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 when approver returns Conflict Error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionConflict{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(time.Now())
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionPreconditionFailed{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when approver returns validation errors", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionUnprocessableEntity{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when approver returns unexpected error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approver := &mocks.ShipmentDiversionApprover{}

		approver.On("ApproveShipmentDiversion", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/approve-diversion", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := ApproveShipmentDiversionHandler{
			handlerConfig,
			approver,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.ApproveShipmentDiversionParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveShipmentDiversionInternalServerError{}, response)
		payload := response.(*shipmentops.ApproveShipmentDiversionInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestRejectShipmentHandler() {
	reason := "reason"

	suite.Run("Returns 200 when all validations pass", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := mtoshipment.NewShipmentRejecter(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}

		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentOK{}, response)
		payload := response.(*shipmentops.RejectShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		uuid := uuid.Must(uuid.NewV4())
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.AssertNumberOfCalls(suite.T(), "RejectShipment", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentForbidden{}, response)
		payload := response.(*shipmentops.RejectShipmentForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 404 when rejecter returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentNotFound{}, response)
		payload := response.(*shipmentops.RejectShipmentNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 when rejecter returns Conflict Error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentConflict{}, response)
		payload := response.(*shipmentops.RejectShipmentConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(time.Now())
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentPreconditionFailed{}, response)
		payload := response.(*shipmentops.RejectShipmentPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when rejecter returns validation errors", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentUnprocessableEntity{}, response)
		payload := response.(*shipmentops.RejectShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when rejecter returns unexpected error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := &mocks.ShipmentRejecter{}

		rejecter.On("RejectShipment", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag, &reason).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body: &ghcmessages.RejectShipment{
				RejectionReason: &reason,
			},
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentInternalServerError{}, response)
		payload := response.(*shipmentops.RejectShipmentInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Requires rejection reason in Body of request", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		rejecter := mtoshipment.NewShipmentRejecter(
			mtoshipment.NewShipmentRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/reject", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RejectShipmentHandler{
			handlerConfig,
			rejecter,
		}
		params := shipmentops.RejectShipmentParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
			Body:        &ghcmessages.RejectShipment{},
		}

		// Validate incoming payload
		suite.Error(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RejectShipmentUnprocessableEntity{}, response)
		payload := response.(*shipmentops.RejectShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}

func (suite *HandlerSuite) TestRequestShipmentCancellationHandler() {
	suite.Run("Returns 200 when all validations pass", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		// valid pickupdate is anytime after the request to cancel date
		actualPickupDate := time.Now().AddDate(0, 0, 1)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					ActualPickupDate: &actualPickupDate,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		canceler := mtoshipment.NewShipmentCancellationRequester(
			mtoshipment.NewShipmentRouter(),
			moveservices.NewMoveRouter(),
		)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}

		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationOK{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		uuid := uuid.Must(uuid.NewV4())
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.AssertNumberOfCalls(suite.T(), "RequestShipmentCancellation", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
			IfMatch:     etag.GenerateEtag(time.Now()),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationForbidden{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 404 when canceler returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.NotFoundError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationNotFound{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 when canceler returns Conflict Error", func() {
		actualPickupDate := time.Now()
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID:               uuid.Must(uuid.NewV4()),
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, mtoshipment.ConflictStatusError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationConflict{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 412 when eTag does not match", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(time.Now())
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.PreconditionFailedError{})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationPreconditionFailed{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when canceler returns validation errors", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationUnprocessableEntity{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when canceler returns unexpected error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		canceler := &mocks.ShipmentCancellationRequester{}

		canceler.On("RequestShipmentCancellation", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, eTag).Return(nil, errors.New("UnexpectedError"))

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-cancellation", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := RequestShipmentCancellationHandler{
			handlerConfig,
			canceler,
			sitstatus.NewShipmentSITStatus(),
		}
		approveParams := shipmentops.RequestShipmentCancellationParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
			IfMatch:     eTag,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentCancellationInternalServerError{}, response)
		payload := response.(*shipmentops.RequestShipmentCancellationInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestRequestShipmentReweighHandler() {
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()

	suite.Run("Returns 200 when all validations pass", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		reweighRequester := mtoshipment.NewShipmentReweighRequester()

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		handlerConfig := suite.HandlerConfig()
		handlerConfig.SetNotificationSender(suite.TestNotificationSender())
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moveservices.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			sitstatus.NewShipmentSITStatus(),
			updater,
		}

		approveParams := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentReweighOK{}, response)
		okResponse := response.(*shipmentops.RequestShipmentReweighOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(strfmt.UUID(shipment.ID.String()), payload.ShipmentID)
		suite.EqualValues(models.ReweighRequesterTOO, payload.RequestedBy)
		suite.WithinDuration(time.Now(), (time.Time)(payload.RequestedAt), 2*time.Second)
		suite.HasWebhookNotification(shipment.ID, traceID)
	})

	suite.Run("Returns a 403 when the office user is not a TOO", func() {
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		uuid := uuid.Must(uuid.NewV4())
		reweighRequester := &mocks.ShipmentReweighRequester{}

		reweighRequester.AssertNumberOfCalls(suite.T(), "RequestShipmentReweigh", 0)

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", uuid.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moveservices.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			sitstatus.NewShipmentSITStatus(),
			updater,
		}
		approveParams := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(uuid),
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.RequestShipmentReweighForbidden{}, response)
		payload := response.(*shipmentops.RequestShipmentReweighForbidden).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 404 when reweighRequester returns NotFoundError", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moveservices.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			sitstatus.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}
		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, apperror.NotFoundError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighNotFound{}, response)
		payload := response.(*shipmentops.RequestShipmentReweighNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Returns 409 when reweighRequester returns Conflict Error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moveservices.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			sitstatus.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, apperror.ConflictError{})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighConflict{}, response)
		payload := response.(*shipmentops.RequestShipmentReweighConflict).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 422 when reweighRequester returns validation errors", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moveservices.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			sitstatus.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}
		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, apperror.InvalidInputError{ValidationErrors: &validate.Errors{}})

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighUnprocessableEntity{}, response)
		payload := response.(*shipmentops.RequestShipmentReweighUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("Returns 500 when reweighRequester returns unexpected error", func() {
		shipment := factory.BuildMTOShipmentMinimal(nil, []factory.Customization{
			{
				Model: models.MTOShipment{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		reweighRequester := &mocks.ShipmentReweighRequester{}

		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/request-reweigh", shipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveRouter := moveservices.NewMoveRouter()
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		updater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)

		handler := RequestShipmentReweighHandler{
			handlerConfig,
			reweighRequester,
			sitstatus.NewShipmentSITStatus(),
			updater,
		}
		params := shipmentops.RequestShipmentReweighParams{
			HTTPRequest: req,
			ShipmentID:  *handlers.FmtUUID(shipment.ID),
		}

		reweighRequester.On("RequestShipmentReweigh", mock.AnythingOfType("*appcontext.appContext"), shipment.ID, models.ReweighRequesterTOO).Return(nil, errors.New("UnexpectedError"))

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&shipmentops.RequestShipmentReweighInternalServerError{}, response)
		payload := response.(*shipmentops.RequestShipmentReweighInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestReviewShipmentAddressUpdateHandler() {
	officeRemarks := "This is a TOO remark"
	status := "APPROVED"

	suite.Run("PATCH Success - 200 OK", func() {

		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)

		newAddress := models.ShipmentAddressUpdate{
			OfficeRemarks: &officeRemarks,
			Status:        models.ShipmentAddressUpdateStatusApproved,
			ID:            addressChange.ID,
		}

		body := shipmentops.ReviewShipmentAddressUpdateBody{
			OfficeRemarks: &officeRemarks,
			Status:        &status,
		}

		req := httptest.NewRequest("PATCH", "/shipments/{mtoShipmentID}/review-shipment-address-update", nil)

		params := shipmentops.ReviewShipmentAddressUpdateParams{
			HTTPRequest: req,
			Body:        body,
			IfMatch:     etag.GenerateEtag(addressChange.Shipment.UpdatedAt),
			ShipmentID:  *handlers.FmtUUID(addressChange.ShipmentID),
		}

		handlerConfig := suite.HandlerConfig()

		mockCreator := mocks.ShipmentAddressUpdateRequester{}

		mockCreator.On("ReviewShipmentAddressChange",
			mock.AnythingOfType("*appcontext.appContext"),
			addressChange.ShipmentID,
			models.ShipmentAddressUpdateStatusApproved,
			"This is a TOO remark",
		).Return(&newAddress, nil)

		handler := ReviewShipmentAddressUpdateHandler{
			handlerConfig,
			&mockCreator,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		okResponse := response.(*shipmentops.ReviewShipmentAddressUpdateOK)
		payload := okResponse.Payload

		suite.IsNotErrResponse(response)
		suite.NotNil(payload)
		suite.Equal(ghcmessages.ShipmentAddressUpdateStatus("APPROVED"), payload.Status)
		suite.Equal("This is a TOO remark", *payload.OfficeRemarks)
	})

	suite.Run("PATCH Failure - 422 Unprocessable Entity error", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)

		mockCreator := mocks.ShipmentAddressUpdateRequester{}

		handlerConfig := suite.HandlerConfig()

		handler := ReviewShipmentAddressUpdateHandler{
			handlerConfig,
			&mockCreator,
		}
		body := shipmentops.ReviewShipmentAddressUpdateBody{
			OfficeRemarks: &officeRemarks,
			Status:        &status,
		}

		req := httptest.NewRequest("PATCH", "/shipments/{mtoShipmentID}/review-shipment-address-update", nil)

		params := shipmentops.ReviewShipmentAddressUpdateParams{
			HTTPRequest: req,
			Body:        body,
			IfMatch:     etag.GenerateEtag(addressChange.Shipment.UpdatedAt),
			ShipmentID:  *handlers.FmtUUID(addressChange.ShipmentID),
		}

		verrs := validate.NewErrors()
		verrs.Add("some key", "some value")
		err := apperror.NewInvalidInputError(uuid.Nil, nil, verrs, "unable to create ShipmentAddressUpdate")

		mockCreator.On("ReviewShipmentAddressChange",
			mock.AnythingOfType("*appcontext.appContext"),
			addressChange.ShipmentID,
			models.ShipmentAddressUpdateStatusApproved,
			"This is a TOO remark",
		).Return(nil, err)

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.ReviewShipmentAddressUpdateUnprocessableEntity{}, response)

		errResponse := response.(*shipmentops.ReviewShipmentAddressUpdateUnprocessableEntity)
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH Failure - 409 Request conflict response error", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)

		mockCreator := mocks.ShipmentAddressUpdateRequester{}

		handlerConfig := suite.HandlerConfig()

		handler := ReviewShipmentAddressUpdateHandler{
			handlerConfig,
			&mockCreator,
		}

		body := shipmentops.ReviewShipmentAddressUpdateBody{
			OfficeRemarks: &officeRemarks,
			Status:        &status,
		}

		req := httptest.NewRequest("PATCH", "/shipments/{mtoShipmentID}/review-shipment-address-update", nil)

		params := shipmentops.ReviewShipmentAddressUpdateParams{
			HTTPRequest: req,
			Body:        body,
			IfMatch:     etag.GenerateEtag(addressChange.Shipment.UpdatedAt),
			ShipmentID:  *handlers.FmtUUID(addressChange.ShipmentID),
		}

		err := apperror.NewConflictError(uuid.Nil, "unable to create ReviewShipmentAddressChange")

		mockCreator.On("ReviewShipmentAddressChange",
			mock.AnythingOfType("*appcontext.appContext"),
			addressChange.ShipmentID,
			models.ShipmentAddressUpdateStatusApproved,
			"This is a TOO remark",
		).Return(nil, err)

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.NotNil(response)
		suite.IsType(&shipmentops.ReviewShipmentAddressUpdateConflict{}, response)
		errResponse := response.(*shipmentops.ReviewShipmentAddressUpdateConflict)

		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH Failure - 404 Not Found response error", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)

		mockCreator := mocks.ShipmentAddressUpdateRequester{}

		handlerConfig := suite.HandlerConfig()

		handler := ReviewShipmentAddressUpdateHandler{
			handlerConfig,
			&mockCreator,
		}

		body := shipmentops.ReviewShipmentAddressUpdateBody{
			OfficeRemarks: &officeRemarks,
			Status:        &status,
		}

		req := httptest.NewRequest("PATCH", "/shipments/{mtoShipmentID}/review-shipment-address-update", nil)

		params := shipmentops.ReviewShipmentAddressUpdateParams{
			HTTPRequest: req,
			Body:        body,
			IfMatch:     etag.GenerateEtag(addressChange.Shipment.UpdatedAt),
			ShipmentID:  *handlers.FmtUUID(addressChange.ShipmentID),
		}

		err := apperror.NewNotFoundError(uuid.Nil, "unable to create ReviewShipmentAddressChange")

		mockCreator.On("ReviewShipmentAddressChange",
			mock.AnythingOfType("*appcontext.appContext"),
			addressChange.ShipmentID,
			models.ShipmentAddressUpdateStatusApproved,
			"This is a TOO remark",
		).Return(nil, err)

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.ReviewShipmentAddressUpdateNotFound{}, response)
		errResponse := response.(*shipmentops.ReviewShipmentAddressUpdateNotFound)
		suite.NoError(errResponse.Payload.Validate(strfmt.Default))
	})

	suite.Run("500 server error", func() {
		addressChange := factory.BuildShipmentAddressUpdate(suite.DB(), nil, nil)

		mockCreator := mocks.ShipmentAddressUpdateRequester{}

		handlerConfig := suite.HandlerConfig()

		handler := ReviewShipmentAddressUpdateHandler{
			handlerConfig,
			&mockCreator,
		}

		body := shipmentops.ReviewShipmentAddressUpdateBody{
			OfficeRemarks: &officeRemarks,
			Status:        &status,
		}

		req := httptest.NewRequest("PATCH", "/shipments/{mtoShipmentID}/review-shipment-address-update", nil)

		params := shipmentops.ReviewShipmentAddressUpdateParams{
			HTTPRequest: req,
			Body:        body,
			IfMatch:     etag.GenerateEtag(addressChange.Shipment.UpdatedAt),
			ShipmentID:  *handlers.FmtUUID(addressChange.ShipmentID),
		}

		err := apperror.NewQueryError("", nil, "unable to reach database")

		mockCreator.On("ReviewShipmentAddressChange",
			mock.AnythingOfType("*appcontext.appContext"),
			addressChange.ShipmentID,
			models.ShipmentAddressUpdateStatusApproved,
			"This is a TOO remark",
		).Return(nil, err)

		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&shipmentops.ReviewShipmentAddressUpdateInternalServerError{}, response)
	})
}

func (suite *HandlerSuite) TestApproveSITExtensionHandler() {
	suite.Run("Returns 200 and updates SIT days allowance when validations pass", func() {
		sitDaysAllowance := 20
		move := factory.BuildApprovalsRequestedMove(suite.DB(), []factory.Customization{
			{
				Model: models.Entitlement{
					StorageInTransit: &sitDaysAllowance,
				},
			},
		}, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					SITDaysAllowance: &sitDaysAllowance,
					Status:           models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		year, month, day := time.Now().Date()
		lastMonthEntry := time.Date(year, month, day-37, 0, 0, 0, 0, time.UTC)
		lastMonthDeparture := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					SITEntryDate:     &lastMonthEntry,
					SITDepartureDate: &lastMonthDeparture,
					Status:           models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		moveRouter := moveservices.NewMoveRouter()
		sitExtensionApprover := sitextension.NewSITExtensionApprover(moveRouter)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-extension/%s/approve", mtoShipment.ID.String(), sitExtension.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
		mockSender := suite.TestNotificationSender()
		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()

		noCheckUpdater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}

		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		sitExtensionShipmentUpdater := shipmentorchestrator.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)

		handler := ApproveSITExtensionHandler{
			handlerConfig,
			sitExtensionApprover,
			sitstatus.NewShipmentSITStatus(),
			sitExtensionShipmentUpdater,
		}
		approvedDays := int64(10)
		requestReason := "AWAITING_COMPLETION_OF_RESIDENCE"
		officeRemarks := "new office remarks"
		approveParams := shipmentops.ApproveSITExtensionParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.ApproveSITExtension{
				ApprovedDays:  &approvedDays,
				RequestReason: requestReason,
				OfficeRemarks: &officeRemarks,
			},
			ShipmentID:     *handlers.FmtUUID(mtoShipment.ID),
			SitExtensionID: *handlers.FmtUUID(sitExtension.ID),
		}

		// Validate incoming payload
		suite.NoError(approveParams.Body.Validate(strfmt.Default))

		response := handler.Handle(approveParams)
		suite.IsType(&shipmentops.ApproveSITExtensionOK{}, response)
		okResponse := response.(*shipmentops.ApproveSITExtensionOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(int64(30), *payload.SitDaysAllowance)
		suite.Equal("APPROVED", payload.SitExtensions[0].Status)
		suite.Require().NotNil(payload.SitExtensions[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitExtensions[0].OfficeRemarks)
	})
}

func (suite *HandlerSuite) TestDenySITExtensionHandler() {
	suite.Run("Returns 200 when validations pass", func() {
		sitDaysAllowance := 20
		move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					SITDaysAllowance: &sitDaysAllowance,
					Status:           models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		sitExtension := factory.BuildSITDurationUpdate(suite.DB(), []factory.Customization{
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		moveRouter := moveservices.NewMoveRouter()
		sitExtensionDenier := sitextension.NewSITExtensionDenier(moveRouter)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-extension/%s/deny", mtoShipment.ID.String(), sitExtension.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := DenySITExtensionHandler{
			handlerConfig,
			sitExtensionDenier,
			sitstatus.NewShipmentSITStatus(),
		}
		officeRemarks := "new office remarks on denial of extension"
		denyParams := shipmentops.DenySITExtensionParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.DenySITExtension{
				OfficeRemarks:            &officeRemarks,
				ConvertToCustomerExpense: models.BoolPointer(false),
			},
			ShipmentID:     *handlers.FmtUUID(mtoShipment.ID),
			SitExtensionID: *handlers.FmtUUID(sitExtension.ID),
		}

		// Validate incoming payload
		suite.NoError(denyParams.Body.Validate(strfmt.Default))

		response := handler.Handle(denyParams)
		suite.IsType(&shipmentops.DenySITExtensionOK{}, response)
		okResponse := response.(*shipmentops.DenySITExtensionOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal("DENIED", payload.SitExtensions[0].Status)
	})
}

func (suite *HandlerSuite) CreateApprovedSITDurationUpdate() {
	suite.Run("Returns 200, creates new SIT extension, and updates SIT days allowance on shipment without an allowance when validations pass", func() {
		mtoShipment := factory.BuildMTOShipment(suite.DB(), nil, nil)

		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approvedSITDurationUpdateCreator := sitextension.NewApprovedSITDurationUpdateCreator()
		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/sit-extension/", mtoShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
		mockSender := suite.TestNotificationSender()
		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()
		moveRouter := moveservices.NewMoveRouter()

		noCheckUpdater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}

		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		sitExtensionShipmentUpdater := shipmentorchestrator.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)

		handler := CreateApprovedSITDurationUpdateHandler{
			handlerConfig,
			approvedSITDurationUpdateCreator,
			sitstatus.NewShipmentSITStatus(),
			sitExtensionShipmentUpdater,
		}
		approvedDays := int64(10)
		officeRemarks := "new office remarks"
		requestReason := "OTHER"
		createParams := shipmentops.CreateApprovedSITDurationUpdateParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.CreateApprovedSITDurationUpdate{
				ApprovedDays:  &approvedDays,
				OfficeRemarks: &officeRemarks,
				RequestReason: &requestReason,
			},
			ShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&shipmentops.CreateApprovedSITDurationUpdateOK{}, response)
		okResponse := response.(*shipmentops.CreateApprovedSITDurationUpdateOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(int64(10), *payload.SitDaysAllowance)
		suite.Equal("APPROVED", payload.SitExtensions[0].Status)
		suite.Require().NotNil(payload.SitExtensions[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitExtensions[0].OfficeRemarks)
	})

	suite.Run("Returns 200, creates new SIT extension, and updates SIT days allowance on shipment that already has an allowance when validations pass", func() {
		sitDaysAllowance := 20
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					SITDaysAllowance: &sitDaysAllowance,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		approvedSITDurationUpdateCreator := sitextension.NewApprovedSITDurationUpdateCreator()
		req := httptest.NewRequest("POST", fmt.Sprintf("/shipments/%s/sit-extension/", mtoShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

		// Get shipment payment request recalculator service
		creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
		statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
		recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
		paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)
		mockSender := suite.TestNotificationSender()
		addressUpdater := address.NewAddressUpdater()
		addressCreator := address.NewAddressCreator()
		moveRouter := moveservices.NewMoveRouter()

		noCheckUpdater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}

		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()

		mobilehomeshipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		sitExtensionShipmentUpdater := shipmentorchestrator.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobilehomeshipmentUpdater)

		handler := CreateApprovedSITDurationUpdateHandler{
			handlerConfig,
			approvedSITDurationUpdateCreator,
			sitstatus.NewShipmentSITStatus(),
			sitExtensionShipmentUpdater,
		}
		approvedDays := int64(10)
		officeRemarks := "new office remarks"
		requestReason := "OTHER"
		createParams := shipmentops.CreateApprovedSITDurationUpdateParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.CreateApprovedSITDurationUpdate{
				ApprovedDays:  &approvedDays,
				OfficeRemarks: &officeRemarks,
				RequestReason: &requestReason,
			},
			ShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&shipmentops.CreateApprovedSITDurationUpdateOK{}, response)
		okResponse := response.(*shipmentops.CreateApprovedSITDurationUpdateOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		suite.Equal(int64(30), *payload.SitDaysAllowance)
		suite.Equal("APPROVED", payload.SitExtensions[0].Status)
		suite.Require().NotNil(payload.SitExtensions[0].OfficeRemarks)
		suite.Equal(officeRemarks, *payload.SitExtensions[0].OfficeRemarks)
	})
}

type createMTOShipmentSubtestData struct {
	builder *query.Builder
	params  mtoshipmentops.CreateMTOShipmentParams
	traceID uuid.UUID
}

func (suite *HandlerSuite) makeCreateMTOShipmentSubtestData() (subtestData *createMTOShipmentSubtestData) {
	subtestData = &createMTOShipmentSubtestData{}

	mto := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
	pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
	mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	mtoShipment.MoveTaskOrderID = mto.ID

	subtestData.builder = query.NewQueryBuilder()

	req := httptest.NewRequest("POST", "/mto-shipments", nil)

	// Set the traceID so we can use it to find the webhook notification
	traceID, err := uuid.NewV4()
	suite.FatalNoError(err, "Error creating a new trace ID.")
	req = req.WithContext(trace.NewContext(req.Context(), traceID))

	subtestData.traceID = traceID
	shipmentType := ghcmessages.MTOShipmentType(mtoShipment.ShipmentType)
	subtestData.params = mtoshipmentops.CreateMTOShipmentParams{
		HTTPRequest: req,
		Body: &ghcmessages.CreateMTOShipment{
			MoveTaskOrderID:     handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
			Agents:              nil,
			CustomerRemarks:     handlers.FmtString("customer remark"),
			CounselorRemarks:    handlers.FmtString("counselor remark"),
			RequestedPickupDate: handlers.FmtDatePtr(mtoShipment.RequestedPickupDate),
			ShipmentType:        &shipmentType,
		},
	}
	subtestData.params.Body.DestinationAddress.Address = ghcmessages.Address{
		City:           &destinationAddress.City,
		PostalCode:     &destinationAddress.PostalCode,
		State:          &destinationAddress.State,
		StreetAddress1: &destinationAddress.StreetAddress1,
		StreetAddress2: destinationAddress.StreetAddress2,
		StreetAddress3: destinationAddress.StreetAddress3,
	}
	subtestData.params.Body.PickupAddress.Address = ghcmessages.Address{
		City:           &pickupAddress.City,
		PostalCode:     &pickupAddress.PostalCode,
		State:          &pickupAddress.State,
		StreetAddress1: &pickupAddress.StreetAddress1,
		StreetAddress2: pickupAddress.StreetAddress2,
		StreetAddress3: pickupAddress.StreetAddress3,
	}

	return subtestData
}

func (suite *HandlerSuite) TestCreateMTOShipmentHandler() {
	moveRouter := moveservices.NewMoveRouter()
	addressCreator := address.NewAddressCreator()

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	suite.Run("Successful POST - Integration Test", func() {
		handlerConfig := suite.HandlerConfig()

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator, boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		createMTOShipmentPayload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(createMTOShipmentPayload.Validate(strfmt.Default))

		suite.Require().Equal(ghcmessages.MTOShipmentStatusSUBMITTED, createMTOShipmentPayload.Status, "MTO Shipment should have been submitted")
		suite.Require().Equal(createMTOShipmentPayload.ShipmentType, ghcmessages.MTOShipmentTypeHHG, "MTO Shipment should be an HHG")
		suite.Equal(int64(models.DefaultServiceMemberSITDaysAllowance), *createMTOShipmentPayload.SitDaysAllowance)
		suite.Equal(string("customer remark"), *createMTOShipmentPayload.CustomerRemarks)
		suite.Equal(string("counselor remark"), *createMTOShipmentPayload.CounselorRemarks)
	})

	suite.Run("POST failure - 500", func() {
		handlerConfig := suite.HandlerConfig()

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		params := subtestData.params

		shipmentCreator := mocks.ShipmentCreator{}
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			&shipmentCreator,
			sitStatus,
		}

		err := errors.New("ServerError")

		shipmentCreator.On("CreateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("*models.MTOShipment"),
		).Return(nil, err)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentInternalServerError{}, response)
		payload := response.(*mtoshipmentops.CreateMTOShipmentInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("POST failure - 422 -- Bad agent IDs set on shipment", func() {
		handlerConfig := suite.HandlerConfig()

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator, boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		badID := params.Body.MoveTaskOrderID
		agent := &ghcmessages.MTOAgent{
			ID:            *badID,
			MtoShipmentID: *badID,
			FirstName:     handlers.FmtString("Mary"),
		}

		paramsBadIDs := params
		paramsBadIDs.Body.Agents = ghcmessages.MTOAgents{agent}

		// Validate incoming payload
		suite.NoError(paramsBadIDs.Body.Validate(strfmt.Default))

		response := handler.Handle(paramsBadIDs)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		suite.NoError(typedResponse.Payload.Validate(strfmt.Default))

		suite.NotEmpty(typedResponse.Payload.InvalidFields)
	})

	suite.Run("POST failure - 422 - invalid input, missing pickup address", func() {
		handlerConfig := suite.HandlerConfig()

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
		boatCreator := boatshipment.NewBoatShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator, boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		badParams := params
		badParams.Body.PickupAddress.Address.StreetAddress1 = nil

		// Validate incoming payload
		suite.NoError(badParams.Body.Validate(strfmt.Default))

		response := handler.Handle(badParams)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentUnprocessableEntity{}, response)
		typedResponse := response.(*mtoshipmentops.CreateMTOShipmentUnprocessableEntity)

		// Validate outgoing payload
		// TODO: Can't validate the response because of the issue noted below. Figure out a way to
		//   either alter the service or relax the swagger requirements.
		// suite.NoError(typedResponse.Payload.Validate(strfmt.Default))
		// CreateShipment is returning apperror.InvalidInputError without any validation errors
		// so InvalidFields won't be added to the payload.
		suite.Empty(typedResponse.Payload.InvalidFields)
	})

	suite.Run("POST failure - 404 -- not found", func() {
		handlerConfig := suite.HandlerConfig()

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder
		params := subtestData.params

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator, boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		uuidString := "d874d002-5582-4a91-97d3-786e8f66c763"
		badParams := params
		badParams.Body.MoveTaskOrderID = handlers.FmtUUID(uuid.FromStringOrNil(uuidString))

		// Validate incoming payload
		suite.NoError(badParams.Body.Validate(strfmt.Default))

		response := handler.Handle(badParams)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentNotFound{}, response)
		payload := response.(*mtoshipmentops.CreateMTOShipmentNotFound).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("POST failure - 400 -- nil body", func() {
		handlerConfig := suite.HandlerConfig()

		subtestData := suite.makeCreateMTOShipmentSubtestData()
		builder := subtestData.builder

		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveRouter, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter, ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveRouter, setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator, boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		paramsNilBody := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
		}

		// Validate incoming payload: nil body (the point of this test)

		response := handler.Handle(paramsNilBody)

		suite.IsType(&mtoshipmentops.CreateMTOShipmentBadRequest{}, response)
		payload := response.(*mtoshipmentops.CreateMTOShipmentBadRequest).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) TestCreateMTOShipmentHandlerUsingPPM() {
	addressCreator := address.NewAddressCreator()

	setUpSignedCertificationCreatorMock := func(returnValue ...interface{}) services.SignedCertificationCreator {
		mockCreator := &mocks.SignedCertificationCreator{}

		mockCreator.On(
			"CreateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
		).Return(returnValue...)

		return mockCreator
	}

	setUpSignedCertificationUpdaterMock := func(returnValue ...interface{}) services.SignedCertificationUpdater {
		mockUpdater := &mocks.SignedCertificationUpdater{}

		mockUpdater.On(
			"UpdateSignedCertification",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.SignedCertification"),
			mock.AnythingOfType("string"),
		).Return(returnValue...)

		return mockUpdater
	}

	suite.Run("Successful POST - Integration Test (PPM, all fields)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		move := factory.BuildMove(suite.DB(), nil, nil)
		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		handlerConfig := suite.HandlerConfig()
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveservices.NewMoveRouter(), addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmCreator := ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator)
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)

		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveservices.NewMoveRouter(), ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveservices.NewMoveRouter(), setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmCreator, boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		sitStatus := sitstatus.NewShipmentSITStatus()
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitStatus,
		}

		shipmentType := ghcmessages.MTOShipmentTypePPM
		expectedDepartureDate := hhgShipment.RequestedPickupDate
		sitExpected := true
		sitLocation := ghcmessages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1700)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3000)
		hasProGear := true
		proGearWeight := unit.Pound(300)
		spouseProGearWeight := unit.Pound(200)
		estimatedIncentive := 654321
		sitEstimatedCost := 67500

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		var pickupAddress ghcmessages.Address
		var secondaryPickupAddress ghcmessages.Address
		var destinationAddress ghcmessages.PPMDestinationAddress
		var secondaryDestinationAddress ghcmessages.Address

		expectedPickupAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		pickupAddress = ghcmessages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedSecondaryPickupAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		secondaryPickupAddress = ghcmessages.Address{
			City:           &expectedSecondaryPickupAddress.City,
			PostalCode:     &expectedSecondaryPickupAddress.PostalCode,
			State:          &expectedSecondaryPickupAddress.State,
			StreetAddress1: &expectedSecondaryPickupAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryPickupAddress.StreetAddress2,
			StreetAddress3: expectedSecondaryPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		destinationAddress = ghcmessages.PPMDestinationAddress{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}

		expectedSecondaryDestinationAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		secondaryDestinationAddress = ghcmessages.Address{
			City:           &expectedSecondaryDestinationAddress.City,
			PostalCode:     &expectedSecondaryDestinationAddress.PostalCode,
			State:          &expectedSecondaryDestinationAddress.State,
			StreetAddress1: &expectedSecondaryDestinationAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryDestinationAddress.StreetAddress2,
			StreetAddress3: expectedSecondaryDestinationAddress.StreetAddress3,
		}

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateMTOShipment{
				MoveTaskOrderID: handlers.FmtUUID(move.ID),
				ShipmentType:    &shipmentType,
				PpmShipment: &ghcmessages.CreatePPMShipment{
					ExpectedDepartureDate:  handlers.FmtDatePtr(expectedDepartureDate),
					PickupAddress:          struct{ ghcmessages.Address }{pickupAddress},
					SecondaryPickupAddress: struct{ ghcmessages.Address }{secondaryPickupAddress},
					DestinationAddress: struct {
						ghcmessages.PPMDestinationAddress
					}{destinationAddress},
					SecondaryDestinationAddress: struct{ ghcmessages.Address }{secondaryDestinationAddress},
					SitExpected:                 &sitExpected,
					SitLocation:                 &sitLocation,
					SitEstimatedWeight:          handlers.FmtPoundPtr(&sitEstimatedWeight),
					SitEstimatedEntryDate:       handlers.FmtDate(sitEstimatedEntryDate),
					SitEstimatedDepartureDate:   handlers.FmtDate(sitEstimatedDepartureDate),
					EstimatedWeight:             handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:                  &hasProGear,
					ProGearWeight:               handlers.FmtPoundPtr(&proGearWeight),
					SpouseProGearWeight:         handlers.FmtPoundPtr(&spouseProGearWeight),
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		// Check MTOShipment fields.
		suite.NotZero(payload.ID)
		suite.NotEqual(uuid.Nil.String(), payload.ID.String())
		suite.Equal(move.ID.String(), payload.MoveTaskOrderID.String())
		suite.Equal(ghcmessages.MTOShipmentTypePPM, payload.ShipmentType)
		suite.Equal(ghcmessages.MTOShipmentStatusSUBMITTED, payload.Status)
		suite.NotZero(payload.CreatedAt)
		suite.NotZero(payload.UpdatedAt)

		// Check PPMShipment fields.
		ppmPayload := payload.PpmShipment
		if suite.NotNil(ppmPayload) {
			suite.NotZero(ppmPayload.ID)
			suite.NotEqual(uuid.Nil.String(), ppmPayload.ID.String())
			suite.EqualDatePtr(expectedDepartureDate, ppmPayload.ExpectedDepartureDate)
			suite.Equal(expectedPickupAddress.PostalCode, *ppmPayload.PickupAddress.PostalCode)
			suite.Equal(&expectedSecondaryPickupAddress.PostalCode, ppmPayload.SecondaryPickupAddress.PostalCode)
			suite.Equal(expectedDestinationAddress.PostalCode, *ppmPayload.DestinationAddress.PostalCode)
			suite.Equal(&expectedSecondaryDestinationAddress.PostalCode, ppmPayload.SecondaryDestinationAddress.PostalCode)
			suite.NotNil(ppmPayload.PickupAddress)
			suite.NotNil(ppmPayload.DestinationAddress)
			suite.NotNil(ppmPayload.SecondaryPickupAddress)
			suite.NotNil(ppmPayload.SecondaryDestinationAddress)
			suite.Equal(sitExpected, *ppmPayload.SitExpected)
			suite.Equal(&sitLocation, ppmPayload.SitLocation)
			suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), ppmPayload.SitEstimatedWeight)
			suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), ppmPayload.SitEstimatedEntryDate)
			suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), ppmPayload.SitEstimatedDepartureDate)
			suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), ppmPayload.EstimatedWeight)
			suite.Equal(&hasProGear, ppmPayload.HasProGear)
			suite.Equal(handlers.FmtPoundPtr(&proGearWeight), ppmPayload.ProGearWeight)
			suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), ppmPayload.SpouseProGearWeight)
			suite.Equal(ghcmessages.PPMShipmentStatusSUBMITTED, ppmPayload.Status)
			suite.Equal(int64(estimatedIncentive), *ppmPayload.EstimatedIncentive)
			suite.Equal(int64(sitEstimatedCost), *ppmPayload.SitEstimatedCost)
			suite.NotZero(ppmPayload.CreatedAt)
		}
	})

	suite.Run("Successful POST - Integration Test (PPM, minimal fields)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		move := factory.BuildMove(suite.DB(), nil, nil)
		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		handlerConfig := suite.HandlerConfig()
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveservices.NewMoveRouter(), addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveservices.NewMoveRouter(), ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveservices.NewMoveRouter(), setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator), boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitstatus.NewShipmentSITStatus(),
		}

		shipmentType := ghcmessages.MTOShipmentTypePPM
		expectedDepartureDate := hhgShipment.RequestedPickupDate
		sitExpected := false
		estimatedWeight := unit.Pound(2450)
		hasProGear := false
		estimatedIncentive := 123456

		var pickupAddress ghcmessages.Address
		var destinationAddress ghcmessages.PPMDestinationAddress

		expectedPickupAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		pickupAddress = ghcmessages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		destinationAddress = ghcmessages.PPMDestinationAddress{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateMTOShipment{
				MoveTaskOrderID: handlers.FmtUUID(move.ID),
				ShipmentType:    &shipmentType,
				PpmShipment: &ghcmessages.CreatePPMShipment{
					ExpectedDepartureDate: handlers.FmtDatePtr(expectedDepartureDate),
					PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
					DestinationAddress: struct {
						ghcmessages.PPMDestinationAddress
					}{destinationAddress},
					SitExpected:     &sitExpected,
					EstimatedWeight: handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:      &hasProGear,
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), nil, nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		// Check MTOShipment fields.
		suite.NotZero(payload.ID)
		suite.NotEqual(uuid.Nil.String(), payload.ID.String())
		suite.Equal(move.ID.String(), payload.MoveTaskOrderID.String())
		suite.Equal(ghcmessages.MTOShipmentTypePPM, payload.ShipmentType)
		suite.Equal(ghcmessages.MTOShipmentStatusSUBMITTED, payload.Status)
		suite.NotZero(payload.CreatedAt)
		suite.NotZero(payload.UpdatedAt)

		// Check PPMShipment fields.
		ppmPayload := payload.PpmShipment
		if suite.NotNil(ppmPayload) {
			suite.NotZero(ppmPayload.ID)
			suite.NotEqual(uuid.Nil.String(), ppmPayload.ID.String())
			suite.EqualDatePtr(expectedDepartureDate, ppmPayload.ExpectedDepartureDate)
			suite.Equal(expectedPickupAddress.PostalCode, *ppmPayload.PickupAddress.PostalCode)
			suite.Equal(expectedDestinationAddress.PostalCode, *ppmPayload.DestinationAddress.PostalCode)
			suite.Nil(ppmPayload.SecondaryPickupAddress)
			suite.Nil(ppmPayload.SecondaryDestinationAddress)
			suite.Equal(sitExpected, *ppmPayload.SitExpected)
			suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), ppmPayload.EstimatedWeight)
			suite.Equal(&hasProGear, ppmPayload.HasProGear)
			suite.Equal(ghcmessages.PPMShipmentStatusSUBMITTED, ppmPayload.Status)
			suite.Equal(int64(estimatedIncentive), *ppmPayload.EstimatedIncentive)
			suite.Nil(ppmPayload.SitEstimatedCost)
			suite.NotZero(ppmPayload.CreatedAt)
		}
	})

	suite.Run("Successful POST and Patch for delete of addresses - Integration Test (PPM, minimal fields)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		move := factory.BuildMove(suite.DB(), nil, nil)
		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		handlerConfig := suite.HandlerConfig()
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		creator := mtoshipment.NewMTOShipmentCreatorV1(builder, fetcher, moveservices.NewMoveRouter(), addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		boatCreator := boatshipment.NewBoatShipmentCreator()
		mobileHomeCreator := mobilehomeshipment.NewMobileHomeShipmentCreator()
		shipmentRouter := mtoshipment.NewShipmentRouter()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		moveTaskOrderUpdater := movetaskorder.NewMoveTaskOrderUpdater(
			builder,
			mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveservices.NewMoveRouter(), ghcrateengine.NewDomesticUnpackPricer(), ghcrateengine.NewDomesticPackPricer(), ghcrateengine.NewDomesticLinehaulPricer(), ghcrateengine.NewDomesticShorthaulPricer(), ghcrateengine.NewDomesticOriginPricer(), ghcrateengine.NewDomesticDestinationPricer(), ghcrateengine.NewFuelSurchargePricer()),
			moveservices.NewMoveRouter(), setUpSignedCertificationCreatorMock(nil, nil), setUpSignedCertificationUpdaterMock(nil, nil), &ppmEstimator,
		)
		shipmentCreator := shipmentorchestrator.NewShipmentCreator(creator, ppmshipment.NewPPMShipmentCreator(&ppmEstimator, addressCreator), boatCreator, mobileHomeCreator, shipmentRouter, moveTaskOrderUpdater)
		handler := CreateMTOShipmentHandler{
			handlerConfig,
			shipmentCreator,
			sitstatus.NewShipmentSITStatus(),
		}

		shipmentType := ghcmessages.MTOShipmentTypePPM
		expectedDepartureDate := hhgShipment.RequestedPickupDate
		sitExpected := false
		estimatedWeight := unit.Pound(2450)
		hasProGear := false
		estimatedIncentive := 123456

		var pickupAddress ghcmessages.Address
		var destinationAddress ghcmessages.PPMDestinationAddress

		expectedPickupAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		pickupAddress = ghcmessages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: expectedPickupAddress.StreetAddress3,
		}

		expectedDestinationAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		destinationAddress = ghcmessages.PPMDestinationAddress{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: expectedDestinationAddress.StreetAddress3,
		}

		req := httptest.NewRequest("POST", "/mto-shipments", nil)

		params := mtoshipmentops.CreateMTOShipmentParams{
			HTTPRequest: req,
			Body: &ghcmessages.CreateMTOShipment{
				MoveTaskOrderID: handlers.FmtUUID(move.ID),
				ShipmentType:    &shipmentType,
				PpmShipment: &ghcmessages.CreatePPMShipment{
					ExpectedDepartureDate: handlers.FmtDatePtr(expectedDepartureDate),
					PickupAddress:         struct{ ghcmessages.Address }{pickupAddress},
					DestinationAddress: struct {
						ghcmessages.PPMDestinationAddress
					}{destinationAddress},
					SitExpected:     &sitExpected,
					EstimatedWeight: handlers.FmtPoundPtr(&estimatedWeight),
					HasProGear:      &hasProGear,
				},
			},
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), nil, nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.CreateMTOShipmentOK{}, response)
		okResponse := response.(*mtoshipmentops.CreateMTOShipmentOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))

		// Check MTOShipment fields.
		suite.NotZero(payload.ID)
		suite.NotEqual(uuid.Nil.String(), payload.ID.String())
		suite.Equal(move.ID.String(), payload.MoveTaskOrderID.String())
		suite.Equal(ghcmessages.MTOShipmentTypePPM, payload.ShipmentType)
		suite.Equal(ghcmessages.MTOShipmentStatusSUBMITTED, payload.Status)
		suite.NotZero(payload.CreatedAt)
		suite.NotZero(payload.UpdatedAt)

		// Check PPMShipment fields.
		ppmPayload := payload.PpmShipment
		if suite.NotNil(ppmPayload) {
			suite.NotZero(ppmPayload.ID)
			suite.NotEqual(uuid.Nil.String(), ppmPayload.ID.String())
			suite.EqualDatePtr(expectedDepartureDate, ppmPayload.ExpectedDepartureDate)
			suite.Equal(expectedPickupAddress.PostalCode, *ppmPayload.PickupAddress.PostalCode)
			suite.Equal(expectedDestinationAddress.PostalCode, *ppmPayload.DestinationAddress.PostalCode)
			suite.Nil(ppmPayload.SecondaryPickupAddress)
			suite.Nil(ppmPayload.SecondaryDestinationAddress)
			suite.Equal(sitExpected, *ppmPayload.SitExpected)
			suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), ppmPayload.EstimatedWeight)
			suite.Equal(&hasProGear, ppmPayload.HasProGear)
			suite.Equal(ghcmessages.PPMShipmentStatusSUBMITTED, ppmPayload.Status)
			suite.Equal(int64(estimatedIncentive), *ppmPayload.EstimatedIncentive)
			suite.Nil(ppmPayload.SitEstimatedCost)
			suite.NotZero(ppmPayload.CreatedAt)
		}
	})

}

func (suite *HandlerSuite) getUpdateShipmentParams(originalShipment models.MTOShipment) mtoshipmentops.UpdateMTOShipmentParams {
	servicesCounselor := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
	servicesCounselor.User.Roles = append(servicesCounselor.User.Roles, roles.Role{
		RoleType: roles.RoleTypeServicesCounselor,
	})
	pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
	pickupAddress.StreetAddress1 = "123 Fake Test St NW"
	destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)
	destinationAddress.StreetAddress1 = "54321 Test Fake Rd SE"
	customerRemarks := "help"
	counselorRemarks := "counselor approved"
	billableWeightCap := int64(8000)
	billableWeightJustification := "Unable to perform reweigh because shipment was already unloaded."
	mtoAgent := factory.BuildMTOAgent(suite.DB(), nil, nil)
	agents := ghcmessages.MTOAgents{&ghcmessages.MTOAgent{
		FirstName: mtoAgent.FirstName,
		LastName:  mtoAgent.LastName,
		Email:     mtoAgent.Email,
		Phone:     mtoAgent.Phone,
		AgentType: string(mtoAgent.MTOAgentType),
	}}

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", originalShipment.MoveTaskOrderID.String(), originalShipment.ID.String()), nil)
	req = suite.AuthenticateOfficeRequest(req, servicesCounselor)

	eTag := etag.GenerateEtag(originalShipment.UpdatedAt)

	now := strfmt.Date(time.Now())
	payload := ghcmessages.UpdateShipment{
		BillableWeightJustification: &billableWeightJustification,
		BillableWeightCap:           &billableWeightCap,
		RequestedPickupDate:         &now,
		RequestedDeliveryDate:       &now,
		ShipmentType:                ghcmessages.MTOShipmentTypeHHG,
		CustomerRemarks:             &customerRemarks,
		CounselorRemarks:            &counselorRemarks,
		Agents:                      agents,
		TacType:                     nullable.NewString("NTS"),
		SacType:                     nullable.NewString(""),
	}
	payload.DestinationAddress.Address = ghcmessages.Address{
		City:           &destinationAddress.City,
		PostalCode:     &destinationAddress.PostalCode,
		State:          &destinationAddress.State,
		StreetAddress1: &destinationAddress.StreetAddress1,
		StreetAddress2: destinationAddress.StreetAddress2,
		StreetAddress3: destinationAddress.StreetAddress3,
	}
	payload.PickupAddress.Address = ghcmessages.Address{
		City:           &pickupAddress.City,
		PostalCode:     &pickupAddress.PostalCode,
		State:          &pickupAddress.State,
		StreetAddress1: &pickupAddress.StreetAddress1,
		StreetAddress2: pickupAddress.StreetAddress2,
		StreetAddress3: pickupAddress.StreetAddress3,
	}

	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest: req,
		ShipmentID:  *handlers.FmtUUID(originalShipment.ID),
		Body:        &payload,
		IfMatch:     eTag,
	}

	return params
}

func (suite *HandlerSuite) TestUpdateShipmentHandler() {
	addressUpdater := address.NewAddressUpdater()
	addressCreator := address.NewAddressCreator()

	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())

	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(planner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

	suite.Run("Successful PATCH - Integration Test", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}

		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)
		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		hhgLOAType := models.LOATypeHHG
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:                    models.MTOShipmentStatusSubmitted,
					UsesExternalVendor:        true,
					TACType:                   &hhgLOAType,
					Diversion:                 true,
					ActualProGearWeight:       models.PoundPointer(1000),
					ActualSpouseProGearWeight: models.PoundPointer(253),
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		suite.Equal(oldShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(oldShipment.ActualProGearWeight, handlers.PoundPtrFromInt64Ptr(updatedShipment.ActualProGearWeight))
		suite.Equal(oldShipment.ActualSpouseProGearWeight, handlers.PoundPtrFromInt64Ptr(updatedShipment.ActualSpouseProGearWeight))
		suite.Equal(params.Body.BillableWeightCap, updatedShipment.BillableWeightCap)
		suite.Equal(params.Body.BillableWeightJustification, updatedShipment.BillableWeightJustification)
		suite.Equal(params.Body.CounselorRemarks, updatedShipment.CounselorRemarks)
		suite.Equal(params.Body.PickupAddress.StreetAddress1, updatedShipment.PickupAddress.StreetAddress1)
		suite.Equal(params.Body.DestinationAddress.StreetAddress1, updatedShipment.DestinationAddress.StreetAddress1)
		suite.Equal(params.Body.RequestedPickupDate.String(), updatedShipment.RequestedPickupDate.String())
		suite.Equal(params.Body.Agents[0].FirstName, updatedShipment.MtoAgents[0].FirstName)
		suite.Equal(params.Body.Agents[0].LastName, updatedShipment.MtoAgents[0].LastName)
		suite.Equal(params.Body.Agents[0].Email, updatedShipment.MtoAgents[0].Email)
		suite.Equal(params.Body.Agents[0].Phone, updatedShipment.MtoAgents[0].Phone)
		suite.Equal(params.Body.Agents[0].AgentType, updatedShipment.MtoAgents[0].AgentType)
		suite.Equal(oldShipment.ID.String(), string(updatedShipment.MtoAgents[0].MtoShipmentID))
		suite.NotEmpty(updatedShipment.MtoAgents[0].ID)
		suite.Equal(params.Body.RequestedDeliveryDate.String(), updatedShipment.RequestedDeliveryDate.String())
		suite.Equal(*params.Body.TacType.Value, string(*updatedShipment.TacType))
		suite.Nil(updatedShipment.SacType)

		// don't update non-nullable booleans if they're not passed in
		suite.Equal(oldShipment.Diversion, updatedShipment.Diversion)
		suite.Equal(oldShipment.UsesExternalVendor, updatedShipment.UsesExternalVendor)
	})

	suite.Run("Successful PATCH - Integration Test (PPM)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		hasProGear := true
		ppmShipment := factory.BuildMinimalPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					HasProGear: &hasProGear,
				},
			},
		}, nil)
		year, month, day := time.Now().Date()
		actualMoveDate := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
		expectedDepartureDate := actualMoveDate.Add(time.Hour * 24 * 2)

		// we expect initial setup data to have no secondary addresses
		suite.Nil(ppmShipment.SecondaryPickupAddress)
		suite.Nil(ppmShipment.SecondaryDestinationAddress)

		expectedPickupAddressStreet3 := "HelloWorld1"
		expectedSecondaryPickupAddressStreet3 := "HelloWorld2"
		expectedDestinationAddressStreet3 := "HelloWorld3"
		expectedSecondaryDestinationAddressStreet3 := "HelloWorld4"

		var pickupAddress ghcmessages.Address
		var secondaryPickupAddress ghcmessages.Address
		var destinationAddress ghcmessages.PPMDestinationAddress
		var secondaryDestinationAddress ghcmessages.Address

		expectedPickupAddress := ppmShipment.PickupAddress
		pickupAddress = ghcmessages.Address{
			City:           &expectedPickupAddress.City,
			PostalCode:     &expectedPickupAddress.PostalCode,
			State:          &expectedPickupAddress.State,
			StreetAddress1: &expectedPickupAddress.StreetAddress1,
			StreetAddress2: expectedPickupAddress.StreetAddress2,
			StreetAddress3: &expectedPickupAddressStreet3,
		}

		expectedSecondaryPickupAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		secondaryPickupAddress = ghcmessages.Address{
			City:           &expectedSecondaryPickupAddress.City,
			PostalCode:     &expectedSecondaryPickupAddress.PostalCode,
			State:          &expectedSecondaryPickupAddress.State,
			StreetAddress1: &expectedSecondaryPickupAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryPickupAddress.StreetAddress2,
			StreetAddress3: &expectedSecondaryPickupAddressStreet3,
		}

		expectedDestinationAddress := ppmShipment.DestinationAddress
		destinationAddress = ghcmessages.PPMDestinationAddress{
			City:           &expectedDestinationAddress.City,
			PostalCode:     &expectedDestinationAddress.PostalCode,
			State:          &expectedDestinationAddress.State,
			StreetAddress1: &expectedDestinationAddress.StreetAddress1,
			StreetAddress2: expectedDestinationAddress.StreetAddress2,
			StreetAddress3: &expectedDestinationAddressStreet3,
		}

		expectedSecondaryDestinationAddress := factory.BuildAddress(nil, []factory.Customization{
			{
				Model: models.Address{
					ID: uuid.Must(uuid.NewV4()),
				},
			},
		}, nil)
		secondaryDestinationAddress = ghcmessages.Address{
			City:           &expectedSecondaryDestinationAddress.City,
			PostalCode:     &expectedSecondaryDestinationAddress.PostalCode,
			State:          &expectedSecondaryDestinationAddress.State,
			StreetAddress1: &expectedSecondaryDestinationAddress.StreetAddress1,
			StreetAddress2: expectedSecondaryDestinationAddress.StreetAddress2,
			StreetAddress3: &expectedSecondaryDestinationAddressStreet3,
		}

		sitExpected := true
		sitLocation := ghcmessages.SITLocationTypeDESTINATION
		sitEstimatedWeight := unit.Pound(1700)
		sitEstimatedEntryDate := expectedDepartureDate.AddDate(0, 0, 5)
		sitEstimatedDepartureDate := sitEstimatedEntryDate.AddDate(0, 0, 20)
		estimatedWeight := unit.Pound(3000)
		proGearWeight := unit.Pound(300)
		spouseProGearWeight := unit.Pound(200)
		estimatedIncentive := 654321
		sitEstimatedCost := 67500

		params := suite.getUpdateShipmentParams(ppmShipment.Shipment)
		params.Body.ShipmentType = ghcmessages.MTOShipmentTypePPM
		params.Body.PpmShipment = &ghcmessages.UpdatePPMShipment{
			ActualMoveDate:         handlers.FmtDatePtr(&actualMoveDate),
			ExpectedDepartureDate:  handlers.FmtDatePtr(&expectedDepartureDate),
			PickupAddress:          struct{ ghcmessages.Address }{pickupAddress},
			SecondaryPickupAddress: struct{ ghcmessages.Address }{secondaryPickupAddress},
			DestinationAddress: struct {
				ghcmessages.PPMDestinationAddress
			}{destinationAddress},
			SecondaryDestinationAddress: struct{ ghcmessages.Address }{secondaryDestinationAddress},
			SitExpected:                 &sitExpected,
			SitEstimatedWeight:          handlers.FmtPoundPtr(&sitEstimatedWeight),
			SitEstimatedEntryDate:       handlers.FmtDatePtr(&sitEstimatedEntryDate),
			SitEstimatedDepartureDate:   handlers.FmtDatePtr(&sitEstimatedDepartureDate),
			SitLocation:                 &sitLocation,
			EstimatedWeight:             handlers.FmtPoundPtr(&estimatedWeight),
			HasProGear:                  &hasProGear,
			ProGearWeight:               handlers.FmtPoundPtr(&proGearWeight),
			SpouseProGearWeight:         handlers.FmtPoundPtr(&spouseProGearWeight),
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))

		suite.Equal(ppmShipment.Shipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(handlers.FmtDatePtr(&actualMoveDate), updatedShipment.PpmShipment.ActualMoveDate)
		suite.Equal(handlers.FmtDatePtr(&expectedDepartureDate), updatedShipment.PpmShipment.ExpectedDepartureDate)
		suite.NotNil(updatedShipment.PpmShipment.PickupAddress)
		suite.NotNil(updatedShipment.PpmShipment.DestinationAddress)

		// expect secondary addresses to be added
		suite.NotNil(updatedShipment.PpmShipment.SecondaryPickupAddress)
		suite.NotNil(updatedShipment.PpmShipment.SecondaryDestinationAddress)
		suite.Equal(expectedPickupAddressStreet3, *updatedShipment.PpmShipment.PickupAddress.StreetAddress3)
		suite.Equal(expectedSecondaryPickupAddressStreet3, *updatedShipment.PpmShipment.SecondaryPickupAddress.StreetAddress3)
		suite.Equal(expectedDestinationAddressStreet3, *updatedShipment.PpmShipment.DestinationAddress.StreetAddress3)
		suite.Equal(expectedSecondaryDestinationAddressStreet3, *updatedShipment.PpmShipment.SecondaryDestinationAddress.StreetAddress3)
		suite.Equal(pickupAddress.PostalCode, updatedShipment.PpmShipment.PickupAddress.PostalCode)
		suite.Equal(secondaryPickupAddress.PostalCode, updatedShipment.PpmShipment.SecondaryPickupAddress.PostalCode)
		suite.Equal(destinationAddress.PostalCode, updatedShipment.PpmShipment.DestinationAddress.PostalCode)
		suite.Equal(secondaryDestinationAddress.PostalCode, updatedShipment.PpmShipment.SecondaryDestinationAddress.PostalCode)

		suite.Equal(sitExpected, *updatedShipment.PpmShipment.SitExpected)
		suite.Equal(&sitLocation, updatedShipment.PpmShipment.SitLocation)
		suite.Equal(handlers.FmtPoundPtr(&sitEstimatedWeight), updatedShipment.PpmShipment.SitEstimatedWeight)
		suite.Equal(handlers.FmtDate(sitEstimatedEntryDate), updatedShipment.PpmShipment.SitEstimatedEntryDate)
		suite.Equal(handlers.FmtDate(sitEstimatedDepartureDate), updatedShipment.PpmShipment.SitEstimatedDepartureDate)
		suite.Equal(int64(sitEstimatedCost), *updatedShipment.PpmShipment.SitEstimatedCost)
		suite.Equal(handlers.FmtPoundPtr(&estimatedWeight), updatedShipment.PpmShipment.EstimatedWeight)
		suite.Equal(int64(estimatedIncentive), *updatedShipment.PpmShipment.EstimatedIncentive)
		suite.Equal(handlers.FmtBool(hasProGear), updatedShipment.PpmShipment.HasProGear)
		suite.Equal(handlers.FmtPoundPtr(&proGearWeight), updatedShipment.PpmShipment.ProGearWeight)
		suite.Equal(handlers.FmtPoundPtr(&spouseProGearWeight), updatedShipment.PpmShipment.SpouseProGearWeight)
	})

	suite.Run("Successful PATCH Delete Addresses - Integration Test (PPM)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		hasProGear := true
		ppmShipment := factory.BuildFullAddressPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					HasProGear: &hasProGear,
				},
			},
		}, nil)

		// we expect initial setup data to have no secondary addresses
		suite.NotNil(ppmShipment.SecondaryPickupAddress)
		suite.NotNil(ppmShipment.SecondaryDestinationAddress)

		estimatedIncentive := 654321
		sitEstimatedCost := 67500

		params := suite.getUpdateShipmentParams(ppmShipment.Shipment)
		params.Body.ShipmentType = ghcmessages.MTOShipmentTypePPM
		params.Body.PpmShipment = &ghcmessages.UpdatePPMShipment{
			HasSecondaryPickupAddress:      models.BoolPointer(false),
			HasSecondaryDestinationAddress: models.BoolPointer(false),
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))
		suite.Equal(ppmShipment.Shipment.ID.String(), updatedShipment.ID.String())

		suite.NotNil(updatedShipment.PpmShipment.PickupAddress)
		suite.NotNil(updatedShipment.PpmShipment.DestinationAddress)
		// expect secondary addresses to be deleted
		suite.Nil(updatedShipment.PpmShipment.SecondaryPickupAddress)
		suite.Nil(updatedShipment.PpmShipment.SecondaryDestinationAddress)

		suite.False(*updatedShipment.PpmShipment.HasSecondaryPickupAddress)
		suite.False(*updatedShipment.PpmShipment.HasSecondaryDestinationAddress)
	})

	suite.Run("Successful PATCH does not Delete Addresses w/o has flag - Integration Test (PPM)", func() {
		// Make a move along with an attached minimal shipment. Shouldn't matter what's in them.
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		hasProGear := true
		ppmShipment := factory.BuildFullAddressPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					HasProGear: &hasProGear,
				},
			},
		}, nil)

		// we expect initial setup data to have no secondary addresses
		suite.NotNil(ppmShipment.SecondaryPickupAddress)
		suite.NotNil(ppmShipment.SecondaryDestinationAddress)

		estimatedIncentive := 654321
		sitEstimatedCost := 67500

		params := suite.getUpdateShipmentParams(ppmShipment.Shipment)
		params.Body.ShipmentType = ghcmessages.MTOShipmentTypePPM
		params.Body.PpmShipment = &ghcmessages.UpdatePPMShipment{
			HasProGear: &hasProGear,
		}

		ppmEstimator.On("EstimateIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(models.CentPointer(unit.Cents(estimatedIncentive)), models.CentPointer(unit.Cents(sitEstimatedCost)), nil).Once()

		ppmEstimator.On("MaxIncentive",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		ppmEstimator.On("FinalIncentiveWithDefaultChecks",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("models.PPMShipment"),
			mock.AnythingOfType("*models.PPMShipment")).
			Return(nil, nil)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)
		updatedShipment := response.(*mtoshipmentops.UpdateMTOShipmentOK).Payload

		// Validate outgoing payload
		suite.NoError(updatedShipment.Validate(strfmt.Default))
		suite.Equal(ppmShipment.Shipment.ID.String(), updatedShipment.ID.String())

		suite.NotNil(updatedShipment.PpmShipment.PickupAddress)
		suite.NotNil(updatedShipment.PpmShipment.DestinationAddress)
		// expect secondary addresses not to be deleted
		suite.NotNil(updatedShipment.PpmShipment.SecondaryPickupAddress)
		suite.NotNil(updatedShipment.PpmShipment.SecondaryDestinationAddress)
		suite.NotNil(updatedShipment.PpmShipment.SecondaryPickupAddress.PostalCode)
		suite.NotNil(updatedShipment.PpmShipment.SecondaryDestinationAddress.PostalCode)

		suite.True(*updatedShipment.PpmShipment.HasSecondaryPickupAddress)
		suite.True(*updatedShipment.PpmShipment.HasSecondaryDestinationAddress)
	})

	suite.Run("PATCH failure - 400 -- nil body", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)
		params.Body = nil

		// Validate incoming payload: nil body (the point of this test)

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		uuidString := handlers.FmtUUID(uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763"))
		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)
		params.ShipmentID = *uuidString

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("PATCH failure - 412 -- etag mismatch", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)
		params.IfMatch = "intentionally-bad-if-match-header-value"

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 412 -- shipment shouldn't be updatable", func() {
		builder := query.NewQueryBuilder()
		fetcher := fetch.NewFetcher(builder)
		mockSender := suite.TestNotificationSender()
		mtoShipmentUpdater := mtoshipment.NewOfficeMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, mockSender, paymentRequestShipmentRecalculator, addressUpdater, addressCreator)
		ppmEstimator := mocks.PPMEstimator{}
		addressCreator := address.NewAddressCreator()
		addressUpdater := address.NewAddressUpdater()
		ppmShipmentUpdater := ppmshipment.NewPPMShipmentUpdater(&ppmEstimator, addressCreator, addressUpdater)

		boatShipmentUpdater := boatshipment.NewBoatShipmentUpdater()
		mobileHomeShipmentUpdater := mobilehomeshipment.NewMobileHomeShipmentUpdater()

		shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(mtoShipmentUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			shipmentUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusDraft,
				},
			},
		}, nil)

		params := suite.getUpdateShipmentParams(oldShipment)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentForbidden{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentForbidden).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 500", func() {
		mockUpdater := mocks.ShipmentUpdater{}
		handler := UpdateShipmentHandler{
			suite.HandlerConfig(),
			&mockUpdater,
			sitstatus.NewShipmentSITStatus(),
		}

		err := errors.New("ServerError")

		mockUpdater.On("UpdateShipment",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, err)

		oldShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		params := suite.getUpdateShipmentParams(oldShipment)

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)

		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)
		payload := response.(*mtoshipmentops.UpdateMTOShipmentInternalServerError).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})
}
func (suite *HandlerSuite) TestUpdateSITServiceItemCustomerExpenseHandler() {
	moveRouter := moveservices.NewMoveRouter()
	builder := query.NewQueryBuilder()
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
	addressCreator := address.NewAddressCreator()

	suite.Run("Successful PATCH - Integration Test", func() {
		// Build shipment with SIT
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dofsit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aMonthAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)

		approvedShipment.MTOServiceItems = models.MTOServiceItems{dofsit}

		eTag := etag.GenerateEtag(approvedShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		updater := mtoserviceitem.NewMTOServiceItemUpdater(planner, builder, moveRouter, shipmentFetcher, addressCreator)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-service-item/convert-to-customer-expense", approvedShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := UpdateSITServiceItemCustomerExpenseHandler{
			handlerConfig,
			updater,
			shipmentFetcher,
			sitstatus.NewShipmentSITStatus(),
		}
		convertToCustomerExpense := true
		customerExpenseReason := "test"
		createParams := shipmentops.UpdateSITServiceItemCustomerExpenseParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.UpdateSITServiceItemCustomerExpense{
				ConvertToCustomerExpense: &convertToCustomerExpense,
				CustomerExpenseReason:    &customerExpenseReason,
			},
			ShipmentID: *handlers.FmtUUID(approvedShipment.ID),
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&shipmentops.UpdateSITServiceItemCustomerExpenseOK{}, response)
		okResponse := response.(*shipmentops.UpdateSITServiceItemCustomerExpenseOK)
		payload := okResponse.Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("PATCH failure - 404 -- not found", func() {
		// Build shipment with SIT
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)

		year, month, day := time.Now().Add(time.Hour * 24 * -30).Date()
		aMonthAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		dofsit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    approvedShipment,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					SITEntryDate: &aMonthAgo,
					Status:       models.MTOServiceItemStatusApproved,
				},
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		}, nil)
		approvedShipment.MTOServiceItems = models.MTOServiceItems{dofsit}

		eTag := etag.GenerateEtag(approvedShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		updater := mtoserviceitem.NewMTOServiceItemUpdater(planner, builder, moveRouter, shipmentFetcher, addressCreator)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-service-item/convert-to-customer-expense", approvedShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := UpdateSITServiceItemCustomerExpenseHandler{
			handlerConfig,
			updater,
			shipmentFetcher,
			sitstatus.NewShipmentSITStatus(),
		}
		convertToCustomerExpense := true
		customerExpenseReason := "test"
		uuidString := handlers.FmtUUID(uuid.FromStringOrNil("d874d002-5582-4a91-97d3-786e8f66c763"))
		createParams := shipmentops.UpdateSITServiceItemCustomerExpenseParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.UpdateSITServiceItemCustomerExpense{
				ConvertToCustomerExpense: &convertToCustomerExpense,
				CustomerExpenseReason:    &customerExpenseReason,
			},
			ShipmentID: *uuidString,
		}

		// Validate incoming payload
		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)

		suite.IsType(&shipmentops.UpdateSITServiceItemCustomerExpenseNotFound{}, response)
		payload := response.(*shipmentops.UpdateSITServiceItemCustomerExpenseNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("500 server error", func() {
		shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
		// Build shipment with SIT
		shipmentSITAllowance := int(90)
		approvedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					SITDaysAllowance: &shipmentSITAllowance,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(approvedShipment.UpdatedAt)
		officeUser := factory.BuildOfficeUserWithRoles(nil, nil, []roles.RoleType{roles.RoleTypeTOO})
		// Mock the updater to return an error.
		mockedUpdater := &mocks.MTOServiceItemUpdater{}
		mockedUpdater.On("ConvertItemToCustomerExpense",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("internal server error"))

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/shipments/%s/sit-service-item/convert-to-customer-expense", approvedShipment.ID.String()), nil)
		req = suite.AuthenticateOfficeRequest(req, officeUser)
		handlerConfig := suite.HandlerConfig()

		handler := UpdateSITServiceItemCustomerExpenseHandler{
			handlerConfig,
			mockedUpdater, // Use the mocked updater here
			shipmentFetcher,
			sitstatus.NewShipmentSITStatus(),
		}
		convertToCustomerExpense := true
		customerExpenseReason := "test"
		createParams := shipmentops.UpdateSITServiceItemCustomerExpenseParams{
			HTTPRequest: req,
			IfMatch:     eTag,
			Body: &ghcmessages.UpdateSITServiceItemCustomerExpense{
				ConvertToCustomerExpense: &convertToCustomerExpense,
				CustomerExpenseReason:    &customerExpenseReason,
			},
			ShipmentID: *handlers.FmtUUID(approvedShipment.ID),
		}

		suite.NoError(createParams.Body.Validate(strfmt.Default))

		response := handler.Handle(createParams)
		suite.IsType(&shipmentops.UpdateSITServiceItemCustomerExpenseInternalServerError{}, response)
	})
}
