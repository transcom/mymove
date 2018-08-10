package main

import (
	"log"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"github.com/transcom/mymove/pkg/testdatagen"
	tdgs "github.com/transcom/mymove/pkg/testdatagen/scenario"
)

// Hey, refactoring self: you can pull the UUIDs from the objects rather than
// querying the db for them again.
func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	rounds := flag.String("rounds", "none", "If not using premade scenarios: Specify none (no awards), full (1 full round of awards), or half (partial round of awards)")
	numTSP := flag.Int("numTSP", 15, "If not using premade scenarios: Specify the number of TSPs you'd like to create")
	scenario := flag.Int("scenario", 0, "Specify which scenario you'd like to run. Current options: 1, 2, 3, 4, 5, 6, 7.")
	namedScenario := flag.String("named-scenario", "", "It's like a scenario, but more descriptive.")
	flag.Parse()

	//DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		log.Panic(err)
	}
	db, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	if *scenario == 1 {
		tdgs.RunAwardQueueScenario1(db)
	} else if *scenario == 2 {
		tdgs.RunAwardQueueScenario2(db)
	} else if *scenario == 4 {
		err = tdgs.RunPPMSITEstimateScenario1(db)
	} else if *scenario == 5 {
		err = tdgs.RunRateEngineScenario1(db)
	} else if *scenario == 6 {
		query := `DELETE FROM transportation_service_provider_performances;
				  DELETE FROM transportation_service_providers;
				  DELETE FROM traffic_distribution_lists;
				  DELETE FROM tariff400ng_zip3s;
				  DELETE FROM tariff400ng_zip5_rate_areas;
				  DELETE FROM tariff400ng_service_areas;
				  DELETE FROM tariff400ng_linehaul_rates;
				  DELETE FROM tariff400ng_shorthaul_rates;
				  DELETE FROM tariff400ng_full_pack_rates;
				  DELETE FROM tariff400ng_full_unpack_rates;`

		err = db.RawQuery(query).Exec()
		if err != nil {
			log.Panic(err)
		}
		err = tdgs.RunRateEngineScenario2(db)
	} else if *scenario == 7 {
		// Create TSPs with shipments divided among them
		numTspUsers := 2
		numShipments := 25
		numShipmentOfferSplit := []int{15, 10}
		status = []string{"DEFAULT", "AWARDED"}
		_, _, _, err := testdatagen.CreateShipmentOfferData(db, numTspUsers, numShipments, numShipmentOfferSplit, status)
		if err != nil {
			log.Panic(err)
		}
	} else if *namedScenario == tdgs.E2eBasicScenario.Name {
		tdgs.E2eBasicScenario.Run(db)
		log.Print("Success! Created e2e test data.")
	} else {
		// Can this be less repetitive without being overly clever?
		testdatagen.MakeTDLData(db)
		testdatagen.MakeTSPs(db, *numTSP)
		testdatagen.MakeShipmentData(db)
		testdatagen.MakeShipmentOfferData(db)
		testdatagen.MakeTSPPerformanceData(db, *rounds)
		testdatagen.MakeBlackoutDateData(db)
		testdatagen.MakePPMData(db)
		testdatagen.MakeReimbursementData(db)
		testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{})
	}
	if err != nil {
		log.Panic(err)
	}
}
