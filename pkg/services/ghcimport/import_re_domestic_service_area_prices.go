package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticServiceAreaPrices(db *pop.Connection) error {
	var stageDomPricingModels []models.StageDomesticServiceAreaPrice

	if err := db.All(&stageDomPricingModels); err != nil {
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
			service string
			price   string
		}{
			{"DSH", stageDomPricingModel.ShorthaulPrice},
			{"DOP", stageDomPricingModel.OriginDestinationPrice},
			{"DDP", stageDomPricingModel.OriginDestinationPrice},
			{"DOFSIT", stageDomPricingModel.OriginDestinationSITFirstDayWarehouse},
			{"DDFSIT", stageDomPricingModel.OriginDestinationSITFirstDayWarehouse},
			{"DOASIT", stageDomPricingModel.OriginDestinationSITAddlDays},
			{"DDASIT", stageDomPricingModel.OriginDestinationSITAddlDays},
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

			verrs, err := db.ValidateAndSave(&domPricingModel)
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
