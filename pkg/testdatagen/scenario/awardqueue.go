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
	tdl, _ := testdatagen.MakeTDL(db, "US13", "15", "2")

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, time.Now(), time.Now(), time.Now(), tdl, sourceGBLOC, &market, nil)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC())
	tsp2, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC())
	tsp3, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC())
	tsp4, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC())
	tsp5, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC())

	// TSPs should be ordered by offer_count first, then BVS.
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0, 0.42, 0.42)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0, 0.33, 0.33)
	testdatagen.MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0, 0.21, 0.21)
	testdatagen.MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0, 0.11, 0.11)
	testdatagen.MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0, 0.05, 0.05)
}

// RunAwardQueueScenario2 creates 9 shipments to be divided between 5 TSPs in 1 TDL and 10 shipments to be divided among 4 TSPs in TDL 2.
// This allows testing against award queue to ensure it behaves as expected. Two TSPs in TDL1 and one TSP in TDL 2 have blackout dates.
func RunAwardQueueScenario2(db *pop.Connection) {
	shipmentsToMake := 9
	shipmentDate := time.Now()

	// Make a TDL to contain our tests
	tdl, _ := testdatagen.MakeTDL(db, "US13", "15", "2")
	tdl2, _ := testdatagen.MakeTDL(db, "US62", "1", "2")

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in first TDL
	for i := 0; i < shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl, sourceGBLOC, &market, nil)
	}
	// Make shipments in second TDL
	for i := 0; i <= shipmentsToMake; i++ {
		testdatagen.MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl2, sourceGBLOC, &market, nil)
	}

	// Make TSPs
	tsp1, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Good TSP with blackout date
	tsp2, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Very good TSP, no blackout date
	tsp3, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Pretty good TSP, no blackout date
	tsp4, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // So-so TSP with blackout date
	tsp5, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Meh TSP, no blackout date
	tsp6, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Sterling TSP with no blackout date
	tsp7, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Decent TSP with blackout date
	tsp8, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // Decent TSP,  no blackout date
	tsp9, _ := testdatagen.MakeTSP(db, testdatagen.RandomSCAC()) // V v bad TSP

	// Put TSPs in 2 TDLs to handle these shipments
	testdatagen.MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0, 0.42, 0.44)
	testdatagen.MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0, 0.31, 0.32)
	testdatagen.MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0, 0.24, 0.25)
	testdatagen.MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0, 0.11, 0.13)
	testdatagen.MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0, 0.05, 0.08)

	testdatagen.MakeTSPPerformance(db, tsp6, tdl2, swag.Int(1), 5, 0, 0.42, 0.44)
	testdatagen.MakeTSPPerformance(db, tsp7, tdl2, swag.Int(2), 4, 0, 0.31, 0.32)
	testdatagen.MakeTSPPerformance(db, tsp8, tdl2, swag.Int(3), 2, 0, 0.11, 0.13)
	testdatagen.MakeTSPPerformance(db, tsp9, tdl2, swag.Int(4), 1, 0, 0.05, 0.08)
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
