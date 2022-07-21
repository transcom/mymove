package supportapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/supportapi/supportoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/supportmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMTOShipmentStatusHandler() {
	setupTestData := func() models.MTOShipment {
		ghcDomesticTransitTime := models.GHCDomesticTransitTime{
			MaxDaysTransitTime: 12,
			WeightLbsLower:     0,
			WeightLbsUpper:     10000,
			DistanceMilesLower: 0,
			DistanceMilesUpper: 10000,
		}
		_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

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

		return mtoShipment
	}

	setupParams := func(shipment models.MTOShipment) mtoshipmentops.UpdateMTOShipmentStatusParams {
		requestUser := testdatagen.MakeStubbedUser(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", shipment.ID.String()), nil)
		req = suite.AuthenticateUserRequest(req, requestUser)

		return mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(shipment.ID),
			Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "APPROVED"},
			IfMatch:       eTag,
		}
	}

	// Used for all tests except 500 error:
	queryBuilder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(queryBuilder)
	moveRouter := moverouter.NewMoveRouter()
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
	planner := &routemocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(1000, nil)
	updater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, siCreator, planner)

	setupHandler := func() UpdateMTOShipmentStatusHandlerFunc {
		return UpdateMTOShipmentStatusHandlerFunc{
			suite.HandlerConfig(),
			fetcher,
			updater,
		}
	}

	suite.Run("Patch failure - 500", func() {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		mockHandler := UpdateMTOShipmentStatusHandlerFunc{
			suite.HandlerConfig(),
			&mockFetcher,
			&mockUpdater,
		}

		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipmentStatus",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := mockHandler.Handle(setupParams(setupTestData()))
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})

	suite.Run("Patch failure - 404", func() {
		params := setupParams(setupTestData())
		notFoundParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: strfmt.UUID(uuid.Nil.String()),
			Body:          params.Body,
			IfMatch:       params.IfMatch,
		}
		response := setupHandler().Handle(notFoundParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusNotFound{}, response)
	})

	suite.Run("Patch failure - 412", func() {
		params := setupParams(setupTestData())
		preconditionParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          params.Body,
			IfMatch:       "eTag",
		}
		response := setupHandler().Handle(preconditionParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusPreconditionFailed{}, response)
	})

	suite.Run("Patch failure - 422", func() {
		params := setupParams(setupTestData())
		invalidInputParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "X"},
			IfMatch:       params.IfMatch,
		}
		// This input error would be caught by Swagger, not the service object:
		suite.Error(invalidInputParams.Body.Validate(strfmt.Default))
	})

	// Second to last because many of the above tests fail because of a conflict error with APPROVED/REJECTED shipments
	// first:
	suite.Run("Successful patch - Integration Test", func() {
		mtoShipment := setupTestData()
		move := mtoShipment.MoveTaskOrder
		move.Status = models.MoveStatusAPPROVED
		_ = suite.DB().Save(&move)

		params := setupParams(mtoShipment)
		response := setupHandler().Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)
	})

	// Last because the shipment has to be either APPROVED or REJECTED before triggering this conflict:
	suite.Run("Patch failure - 409", func() {
		params := setupParams(setupTestData())
		conflictParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "SUBMITTED"},
			IfMatch:       params.IfMatch,
		}
		response := setupHandler().Handle(conflictParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusConflict{}, response)
	})

	// Test Successful Cancellation Request
	suite.Run("Successful patch - Integration Test for CANCELLATION_REQUESTED", func() {
		// Under test: updateMTOShipmentHandler function
		// Mocked:     none
		// Setup: We create a new mtoShipment, then try to update the status from Approved to Cancellation_Requested
		// Expected outcome:
		//             Successfully updated status to CANCELLATION_REQUESTED
		mto := testdatagen.MakeMove(suite.DB(), testdatagen.Assertions{Move: models.Move{Status: models.MoveStatusAPPROVED}})
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
			MTOShipment: models.MTOShipment{
				Status:       models.MTOShipmentStatusApproved,
				ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
			},
		})
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

		baseParams := setupParams(mtoShipment)
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   baseParams.HTTPRequest,
			MtoShipmentID: *handlers.FmtUUID(mtoShipment.ID),
			Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "CANCELLATION_REQUESTED"},
			IfMatch:       eTag,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := setupHandler().Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusOK)
		suite.Equal(supportmessages.UpdateMTOShipmentStatusStatusCANCELLATIONREQUESTED, okResponse.Payload.Status)
		suite.NotZero(okResponse.Payload.ETag)
	})
}
