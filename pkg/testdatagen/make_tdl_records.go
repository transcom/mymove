package testdatagen

import (
	"fmt"
	"github.com/markbates/pop"
	"github.com/transcom/mymove/pkg/models"
	"log"
)

// MakeTDLData creates three TDL records
func MakeTDLData(dbConnection *pop.Connection) {
	// It would be nice to make this less repetitive
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

	_, err := dbConnection.ValidateAndSave(&tdl1)
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
}
