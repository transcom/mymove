package main

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
	// "github.com/satori/go.uuid"
	"github.com/namsral/flag"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// Add three TDL table records
	tdl1 := models.TrafficDistributionList{
		SourceRateArea:    "california",
		DestinationRegion: "90210",
		CodeOfService:     "2"}

	tdl2 := models.TrafficDistributionList{
		SourceRateArea:    "north carolina",
		DestinationRegion: "27007",
		CodeOfService:     "4"}

	tdl3 := models.TrafficDistributionList{
		SourceRateArea:    "washington",
		DestinationRegion: "98310",
		CodeOfService:     "1"}

	_, err = dbConnection.ValidateAndSave(&tdl1)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&tdl2)
	if err != nil {
		log.Panic(err)
	}

	_, err = dbConnection.ValidateAndSave(&tdl3)
	if err != nil {
		log.Panic(err)
	}

	// Query for newly made records and print IDs in terminal
	tdls := []models.TrafficDistributionList{}
	err = dbConnection.All(&tdls)
	if err != nil {
		fmt.Print("Error!\n")
		fmt.Printf("%v\n", err)
	} else {
		for _, v := range tdls {
			fmt.Print(v.ID)
		}
	}
}
