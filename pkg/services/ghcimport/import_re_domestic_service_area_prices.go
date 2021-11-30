package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticServiceAreaPrices(appCtx appcontext.AppContext) error {
	var stageDomPricingModels []models.StageDomesticServiceAreaPrice

	if err := appCtx.DB().All(&stageDomPricingModels); err != nil {
		return fmt.Errorf("error looking up StageDomesticServiceAreaPrice data: %w", err)
	}

	for _, stageDomPricingModel := range stageDomPricingModels {
		isPeakPeriod, ippErr := isPeakPeriod(stageDomPricingModel.Season)
		if ippErr != nil {
			return ippErr
		}

		serviceAreaNumber, csaErr := cleanServiceAreaNumber(stageDomPricingModel.ServiceAreaNumber)
		if csaErr != nil {
			return csaErr
		}

		serviceAreaID, found := gre.serviceAreaToIDMap[serviceAreaNumber]
		if !found {
			return fmt.Errorf("could not find service area [%s] in map", serviceAreaNumber)
		}

		servicesToInsert := []struct {
			service models.ReServiceCode
			price   string
		}{
			{models.ReServiceCodeDSH, stageDomPricingModel.ShorthaulPrice},
			{models.ReServiceCodeDOP, stageDomPricingModel.OriginDestinationPrice},
			{models.ReServiceCodeDDP, stageDomPricingModel.OriginDestinationPrice},
			{models.ReServiceCodeDOFSIT, stageDomPricingModel.OriginDestinationSITFirstDayWarehouse},
			{models.ReServiceCodeDDFSIT, stageDomPricingModel.OriginDestinationSITFirstDayWarehouse},
			{models.ReServiceCodeDOASIT, stageDomPricingModel.OriginDestinationSITAddlDays},
			{models.ReServiceCodeDDASIT, stageDomPricingModel.OriginDestinationSITAddlDays},
		}

		for _, serviceToInsert := range servicesToInsert {
			service := serviceToInsert.service
			price := serviceToInsert.price

			serviceID, found := gre.serviceToIDMap[service]
			if !found {
				return fmt.Errorf("missing service [%s] in map of services", service)
			}

			cents, convErr := priceToCents(price)
			if convErr != nil {
				return fmt.Errorf("failed to parse price for service code %s: %+v error: %w", service, price, convErr)
			}

			domPricingModel := models.ReDomesticServiceAreaPrice{
				ContractID:            gre.ContractID,
				ServiceID:             serviceID,
				IsPeakPeriod:          isPeakPeriod,
				DomesticServiceAreaID: serviceAreaID,
				PriceCents:            unit.Cents(cents),
			}

			verrs, err := appCtx.DB().ValidateAndSave(&domPricingModel)
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReDomesticServiceAreaPrices: %+v with validation errors: %w", domPricingModel, verrs)
			}
			if err != nil {
				return fmt.Errorf("error saving ReDomesticServiceAreaPrices: %+v with error: %w", domPricingModel, err)
			}
		}
	}

	return nil
}
