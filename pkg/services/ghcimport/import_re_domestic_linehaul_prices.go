package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticLinehaulPrices(dbTx *pop.Connection) error {
	// Read all the staged prices
	var stageDomesticLinehaulPrices []models.StageDomesticLinehaulPrice
	err := dbTx.All(&stageDomesticLinehaulPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic linehaul prices: %w", err)
	}

	for _, stagePrice := range stageDomesticLinehaulPrices {
		weightLowerInt, err := stringToInteger(stagePrice.WeightLower)
		if err != nil {
			return fmt.Errorf("could not process weight lower [%s]: %w", stagePrice.WeightLower, err)
		}

		weightUpperInt, err := stringToInteger(stagePrice.WeightUpper)
		if err != nil {
			return fmt.Errorf("could not process weight upper [%s]: %w", stagePrice.WeightUpper, err)
		}

		milesLowerInt, err := stringToInteger(stagePrice.MilesLower)
		if err != nil {
			return fmt.Errorf("could not process miles lower [%s]: %w", stagePrice.MilesLower, err)
		}

		milesUpperInt, err := stringToInteger(stagePrice.MilesUpper)
		if err != nil {
			return fmt.Errorf("could not process miles upper [%s]: %w", stagePrice.MilesUpper, err)
		}

		isPeakPeriod, err := isPeakPeriod(stagePrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stagePrice.Season, err)
		}

		serviceArea, err := cleanServiceAreaNumber(stagePrice.ServiceAreaNumber)
		if err != nil {
			return fmt.Errorf("could not process service area number [%s]: %w", stagePrice.ServiceAreaNumber, err)
		}
		serviceAreaID, found := gre.serviceAreaToIDMap[serviceArea]
		if !found {
			return fmt.Errorf("could not find service area [%s] in map", stagePrice.ServiceAreaNumber)
		}

		priceMillicents, err := priceToMillicents(stagePrice.Rate)
		if err != nil {
			return fmt.Errorf("could not process rate [%s]: %w", stagePrice.Rate, err)
		}

		domesticLinehaulPrice := models.ReDomesticLinehaulPrice{
			ContractID:            gre.ContractID,
			WeightLower:           unit.Pound(weightLowerInt),
			WeightUpper:           unit.Pound(weightUpperInt),
			MilesLower:            milesLowerInt,
			MilesUpper:            milesUpperInt,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceMillicents:       unit.Millicents(priceMillicents),
		}

		verrs, err := dbTx.ValidateAndSave(&domesticLinehaulPrice)
		if verrs.HasAny() {
			return fmt.Errorf("validation errors when saving domestic linehaul price [%+v]: %w", domesticLinehaulPrice, verrs)
		}
		if err != nil {
			return fmt.Errorf("could not save domestic linehaul price [%+v]: %w", domesticLinehaulPrice, err)
		}
	}

	return nil
}
