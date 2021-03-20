package ghcapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/gofrs/uuid"

	routemocks "github.com/transcom/mymove/pkg/route/mocks"

	"github.com/transcom/mymove/pkg/services"

	"github.com/gobuffalo/validate/v3"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestListMTOShipmentsHandler() {
	mto := testdatagen.MakeDefaultMove(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
	})
	mtoAgent := testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipmentID: mtoShipment.ID,
		},
	})
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			MTOShipmentID: &mtoShipment.ID,
		},
	})

	shipments := models.MTOShipments{mtoShipment}
	requestUser := testdatagen.MakeStubbedUser(suite.DB())

	req := httptest.NewRequest("GET", fmt.Sprintf("/move_task_orders/%s/mto_shipments", mto.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoshipmentops.ListMTOShipmentsParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
	}

	suite.T().Run("Successful list fetch - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		listFetcher := fetch.NewListFetcher(queryBuilder)
		fetcher := fetch.NewFetcher(queryBuilder)
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			listFetcher,
			fetcher,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsOK{}, response)

		okResponse := response.(*mtoshipmentops.ListMTOShipmentsOK)
		suite.Len(okResponse.Payload, 1)
		suite.Equal(shipments[0].ID.String(), okResponse.Payload[0].ID.String())
		suite.Equal(mtoAgent.ID.String(), okResponse.Payload[0].MtoAgents[0].ID.String())
		suite.Equal(mtoServiceItem.ID.String(), okResponse.Payload[0].MtoServiceItems[0].ID.String())
	})

	suite.T().Run("Failure list fetch - Internal Server Error", func(t *testing.T) {
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		internalServerErr := errors.New("ServerError")

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(nil)

		mockListFetcher.On("FetchRecordList",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(internalServerErr)

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsInternalServerError{}, response)
	})

	suite.T().Run("Failure list fetch - 404 Not Found - Move Task Order ID", func(t *testing.T) {
		mockListFetcher := mocks.ListFetcher{}
		mockFetcher := mocks.Fetcher{}
		handler := ListMTOShipmentsHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockListFetcher,
			&mockFetcher,
		}

		notfound := errors.New("Not found error")

		mockFetcher.On("FetchRecord",
			mock.Anything,
			mock.Anything,
		).Return(notfound)

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.ListMTOShipmentsNotFound{}, response)
	})
}

func (suite *HandlerSuite) TestPatchMTOShipmentHandler() {
	mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusAPPROVED}})
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: mto,
		MTOShipment: models.MTOShipment{
			Status:       models.MTOShipmentStatusSubmitted,
			ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
		},
	})
	// Populate the reServices table with codes needed by the
	// HHG_LONGHAUL_DOMESTIC shipment type
	reServiceCodes := []models.ReServiceCode{
		models.ReServiceCodeDLH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDPK,
		models.ReServiceCodeDUPK,
	}
	for _, serviceCode := range reServiceCodes {
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code:      serviceCode,
				Name:      "test",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}

	requestUser := testdatagen.MakeStubbedUser(suite.DB())
	eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", mto.ID.String(), mtoShipment.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoshipmentops.PatchMTOShipmentStatusParams{
		HTTPRequest:     req,
		MoveTaskOrderID: *handlers.FmtUUID(mtoShipment.MoveTaskOrderID),
		ShipmentID:      *handlers.FmtUUID(mtoShipment.ID),
		Body:            &ghcmessages.PatchMTOShipmentStatus{Status: ghcmessages.MTOShipmentStatusAPPROVED},
		IfMatch:         eTag,
	}

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	planner := &routemocks.Planner{}
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(1000, nil)

	suite.T().Run("Successful patch - Integration Test", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

		updater := mtoshipment.NewMTOShipmentStatusUpdater(suite.DB(), queryBuilder, siCreator, planner)
		handler := PatchShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			fetcher,
			updater,
		}
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.PatchMTOShipmentStatusOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)
	})

	suite.T().Run("Patch failure - 500", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		handler := PatchShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
		}

		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusInternalServerError{}, response)
	})

	suite.T().Run("Patch failure - 404", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		handler := PatchShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.NotFoundError{})

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusNotFound{}, response)
	})

	suite.T().Run("Patch failure - 422", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		handler := PatchShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.InvalidInputError{ValidationErrors: &validate.Errors{}})

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusUnprocessableEntity{}, response)
	})

	suite.T().Run("Patch failure - 412", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		handler := PatchShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, services.PreconditionFailedError{})

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusPreconditionFailed{}, response)
	})

	suite.T().Run("Patch failure - 409", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		handler := PatchShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockFetcher,
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, mtoshipment.ConflictStatusError{})

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusConflict{}, response)
	})

	suite.T().Run("Successful patch with webhook notification - On a submitted move", func(t *testing.T) {
		// Create mock fetcher and updater
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}

		// Create an mto shipment on a move that is available to prime
		now := time.Now()
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
			MTOShipment: models.MTOShipment{
				Status:       models.MTOShipmentStatusSubmitted,
				ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
			},
		})

		// Set the traceID so we can use it to find the webhook notification
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handlerContext.SetTraceID(uuid.Must(uuid.NewV4()))

		handler := PatchShipmentHandler{
			handlerContext,
			&mockFetcher,
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&mtoShipment, nil)

		// Call the handler
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.PatchMTOShipmentStatusOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)

		// Check that webhook notification was stored
		suite.HasWebhookNotification(mtoShipment.ID, handlerContext.GetTraceID())

	})

	suite.T().Run("Successful patch with no webhook notification - On an unsubmitted move", func(t *testing.T) {
		// Create mock fetcher and updater
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}

		// Create an mto shipment on a move that is available to prime
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:       models.MTOShipmentStatusSubmitted,
				ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
			},
		})

		// Set the traceID so we can use it to find the webhook notification
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handlerContext.SetTraceID(uuid.Must(uuid.NewV4()))

		handler := PatchShipmentHandler{
			handlerContext,
			&mockFetcher,
			&mockUpdater,
		}

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(&mtoShipment, nil)

		// Call the handler
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.PatchMTOShipmentStatusOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)

		// Check that webhook notification was stored
		suite.HasNoWebhookNotification(mtoShipment.ID, handlerContext.GetTraceID())

	})

	suite.T().Run("Successful patch to CANCELLATION_REQUESTED status", func(t *testing.T) {
		queryBuilder := query.NewQueryBuilder(suite.DB())
		fetcher := fetch.NewFetcher(queryBuilder)
		siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)
		updater := mtoshipment.NewMTOShipmentStatusUpdater(suite.DB(), queryBuilder, siCreator, planner)

		// Create an APPROVED shipment on a move that is available to prime
		approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: testdatagen.MakeAvailableMove(suite.DB()),
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})
		eTag := etag.GenerateEtag(approvedShipment.UpdatedAt)

		// Set the traceID so we can use it to find the webhook notification
		handlerContext := handlers.NewHandlerContext(suite.DB(), suite.TestLogger())
		handlerContext.SetTraceID(uuid.Must(uuid.NewV4()))

		handler := PatchShipmentHandler{
			handlerContext,
			fetcher,
			updater,
		}
		cancellationParams := mtoshipmentops.PatchMTOShipmentStatusParams{
			HTTPRequest:     req,
			MoveTaskOrderID: *handlers.FmtUUID(approvedShipment.MoveTaskOrderID),
			ShipmentID:      *handlers.FmtUUID(approvedShipment.ID),
			Body:            &ghcmessages.PatchMTOShipmentStatus{Status: ghcmessages.MTOShipmentStatusCANCELLATIONREQUESTED},
			IfMatch:         eTag,
		}
		suite.NoError(cancellationParams.Body.Validate(strfmt.Default))

		// Call the handler
		response := handler.Handle(cancellationParams)
		suite.IsType(&mtoshipmentops.PatchMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.PatchMTOShipmentStatusOK)
		suite.Equal(approvedShipment.ID.String(), okResponse.Payload.ID.String())
		suite.Equal(string(okResponse.Payload.Status), string(models.MTOShipmentStatusCancellationRequested))
		suite.NotNil(okResponse.Payload.ETag)

		// Check that webhook notification was stored
		suite.HasWebhookNotification(approvedShipment.ID, handlerContext.GetTraceID())

	})
}
