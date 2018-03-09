package main

import (
	"log"
	"time"

	"github.com/go-openapi/swag"
	"github.com/markbates/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/testdatagen"
)

const mps int = 15

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	rounds := flag.String("rounds", "none", "Choose none (no awards), full (1 full round of awards), or half (partial round of awards)")
	numTSP := flag.Int("numTSP", 15, "The number of TSPs you'd like to create")
	fullRound := flag.Bool("fullRound", false, "If you want to see how a full round of shipments are awarded.")
	halfRound := flag.Bool("halfRound", false, "If you want to see how half a round of shipments are awarded.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	if *fullRound == true {
		shipmentsToMake := 17

		// Make a TDL to contain our tests
		tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")

		// Make shipments in this TDL
		for i := 0; i < shipmentsToMake; i++ {
			testdatagen.MakeShipment(db, time.Now(), time.Now(), tdl)
		}

		// Make TSPs in the same TDL to handle these shipments
		tsp1, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 1", "TSP1")
		tsp2, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 1", "TSP2")
		tsp3, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 2", "TSP3")
		tsp4, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 3", "TSP4")
		tsp5, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 4", "TSP5")

		// TSPs should be orderd by award_count first, then BVS.
		testdatagen.MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), mps+5, 0)
		testdatagen.MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), mps+4, 0)
		testdatagen.MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), mps+3, 0)
		testdatagen.MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), mps+2, 0)
		testdatagen.MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), mps+1, 0)
	} else if *halfRound == true {
		shipmentsToMake := 9

		// Make a TDL to contain our tests
		tdl, _ := testdatagen.MakeTDL(db, "california", "90210", "2")
		tdl2, _ := testdatagen.MakeTDL(db, "New York", "10024", "2")

		// Make shipments in this TDL
		for i := 0; i < shipmentsToMake; i++ {
			testdatagen.MakeShipment(db, time.Now(), time.Now(), tdl)
		}

		// Make TSPs in the same TDL to handle these shipments
		tsp1, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 1", "TSP1")
		tsp2, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 1", "TSP2")
		tsp3, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 2", "TSP3")
		tsp4, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 3", "TSP4")
		tsp5, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 4", "TSP5")
		tsp6, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 1", "TSP6")
		tsp7, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 2", "TSP7")
		tsp8, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 3", "TSP8")
		tsp9, _ := testdatagen.MakeTSP(db, "Test TSP Quality Band 4", "TSP9")

		// TSPs should be orderd by award_count first, then BVS.
		testdatagen.MakeTSPPerformance(db, tsp1, tdl, swag.Int(1), mps+5, 0)
		testdatagen.MakeTSPPerformance(db, tsp2, tdl, swag.Int(1), mps+4, 0)
		testdatagen.MakeTSPPerformance(db, tsp3, tdl, swag.Int(2), mps+3, 0)
		testdatagen.MakeTSPPerformance(db, tsp4, tdl, swag.Int(3), mps+2, 0)
		testdatagen.MakeTSPPerformance(db, tsp5, tdl, swag.Int(4), mps+1, 0)
		testdatagen.MakeTSPPerformance(db, tsp6, tdl2, swag.Int(1), mps+5, 0)
		testdatagen.MakeTSPPerformance(db, tsp7, tdl2, swag.Int(2), mps+4, 0)
		testdatagen.MakeTSPPerformance(db, tsp8, tdl2, swag.Int(3), mps+2, 0)
		testdatagen.MakeTSPPerformance(db, tsp9, tdl2, swag.Int(4), mps+1, 0)
	} else {
		// Can this be less repetitive without being overly clever?
		testdatagen.MakeTDLData(db)
		testdatagen.MakeTSPs(db, *numTSP)
		testdatagen.MakeShipmentData(db)
		testdatagen.MakeShipmentData(db)
		testdatagen.MakeShipmentData(db)
		testdatagen.MakeShipmentAwardData(db)
		testdatagen.MakeTSPPerformanceData(db, *rounds)
		testdatagen.MakeBlackoutDateData(db)
	}
}
