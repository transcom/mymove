package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
)

func main() {
	scenarioNumber := flag.Int("scenario", 1, "Specify which scenario you'd like to run. Current options: 1, 2.")
	flag.Parse()

	pop.AddLookupPaths("config")
	db, err := pop.Connect("test")
	if err != nil {
		log.Fatalf("could not connect to database: %+v", err)
	}
	db.TruncateAll()

	logger := zap.NewNop()
	planner := route.NewTestingPlanner(362)
	engine := rateengine.NewRateEngine(db, logger, planner)

	fmt.Printf("Running scenario %d: New Smyrna Beach, FL -> Awendaw, SC\n", *scenarioNumber)
	if *scenarioNumber == 1 {
		if err := scenario.RunRateEngineScenario1(db); err != nil {
			log.Fatalf("failed to run scenario 1.")
		}

		weight := unit.Pound(4000)
		originZip5 := "32168"
		destinationZip5 := "29429"
		date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
		inverseDiscount := 0.33

		cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
		if err != nil {
			log.Fatalf("could not compute PPM: %+v", err)
		}
		fmt.Printf("%+v", cost)

	} else if *scenarioNumber == 2 {
		scenario.RunRateEngineScenario2(db)

	}
}
