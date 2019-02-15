package models_test

import (
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/dates"
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
	calendar := dates.NewUSCalendar()
	pickupDate := dates.NextWorkday(*calendar, time.Date(testdatagen.TestYear, time.January, 28, 0, 0, 0, 0, time.UTC))
	deliveryDate := dates.NextWorkday(*calendar, pickupDate)
	tdl := testdatagen.MakeDefaultTDL(suite.DB())
	tsp := testdatagen.MakeDefaultTSP(suite.DB())
	tspp := testdatagen.MakeDefaultTSPPerformance(suite.DB())

	sourceGBLOC := "KKFA"
	destinationGBLOC := "HAFC"
	market := "dHHG"

	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: Shipment{
			RequestedPickupDate:     &pickupDate,
			ActualPickupDate:        &pickupDate,
			ActualDeliveryDate:      &deliveryDate,
			TrafficDistributionList: &tdl,
			SourceGBLOC:             &sourceGBLOC,
			DestinationGBLOC:        &destinationGBLOC,
			Market:                  &market,
		},
	})

	shipmentOffer, err := CreateShipmentOffer(suite.DB(), shipment.ID, tsp.ID, tspp.ID, false)
	suite.Nil(err, "error making ShipmentOffer")

	expectedShipmentOffer := ShipmentOffer{}
	if err := suite.DB().Find(&expectedShipmentOffer, shipmentOffer.ID); err != nil {
		t.Fatalf("could not find shipmentOffer: %v", err)
	}
}

func (suite *ModelSuite) TestShipmentOfferStateMachine() {
	// Try to accept an offer
	shipmentOffer := testdatagen.MakeDefaultShipmentOffer(suite.DB())
	suite.Nil(shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	err := shipmentOffer.Accept()
	suite.Nil(err)
	suite.True(*shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	// Try to Reject an offer
	shipmentOffer = testdatagen.MakeDefaultShipmentOffer(suite.DB())
	suite.Nil(shipmentOffer.Accepted)
	suite.Nil(shipmentOffer.RejectionReason)

	err = shipmentOffer.Reject("DO NOT WANT")
	suite.Nil(err)
	suite.False(*shipmentOffer.Accepted)
	suite.Equal("DO NOT WANT", *shipmentOffer.RejectionReason)
}

func (suite *ModelSuite) TestAccepted() {
	// Trying a nil slice of shipment offers.
	var shipmentOffers ShipmentOffers
	shipmentOffers, err := shipmentOffers.Accepted()
	suite.Nil(err)
	suite.Nil(shipmentOffers)

	// Make a default shipment offer (which shouldn't be accepted).
	unacceptedOffer := testdatagen.MakeDefaultShipmentOffer(suite.DB())
	shipmentOffers = ShipmentOffers{unacceptedOffer}
	shipmentOffers, err = shipmentOffers.Accepted()
	suite.Nil(err)
	suite.Nil(shipmentOffers)

	// Add an accepted shipment to our slice.
	acceptedOffer := testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: ShipmentOffer{
			Shipment: unacceptedOffer.Shipment,
			Accepted: swag.Bool(true),
		},
	})
	shipmentOffers = append(shipmentOffers, acceptedOffer)
	acceptedOffers, _ := shipmentOffers.Accepted()
	suite.Len(acceptedOffers, 1)
	suite.Equal(acceptedOffer.ID, acceptedOffers[0].ID)
}
