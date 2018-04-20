package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"github.com/pkg/errors"
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

	if *scenarioNumber == 1 {
		fmt.Println("Running scenario 1: New Smyrna Beach, FL -> Awendaw, SC")
		fmt.Println("")

		if err := scenario.RunRateEngineScenario1(db); err != nil {
			log.Fatalf("failed to run scenario 1.")
		}

		planner := route.NewTestingPlanner(362)
		engine := rateengine.NewRateEngine(db, logger, planner)

		weight := unit.Pound(4000)
		originZip5 := "32168"
		destinationZip5 := "29429"
		date := time.Date(2018, time.June, 18, 0, 0, 0, 0, time.UTC)
		inverseDiscount := 0.33

		cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
		if err != nil {
			log.Fatalf("could not compute PPM: %+v", err)
		}
		fmt.Printf("%-30s%s\n", "Base linehaul (non-disc'd):", cost.BaseLinehaul.ToDollarString())
		fmt.Printf("%-30s%s\n", "Origin linehaul factor:", cost.OriginLinehaulFactor.ToDollarString())
		fmt.Printf("%-30s%s\n", "Destination linehaul factor:", cost.DestinationLinehaulFactor.ToDollarString())
		fmt.Printf("%-30s%s\n", "Shorthaul charge:", cost.ShorthaulCharge.ToDollarString())
		fmt.Printf("%-30s%s\n", "Linehaul total (trans. cost):", cost.LinehaulChargeTotal.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Origin service fee:", cost.OriginServiceFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Destination service fee:", cost.DestinationServiceFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Full Pack/Unpack fee:", cost.FullPackUnpackFee.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Government Constructive Cost:", cost.GCC.ToDollarString())

	} else if *scenarioNumber == 2 {
		fmt.Println("Running scenario 2: Hayward, CA -> Georgetown, TX")
		fmt.Println("")

		if err := scenario.RunRateEngineScenario2(db); err != nil {
			log.Fatalf("failed to run scenario 2.")
		}

		planner := route.NewTestingPlanner(1693)
		engine := rateengine.NewRateEngine(db, logger, planner)

		weight := unit.Pound(7500)
		originZip5 := "94540"
		destinationZip5 := "78626"
		date := time.Date(2018, time.December, 5, 0, 0, 0, 0, time.UTC)
		inverseDiscount := 0.33

		cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, inverseDiscount)
		if err != nil {
			log.Fatalf("could not compute PPM: %+v", errors.Cause(err))
		}
		fmt.Printf("%-30s%s\n", "Base linehaul (non-disc'd):", cost.BaseLinehaul.ToDollarString())
		fmt.Printf("%-30s%s\n", "Origin linehaul factor:", cost.OriginLinehaulFactor.ToDollarString())
		fmt.Printf("%-30s%s\n", "Destination linehaul factor:", cost.DestinationLinehaulFactor.ToDollarString())
		fmt.Printf("%-30s%s\n", "Shorthaul charge:", cost.ShorthaulCharge.ToDollarString())
		fmt.Printf("%-30s%s\n", "Linehaul total (trans. cost):", cost.LinehaulChargeTotal.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Origin service fee:", cost.OriginServiceFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Destination service fee:", cost.DestinationServiceFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Full Pack/Unpack fee:", cost.FullPackUnpackFee.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Government Constructive Cost:", cost.GCC.ToDollarString())
	}
}
