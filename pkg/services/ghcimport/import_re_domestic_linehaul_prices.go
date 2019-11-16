package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticLinehaulPrices(dbTx *pop.Connection) error {
	// Read all the staged prices
	var stageDomesticLinehaulPrices []models.StageDomesticLinehaulPrice
	err := dbTx.All(&stageDomesticLinehaulPrices)
	if err != nil {
		return errors.Wrap(err, "Could not read staged domestic linehaul prices")
	}

	for _, stagePrice := range stageDomesticLinehaulPrices {
		weightLowerInt, err := stringToInteger(stagePrice.WeightLower)
		if err != nil {
			return errors.Wrapf(err, "Could not process weight lower [%s]", stagePrice.WeightLower)
		}

		weightUpperInt, err := stringToInteger(stagePrice.WeightUpper)
		if err != nil {
			return errors.Wrapf(err, "Could not process weight upper [%s]", stagePrice.WeightUpper)
		}

		milesLowerInt, err := stringToInteger(stagePrice.MilesLower)
		if err != nil {
			return errors.Wrapf(err, "Could not process miles lower [%s]", stagePrice.MilesLower)
		}

		milesUpperInt, err := stringToInteger(stagePrice.MilesUpper)
		if err != nil {
			return errors.Wrapf(err, "Could not process miles lower [%s]", stagePrice.MilesUpper)
		}

		isPeakPeriod, err := isPeakPeriod(stagePrice.Season)
		if err != nil {
			return errors.Wrapf(err, "Could not process season [%s]", stagePrice.Season)
		}

		serviceArea, err := cleanServiceAreaNumber(stagePrice.ServiceAreaNumber)
		if err != nil {
			return errors.Wrapf(err, "Could not process service area number [%s]", stagePrice.ServiceAreaNumber)
		}
		serviceAreaID, found := gre.serviceAreaToIDMap[serviceArea]
		if !found {
			return errors.New(fmt.Sprintf("Could not find service area [%s] in map", stagePrice.ServiceAreaNumber))
		}

		priceMillicents, err := priceToMillicents(stagePrice.Rate)
		if err != nil {
			return errors.Wrapf(err, "Could not process rate [%s]", stagePrice.MilesUpper)
		}

		domesticLinehaulPrice := models.ReDomesticLinehaulPrice{
			ContractID:            gre.contractID,
			WeightLower:           unit.Pound(weightLowerInt),
			WeightUpper:           unit.Pound(weightUpperInt),
			MilesLower:            milesLowerInt,
			MilesUpper:            milesUpperInt,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceMillicents:       unit.Millicents(priceMillicents),
		}

		verrs, err := dbTx.ValidateAndSave(&domesticLinehaulPrice)
		if err != nil {
			return errors.Wrapf(err, "Could not save domestic linehaul price: %+v", domesticLinehaulPrice)
		}
		if verrs.HasAny() {
			return errors.Wrapf(verrs, "Validation errors when saving domestic linehaul price: %+v", domesticLinehaulPrice)
		}
	}

	return nil
}
