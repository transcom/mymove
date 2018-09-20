package scenario

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// RunAwardQueueScenario1 creates 17 shipments and 5 TSPs in 1 TDL. This allows testing against
// award queue to ensure it behaves as expected. This doesn't track blackout dates.
func RunAwardQueueScenario1(db *pop.Connection) {
	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl := testdatagen.MakeTDL(db, testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US13",
			DestinationRegion: "5",
			CodeOfService:     "2",
		},
	})

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		now := time.Now()
		testdatagen.MakeShipment(db, testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate:     &now,
				ActualPickupDate:        &now,
				DeliveryDate:            &now,
				TrafficDistributionList: &tdl,
				SourceGBLOC:             &sourceGBLOC,
				Market:                  &market,
			},
		})
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1 := testdatagen.MakeDefaultTSP(db)
	tsp2 := testdatagen.MakeDefaultTSP(db)
	tsp3 := testdatagen.MakeDefaultTSP(db)
	tsp4 := testdatagen.MakeDefaultTSP(db)
	tsp5 := testdatagen.MakeDefaultTSP(db)

	// TSPs should be ordered by offer_count first, then BVS.
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp1, tdl, swag.Int(1), 5, 0, 0.42, 0.42)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp2, tdl, swag.Int(1), 4, 0, 0.33, 0.33)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp3, tdl, swag.Int(2), 3, 0, 0.21, 0.21)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp4, tdl, swag.Int(3), 2, 0, 0.11, 0.11)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp5, tdl, swag.Int(4), 1, 0, 0.05, 0.05)
}

// RunAwardQueueScenario2 creates 9 shipments to be divided between 5 TSPs in 1 TDL and 10 shipments to be divided among 4 TSPs in TDL 2.
// This allows testing against award queue to ensure it behaves as expected. Two TSPs in TDL1 and one TSP in TDL 2 have blackout dates.
func RunAwardQueueScenario2(db *pop.Connection) {
	shipmentsToMake := 9
	shipmentDate := time.Now()

	// Make a TDL to contain our tests
	tdl := testdatagen.MakeTDL(db, testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US13",
			DestinationRegion: "15",
			CodeOfService:     "2",
		},
	})
	tdl2 := testdatagen.MakeTDL(db, testdatagen.Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US62",
			DestinationRegion: "1",
			CodeOfService:     "2",
		},
	})

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in first TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate:     &shipmentDate,
				ActualPickupDate:        &shipmentDate,
				DeliveryDate:            &shipmentDate,
				TrafficDistributionList: &tdl,
				SourceGBLOC:             &sourceGBLOC,
				Market:                  &market,
			},
		})
	}
	// Make shipments in second TDL
	for i := 0; i <= shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, testdatagen.Assertions{
			Shipment: models.Shipment{
				RequestedPickupDate:     &shipmentDate,
				ActualPickupDate:        &shipmentDate,
				DeliveryDate:            &shipmentDate,
				TrafficDistributionList: &tdl2,
				SourceGBLOC:             &sourceGBLOC,
				Market:                  &market,
			},
		})
	}

	// Make TSPs
	tsp1 := testdatagen.MakeDefaultTSP(db) // Good TSP with blackout date
	tsp2 := testdatagen.MakeDefaultTSP(db) // Very good TSP, no blackout date
	tsp3 := testdatagen.MakeDefaultTSP(db) // Pretty good TSP, no blackout date
	tsp4 := testdatagen.MakeDefaultTSP(db) // So-so TSP with blackout date
	tsp5 := testdatagen.MakeDefaultTSP(db) // Meh TSP, no blackout date
	tsp6 := testdatagen.MakeDefaultTSP(db) // Sterling TSP with no blackout date
	tsp7 := testdatagen.MakeDefaultTSP(db) // Decent TSP with blackout date
	tsp8 := testdatagen.MakeDefaultTSP(db) // Decent TSP,  no blackout date
	tsp9 := testdatagen.MakeDefaultTSP(db) // V v bad TSP

	// Put TSPs in 2 TDLs to handle these shipments
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp1, tdl, swag.Int(1), 5, 0, 0.42, 0.44)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp2, tdl, swag.Int(1), 4, 0, 0.31, 0.32)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp3, tdl, swag.Int(2), 3, 0, 0.24, 0.25)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp4, tdl, swag.Int(3), 2, 0, 0.11, 0.13)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp5, tdl, swag.Int(4), 1, 0, 0.05, 0.08)

	testdatagen.MakeTSPPerformanceDeprecated(db, tsp6, tdl2, swag.Int(1), 5, 0, 0.42, 0.44)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp7, tdl2, swag.Int(2), 4, 0, 0.31, 0.32)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp8, tdl2, swag.Int(3), 2, 0, 0.11, 0.13)
	testdatagen.MakeTSPPerformanceDeprecated(db, tsp9, tdl2, swag.Int(4), 1, 0, 0.05, 0.08)
	// Add blackout dates
	blackoutStart := shipmentDate.AddDate(0, 0, -3)
	blackoutEnd := shipmentDate.AddDate(0, 0, 3)

	gbloc := "BKAS"
	testdatagen.MakeBlackoutDate(db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp1.ID,
			StartBlackoutDate:               blackoutStart,
			EndBlackoutDate:                 blackoutEnd,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &gbloc,
			Market:                          &market,
		},
	})
	testdatagen.MakeBlackoutDate(db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp4.ID,
			StartBlackoutDate:               blackoutStart,
			EndBlackoutDate:                 blackoutEnd,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &gbloc,
			Market:                          &market,
		},
	})
	testdatagen.MakeBlackoutDate(db, testdatagen.Assertions{
		BlackoutDate: models.BlackoutDate{
			TransportationServiceProviderID: tsp7.ID,
			StartBlackoutDate:               blackoutStart,
			EndBlackoutDate:                 blackoutEnd,
			TrafficDistributionListID:       &tdl.ID,
			SourceGBLOC:                     &gbloc,
			Market:                          &market,
		},
	})
}
