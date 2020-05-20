package ghcimport

import (
	"fmt"
	"os"

	"github.com/gobuffalo/pop"
	"github.com/gocarina/gocsv"

	"github.com/transcom/mymove/pkg/models"
)

// Zip3Fixture stores zip and rate and area information from tariff400ng_zip3s_fixture.csv
type Zip3Fixture struct {
	Zip      string `csv:"zip3"`
	RateArea string `csv:"rate_area"`
}

// Zip5Fixture stores zip and rate and area information from tariff400ng_zip5_rate_areas_fixture.csv
type Zip5Fixture struct {
	Zip      string `csv:"zip5"`
	RateArea string `csv:"rate_area"`
}

func (gre *GHCRateEngineImporter) mapZipsToRateAreas(dbTx *pop.Connection, zip3FixturePath string, zip5RateAreasFixturePath string) error {
	// Maps re_zip3s to re_rate_areas based on tariff400ng_zip3s_fixture.csv
	err := gre.mapZip3s(dbTx, zip3FixturePath)
	if err != nil {
		return fmt.Errorf("mapZip3s failed: %w", err)
	}

	// Creates re_zip5_rate_areas records from tariff400ng_zip5_rate_areas_fixture.csv
	err = gre.createZip5s(dbTx, zip5RateAreasFixturePath)
	if err != nil {
		return fmt.Errorf("createZip5s failed: %w", err)
	}

	return nil
}

func (gre *GHCRateEngineImporter) mapZip3s(dbTx *pop.Connection, fixturePath string) error {
	csvFile, err := os.Open(fixturePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer csvFile.Close()

	var zip3RateAreas []*Zip3Fixture

	err = gocsv.UnmarshalFile(csvFile, &zip3RateAreas)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}

	for _, zip3RateArea := range zip3RateAreas {
		var reZip3 models.ReZip3

		err = dbTx.Where("zip3 = ?", zip3RateArea.Zip).Where("contract_id = ?", gre.ContractID).First(&reZip3)
		if err != nil {
			return fmt.Errorf("failed to find ReZip3 record with zip %s: %w", zip3RateArea.Zip, err)
		}

		if zip3RateArea.RateArea == "ZIP" {
			reZip3.RateAreaID = nil
			reZip3.HasMultipleRateAreas = true

			verrs, err := dbTx.ValidateAndUpdate(&reZip3)
			if err != nil {
				return fmt.Errorf("failed to update ReZip3 %v: %w", reZip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate ReZip3 %v: %w", reZip3, verrs)
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
				return fmt.Errorf("failed to update ReZip3: %v: %w", reZip3.Zip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate ReZip3: %v: %w", reZip3.Zip3, verrs)
			}
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) createZip5s(dbTx *pop.Connection, fixturePath string) error {
	csvFile, err := os.Open(fixturePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer csvFile.Close()

	var zip5RateAreas []*Zip5Fixture

	err = gocsv.UnmarshalFile(csvFile, &zip5RateAreas)
	if err != nil {
		return fmt.Errorf("failed to unmarshal file: %w", err)
	}

	for _, zip5RateArea := range zip5RateAreas {
		var reZip5 models.ReZip5RateArea

		rateAreaID, found := gre.domesticRateAreaToIDMap[zip5RateArea.RateArea]
		if !found {
			return fmt.Errorf("failed to map %s rate area to ID", zip5RateArea.RateArea)
		}

		reZip5.ContractID = gre.ContractID
		reZip5.Zip5 = zip5RateArea.Zip
		reZip5.RateAreaID = rateAreaID

		verrs, err := dbTx.ValidateAndCreate(&reZip5)
		if err != nil {
			return fmt.Errorf("failed to update ReZip5RateArea: %v: %w", reZip5.Zip5, err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ReZip5RateArea: %v: %w", reZip5.Zip5, verrs)
		}
	}

	return nil
}