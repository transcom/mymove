package models_test

import (
	"time"

	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_ShipmentOfferValidations() {
	sa := &ShipmentOffer{}

	var expErrors = map[string][]string{
		"shipment_id":                                    []string{"ShipmentID can not be blank."},
		"transportation_service_provider_id":             []string{"TransportationServiceProviderID can not be blank."},
		"transportation_service_provider_performance_id": []string{"TransportationServiceProviderPerformanceID can not be blank."},
	}

	suite.verifyValidationErrors(sa, expErrors)
}

// Test_CreateShipmentOffer tests that a shipment is created when expected
func (suite *ModelSuite) Test_CreateShipmentOffer() {
	t := suite.T()
	pickupDate := time.Now()
	deliveryDate := time.Now().AddDate(0, 0, 1)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	tsp := testdatagen.MakeDefaultTSP(suite.db)
	tspp := testdatagen.MakeDefaultTSPPerformance(suite.db)

	sourceGBLOC := "OHAI"
	market := "dHHG"

	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			Market:                  &market,
		},
	})

	shipmentOffer, err := CreateShipmentOffer(suite.db, shipment.ID, tsp.ID, tspp.ID, false)
	suite.Nil(err, "error making ShipmentOffer")

	expectedShipmentOffer := ShipmentOffer{}
	if err := suite.db.Find(&expectedShipmentOffer, shipmentOffer.ID); err != nil {
		t.Fatalf("could not find shipmentOffer: %v", err)
	}
}

func (suite *ModelSuite) TestShipmentOfferStateMachine() {
	// Try to accept an offer
	shipmentOffer := testdatagen.MakeDefaultShipmentOffer(suite.db)
	suite.Nil(shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	err := shipmentOffer.Accept()
	suite.Nil(err)
	suite.True(*shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	// Try to Reject an offer
	shipmentOffer = testdatagen.MakeDefaultShipmentOffer(suite.db)
	suite.Nil(shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	err = shipmentOffer.Reject("DO NOT WANT")
	suite.Nil(err)
	suite.False(*shipmentOffer.Accepted)
	suite.Equal("DO NOT WANT", *shipmentOffer.RejectionReason)
}
