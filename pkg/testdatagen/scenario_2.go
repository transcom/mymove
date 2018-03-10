package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/markbates/pop"
)

// RunScenarioTwo creates 17 shipments and 5 TSPs in 1 TDL. This allows testing against
// award queue to ensure it behaves as expected. This doesn't track blackout dates.
func RunScenarioTwo(db *pop.Connection) {
	shipmentsToMake := 9

	// Make a TDL to contain our tests
	tdl, _ := MakeTDL(db, "california", "90210", "2")
	tdl2, _ := MakeTDL(db, "New York", "10024", "2")

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		MakeShipment(db, time.Now(), time.Now(), tdl)
	}

	// Make TSPs in the same TDL to handle these shipments
	tsp1, _ := MakeTSP(db, "Test TSP Quality Band 1", "TSP1")
	tsp2, _ := MakeTSP(db, "Test TSP Quality Band 1", "TSP2")
	tsp3, _ := MakeTSP(db, "Test TSP Quality Band 2", "TSP3")
	tsp4, _ := MakeTSP(db, "Test TSP Quality Band 3", "TSP4")
	tsp5, _ := MakeTSP(db, "Test TSP Quality Band 4", "TSP5")
	tsp6, _ := MakeTSP(db, "Test TSP Quality Band 1", "TSP6")
	tsp7, _ := MakeTSP(db, "Test TSP Quality Band 2", "TSP7")
	tsp8, _ := MakeTSP(db, "Test TSP Quality Band 3", "TSP8")
	tsp9, _ := MakeTSP(db, "Test TSP Quality Band 4", "TSP9")

	// TSPs should be orderd by award_count first, then BVS.
	MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0)
	MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0)
	MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0)
	MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0)
	MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0)
	MakeTSPPerformance(db, tsp6, tdl2, swag.Int(1), 5, 0)
	MakeTSPPerformance(db, tsp7, tdl2, swag.Int(2), 4, 0)
	MakeTSPPerformance(db, tsp8, tdl2, swag.Int(3), 2, 0)
	MakeTSPPerformance(db, tsp9, tdl2, swag.Int(4), 1, 0)
}
