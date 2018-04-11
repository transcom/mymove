package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
)

// RunScenarioOne creates 17 shipments and 5 TSPs in 1 TDL. This allows testing against
// award queue to ensure it behaves as expected. This doesn't track blackout dates.
func RunScenarioOne(db *pop.Connection) {
	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl, _ := MakeTDL(db, "california", "90210", "2")

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		MakeShipment(db, time.Now(), time.Now(), time.Now(), tdl, sourceGBLOC, &market)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := MakeTSP(db, "Excellent TSP", RandomSCAC())
	tsp2, _ := MakeTSP(db, "Pretty Good TSP", RandomSCAC())
	tsp3, _ := MakeTSP(db, "Good TSP", RandomSCAC())
	tsp4, _ := MakeTSP(db, "OK TSP", RandomSCAC())
	tsp5, _ := MakeTSP(db, "Bad TSP", RandomSCAC())

	// TSPs should be orderd by offer_count first, then BVS.
	MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0)
	MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0)
	MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0)
	MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0)
	MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0)
}
