package main

import (
	"fmt"
	"log"
	"os"
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

// This executable is used to demonstrate the rate engine and as a diagnostic tool to
// easily see what values it is computing for a known scenario.
//
// Run using go run cmd/demo/rateengine.go -scenario=n, where n is either 1 or 2
func main() {
	scenarioNumber := flag.Int("scenario", 1, "Specify which scenario you'd like to run. Current options: 1, 2.")
	flag.Parse()

	err := pop.AddLookupPaths("config")
	if err != nil {
		log.Fatalf("failed to add config to pop paths: %+v", err)
	}
	db, err := pop.Connect("development")
	if err != nil {
		log.Fatalf("could not connect to database: %+v", err)
	}
	err = db.TruncateAll()
	if err != nil {
		log.Fatalf("could not truncate the database: %+v", err)
	}

	var input string
	for {
		fmt.Println("Running this tool will delete everything in your development database.")
		fmt.Print("Do you wish to proceed? (y)es or (n)o: ")
		count, err := fmt.Scanln(&input)
		if err != nil || count == 0 || input == "n" || input == "no" {
			os.Exit(1)
		} else if input == "y" || input == "yes" {
			break
		}
		fmt.Println("")
	}

	logger := zap.NewNop()

	if *scenarioNumber == 1 {
		fmt.Println("Running scenario 1")
		fmt.Println("Origin: New Smyrna Beach, FL, 32168")
		fmt.Println("Destination: Awendaw, SC, 29429")
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
		lhDiscount := unit.DiscountRate(0.67)

		cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, 0, lhDiscount, 0)
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
		fmt.Printf("%-30s%s\n", "Full Pack fee:", cost.PackFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Full Unpack fee:", cost.UnpackFee.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Government Constructed Cost:", cost.GCC.ToDollarString())

	} else if *scenarioNumber == 2 {
		fmt.Println("Running scenario 2")
		fmt.Println("Origin: Hayward, CA, 94540")
		fmt.Println("Destination: Georgetown, TX, 78626")
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
		lhDiscount := unit.DiscountRate(0.67)

		cost, err := engine.ComputePPM(weight, originZip5, destinationZip5, date, 0, lhDiscount, 0)
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
		fmt.Printf("%-30s%s\n", "Full Pack fee:", cost.PackFee.ToDollarString())
		fmt.Printf("%-30s%s\n", "Full Unpack fee:", cost.UnpackFee.ToDollarString())
		fmt.Println("")
		fmt.Printf("%-30s%s\n", "Government Constructed Cost:", cost.GCC.ToDollarString())
	}
}
