package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
)

// RunScenarioThree creates 5? shipments and 5? TSPs in 3? TDLs. This allows
// testing against award queue to ensure it behaves as expected. This
// doesn't track blackout dates.
func RunScenarioThree(db *pop.Connection) {
	shipmentsToMake := 17

	// Make a TDL to contain our tests
	tdl, _ := MakeTDL(db, "california", "90210", "2")

	// Make shipments in this TDL
	for i := 0; i < shipmentsToMake; i++ {
		MakeShipment(db, time.Now(), time.Now(), time.Now(), tdl)
	}

	// Make TSPs in two TDLs to handle these shipments
	tspA1, _ := MakeTSP(db, "Excellent TSP", "TSP1")
	tspA2, _ := MakeTSP(db, "Pretty Good TSP", "TSP2")
	tspB1, _ := MakeTSP(db, "Good TSP", "TSP3")
	tspB2, _ := MakeTSP(db, "OK TSP", "TSP4")
	tspB3, _ := MakeTSP(db, "Bad TSP", "TSP5")

	// TSPs should be orderd by offer_count first, then BVS.
	MakeTSPPerformance(db, tspA1, tdl, swag.Int(1), 5, 0)
	MakeTSPPerformance(db, tspA2, tdl, swag.Int(1), 4, 0)
	MakeTSPPerformance(db, tspB1, tdl, swag.Int(2), 3, 0)
	MakeTSPPerformance(db, tspB2, tdl, swag.Int(3), 2, 0)
	MakeTSPPerformance(db, tspB3, tdl, swag.Int(4), 1, 0)
}
