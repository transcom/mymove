package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
)

// RunScenarioTwo creates 9 shipments to be divided between 5 TSPs in 1 TDL and 10 shipments to be divided among 4 TSPs in TDL 2.
// This allows testing against award queue to ensure it behaves as expected. Two TSPs in TDL1 and one TSP in TDL 2 have blackout dates.
func RunScenarioTwo(db *pop.Connection) {
	shipmentsToMake := 9
	shipmentDate := time.Now()

	// Make a TDL to contain our tests
	tdl, _ := MakeTDL(db, "california", "90210", "2")
	tdl2, _ := MakeTDL(db, "New York", "10024", "2")

	// Make a market
	market := "dHHG"

	// Make a source GBLOC
	sourceGBLOC := "OHAI"

	// Make shipments in first TDL
	for i := 0; i < shipmentsToMake; i++ {
		MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl, sourceGBLOC, market)
	}
	// Make shipments in second TDL
	for i := 0; i <= shipmentsToMake; i++ {
		MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl2, sourceGBLOC, market)
	}

	// Make TSPs
	tsp1, _ := MakeTSP(db, "Excellent TSP with Blackout Date", "TSP1")
	tsp2, _ := MakeTSP(db, "Very Good TSP", "TSP2")
	tsp3, _ := MakeTSP(db, "Pretty Good TSP", "TSP3")
	tsp4, _ := MakeTSP(db, "OK TSP with Blackout Date", "TSP4")
	tsp5, _ := MakeTSP(db, "Are you even trying TSP", "TSP5")
	tsp6, _ := MakeTSP(db, "Excellent TSP", "TSP6")
	tsp7, _ := MakeTSP(db, "Pretty Good TSP with Blackout Date", "TSP7")
	tsp8, _ := MakeTSP(db, "OK TSP", "TSP8")
	tsp9, _ := MakeTSP(db, "Going out of business TSP", "TSP9")

	// Put TSPs in 2 TDLs to handle these shipments
	MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0)
	MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0)
	MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0)
	MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0)
	MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0)

	MakeTSPPerformance(db, tsp6, tdl2, swag.Int(1), 5, 0)
	MakeTSPPerformance(db, tsp7, tdl2, swag.Int(2), 4, 0)
	MakeTSPPerformance(db, tsp8, tdl2, swag.Int(3), 2, 0)
	MakeTSPPerformance(db, tsp9, tdl2, swag.Int(4), 1, 0)
	// Add blackout dates
	blackoutStart := shipmentDate.AddDate(0, 0, -3)
	blackoutEnd := shipmentDate.AddDate(0, 0, 3)

	gbloc := "BKAS"
	MakeBlackoutDate(db,
		tsp1,
		blackoutStart,
		blackoutEnd,
		&tdl,
		&gbloc,
		&market)
	MakeBlackoutDate(db,
		tsp4,
		blackoutStart,
		blackoutEnd,
		&tdl,
		&gbloc,
		&market)
	MakeBlackoutDate(db,
		tsp7,
		blackoutStart,
		blackoutEnd,
		&tdl,
		&gbloc,
		&market)
}
