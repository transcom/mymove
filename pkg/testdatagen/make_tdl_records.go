package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTDL finds or makes a single traffic_distribution_list record
func MakeTDL(db *pop.Connection, assertions Assertions) models.TrafficDistributionList {

	source := assertions.TrafficDistributionList.SourceRateArea
	if source == "" {
		source = DefaultSrcRateArea
	}
	dest := assertions.TrafficDistributionList.DestinationRegion
	if dest == "" {
		dest = DefaultDstRegion
	}
	cos := assertions.TrafficDistributionList.CodeOfService
	if cos == "" {
		cos = DefaultCOS
	}

	tdls := []models.TrafficDistributionList{}
	query := db.Where(fmt.Sprintf("source_rate_area = '%s' AND destination_region = '%s' AND code_of_service = '%s'", source, dest, cos))
	err := query.All(&tdls)
	if err != nil {
		log.Panic(err)
	}

	// Create a new record if none are found
	if len(tdls) == 0 {
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
		return tdl
	}
	return tdls[0]
}

// MakeDefaultTDL makes a TDL with default values
func MakeDefaultTDL(db *pop.Connection) models.TrafficDistributionList {
	return MakeTDL(db, Assertions{})
}

// MakeTDLData creates three TDL records
func MakeTDLData(db *pop.Connection) {
	// It would be nice to make this less repetitive
	MakeTDL(db, Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US1",
			DestinationRegion: "2",
			CodeOfService:     "2",
		},
	})
	MakeTDL(db, Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US10",
			DestinationRegion: "9",
			CodeOfService:     "2",
		},
	})
	MakeTDL(db, Assertions{
		TrafficDistributionList: models.TrafficDistributionList{
			SourceRateArea:    "US4964400",
			DestinationRegion: "4",
			CodeOfService:     "D",
		},
	})
}
