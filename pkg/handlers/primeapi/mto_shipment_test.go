package primeapi

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/mocks"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// ClearNonUpdateFields clears out the MTOShipment payload fields that CANNOT be sent in for a successful update
func ClearNonUpdateFields(mtoShipment *models.MTOShipment) *primemessages.MTOShipment {
	mtoShipment.MoveTaskOrderID = uuid.FromStringOrNil("")
	mtoShipment.CreatedAt = time.Time{}
	mtoShipment.UpdatedAt = time.Time{}
	mtoShipment.PrimeEstimatedWeightRecordedDate = &time.Time{}
	mtoShipment.RequiredDeliveryDate = &time.Time{}
	mtoShipment.ApprovedDate = &time.Time{}
	mtoShipment.Status = ""
	mtoShipment.RejectionReason = nil
	mtoShipment.CustomerRemarks = nil
	mtoShipment.MTOAgents = nil

	return payloads.MTOShipment(mtoShipment)
}

func (suite *HandlerSuite) TestUpdateMTOShipmentHandler() {
	primeEstimatedWeight := unit.Pound(500)
	primeEstimatedWeightDate := testdatagen.DateInsidePeakRateCycle
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	mtoShipment.PrimeEstimatedWeight = &primeEstimatedWeight
	mtoShipment.PrimeEstimatedWeightRecordedDate = &primeEstimatedWeightDate

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	_, _ = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment,
			MTOShipmentID: mtoShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	builder := query.NewQueryBuilder(suite.DB())
	fetcher := fetch.NewFetcher(builder)

	req := httptest.NewRequest("PUT", fmt.Sprintf("/mto-shipments/%s", mtoShipment.ID.String()), nil)
	eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)
	params := mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req,
		MtoShipmentID: *handlers.FmtUUID(mtoShipment.ID),
		Body:          ClearNonUpdateFields(&mtoShipment),
		IfMatch:       eTag,
	}
	// used for all tests except the 500 server error:
	updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, route.NewTestingPlanner(400))
	handler := UpdateMTOShipmentHandler{
		handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
		updater,
	}

	suite.T().Run("Successful PUT - Integration Test", func(t *testing.T) {
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment.ID.String(), okResponse.Payload.ID.String())
	})

	suite.T().Run("PUT failure - 500", func(t *testing.T) {
		mockUpdater := mocks.MTOShipmentUpdater{}
		mockHandler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			&mockUpdater,
		}
		internalServerErr := errors.New("ServerError")

		mockUpdater.On("UpdateMTOShipment",
			mock.Anything,
			mock.Anything,
		).Return(nil, internalServerErr)

		response := mockHandler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentInternalServerError{}, response)
	})

	suite.T().Run("PUT failure - 404", func(t *testing.T) {
		notFoundParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: strfmt.UUID(uuid.Nil.String()),
			Body:          &primemessages.MTOShipment{},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(notFoundParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentNotFound{}, response)
	})

	suite.T().Run("PUT failure - 409", func(t *testing.T) {
		remarks := fmt.Sprintf("test conflict %s", time.Now())
		conflictParams := mtoshipmentops.UpdateMTOShipmentParams{
			HTTPRequest:   params.HTTPRequest,
			MtoShipmentID: params.MtoShipmentID,
			Body:          &primemessages.MTOShipment{CustomerRemarks: &remarks},
			IfMatch:       params.IfMatch,
		}
		response := handler.Handle(conflictParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentConflict{}, response)
	})

	suite.T().Run("PUT failure - 412", func(t *testing.T) {
		staleParams := params
		// no need to update IfMatch or eTag value because it's still the old value from before the successful PUT
		staleParams.Body.PrimeEstimatedWeight = 0 // causes this test to fail because any updates to this are invalid
		// (and input validation happens first before this check)
		response := handler.Handle(staleParams)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentPreconditionFailed{}, response)
	})

	suite.T().Run("PUT failure - 422", func(t *testing.T) {
		params.Body.PrimeEstimatedWeight = 1 // cannot update once initial value has been set
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentUnprocessableEntity{}, response)
	})

	mto2 := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	mtoShipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOAgent(suite.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	payload := primemessages.MTOShipment{
		ID: strfmt.UUID(mtoShipment2.ID.String()),
	}

	req2 := httptest.NewRequest("PUT", fmt.Sprintf("/move_task_orders/%s/mto_shipments/%s", mto2.ID.String(), mtoShipment2.ID.String()), nil)

	eTag = etag.GenerateEtag(mtoShipment2.UpdatedAt)
	params = mtoshipmentops.UpdateMTOShipmentParams{
		HTTPRequest:   req2,
		MtoShipmentID: *handlers.FmtUUID(mtoShipment2.ID),
		Body:          &payload,
		IfMatch:       eTag,
	}

	suite.T().Run("Successful PUT - Integration Test with Only Required Fields in Payload", func(t *testing.T) {
		updater := mtoshipment.NewMTOShipmentUpdater(suite.DB(), builder, fetcher, route.NewTestingPlanner(400))
		handler := UpdateMTOShipmentHandler{
			handlers.NewHandlerContext(suite.DB(), suite.TestLogger()),
			updater,
		}

		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateMTOShipmentOK{}, response)

		okResponse := response.(*mtoshipmentops.UpdateMTOShipmentOK)
		suite.Equal(mtoShipment2.ID.String(), okResponse.Payload.ID.String())
	})
}
