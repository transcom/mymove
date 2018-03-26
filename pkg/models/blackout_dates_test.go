package models_test

import (
  "time"

  "github.com/go-openapi/swag"

  . "github.com/transcom/mymove/pkg/models"
  "github.com/transcom/mymove/pkg/models"
  "github.com/transcom/mymove/pkg/testdatagen"
  )

func (suite *ModelSuite) Test_FetchTSPBlackoutDates() {
	t := suite.T()
	// TODO: add test for this model
  // Use FetchTSPBlackoutDates on two queries: one that should use a market value and one that doesn't.
  // Create one blackout date object with a market.
  tsp, err := testdatagen.MakeTSP(suite.db, "A Very Excellent TSP", "XYZA")
  tdl, err := testdatagen.MakeTDL(suite.db, "Oklahoma", "62240", "5")
  blackoutStartDate := time.Now()
  blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
  pickupDate := blackoutStartDate.Add(time.Hour)
  deliverDate := blackoutStartDate.Add(time.Hour * 24 * 60)
  market := "dHHG"
  testdatagen.MakeBlackoutDate(suite.db, tsp, blackoutStartDate, blackoutEndDate, &tdl, nil, &market)

  // Create a time within the blackout date.
  dateDuring := blackoutStartDate.Add(time.Hour)

  // Create two shipments, one with market and one without.
  shipmentWithMarket, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl)
  shipmentWithMarket.Market = &market
  shipmentWithoutMarket, _ := testdatagen.MakeShipment(suite.db, pickupDate, pickupDate, deliverDate, tdl)

  // Create two ShipmentWithOffers, one using the first set of times and a market, the other using the same times but without a market.
  shipmentWithOfferWithMarket := models.ShipmentWithOffer{
    ID:                              shipmentWithMarket.ID,
    TrafficDistributionListID:       tdl.ID,
    PickupDate:                      pickupDate,
    TransportationServiceProviderID: nil,
    Accepted:                        nil,
    RejectionReason:                 nil,
    AdministrativeShipment:          swag.Bool(false),
    BookDate:                        testdatagen.DateInsidePerformancePeriod,
  }

  shipmentWithOfferWithoutMarket := models.ShipmentWithOffer{
    ID:                              shipmentWithoutMarket.ID,
    TrafficDistributionListID:       tdl.ID,
    PickupDate:                      pickupDate,
    TransportationServiceProviderID: nil,
    Accepted:                        nil,
    RejectionReason:                 nil,
    AdministrativeShipment:          swag.Bool(false),
    BookDate:                        testdatagen.DateInsidePerformancePeriod,
  }

  fetchWithMarket, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentWithMarket)
  if err != nil {
    t.Errorf("Error fetching blackout dates.")
  } else if fetchWithMarket == false {
    t.Errorf("Blackout dates query erroneously returned false.")
  }

  fetchWithoutMarket, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentWithoutMarket)
  if err != nil {
    t.Errorf("Error fetching blackout dates.")
  } else if fetchWithoutMarket == false {
    t.Errorf("Blackout dates query erroneously returned false.")
  }
}
