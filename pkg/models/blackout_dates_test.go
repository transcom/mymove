package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_FetchTSPBlackoutDates() {
	t := suite.T()
	// Use FetchTSPBlackoutDates on two queries: one that should use a market value and one that doesn't.
	// Create one blackout date object with a market.
	tsp := testdatagen.MakeDefaultTSP(suite.db)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	blackoutStartDate := time.Now()
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	pickupDate := blackoutStartDate.Add(time.Hour)
	market1 := "dHHG"
	sourceGBLOC := "KKFA"
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

	shipmentDomesticMarket := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			ActualPickupDate: &pickupDate,
			BookDate:         &testdatagen.DateInsidePerformancePeriod,
			Status:           models.ShipmentStatusSUBMITTED,
		},
	})

	shipmentInternationalMarket := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			ActualPickupDate: &pickupDate,
			BookDate:         &testdatagen.DateInsidePerformancePeriod,
			Status:           models.ShipmentStatusSUBMITTED,
		},
	})

	fetchWithDomesticMarket, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentDomesticMarket)
	if err != nil {
		t.Errorf("Error fetching blackout dates.")
	} else if len(fetchWithDomesticMarket) == 0 {
		t.Errorf("Blackout dates query should have returned one result but returned zero instead.")
	}

	fetchWithInternationalMarket, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentInternationalMarket)
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
	tsp := testdatagen.MakeDefaultTSP(suite.db)
	tdl := testdatagen.MakeDefaultTDL(suite.db)
	blackoutStartDate := time.Now()
	blackoutEndDate := blackoutStartDate.Add(time.Hour * 24 * 2)
	pickupDate := blackoutStartDate.Add(time.Hour)
	market1 := "dHHG"
	sourceGBLOC1 := "KKFA"
	destinationGBLOC1 := "HAFC"
	sourceGBLOC2 := "KKNO"
	destinationGBLOC2 := "HANO"
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

	shipmentInGBLOC1 := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			ActualPickupDate: &pickupDate,
			SourceGBLOC:      &sourceGBLOC1,
			DestinationGBLOC: &destinationGBLOC1,
			Market:           &market1,
			BookDate:         &testdatagen.DateInsidePerformancePeriod,
			Status:           models.ShipmentStatusSUBMITTED,
		},
	})

	shipmentInGBLOC2 := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			ActualPickupDate: &pickupDate,
			SourceGBLOC:      &sourceGBLOC2,
			DestinationGBLOC: &destinationGBLOC2,
			BookDate:         &testdatagen.DateInsidePerformancePeriod,
			Status:           models.ShipmentStatusSUBMITTED,
		},
	})

	fetchWithMatchingGBLOC, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentInGBLOC1)
	if err != nil {
		t.Errorf("Error fetching blackout dates.")
	} else if len(fetchWithMatchingGBLOC) != 1 {
		t.Errorf("Blackout dates query should have returned one result but returned zero instead.")
	}

	fetchWithMismatchedGBLOC, err := FetchTSPBlackoutDates(suite.db, tsp.ID, shipmentInGBLOC2)
	if err != nil {
		t.Errorf("Error fetching blackout dates: %s.", err)
	} else if len(fetchWithMismatchedGBLOC) != 0 {
		t.Errorf("Blackout dates query should have returned no results but returned one instead.")
	}
}
