package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTDL makes a single traffic_distribution_list record
func MakeTDL(db *pop.Connection, source string, dest string, cos string) (models.TrafficDistributionList, error) {

	tdl := models.TrafficDistributionList{
		SourceRateArea:    source,
		DestinationRegion: dest,
		CodeOfService:     cos,
	}

	verrs, err := db.ValidateAndSave(&tdl)
	if verrs.HasAny() {
		err = fmt.Errorf("TDL validation errors: %v", verrs)
	}
	if err != nil {
		log.Panic(err)
	}

	return tdl, err
}

// MakeTDLData creates three TDL records
func MakeTDLData(db *pop.Connection) {
	// It would be nice to make this less repetitive
	MakeTDL(db, "US1", "2", "2")
	MakeTDL(db, "US10", "9", "2")
	MakeTDL(db, "US4964400", "4", "D")
}
