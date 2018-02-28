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
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// Can this be less repetitive without being overly clever?
	testdatagen.MakeTDLData(db)
	testdatagen.MakeTSPData(db)
	testdatagen.MakeShipmentData(db)
	testdatagen.MakeShipmentAwardData(db)
	testdatagen.MakeTSPPerformanceData(db)
}
