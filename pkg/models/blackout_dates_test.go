package models_test

import (
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_FetchTSPBlackoutDates() {
	t := suite.T()
	// Use FetchTSPBlackoutDates on two queries: one that should use a market value and one that doesn't.
	// Create one blackout date object with a market.
	tsp, _ := testdatagen.MakeTSP(suite.db, "A Very Excellent TSP", "XYZA")
	tdl, _ := testdatagen.MakeTDL(suite.db, "Oklahoma", "62240", "5")
	blackoutStartDate := time.Now()
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	pickupDate := blackoutStartDate.Add(time.Hour)
	deliveryDate := blackoutStartDate.Add(time.Hour * 24 * 60)
	market1 := "dHHG"
	sourceGBLOC := "OHAI"
	testdatagen.MakeBlackoutDate(suite.db, tsp, blackoutStartDate, blackoutEndDate, &tdl, &sourceGBLOC, &market1)

	// Create two shipments, one with market and one without.
	shipmentWithDomesticMarket, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliveryDate, tdl, sourceGBLOC, market1)
	shipmentWithInternationalMarket, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliveryDate, tdl, sourceGBLOC, market1)

	// Create two ShipmentWithOffers, one using the first set of times and a market, the other using the same times but without a market.
	shipmentWithOfferWithDomesticMarket := models.ShipmentWithOffer{
		ID: shipmentWithDomesticMarket.ID,
		TrafficDistributionListID:       tdl.ID,
		PickupDate:                      pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        testdatagen.DateInsidePerformancePeriod,
	}

	shipmentWithOfferWithInternationalMarket := models.ShipmentWithOffer{
		ID: shipmentWithInternationalMarket.ID,
		TrafficDistributionListID:       tdl.ID,
		PickupDate:                      pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        testdatagen.DateInsidePerformancePeriod,
	}

	fetchWithDomesticMarket, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentWithOfferWithDomesticMarket)
	if err != nil {
		t.Errorf("Error fetching blackout dates.")
	} else if len(fetchWithDomesticMarket) == 0 {
		t.Errorf("Blackout dates query should have returned one result but returned zero instead.")
	}

	fetchWithInternationalMarket, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentWithOfferWithInternationalMarket)
	if err != nil {
		t.Errorf("Error fetching blackout dates.")
	} else if len(fetchWithInternationalMarket) == 0 {
		t.Errorf("Blackout dates query should have returned one result but returned zero instead.")
	}
}
