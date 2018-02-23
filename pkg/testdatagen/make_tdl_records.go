package testdatagen

import (
	"log"

	"github.com/markbates/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTDL makes a single traffic_distribution_list record
func MakeTDL(db *pop.Connection, source string, dest string, cos string) (models.TrafficDistributionList, error) {

	tdl := models.TrafficDistributionList{
		SourceRateArea:    source,
		DestinationRegion: dest,
		CodeOfService:     cos,
	}

	_, err := db.ValidateAndSave(&tdl)
	if err != nil {
		log.Panic(err)
	}

	return tdl, err
}

// MakeTDLData creates three TDL records
func MakeTDLData(db *pop.Connection) {
	// It would be nice to make this less repetitive
	MakeTDL(db, "california", "90210", "2")
	MakeTDL(db, "north carolina", "27007", "4")
	MakeTDL(db, "washington", "98310", "1")
}
