package ghcimport

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/gocarina/gocsv"

	"github.com/transcom/mymove/pkg/models"
)

//Zip3Fixture stores zip and rate and area information from tariff400ng_zip3s_fixture.csv
type Zip3Fixture struct {
	Zip      string `csv:"zip"`
	RateArea string `csv:"rate_area"`
}

//Zip5Fixture stores zip and rate and area information from tariff400ng_zip5_rate_areas_fixture.csv
type Zip5Fixture struct {
	Zip      string `csv:"zip"`
	RateArea string `csv:"rate_area"`
}

func (gre *GHCRateEngineImporter) mapZipsToRateAreas(dbTx *pop.Connection) error {
	// Maps re_zip3s to re_rate_areas based on tariff400ng_zip3s_fixture.csv
	err := gre.mapZip3s(dbTx)
	if err != nil {
		return fmt.Errorf("mapZip3s failed: %w", err)
	}

	// Creates re_zip5_rate_areas records from tariff400ng_zip5_rate_areas_fixture.csv
	err = gre.createZip5s(dbTx)
	if err != nil {
		return fmt.Errorf("createZip5s failed: %w", err)
	}

	return nil
}

func (gre *GHCRateEngineImporter) mapZip3s(dbTx *pop.Connection) error {
	csvFile, err := os.OpenFile("pkg/services/ghcimport/fixtures/tariff400ng_zip3s_fixture.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer csvFile.Close()

	zip3RateAreas := []*Zip3Fixture{}

	err = gocsv.UnmarshalFile(csvFile, &zip3RateAreas)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}

	for _, zip3RateArea := range zip3RateAreas {
		var reZip3 models.ReZip3
		err = dbTx.Where("zip3 = ?", zip3RateArea.Zip).First(&reZip3)
		if err != nil {
			return fmt.Errorf("failed to find ReZip3 record with zip %s: %w", zip3RateArea.Zip, err)
		}

		if zip3RateArea.RateArea == "ZIP" {
			reZip3.RateAreaID = nil
			reZip3.HasMultipleRateAreas = true
			verrs, err := dbTx.ValidateAndUpdate(&reZip3)
			if err != nil {
				return fmt.Errorf("failed to update %v: %v", reZip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate %v: %v", reZip3, verrs)
			}
		} else {
			rateAreaID, found := gre.domesticRateAreaToIDMap[zip3RateArea.RateArea]
			if !found {
				return fmt.Errorf("failed to map %s rate area to ID", zip3RateArea.RateArea)
			}
			reZip3.RateAreaID = &rateAreaID
			reZip3.HasMultipleRateAreas = false
			verrs, err := dbTx.ValidateAndUpdate(&reZip3)
			if err != nil {
				return fmt.Errorf("failed to update %v: %v", reZip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate %v: %v", reZip3, verrs)
			}
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) createZip5s(db *pop.Connection) error {
	//TODO 1: Load /fixtures/tariff400ng_zip5s_rate_areas_fixture.csv
	//TODO 2: Iterate over each row in CSV fixture
	//TODO 3: Store zip5 and rate_area value from each row of the CSV
	//TODO 4: Create a new record in re_zip5s_rate_areas table for each row of the CSV fixture
	//TODO 5: Store zip5 value in record
	//TODO 6: Find the corresponding re_rate_areas record and associate it with the new re_zip5s record

	return nil
}