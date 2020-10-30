package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) mapZipCodesToRERateAreas(dbTx *pop.Connection) error {
	err := gre.mapREZip3sToRERateAreas(dbTx)
	if err != nil {
		return fmt.Errorf("mapREZip3sToRERateAreas failed: %w", err)
	}

	err = gre.createAndMapREZip5sToRERateAreas(dbTx)
	if err != nil {
		return fmt.Errorf("createAndMapREZip5sToRERateAreas failed: %w", err)
	}

	return nil
}

func (gre *GHCRateEngineImporter) mapREZip3sToRERateAreas(dbTx *pop.Connection) error {
	var reZip3s []models.ReZip3

	err := dbTx.Where("contract_id = ?", gre.ContractID).All(&reZip3s)
	if err != nil {
		return fmt.Errorf("failed to collect all ReZip3 records: %w", err)
	}

	for _, reZip3 := range reZip3s {
		rateArea, found := zip3ToRateAreaMappings[reZip3.Zip3]
		if !found {
			return fmt.Errorf("failed to find rate area map for zip3 %s in zip3ToRateAreaMappings", reZip3.Zip3)
		}

		if rateArea == "ZIP" {
			reZip3.RateAreaID = nil
			reZip3.HasMultipleRateAreas = true

			verrs, err := dbTx.ValidateAndUpdate(&reZip3)
			if err != nil {
				return fmt.Errorf("failed to update ReZip3 %v: %w", reZip3.Zip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate ReZip3 %v: %w", reZip3.Zip3, verrs)
			}
		} else {
			rateAreaID, found := gre.domesticRateAreaToIDMap[rateArea]
			if !found {
				return fmt.Errorf("failed to find ID for rate area %s in domesticRateAreaToIDMap", rateArea)
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

func (gre *GHCRateEngineImporter) createAndMapREZip5sToRERateAreas(dbTx *pop.Connection) error {
	for zip5, rateArea := range zip5ToRateAreaMappings {
		var reZip5RateArea models.ReZip5RateArea

		rateAreaID, found := gre.domesticRateAreaToIDMap[rateArea]
		if !found {
			return fmt.Errorf("failed to find ID for rate area %s in domesticRateAreaToIDMap", rateArea)
		}

		reZip5RateArea.ContractID = gre.ContractID
		reZip5RateArea.Zip5 = zip5
		reZip5RateArea.RateAreaID = rateAreaID

		verrs, err := dbTx.ValidateAndCreate(&reZip5RateArea)
		if err != nil {
			return fmt.Errorf("failed to update ReZip5RateArea: %v: %w", reZip5RateArea.Zip5, err)
		}
		if verrs.HasAny() {
			return fmt.Errorf("failed to validate ReZip5RateArea: %v: %w", reZip5RateArea.Zip5, verrs)
		}
	}

	return nil
}
