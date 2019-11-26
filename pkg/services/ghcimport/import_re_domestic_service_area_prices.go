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
		return fmt.Errorf("Error looking up StageDomesticServiceAreaPrice data: %w", err)
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
		service, err := models.FetchReServiceItem(db, "DSH")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DSH: %w", err)
		}

		cents, convErr := priceToCents(stageDomPricingModel.ShorthaulPrice)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for Shorthaul data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDSH := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDSH)

		//DODP - OriginDestinationPrice
		service, err = models.FetchReServiceItem(db, "DODP")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DODP: %w", err)
		}

		cents, convErr = priceToCents(stageDomPricingModel.OriginDestinationPrice)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for OriginDestinationPrice data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDODP := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDODP)

		//DFSIT - OriginDestinationSITFirstDayWarehouse
		service, err = models.FetchReServiceItem(db, "DFSIT")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DFSIT: %w", err)
		}

		cents, convErr = priceToCents(stageDomPricingModel.OriginDestinationSITFirstDayWarehouse)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for OriginDestinationSITFirstDayWarehouse data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDFSIT := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDFSIT)

		//DASIT - OriginDestinationSITAddlDays
		service, err = models.FetchReServiceItem(db, "DASIT")
		if err != nil {
			return fmt.Errorf("Failed importing re_service from StageDomesticServiceAreaPrice with code DASIT: %w", err)
		}

		cents, convErr = priceToCents(stageDomPricingModel.OriginDestinationSITAddlDays)
		if convErr != nil {
			return fmt.Errorf("Failed to parse price for OriginDestinationSITAddlDays data: %+v error: %w", stageDomPricingModel, convErr)
		}

		domPricingModelDASIT := models.ReDomesticServiceAreaPrice{
			ContractID:            gre.contractID,
			ServiceID:             service.ID,
			IsPeakPeriod:          isPeakPeriod,
			DomesticServiceAreaID: serviceAreaID,
			PriceCents:            unit.Cents(cents),
		}

		domPricingModels = append(domPricingModels, domPricingModelDASIT)

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
