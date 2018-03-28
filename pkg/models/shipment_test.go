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
		"source_gbloc":                 []string{"source_gbloc can not be blank."},
	}

	suite.verifyValidationErrors(shipment, expErrors)
}

// Test_FetchAllShipments tests that a shipment is returned when we fetch shipments with their offers.
func (suite *ModelSuite) Test_FetchAllShipments() {
	t := suite.T()
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl)
	shipment2, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl)
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	CreateShipmentOffer(suite.db, shipment.ID, tsp.ID, false)
	shipments, err := FetchShipments(suite.db, false)

	// Expect both shipments returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if shipments[0].ID != shipment.ID || shipments[1].ID != shipment2.ID {
		t.Errorf("Failed to return correct shipments. Expected shipments %v and %v, got %v and %v",
			shipment.ID, shipment2.ID, shipments[0].ID, shipments[1].ID)
	}
}

// Test_FetchUnassignedShipments tests that a shipment is returned when we fetch shipments with offers.
func (suite *ModelSuite) Test_FetchUnassignedShipments() {
	t := suite.T()
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl)
	shipment2, _ := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl)
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	CreateShipmentOffer(suite.db, shipment.ID, tsp.ID, false)
	shipments, err := FetchShipments(suite.db, true)

	// Expect only unassigned shipment returned
	if err != nil {
		t.Errorf("Failed to find Shipments: %v", err)
	} else if len(shipments) != 1 {
		t.Errorf("Returned too many shipments. Expected %v, got %v", shipment2.ID, shipments)
	}
}

func equalShipmentsSlice(a []ShipmentWithOffer, b []ShipmentWithOffer) bool {
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
