package ghcapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	boatshipment "github.com/transcom/mymove/pkg/services/boat_shipment"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mobilehomeshipment "github.com/transcom/mymove/pkg/services/mobile_home_shipment"
	"github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	shipmentorchestrator "github.com/transcom/mymove/pkg/services/orchestrators/shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	ppmshipment "github.com/transcom/mymove/pkg/services/ppmshipment"
	"github.com/transcom/mymove/pkg/services/query"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/trace"
)

func (suite *HandlerSuite) TestListMTOServiceItemHandler() {
	reServiceID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	mtoShipmentID, _ := uuid.NewV4()
	var mtoID uuid.UUID

	setupTestData := func() (models.User, models.MTOServiceItems) {
		mto := factory.BuildMove(suite.DB(), nil, nil)
		mtoID = mto.ID
		reService := factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					ID:   reServiceID,
					Code: "TEST10000",
				},
			},
		}, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{ID: mtoShipmentID},
			},
		}, nil)
		requestUser := factory.BuildUser(nil, nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					ID: serviceItemID,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		year, month, day := time.Now().Date()
		aWeekAgo := time.Date(year, month, day-7, 0, 0, 0, 0, time.UTC)
		departureDate := aWeekAgo.Add(time.Hour * 24 * 30)
		originSit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
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
					Code: models.ReServiceCodeDOPSIT,
				},
			},
		}, nil)

		customerContact := testdatagen.MakeMTOServiceItemCustomerContact(suite.DB(), testdatagen.Assertions{})
		destinationSit := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact},
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
					Code: models.ReServiceCodeDDFSIT,
					Name: "Destination 1st Day SIT",
				},
			},
		}, nil)

		serviceItems := models.MTOServiceItems{serviceItem, originSit, destinationSit}

		return requestUser, serviceItems
	}

	suite.Run("Successful list fetch - Integration Test", func() {
		requestUser, serviceItems := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItems[0].MoveTaskOrderID),
		}

		serviceItem1 := serviceItems[0]

		serviceRequestDocumentUpload := factory.BuildServiceRequestDocumentUpload(suite.DB(), []factory.Customization{
			{
				Model:    serviceItem1,
				LinkOnly: true,
			},
		}, nil)

		serviceItem1.ServiceRequestDocuments = models.ServiceRequestDocuments{serviceRequestDocumentUpload.ServiceRequestDocument}

		queryBuilder := query.NewQueryBuilder()
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		counselingPricer := ghcrateengine.NewCounselingServicesPricer()
		moveManagementPricer := ghcrateengine.NewManagementServicesPricer()
		handler := ListMTOServiceItemsHandler{
			suite.createS3HandlerConfig(),
			listFetcher,
			fetcher,
			counselingPricer,
			moveManagementPricer,
		}

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsOK{}, response)
		okResponse := response.(*mtoserviceitemop.ListMTOServiceItemsOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))
		fmt.Println(okResponse.Payload)

		suite.Len(okResponse.Payload, 3)
		for _, serviceItem := range serviceItems {
			for _, payload := range okResponse.Payload {
				// Validate that the Customer Contacts were included in the payload
				if len(serviceItem.CustomerContacts) > 0 {
					if len(payload.CustomerContacts) > 0 {
						suite.Equal(serviceItem.ID.String(), payload.ID.String())
						suite.Len(payload.CustomerContacts, 1)
					}
				}

				//Validate that Service Request Document upload was included in payload
				if len(serviceItem.ServiceRequestDocuments) == 1 && suite.Len(payload.ServiceRequestDocuments, 1) {
					if len(serviceItem.ServiceRequestDocuments[0].ServiceRequestDocumentUploads) == 1 && suite.Len(payload.ServiceRequestDocuments[0].Uploads, 1) {
						upload := serviceItem.ServiceRequestDocuments[0].ServiceRequestDocumentUploads[0].Upload
						uploadPayload := payload.ServiceRequestDocuments[0].Uploads[0]
						suite.Equal(upload.ID.String(), uploadPayload.ID.String())
						suite.NotEqual(string(uploadPayload.URL), "")
					}
				}
			}
		}
	})

	suite.Run("Failure list fetch - Internal Server Error", func() {
		requestUser, serviceItems := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItems[0].MoveTaskOrderID),
		}
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		mockCounselingPricer := mocks.CounselingServicesPricer{}
		mockMoveManagementPricer := mocks.ManagementServicesPricer{}
		handler := ListMTOServiceItemsHandler{
			suite.HandlerConfig(),
			&mockListFetcher,
			&mockFetcher,
			&mockCounselingPricer,
			&mockMoveManagementPricer,
		}

		internalServerErr := errors.New("ServerError")

		mockFetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockListFetcher.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(internalServerErr)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsInternalServerError{}, response)
		payload := response.(*mtoserviceitemop.ListMTOServiceItemsInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("Failure list fetch - 404 Not Found - Move Task Order ID", func() {
		requestUser, serviceItems := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItems[0].MoveTaskOrderID),
		}

		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		mockCounselingPricer := mocks.CounselingServicesPricer{}
		mockMoveManagementPricer := mocks.ManagementServicesPricer{}
		handler := ListMTOServiceItemsHandler{
			suite.HandlerConfig(),
			&mockListFetcher,
			&mockFetcher,
			&mockCounselingPricer,
			&mockMoveManagementPricer,
		}

		notfound := errors.New("Not found error")

		mockFetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(notfound)

		// Validate incoming payload: no body to validate

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsNotFound{}, response)
		payload := response.(*mtoserviceitemop.ListMTOServiceItemsNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})
}

func (suite *HandlerSuite) createServiceItem() (models.MTOServiceItem, models.Move) {
	move := factory.BuildApprovalsRequestedMove(suite.DB(), nil, nil)
	serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	return serviceItem, move
}

func (suite *HandlerSuite) TestUpdateMTOServiceItemStatusHandler() {

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
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()

	moveTaskOrderID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	var requestUser models.User

	setupTestData := func() mtoserviceitemop.UpdateMTOServiceItemStatusParams {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status",
			moveTaskOrderID, serviceItemID), nil)

		req = suite.AuthenticateUserRequest(req, requestUser)
		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(time.Now()),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "APPROVED"},
			MoveTaskOrderID:  moveTaskOrderID.String(),
			MtoServiceItemID: serviceItemID.String(),
		}
		return params
	}

	// With this first set of tests we'll use mocked service object responses so that we can make sure the handler
	// is returning the right HTTP code given a set of circumstances.
	suite.Run("404 - not found response", func() {
		params := setupTestData()
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(errors.New("Not found error")).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusNotFound{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusNotFound).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("200 - success response", func() {
		params := setupTestData()
		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("ApproveOrRejectServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&models.MTOServiceItem{ID: serviceItemID}, nil).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("412 - precondition failed response", func() {
		params := setupTestData()

		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("ApproveOrRejectServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, apperror.NewPreconditionFailedError(serviceItemID, errors.New("oh no"))).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusPreconditionFailed{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusPreconditionFailed).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	suite.Run("500 - internal server error response", func() {
		params := setupTestData()

		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		serviceItemStatusUpdater.On("ApproveOrRejectServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, errors.New("oh no")).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusInternalServerError{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.Nil(payload)
	})

	suite.Run("422 - unprocessable entity response", func() {
		params := setupTestData()

		serviceItemStatusUpdater := mocks.MTOServiceItemUpdater{}
		fetcher := mocks.Fetcher{}
		params.MtoServiceItemID = ""
		fetcher.On("FetchRecord",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(nil).Once()

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: &serviceItemStatusUpdater,
			Fetcher:               &fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusUnprocessableEntity{}, response)
		payload := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusUnprocessableEntity).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// With this we'll do a happy path integration test to ensure that the use of the service object
	// by the handler is working as expected.
	suite.Run("Successful rejected status update - Integration test", func() {

		queryBuilder := query.NewQueryBuilder()
		mtoServiceItem, move := suite.createServiceItem()
		requestUser := factory.BuildUser(nil, nil, nil)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status",
			moveTaskOrderID, serviceItemID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		rejectionReason := "No justification given"
		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "REJECTED", RejectionReason: &rejectionReason},
			MoveTaskOrderID:  move.ID.String(),
			MtoServiceItemID: mtoServiceItem.ID.String(),
		}

		fetcher := fetch.NewFetcher(queryBuilder)
		moveRouter := moveservices.NewMoveRouter()
		shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
		addressCreator := address.NewAddressCreator()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		mtoServiceItemStatusUpdater := mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator)

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: mtoServiceItemStatusUpdater,
			Fetcher:               fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		okResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(string(models.MTOServiceItemStatusRejected), string(okResponse.Payload.Status))
		suite.NotNil(okResponse.Payload.RejectedAt)
		suite.Equal(rejectionReason, *okResponse.Payload.RejectionReason)
	})

	// With this we'll do a happy path integration test to ensure that the use of the service object
	// by the handler is working as expected.
	suite.Run("Successful status update of MTO service item and event trigger", func() {
		queryBuilder := query.NewQueryBuilder()
		moveRouter := moveservices.NewMoveRouter()
		shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()
		mtoServiceItem, availableMove := suite.createServiceItem()
		requestUser := factory.BuildUser(nil, nil, nil)
		availableMoveID := availableMove.ID
		mtoServiceItemID := mtoServiceItem.ID

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_service_items/%s/status", availableMoveID, mtoServiceItemID), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		traceID, err := uuid.NewV4()
		suite.FatalNoError(err, "Error creating a new trace ID.")
		req = req.WithContext(trace.NewContext(req.Context(), traceID))

		params := mtoserviceitemop.UpdateMTOServiceItemStatusParams{
			HTTPRequest:      req,
			IfMatch:          etag.GenerateEtag(mtoServiceItem.UpdatedAt),
			Body:             &ghcmessages.PatchMTOServiceItemStatusPayload{Status: "APPROVED"},
			MoveTaskOrderID:  availableMoveID.String(),
			MtoServiceItemID: mtoServiceItemID.String(),
		}

		fetcher := fetch.NewFetcher(queryBuilder)
		addressCreator := address.NewAddressCreator()
		planner := &routemocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(400, nil)
		mtoServiceItemStatusUpdater := mtoserviceitem.NewMTOServiceItemUpdater(planner, queryBuilder, moveRouter, shipmentFetcher, addressCreator)

		handler := UpdateMTOServiceItemStatusHandler{
			HandlerConfig:         suite.HandlerConfig(),
			MTOServiceItemUpdater: mtoServiceItemStatusUpdater,
			Fetcher:               fetcher,
			ShipmentSITStatus:     sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:    shipmentFetcher,
			ShipmentUpdater:       shipmentUpdater,
		}

		// Validate incoming payload
		suite.NoError(params.Body.Validate(strfmt.Default))

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateMTOServiceItemStatusOK{}, response)
		okResponse := response.(*mtoserviceitemop.UpdateMTOServiceItemStatusOK)

		// Validate outgoing payload
		suite.NoError(okResponse.Payload.Validate(strfmt.Default))

		suite.Equal(string(models.MTOServiceItemStatusApproved), string(okResponse.Payload.Status))
		suite.NotNil(okResponse.Payload.ApprovedAt)
		suite.HasWebhookNotification(mtoServiceItemID, traceID)

		impactedMove := models.Move{}
		_ = suite.DB().Find(&impactedMove, okResponse.Payload.MoveTaskOrderID)
		suite.Equal(models.MoveStatusAPPROVED, impactedMove.Status)
	})
}

func (suite *HandlerSuite) TestGetMTOServiceItemHandler() {
	serviceItemID := uuid.Must(uuid.FromString("f7b4b9e2-04e8-4c34-827a-df917e69caf4"))
	moveTaskOrderID := uuid.Must(uuid.FromString("f7b4b9e2-04e8-4c34-1234-df917e69caf4"))
	var requestUser models.User

	setupTestData := func() mtoserviceitemop.GetMTOServiceItemParams {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/service_items/%s",
			moveTaskOrderID, serviceItemID), nil)

		req = suite.AuthenticateUserRequest(req, requestUser)
		params := mtoserviceitemop.GetMTOServiceItemParams{
			HTTPRequest:      req,
			MoveTaskOrderID:  moveTaskOrderID.String(),
			MtoServiceItemID: serviceItemID.String(),
		}
		return params
	}

	suite.Run("200 - success response", func() {
		// setting up test data
		params := setupTestData()
		// mock function
		serviceItemFetcher := mocks.MTOServiceItemFetcher{}

		// creating struct that returns from mock function & updating values
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)
		mtoServiceItem.ID = serviceItemID
		// calling mock function, passing in what we want, and saying "return this"
		serviceItemFetcher.On("GetServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			serviceItemID,
		).Return(&mtoServiceItem, nil).Once()

		handler := GetMTOServiceItemHandler{
			HandlerConfig:         suite.HandlerConfig(),
			mtoServiceItemFetcher: &serviceItemFetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.GetMTOServiceItemOK{}, response)
		payload := response.(*mtoserviceitemop.GetMTOServiceItemOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// With this first set of tests we'll use mocked service object responses so that we can make sure the handler
	// is returning the right HTTP code given a set of circumstances.
	suite.Run("404 - not found response", func() {
		params := setupTestData()
		serviceItemFetcher := mocks.MTOServiceItemFetcher{}
		serviceItemFetcher.On("GetServiceItem",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, errors.New("Not found error")).Once()

		handler := GetMTOServiceItemHandler{
			HandlerConfig:         suite.HandlerConfig(),
			mtoServiceItemFetcher: &serviceItemFetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.GetMTOServiceItemInternalServerError{}, response)
		payload := response.(*mtoserviceitemop.GetMTOServiceItemInternalServerError).Payload

		// Validate outgoing payload: nil payload
		suite.IsType(payload, &ghcmessages.Error{})
	})
}

func (suite *HandlerSuite) TestUpdateServiceItemSitEntryDateHandler() {
	serviceItemID := uuid.Must(uuid.FromString("f7b4b9e2-04e8-4c34-827a-df917e69caf4"))
	var requestUser models.User
	newSitEntryDate := time.Date(2023, time.October, 10, 10, 10, 0, 0, time.UTC)

	sitEntryDateParamsBody := ghcmessages.ServiceItemSitEntryDate{
		ID:           *handlers.FmtUUID(serviceItemID),
		SitEntryDate: handlers.FmtDateTime(newSitEntryDate),
	}

	setupTestData := func() mtoserviceitemop.UpdateServiceItemSitEntryDateParams {
		requestUser = factory.BuildUser(nil, nil, nil)
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/service_item/%s/entry_date_update",
			serviceItemID), nil)

		req = suite.AuthenticateUserRequest(req, requestUser)
		params := mtoserviceitemop.UpdateServiceItemSitEntryDateParams{
			HTTPRequest:      req,
			Body:             &sitEntryDateParamsBody,
			MtoServiceItemID: serviceItemID.String(),
		}
		return params
	}

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
	shipmentUpdater := shipmentorchestrator.NewShipmentUpdater(noCheckUpdater, ppmShipmentUpdater, boatShipmentUpdater, mobileHomeShipmentUpdater)
	shipmentFetcher := mtoshipment.NewMTOShipmentFetcher()

	suite.Run("200 - success response", func() {
		// setting up test data
		params := setupTestData()
		// mock function
		sitEntryDateUpdater := mocks.SitEntryDateUpdater{}

		// setting up data to be used for updating sit entry date
		newSitEntryDate := time.Date(2023, time.October, 10, 10, 10, 0, 0, time.UTC)
		expectedUUID := serviceItemID
		// creating struct that passes into mock function
		sitEntryDateUpdateModel := models.SITEntryDateUpdate{
			ID: expectedUUID, SITEntryDate: &newSitEntryDate,
		}
		// creating struct that returns from mock function & updating values
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)
		mtoServiceItem.ID = sitEntryDateUpdateModel.ID
		mtoServiceItem.SITEntryDate = sitEntryDateUpdateModel.SITEntryDate
		// calling mock function, passing in what we want, and saying "return this"
		sitEntryDateUpdater.On("UpdateSitEntryDate",
			mock.AnythingOfType("*appcontext.appContext"),
			&sitEntryDateUpdateModel,
		).Return(&mtoServiceItem, nil).Once()

		handler := UpdateServiceItemSitEntryDateHandler{
			HandlerConfig:       suite.HandlerConfig(),
			sitEntryDateUpdater: &sitEntryDateUpdater,
			ShipmentSITStatus:   sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:  shipmentFetcher,
			ShipmentUpdater:     shipmentUpdater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateServiceItemSitEntryDateOK{}, response)
		payload := response.(*mtoserviceitemop.UpdateServiceItemSitEntryDateOK).Payload

		// Validate outgoing payload
		suite.NoError(payload.Validate(strfmt.Default))
	})

	// With this first set of tests we'll use mocked service object responses so that we can make sure the handler
	// is returning the right HTTP code given a set of circumstances.
	suite.Run("404 - not found response", func() {
		params := setupTestData()
		sitEntryDateUpdater := mocks.SitEntryDateUpdater{}
		sitEntryDateUpdater.On("UpdateSitEntryDate",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
		).Return(nil, errors.New("Not found error")).Once()

		handler := UpdateServiceItemSitEntryDateHandler{
			HandlerConfig:       suite.HandlerConfig(),
			sitEntryDateUpdater: &sitEntryDateUpdater,
			ShipmentSITStatus:   sitstatus.NewShipmentSITStatus(),
			MTOShipmentFetcher:  shipmentFetcher,
			ShipmentUpdater:     shipmentUpdater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.UpdateServiceItemSitEntryDateUnprocessableEntity{}, response)
		payload := response.(*mtoserviceitemop.UpdateServiceItemSitEntryDateUnprocessableEntity).Payload

		// Validate outgoing payload: nil payload
		suite.IsType(payload, &ghcmessages.ValidationError{})
	})
}

func (suite *HandlerSuite) TestListMTOServiceItemsHandlerWithICRTandIUCRT() {
	reServiceID, _ := uuid.NewV4()
	serviceItemID, _ := uuid.NewV4()
	serviceItemID2, _ := uuid.NewV4()
	mtoShipmentID, _ := uuid.NewV4()
	var mtoID uuid.UUID

	setupTestData := func() (models.User, models.MTOServiceItems) {
		mto := factory.BuildMove(suite.DB(), nil, nil)
		mtoID = mto.ID
		reServiceICRT := factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					ID:   reServiceID,
					Code: models.ReServiceCodeICRT,
				},
			},
		}, nil)
		reServiceIUCRT := factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					ID:   reServiceID,
					Code: models.ReServiceCodeIUCRT,
				},
			},
		}, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{ID: mtoShipmentID},
			},
		}, nil)
		requestUser := factory.BuildUser(nil, nil, nil)
		serviceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					ID: serviceItemID,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    reServiceICRT,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)
		serviceItem2 := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					ID: serviceItemID2,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    reServiceIUCRT,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
		}, nil)

		serviceItems := models.MTOServiceItems{serviceItem, serviceItem2}

		return requestUser, serviceItems
	}

	suite.Run("200 - successfully loads PickupAddress and DestinationAddress for intl crating", func() {
		requestUser, serviceItem := setupTestData()
		req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_service_items", mtoID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		params := mtoserviceitemop.ListMTOServiceItemsParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(serviceItem[0].MoveTaskOrderID),
		}

		// Create the addresses
		pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)
		destinationAddress := factory.BuildAddress(suite.DB(), nil, nil)

		// Create the MTOShipment with populated PickupAddress and DestinationAddress
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					PickupAddressID:      &pickupAddress.ID,
					DestinationAddressID: &destinationAddress.ID,
				},
			},
		}, nil)

		// Create service items with references to the MTOShipment
		serviceItems := models.MTOServiceItems{
			{
				ID:            uuid.Must(uuid.NewV4()),
				ReService:     models.ReService{Code: models.ReServiceCodeICRT},
				MTOShipmentID: &mtoShipment.ID,
				MTOShipment:   mtoShipment,
			},
			{
				ID:            uuid.Must(uuid.NewV4()),
				ReService:     models.ReService{Code: models.ReServiceCodeIUCRT},
				MTOShipmentID: &mtoShipment.ID,
				MTOShipment:   mtoShipment,
			},
		}

		// Mock Load function for PickupAddress and DestinationAddress
		mockLoad := func(item interface{}, associations ...string) error {
			mtoServiceItem, ok := item.(*models.MTOServiceItem)
			if !ok {
				return fmt.Errorf("unexpected type for item: %T", item)
			}
			if len(associations) == 2 && associations[0] == "MTOShipment.PickupAddress" && associations[1] == "MTOShipment.DestinationAddress" {
				mtoServiceItem.MTOShipment.PickupAddress = &pickupAddress
				mtoServiceItem.MTOShipment.DestinationAddress = &destinationAddress
				return nil
			}
			return fmt.Errorf("unexpected association: %v", associations)
		}

		// Inject mockLoad behavior
		for i := range serviceItems {
			if serviceItems[i].ReService.Code == models.ReServiceCodeICRT || serviceItems[i].ReService.Code == models.ReServiceCodeIUCRT {
				err := mockLoad(&serviceItems[i], "MTOShipment.PickupAddress", "MTOShipment.DestinationAddress")
				suite.NoError(err, "Expected no error when loading Pickup and Destination Addresses for ICRT/IUCRT codes")
			}
		}

		// Mock ListFetcher
		listFetcherMock := mocks.ListFetcher{}
		listFetcherMock.On("FetchRecordList",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Run(func(args mock.Arguments) {
			arg := args.Get(1).(*models.MTOServiceItems)
			*arg = serviceItems
		}).Return(nil)

		queryBuilder := query.NewQueryBuilder()
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		counselingPricer := ghcrateengine.NewCounselingServicesPricer()
		moveManagementPricer := ghcrateengine.NewManagementServicesPricer()

		// Configure the handler with mocks
		handler := ListMTOServiceItemsHandler{
			suite.createS3HandlerConfig(),
			listFetcher,
			fetcher,
			counselingPricer,
			moveManagementPricer,
		}

		// Run the handler
		response := handler.Handle(params)
		suite.IsType(&mtoserviceitemop.ListMTOServiceItemsOK{}, response)
		okResponse := response.(*mtoserviceitemop.ListMTOServiceItemsOK)

		// Validate the response
		suite.Len(okResponse.Payload, len(serviceItems))
		for _, payload := range okResponse.Payload {
			if *payload.ReServiceCode == string(models.ReServiceCodeICRT) {
				suite.NotNil(payload.Market, "Expected Market to be set for ICRT")
			} else if *payload.ReServiceCode == string(models.ReServiceCodeIUCRT) {
				suite.NotNil(payload.Market, "Expected Market to be set for IUCRT")
			}
		}
	})
}
