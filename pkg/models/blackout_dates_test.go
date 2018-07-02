package models_test

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_FetchTSPBlackoutDates() {
	t := suite.T()
	// Use FetchTSPBlackoutDates on two queries: one that should use a market value and one that doesn't.
	// Create one blackout date object with a market.
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, testdatagen.DefaultCOS)
	blackoutStartDate := time.Now()
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	pickupDate := blackoutStartDate.Add(time.Hour)
	market1 := "dHHG"
	sourceGBLOC := "OHAI"
	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp.ID,
			StartBlackoutDate:               blackoutStartDate,
			EndBlackoutDate:                 blackoutEndDate,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &sourceGBLOC,
			Market:                          &market1,
		},
	})

	// Create two ShipmentWithOffers, one using the first set of times and a market, the other using the same times but without a market.
	shipmentWithOfferWithDomesticMarket := models.ShipmentWithOffer{
		ID: uuid.Must(uuid.NewV4()),
		TrafficDistributionListID:       tdl.ID,
		PickupDate:                      pickupDate,
		TransportationServiceProviderID: nil,
		Accepted:                        nil,
		RejectionReason:                 nil,
		AdministrativeShipment:          swag.Bool(false),
		BookDate:                        testdatagen.DateInsidePerformancePeriod,
	}

	shipmentWithOfferWithInternationalMarket := models.ShipmentWithOffer{
		ID: uuid.Must(uuid.NewV4()),
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

// Need a test where the GBLOC does and doesn't match
func (suite *ModelSuite) Test_FetchTSPBlackoutDatesWithGBLOC() {
	t := suite.T()
	// Use FetchTSPBlackoutDates on two queries: one that should use a market value and one that doesn't.
	// Create one blackout date object with a market.
	tsp, _ := testdatagen.MakeTSP(suite.db, testdatagen.RandomSCAC())
	tdl, _ := testdatagen.MakeTDL(suite.db, testdatagen.DefaultSrcRateArea, testdatagen.DefaultDstRegion, testdatagen.DefaultCOS)
	blackoutStartDate := time.Now()
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	pickupDate := blackoutStartDate.Add(time.Hour)
	market1 := "dHHG"
	sourceGBLOC1 := "OHAI"
	sourceGBLOC2 := "OHNO"
	testdatagen.MakeBlackoutDate(suite.db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp.ID,
			StartBlackoutDate:               blackoutStartDate,
			EndBlackoutDate:                 blackoutEndDate,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &sourceGBLOC1,
			Market:                          &market1,
		},
	})

	// Create two ShipmentWithOffers, one using the first set of times and a market, the other using the same times but without a market.
	shipmentWithOfferInGBLOC1 := models.ShipmentWithOffer{
		ID: uuid.Must(uuid.NewV4()),
		TrafficDistributionListID: tdl.ID,
		PickupDate:                pickupDate,
		SourceGBLOC:               &sourceGBLOC1,
		Market:                    &market1,
		AdministrativeShipment:    swag.Bool(false),
		BookDate:                  testdatagen.DateInsidePerformancePeriod,
	}

	shipmentWithOfferInGBLOC2 := models.ShipmentWithOffer{
		ID: uuid.Must(uuid.NewV4()),
		TrafficDistributionListID: tdl.ID,
		PickupDate:                pickupDate,
		SourceGBLOC:               &sourceGBLOC2,
		Market:                    nil,
		AdministrativeShipment:    swag.Bool(false),
		BookDate:                  testdatagen.DateInsidePerformancePeriod,
	}

	fetchWithMatchingGBLOC, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentWithOfferInGBLOC1)
	if err != nil {
		t.Errorf("Error fetching blackout dates.")
	} else if len(fetchWithMatchingGBLOC) != 1 {
		t.Errorf("Blackout dates query should have returned one result but returned zero instead.")
	}

	fetchWithMismatchedGBLOC, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentWithOfferInGBLOC2)
	if err != nil {
		t.Errorf("Error fetching blackout dates: %s.", err)
	} else if len(fetchWithMismatchedGBLOC) != 0 {
		t.Errorf("Blackout dates query should have returned no results but returned one instead.")
	}
}
