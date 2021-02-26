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
		copyOfReZip3 := reZip3 // Make copy to avoid implicit memory aliasing of items from a range statement.
		rateArea, found := zip3ToRateAreaMappings[copyOfReZip3.Zip3]
		if !found {
			return fmt.Errorf("failed to find rate area map for zip3 %s in zip3ToRateAreaMappings", copyOfReZip3.Zip3)
		}

		if rateArea == "ZIP" {
			copyOfReZip3.RateAreaID = nil
			copyOfReZip3.HasMultipleRateAreas = true

			verrs, err := dbTx.ValidateAndUpdate(&copyOfReZip3)
			if err != nil {
				return fmt.Errorf("failed to update ReZip3 %v: %w", copyOfReZip3.Zip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate ReZip3 %v: %w", copyOfReZip3.Zip3, verrs)
			}
		} else {
			rateAreaID, found := gre.domesticRateAreaToIDMap[rateArea]
			if !found {
				return fmt.Errorf("failed to find ID for rate area %s in domesticRateAreaToIDMap", rateArea)
			}

			copyOfReZip3.RateAreaID = &rateAreaID
			copyOfReZip3.HasMultipleRateAreas = false

			verrs, err := dbTx.ValidateAndUpdate(&copyOfReZip3)
			if err != nil {
				return fmt.Errorf("failed to update ReZip3: %v: %w", copyOfReZip3.Zip3, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("failed to validate ReZip3: %v: %w", copyOfReZip3.Zip3, verrs)
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
