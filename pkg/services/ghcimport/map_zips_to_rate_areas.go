package ghcimport

import (
	"fmt"
	"github.com/transcom/mymove/pkg/models"

	"github.com/gobuffalo/pop"
	//"github.com/gobuffalo/validate"
	//"github.com/gofrs/uuid"
	//"github.com/pkg/errors"
	//
	//"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) mapZipsToRateAreas(db *pop.Connection) error {
	//fmt.Println("Test mapZipsToRateAreas function")
	var err error
	err = gre.mapZip3s(dbTx)
	if err != nil {
		return fmt.Errorf("mapZipsToRateAreas failed to map: %w", err)
	}

	err = gre.mapZip5s(dbTx)
	if err != nil {
		return fmt.Errorf("mapZipsToRateAreas failed to map: %w", err)
	}
	return nil
}

func (gre *GHCRateEngineImporter) mapZip3s(db *pop.Connection) error {

	gre.domesticRateAreaToIDMap


	// Load /fixtures/tariff400ng_zip3s_fixture.csv and for each row
		// var zip3 = fixture zip3
		// var rate_area = fixture rate_area

		// Fetch the existing zip3 record for the current zip3 in this CSV line
		var reZip3 models.ReZip3
		err := db.Where("zip3 = $1", zip3FromCsv).First(&reZip3)
		// TODO: Update rate_area_id and has_multiple_rate_areas
		reZip3.RateAreaID = // look up the rate area id in the domesticRateAreaToIDMap for the rate area in this CSV line
		reZip3.HasMultipleRateAreas = // true or false depending on whether "ZIP" is in there
		verrs, err := db.ValidateAndUpdate(&reZip3)


	// Load re_zip3s and find record where re_zip3s zip3 == fixture zip3
			// if fixture rate_area == 'ZIP' set re_zip3s has_multiple_rate_areas to true
			// else
				// load re_rate_areas and find re_rate_areas code == fixture rate_area
				// set re_zip3s rate_area_id to re_rate_areas id
				// set re_zip3s has_multiple_rate_areas to false
	return nil
}

func (gre *GHCRateEngineImporter) mapZip5s(db *pop.Connection) error {
	// Load /fixtures/tariff400ng_zip5s_rate_areas_fixture.csv and for each row
	// var zip5 = fixture zip5
	// var rate_area = fixture rate_area

	// Create a new record in re_zip5_rate_areas where re_zip5_rate_areas zip5 == fixture zip5
	// load re_rate_areas and find re_rate_areas code == fixture rate_area
	// set re_zip5_rate_areas rate_area_id to re_rate_areas id
	return nil
}