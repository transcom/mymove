package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) appendDomesticServiceAreaPrice(
	domPricingModels models.ReDomesticServiceAreaPrices,
	code string,
	contractID uuid.UUID,
	isPeakPeriod bool,
	serviceAreaID uuid.UUID,
	price string,
) (models.ReDomesticServiceAreaPrices, error) {
	serviceID, found := gre.serviceToIDMap[code]
	if !found {
		return domPricingModels, fmt.Errorf("missing service [%s] in map of services", code)
	}

	cents, convErr := priceToCents(price)
	if convErr != nil {
		return domPricingModels, fmt.Errorf("failed to parse price for service code %s: %+v error: %w", code, price, convErr)
	}

	domPricingModel := models.ReDomesticServiceAreaPrice{
		ContractID:            contractID,
		ServiceID:             serviceID,
		IsPeakPeriod:          isPeakPeriod,
		DomesticServiceAreaID: serviceAreaID,
		PriceCents:            unit.Cents(cents),
	}

	return append(domPricingModels, domPricingModel), nil
}

func (gre *GHCRateEngineImporter) importREDomesticServiceAreaPrices(db *pop.Connection) error {
	var stageDomPricingModels []models.StageDomesticServiceAreaPrice

	if err := db.All(&stageDomPricingModels); err != nil {
		return fmt.Errorf("error looking up StageDomesticServiceAreaPrice data: %w", err)
	}

	for _, stageDomPricingModel := range stageDomPricingModels {
		var domPricingModels models.ReDomesticServiceAreaPrices

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

		//DSH - ShorthaulPrice
		var err error
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DSH", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.ShorthaulPrice)
		if err != nil {
			return err
		}

		//DOP - OriginPrice
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DOP", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.OriginDestinationPrice)
		if err != nil {
			return err
		}

		//DDP - DestinationPrice
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DDP", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.OriginDestinationPrice)
		if err != nil {
			return err
		}

		//DOFSIT - OriginSITFirstDayWarehouse
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DOFSIT", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.OriginDestinationSITFirstDayWarehouse)
		if err != nil {
			return err
		}

		//DDFSIT - DestinationSITFirstDayWarehouse
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DDFSIT", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.OriginDestinationSITFirstDayWarehouse)
		if err != nil {
			return err
		}

		//DOASIT - OriginSITAddlDays
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DOASIT", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.OriginDestinationSITAddlDays)
		if err != nil {
			return err
		}

		//DDASIT - DestinationSITAddlDays
		domPricingModels, err = gre.appendDomesticServiceAreaPrice(domPricingModels, "DDASIT", gre.contractID, isPeakPeriod, serviceAreaID, stageDomPricingModel.OriginDestinationSITAddlDays)
		if err != nil {
			return err
		}

		for _, model := range domPricingModels {
			verrs, err := db.ValidateAndSave(&model)
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReDomesticServiceAreaPrices: %+v with validation errors: %w", model, verrs)
			}
			if err != nil {
				return fmt.Errorf("error saving ReDomesticServiceAreaPrices: %+v with error: %w", model, err)
			}
		}
	}

	return nil
}
