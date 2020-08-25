package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestMTOShipmentValidation() {
	suite.T().Run("test valid MTOShipment", func(t *testing.T) {
		// mock weights
		estimatedWeight := unit.Pound(1000)
		actualWeight := unit.Pound(980)
		validMTOShipment := models.MTOShipment{
			MoveTaskOrderID:      uuid.Must(uuid.NewV4()),
			Status:               models.MTOShipmentStatusApproved,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOShipment, expErrors)
	})

	suite.T().Run("test empty MTOShipment", func(t *testing.T) {
		emptyMTOShipment := models.MTOShipment{}
		expErrors := map[string][]string{
			"move_task_order_id": {"MoveTaskOrderID can not be blank."},
			"status":             {"Status is not in the list [APPROVED, REJECTED, SUBMITTED, DRAFT]."},
		}
		suite.verifyValidationErrors(&emptyMTOShipment, expErrors)
	})

	suite.T().Run("test rejected MTOShipment", func(t *testing.T) {
		rejectionReason := "bad shipment"
		rejectedMTOShipment := models.MTOShipment{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			Status:          models.MTOShipmentStatusRejected,
			RejectionReason: &rejectionReason,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&rejectedMTOShipment, expErrors)
	})

	suite.T().Run("test validation failures", func(t *testing.T) {
		// mock weights
		estimatedWeight := unit.Pound(-1000)
		actualWeight := unit.Pound(-980)
		invalidMTOShipment := models.MTOShipment{
			MoveTaskOrderID:      uuid.Must(uuid.NewV4()),
			Status:               models.MTOShipmentStatusRejected,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
		}
		expErrors := map[string][]string{
			"prime_estimated_weight": {"-1000 is not greater than -1."},
			"prime_actual_weight":    {"-980 is not greater than -1."},
			"rejection_reason":       {"RejectionReason can not be blank."},
		}
		suite.verifyValidationErrors(&invalidMTOShipment, expErrors)
	})
}
