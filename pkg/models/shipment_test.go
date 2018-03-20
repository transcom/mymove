package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
	// Make a shipment to award
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now.AddDate(0, 0, 1), tdl)
	// Make a shipment not to award
	shipment2, _ := testdatagen.MakeShipment(suite.db, now, now.AddDate(0, 0, 1), tdl)
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	// Award one of the shipments
	CreateShipmentAward(suite.db, shipment.ID, tsp.ID, false)
	// Run FetchPossiblyAwardedShipments
	shipments, err := FetchPossiblyAwardedShipments(suite.db)

	// Expect both shipments returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if shipments[0].ID != shipment.ID || shipments[1].ID != shipment2.ID {
		t.Errorf("Failed to return correct shipments. Expected shipments %v and %v, got %v and %v",
			shipment.ID, shipment2.ID, shipments[0].ID, shipments[1].ID)
	}
}

// Test_FetchUnawardedShipments tests that a shipment is returned when we fetch possibly awarded shipments
func (suite *ModelSuite) Test_FetchUnawardedShipments() {
	t := suite.T()
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now.AddDate(0, 0, 1), tdl)
	shipment2, _ := testdatagen.MakeShipment(suite.db, now, now.AddDate(0, 0, 1), tdl)
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	// Award one of the shipments
	CreateShipmentAward(suite.db, shipment.ID, tsp.ID, false)
	// Run FetchUnawardedShipments
	shipments, err := FetchUnawardedShipments(suite.db)
	// Expect only unawarded shipment returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if len(shipments) != 1 {
		t.Errorf("Returned too many shipments. Expected %v, got %v", shipment2.ID, shipments)
	}
}

func equalShipmentsSlice(a []PossiblyAwardedShipment, b []PossiblyAwardedShipment) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
