package main

import (
	"log"

	"github.com/markbates/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	rounds := flag.String("rounds", "none", "Choose none (no awards), full (1 full round of awards), or half (partial round of awards)")
	numTSP := flag.Int("numTSP", 15, "The number of TSPs you'd like to create")
	scenarioOne := flag.Bool("scenarioOne", false, "To run test data generation scenario 1.")
	scenarioTwo := flag.Bool("scenarioTwo", false, "To run test data generation scenario 2.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	if *scenarioOne == true {
		testdatagen.RunScenarioOne(db)
	} else if *scenarioTwo == true {
		testdatagen.RunScenarioTwo(db)
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
