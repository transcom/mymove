package primeapi

import (
	"fmt"
	"net/http/httptest"
	"time"

	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/gofrs/uuid"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/etag"
	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	reweighservice "github.com/transcom/mymove/pkg/services/reweigh"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const (
	recalculateTestPickupZip      = "30907"
	recalculateTestDestinationZip = "78234"
	recalculateTestZip3Distance   = 1234
)

func (suite *HandlerSuite) TestUpdateReweighHandler() {

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("ZipTransitDistance",
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

	setupTestData := func() (UpdateReweighHandler, models.Reweigh) {

		// Create handler
		handler := UpdateReweighHandler{
			handlers.NewHandlerConfig(suite.DB(), suite.Logger()),
			reweighservice.NewReweighUpdater(movetaskorder.NewMoveTaskOrderChecker(), paymentRequestShipmentRecalculator),
		}
		// Make an available MTO
		mto := testdatagen.MakeAvailableMove(suite.DB())

		// Make a shipment on the available MTO
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: mto,
		})

		// Make Reweigh Request
		reweigh := testdatagen.MakeReweighWithNoWeightForShipment(
			suite.DB(),
			testdatagen.Assertions{},
			shipment,
		)
		return handler, reweigh
	}

	suite.Run("Success 200 - Update reweigh weight", func() {
		// Testcase:   reweigh us updated with the new weight of the shipment
		// Expected:   Success response 200
		handler, reweigh := setupTestData()

		// Update with weights
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s/rewighs/%s", reweigh.ShipmentID.String(), reweigh.ID.String()), nil)
		weight := int64(8000)
		params := mtoshipmentops.UpdateReweighParams{
			HTTPRequest:   req,
			ReweighID:     *handlers.FmtUUID(reweigh.ID),
			MtoShipmentID: *handlers.FmtUUID(reweigh.ShipmentID),
			IfMatch:       etag.GenerateEtag(reweigh.UpdatedAt),
			Body: &primemessages.UpdateReweigh{
				Weight: &weight,
			},
		}
		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateReweighOK{}, response)

		// Check values
		reweighOk := response.(*mtoshipmentops.UpdateReweighOK)
		suite.Equal(&weight, reweighOk.Payload.Weight)
	})

	suite.Run("Success 200 - Update reweigh verification reason", func() {
		// Testcase:   reweigh is updated with the verification reason
		// Expected:   Success response 200
		handler, reweigh := setupTestData()

		// Update with verification reason
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s/rewighs/%s", reweigh.ShipmentID.String(), reweigh.ID.String()), nil)
		reason := "The shipment was already delivered."
		params := mtoshipmentops.UpdateReweighParams{
			HTTPRequest:   req,
			ReweighID:     *handlers.FmtUUID(reweigh.ID),
			MtoShipmentID: *handlers.FmtUUID(reweigh.ShipmentID),
			IfMatch:       etag.GenerateEtag(reweigh.UpdatedAt),
			Body: &primemessages.UpdateReweigh{
				VerificationReason: &reason,
			},
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateReweighOK{}, response)

		// Check values
		reweighOk := response.(*mtoshipmentops.UpdateReweighOK)

		suite.Equal(&reason, reweighOk.Payload.VerificationReason)
	})

	suite.Run("Failure 422 - Failed to update reweigh weight due to bad request - zero reweigh value", func() {
		// Testcase:   reweigh us updated with the new weight of the shipment
		// Expected:   Failure 422
		_, reweigh := setupTestData()

		// Update with weights

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s/rewighs/%s", reweigh.ShipmentID.String(), reweigh.ID.String()), nil)
		weight := int64(0)

		params := mtoshipmentops.UpdateReweighParams{
			HTTPRequest:   req,
			ReweighID:     *handlers.FmtUUID(reweigh.ID),
			MtoShipmentID: *handlers.FmtUUID(reweigh.ShipmentID),
			IfMatch:       etag.GenerateEtag(reweigh.UpdatedAt),
			Body: &primemessages.UpdateReweigh{
				Weight: &weight,
			},
		}

		// Run swagger validations
		err := params.Body.Validate(strfmt.Default)
		suite.Equal("validation failure list:\nweight in body should be greater than or equal to 1", err.Error())
	})

	suite.Run("Failure 422 - Failed to update reweigh weight due to bad request - negative reweigh value", func() {
		// Testcase:   reweigh us updated with the new weight of the shipment
		// Expected:   Failure response 422
		_, reweigh := setupTestData()

		// Update with weights

		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s/rewighs/%s", reweigh.ShipmentID.String(), reweigh.ID.String()), nil)
		weight := int64(-10)

		params := mtoshipmentops.UpdateReweighParams{
			HTTPRequest:   req,
			ReweighID:     *handlers.FmtUUID(reweigh.ID),
			MtoShipmentID: *handlers.FmtUUID(reweigh.ShipmentID),
			IfMatch:       etag.GenerateEtag(reweigh.UpdatedAt),
			Body: &primemessages.UpdateReweigh{
				Weight: &weight,
			},
		}

		// Run swagger validations
		err := params.Body.Validate(strfmt.Default)
		suite.Equal("validation failure list:\nweight in body should be greater than or equal to 1", err.Error())
	})

	suite.Run("Failure 404 - Reweigh not found", func() {
		// Testcase:   Reweigh ID is not found
		// Expected:   Failure response 404
		handler, reweigh := setupTestData()

		// Update with verification reason
		badID, _ := uuid.NewV4()
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s/rewighs/%s", reweigh.ShipmentID.String(), badID.String()), nil)
		reason := "The shipment was already delivered."
		params := mtoshipmentops.UpdateReweighParams{
			HTTPRequest:   req,
			ReweighID:     *handlers.FmtUUID(badID),
			MtoShipmentID: *handlers.FmtUUID(reweigh.ShipmentID),
			IfMatch:       etag.GenerateEtag(reweigh.UpdatedAt),
			Body: &primemessages.UpdateReweigh{
				VerificationReason: &reason,
			},
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateReweighNotFound{}, response)
	})

	suite.Run("Fail - PreconditionFailed due to wrong etag", func() {
		// Testcase:   etag for reweigh is wrong
		// Expected:   PreconditionFailed error is returned
		handler, reweigh := setupTestData()

		// Update with reweigh with a bad etag
		// Testcase:   Reweigh updated with incorrect etag
		// Expected:   Failure response 404

		// Update with verification reason
		req := httptest.NewRequest("PATCH", fmt.Sprintf("/mto-shipments/%s/rewighs/%s", reweigh.ShipmentID.String(), reweigh.ID.String()), nil)
		reason := "The shipment was already delivered."
		params := mtoshipmentops.UpdateReweighParams{
			HTTPRequest:   req,
			ReweighID:     *handlers.FmtUUID(reweigh.ID),
			MtoShipmentID: *handlers.FmtUUID(reweigh.ShipmentID),
			IfMatch:       etag.GenerateEtag(time.Now()),
			Body: &primemessages.UpdateReweigh{
				VerificationReason: &reason,
			},
		}

		// Run swagger validations
		suite.NoError(params.Body.Validate(strfmt.Default))

		// Run handler and check response
		response := handler.Handle(params)
		suite.IsType(&mtoshipmentops.UpdateReweighPreconditionFailed{}, response)
	})

}
