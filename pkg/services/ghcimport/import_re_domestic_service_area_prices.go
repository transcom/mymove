package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

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
		var serviceDSH models.ReService
		err := db.Where("code = 'DSH'").First(&serviceDSH)
		if err != nil {
			return fmt.Errorf("failed importing re_service from StageDomesticServiceAreaPrice with code DSH: %w", err)
		}

		cents, convErr := priceToCents(stageDomPricingModel.ShorthaulPrice)
		if convErr != nil {
			return fmt.Errorf("failed to parse price for Shorthaul data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceDSH.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDSH)

		//DOP - OriginDestinationPrice
		var serviceDOP models.ReService
		err = db.Where("code = 'DOP'").First(&serviceDOP)
		if err != nil {
			return fmt.Errorf("failed importing re_service from StageDomesticServiceAreaPrice with code DOP: %w", err)
		}

		cents, convErr = priceToCents(stageDomPricingModel.OriginDestinationPrice)
		if convErr != nil {
			return fmt.Errorf("failed to parse price for OriginDestinationPrice data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDOP := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceDOP.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDOP)

		//DOFSIT - OriginDestinationSITFirstDayWarehouse
		var serviceDOFSIT models.ReService
		err = db.Where("code = 'DOFSIT'").First(&serviceDOFSIT)
		if err != nil {
			return fmt.Errorf("failed importing re_service from StageDomesticServiceAreaPrice with code DOFSIT: %w", err)
		}

		cents, convErr = priceToCents(stageDomPricingModel.OriginDestinationSITFirstDayWarehouse)
		if convErr != nil {
			return fmt.Errorf("failed to parse price for OriginDestinationSITFirstDayWarehouse data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDOFSIT := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceDOFSIT.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDOFSIT)

		//DOASIT - OriginDestinationSITAddlDays
		var serviceDOASIT models.ReService
		err = db.Where("code = 'DOASIT'").First(&serviceDOASIT)
		if err != nil {
			return fmt.Errorf("failed importing re_service from StageDomesticServiceAreaPrice with code DOASIT: %w", err)
		}

		cents, convErr = priceToCents(stageDomPricingModel.OriginDestinationSITAddlDays)
		if convErr != nil {
			return fmt.Errorf("failed to parse price for OriginDestinationSITAddlDays data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDOASIT := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceDOASIT.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDOASIT)

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
