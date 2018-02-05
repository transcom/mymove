package main

import (
	"fmt"
	"log"

	"github.com/markbates/pop"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags

	"github.com/transcom/mymove/models"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presenetly.")
	flag.Parse()

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	tdl := models.TrafficDistributionList{
		SourceRateArea:    "california",
		DestinationRegion: "90210",
		CodeOfService:     "2"}

	_, err = dbConnection.ValidateAndSave(&tdl)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Hello, TSP Award Queue!")
}
