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
	query := db.Where("source_rate_area = ?", source).
		Where("destination_region = ?", dest).
		Where("code_of_service = ?", cos)
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

		verrs, err := db.ValidateAndCreate(&tdl)
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
