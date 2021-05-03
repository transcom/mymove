package supportapi

import (
	"fmt"
	"net/http/httptest"
	"testing"
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
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
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

	req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s", mtoShipment.ID.String()), nil)
	req = suite.AuthenticateUserRequest(req, requestUser)

	params := mtoshipmentops.UpdateMTOShipmentStatusParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "APPROVED"},
		IfMatch:       eTag,
	}

	// Used for all tests except 500 error:
	queryBuilder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(queryBuilder)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)
	planner := &routemocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	planner.On("TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(1000, nil)

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	updater := mtoshipment.NewMTOShipmentStatusUpdater(suite.DB(), queryBuilder, siCreator, planner)
	handler := UpdateMTOShipmentStatusHandlerFunc{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		fetcher,
		updater,
	}

	suite.T().Run("Patch failure - 500", func(t *testing.T) {
		mockFetcher := mocks.Fetcher{}
		mockUpdater := mocks.MTOShipmentStatusUpdater{}
		mockHandler := UpdateMTOShipmentStatusHandlerFunc{
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

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusInternalServerError{}, response)

		errResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusInternalServerError)
		suite.Equal(handlers.InternalServerErrMessage, string(*errResponse.Payload.Title), "Payload title is wrong")

	})

	suite.T().Run("Patch failure - 404", func(t *testing.T) {
		notFoundParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: strfmt.UUID(uuid.Nil.String()),
			Body:          params.Body,
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(notFoundParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusNotFound{}, response)
	})

	suite.T().Run("Patch failure - 412", func(t *testing.T) {
		preconditionParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          params.Body,
			IfMatch:       "eTag",
		}
		response := handler.Handle(preconditionParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusPreconditionFailed{}, response)
	})

	suite.T().Run("Patch failure - 422", func(t *testing.T) {
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
	suite.T().Run("Successful patch - Integration Test", func(t *testing.T) {
		move := mtoShipment.MoveTaskOrder
		move.Status = models.MoveStatusAPPROVED
		_ = suite.DB().Save(&move)
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
		suite.NotNil(okResponse.Payload.ETag)
	})

	// Last because the shipment has to be either APPROVED or REJECTED before triggering this conflict:
	suite.T().Run("Patch failure - 409", func(t *testing.T) {
		conflictParams := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "SUBMITTED"},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(conflictParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusConflict{}, response)
	})

	// Test Successful Cancellation Request
	suite.T().Run("Successful patch - Integration Test for CANCELLATION_REQUESTED", func(t *testing.T) {
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
		params := mtoshipmentops.UpdateMTOShipmentStatusParams{
			HTTPRequest:   req,
			MtoShipmentID: *handlers.FmtUUID(mtoShipment.ID),
			Body:          &supportmessages.UpdateMTOShipmentStatus{Status: "CANCELLATION_REQUESTED"},
			IfMatch:       eTag,
		}

		suite.NoError(params.Body.Validate(strfmt.Default))
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentStatusOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentStatusOK)
		suite.Equal(supportmessages.UpdateMTOShipmentStatusStatusCANCELLATIONREQUESTED, okResponse.Payload.Status)
		suite.NotZero(okResponse.Payload.ETag)
	})
}
