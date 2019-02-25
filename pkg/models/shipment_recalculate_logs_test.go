package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/testdatagen"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestShipmentRecalculateLog() {

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{})
	shipmentID := shipment.ID

	testCases := map[string]struct {
		recalculateLog models.ShipmentRecalculateLog
		expectedErrs   map[string][]string
	}{
		"Successful Create": {
			recalculateLog: models.ShipmentRecalculateLog{
				ID:         uuid.Must(uuid.NewV4()),
				ShipmentID: shipmentID,
				Shipment:   shipment,
			},
			expectedErrs: nil,
		},

		"Empty Fields": {
			recalculateLog: models.ShipmentRecalculateLog{},
			expectedErrs: map[string][]string{
				"shipment_id": {"ShipmentID can not be blank."},
			},
		},
	}

	for name, test := range testCases {
		suite.T().Run(name, func(t *testing.T) {
			suite.verifyValidationErrors(&test.recalculateLog, test.expectedErrs)
		})
	}
}
