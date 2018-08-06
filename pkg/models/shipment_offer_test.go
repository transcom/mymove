package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_ShipmentOfferValidations() {
	sa := &ShipmentOffer{}

	var expErrors = map[string][]string{
		"shipment_id":                        []string{"ShipmentID can not be blank."},
		"transportation_service_provider_id": []string{"TransportationServiceProviderID can not be blank."},
	}

	suite.verifyValidationErrors(sa, expErrors)
}

// Test_CreateShipmentOffer tests that a shipment is created when expected
func (suite *ModelSuite) Test_CreateShipmentOffer() {
	t := suite.T()
	now := time.Now()
	tdl, err := testdatagen.MakeTDL(suite.db,
		testdatagen.DefaultSrcRateArea,
		testdatagen.DefaultDstRegion,
		testdatagen.DefaultCOS)
	suite.Nil(err, "error making TDL")

	tsp, err := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	suite.Nil(err, "error making TSP")

	market := "dHHG"
	shipment, err := testdatagen.MakeShipment(suite.db, now, now, now.AddDate(0, 0, 1), tdl, "OHAI", &market, nil)
	suite.Nil(err, "error making Shipment")
	shipmentOffer, err := CreateShipmentOffer(suite.db, shipment.ID, tsp.ID, false)
	suite.Nil(err, "error making ShipmentOffer")

	expectedShipmentOffer := ShipmentOffer{}
	if err := suite.db.Find(&expectedShipmentOffer, shipmentOffer.ID); err != nil {
		t.Fatalf("could not find shipmentOffer: %v", err)
	}
}
