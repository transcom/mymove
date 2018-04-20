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
		log.Fatalf("count not connect to database: %+v", err)
	}
	db.TruncateAll()

	logger := zap.NewNop()
	planner := route.NewTestingPlanner(362)
	engine := rateengine.NewRateEngine(db, logger, planner)

	fmt.Printf("Running scenario %d: New Smyrna Beach, FL -> Awendaw, SC\n\n", *scenarioNumber)
	if *scenarioNumber == 1 {
		scenario.RunRateEngineScenario1(db)

		weight := unit.Pound(4000)
		originZip5 := "32168"
		destinationZip5 := "29429"
		date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
		inverseDiscount := 0.33

		cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
		if err != nil {
			log.Fatalf("cound not compute PPM: %+v", err)
		}

		fmt.Printf("%-30s%s\n", "Base linehaul (non-disc'd):", cost.BaseLinehaul.ToDollarString())
		fmt.Printf("%-30s%s\n", "Origin linehaul factor:", cost.OriginLinehaulFactor.ToDollarString())
		fmt.Printf("%-30s%s\n", "Destination linehaul factor:", cost.DestinationLinehaulFactor.ToDollarString())
		fmt.Printf("%-30s%s\n", "Shorthaul chargeg:", cost.ShorthaulCharge.ToDollarString())
		fmt.Printf("%-30s%s\n", "Linehaul total (trans. cost):", cost.LinehaulChargeTotal.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Origin service fee:", cost.OriginServiceFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Destination service fee:", cost.DestinationServiceFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Full Pack charge:", cost.PackFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Full Unpack charge:", cost.UnpackFee.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Government Constructive Cost:", cost.GCC.ToDollarString())

	} else if *scenarioNumber == 2 {
		scenario.RunRateEngineScenario2(db)

	}
}
