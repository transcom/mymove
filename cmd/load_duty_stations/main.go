package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"

	"github.com/transcom/mymove/internal/pkg/dutystationsloader"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	validate := flag.Bool("validate", false, "Only run file validations")
	output := flag.String("output", "", "Where to output the migration file")
	stationsPath := flag.String("stations", "", "Input file for duty stations")
	officesPath := flag.String("offices", "", "Input file for transportation offices")
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

	// Parse duty stations
	stationRows, err := dutystationsloader.ParseStations(*stationsPath)
	if err != nil {
		log.Panic(err)
	}

	// Parse transportation offices
	officeRows, err := dutystationsloader.ParseOffices(*officesPath)
	if err != nil {
		log.Panic(err)
	}

	// Searches the database for existing stations/offices
	stationDupes, officeDupes, err := dutystationsloader.CheckDatabaseForDuplicates(db, stationRows, officeRows)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Found %v duty station duplicates!\n", len(stationDupes))
	for _, d := range stationDupes {
		fmt.Println(d.Name)
	}
	fmt.Printf("Found %v transportation office duplicates!\n", len(officeDupes))
	for _, d := range officeDupes {
		fmt.Println(d.Name)
	}

	// Attempts to pair transportation offices with duty stations
	pairs, unpairedStations, unpairedOffices := dutystationsloader.PairStationsAndOffices(stationRows, officeRows)
	fmt.Printf("Found %v pairs!\n", len(pairs))
	fmt.Printf("%v duty stations were left unpaired\n", len(unpairedStations))
	for _, s := range unpairedStations {
		fmt.Printf("Station: %v, Office name: %v", s.Name, s.TransportationOffice)
	}
	fmt.Printf("%v transportation offices were left unpaired\n", len(unpairedOffices))
	for _, s := range unpairedOffices {
		fmt.Printf("Office: %v", s.Name)
	}

	// If we just want to validate files we can exit
	if validate != nil && *validate {
		os.Exit(0)
	}

	var migration strings.Builder
	migration.WriteString("-- Migration generated using cmd/load_duty_stations\n")
	migration.WriteString(fmt.Sprintf("-- Duty stations file: %v\n", *stationsPath))
	migration.WriteString(fmt.Sprintf("-- Transportation offices file: %v\n", *officesPath))
	migration.WriteString("\n")
	migration.WriteString(dutystationsloader.GenerateMigrationString(pairs))

	f, err := os.OpenFile(*output, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	_, err = f.WriteString(migration.String())
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Complete! Migration written to %v\n", *output)
}
