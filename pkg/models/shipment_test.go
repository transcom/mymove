package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"time"
)

func (suite *ModelSuite) Test_ShipmentValidations() {
	shipment := &Shipment{}

	expErrors := map[string][]string{
		"traffic_distribution_list_id": []string{"traffic_distribution_list_id can not be blank."},
		"gbloc": []string{"gbloc can not be blank."},
	}

	suite.verifyValidationErrors(shipment, expErrors)
}

// Test_FetchPossiblyAwardedShipments tests that a shipment is returned when we fetch possibly awarded shipments
func (suite *ModelSuite) Test_FetchPossiblyAwardedShipments() {
	t := suite.T()
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now.AddDate(0, 0, 1), tdl)

	shipments, err := FetchPossiblyAwardedShipments(suite.db)

	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if shipments[0].ID != shipment.ID {
		t.Errorf("Wrong shipment returned: expected %s, got %s",
			shipment.ID, shipments[0].ID)
	}
}
