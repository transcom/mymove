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
		pickupAddressUUID := uuid.Must(uuid.NewV4())
		destinationAddressUUID := uuid.Must(uuid.NewV4())
		validMTOShipment := models.MTOShipment{
			MoveTaskOrderID:      uuid.Must(uuid.NewV4()),
			PickupAddressID:      &pickupAddressUUID,
			DestinationAddressID: &destinationAddressUUID,
			Status:               models.MTOShipmentStatusApproved,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOShipment, expErrors)
	})
}
