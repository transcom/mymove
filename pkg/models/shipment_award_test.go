package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_ShipmentAwardValidations() {
	sa := &ShipmentAward{}

	var expErrors = map[string][]string{
		"shipment_id":                        []string{"ShipmentID can not be blank."},
		"transportation_service_provider_id": []string{"TransportationServiceProviderID can not be blank."},
	}

	suite.verifyValidationErrors(sa, expErrors)
}

// Test_CreateShipmentAward tests that a shipment is created when expected
func (suite *ModelSuite) Test_CreateShipmentAward() {
	t := suite.T()
	now := time.Now()
	tdl, _ := testdatagen.MakeTDL(suite.db, "california", "90210", "2")
	tsp, _ := testdatagen.MakeTSP(suite.db, "Test TSP 1", "TSP1")
	shipment, _ := testdatagen.MakeShipment(suite.db, now, now.AddDate(0, 0, 1), tdl)
	shipmentAward, err := CreateShipmentAward(suite.db, shipment.ID, tsp.ID, false)

	if err != nil {
		t.Errorf("Failed to create Shipment Award: %v", err)
	}
	expectedShipmentAward := ShipmentAward{}
	if err := suite.db.Find(&expectedShipmentAward, shipmentAward.ID); err != nil {
		t.Fatalf("could not find shipmentAward: %v", err)
	}
}
