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
		MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl, sourceGBLOC, &market)
	}
	// Make shipments in second TDL
	for i := 0; i <= shipmentsToMake; i++ {
		MakeShipment(db, shipmentDate, shipmentDate, shipmentDate, tdl2, sourceGBLOC, &market)
	}

	// Make TSPs
	tsp1, _ := MakeTSP(db, "Excellent TSP with Blackout Date", RandomSCAC())
	tsp2, _ := MakeTSP(db, "Very Good TSP", RandomSCAC())
	tsp3, _ := MakeTSP(db, "Pretty Good TSP", RandomSCAC())
	tsp4, _ := MakeTSP(db, "OK TSP with Blackout Date", RandomSCAC())
	tsp5, _ := MakeTSP(db, "Are you even trying TSP", RandomSCAC())
	tsp6, _ := MakeTSP(db, "Excellent TSP", RandomSCAC())
	tsp7, _ := MakeTSP(db, "Pretty Good TSP with Blackout Date", RandomSCAC())
	tsp8, _ := MakeTSP(db, "OK TSP", RandomSCAC())
	tsp9, _ := MakeTSP(db, "Going out of business TSP", RandomSCAC())

	// Put TSPs in 2 TDLs to handle these shipments
	MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), 5, 0, 4.2, 4.4)
	MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), 4, 0, 3.1, 3.2)
	MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), 3, 0, 2.4, 2.5)
	MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), 2, 0, 1.1, 1.3)
	MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), 1, 0, .5, .8)

	MakeTSPPerformance(db, tsp6, tdl2, swag.Int(1), 5, 0, 4.2, 4.4)
	MakeTSPPerformance(db, tsp7, tdl2, swag.Int(2), 4, 0, 3.1, 3.2)
	MakeTSPPerformance(db, tsp8, tdl2, swag.Int(3), 2, 0, 1.1, 1.3)
	MakeTSPPerformance(db, tsp9, tdl2, swag.Int(4), 1, 0, .5, .8)
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
