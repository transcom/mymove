package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREInternationalPrices(dbTx *pop.Connection) error {
	//tab 3a) OCONUS to OCONUS data
	var oconusToOconusPrices []models.StageOconusToOconusPrice
	err := dbTx.All(&oconusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged OCONUS to OCONUS prices: %w", err)
	}

	//Int'l O->O Shipping & LH
	var serviceIOOLH models.ReService
	err = dbTx.Where("code = 'IOOLH'").First(&serviceIOOLH)
	if err != nil {
		return fmt.Errorf("failed importing re_intl_prices from StageOconousToOconus with code IOOLH: %w", err)
	}

	//Int'l O->O UB
	var serviceIOOUB models.ReService
	err = dbTx.Where("code = 'IOOUB'").First(&serviceIOOUB)
	if err != nil {
		return fmt.Errorf("failed importing re_intl_prices from StageOconousToOconus with code IOOUB: %w", err)
	}

	//loop through the OCONUS to OCONUS data and store in db
	for _, stageOconusToOconusPrice := range oconusToOconusPrices {
		var intlPricingModels models.ReIntlPrices
		var peakPeriod bool
		peakPeriod, err = isPeakPeriod(stageOconusToOconusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process sesason [%s]: %w", stageOconusToOconusPrice.Season, err)
		}

		originRateArea := stageOconusToOconusPrice.OriginIntlPriceAreaID
		originRateAreaID, found := gre.internationalRateAreaToIDMap[originRateArea]
		if !found {
			return fmt.Errorf("could not find origin service area [%s] in map", stageOconusToOconusPrice.OriginIntlPriceAreaID)
		}

		destinationRateArea := stageOconusToOconusPrice.DestinationIntlPriceAreaID
		destinationRateAreaID, found := gre.internationalRateAreaToIDMap[destinationRateArea]
		if !found {
			return fmt.Errorf("could not find destination service area [%s] in map", stageOconusToOconusPrice.DestinationIntlPriceAreaID)
		}

		var perUnitCentsHHG int
		perUnitCentsHHG, err = priceToCents(stageOconusToOconusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		var perUnitCentsUB int
		perUnitCentsUB, err = priceToCents(stageOconusToOconusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		intlPricingModelIOOLH := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOOLH.ID,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelIOOLH)

		intlPricingModelIOOUB := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOOUB.ID,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelIOOUB)

		for _, model := range intlPricingModels {
			verrs, dbErr := dbTx.ValidateAndSave(&model)
			if dbErr != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, dbErr)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	//tab 3b CONUS to OCONUS data
	var conusToOconusPrices []models.StageConusToOconusPrice
	err = dbTx.All(&conusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged CONUS to OCONUS prices: %w", err)
	}

	//Int'l C->O Shipping & LH
	var serviceICOLH models.ReService
	err = dbTx.Where("code = 'ICOLH'").First(&serviceICOLH)
	if err != nil {
		return fmt.Errorf("failed importing re_intl_prices from StageConousToOconus with code ICOLH: %w", err)
	}

	//Int'l C->O UB
	var serviceICOUB models.ReService
	err = dbTx.Where("code = 'ICOUB'").First(&serviceICOUB)
	if err != nil {
		return fmt.Errorf("failed importing re_intl_prices from StageConousToOconus with code ICOUB: %w", err)
	}

	//loop through the CONUS to OCONUS data and store in db
	for _, stageConusToOconusPrice := range conusToOconusPrices {
		var intlPricingModels models.ReIntlPrices
		var peakPeriod bool
		peakPeriod, err = isPeakPeriod(stageConusToOconusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process sesason [%s]: %w", stageConusToOconusPrice.Season, err)
		}

		originRateArea := stageConusToOconusPrice.OriginDomesticPriceAreaCode
		originRateAreaID, found := gre.domesticRateAreaToIDMap[originRateArea]
		if !found {
			return fmt.Errorf("could not find service [%s] in map", stageConusToOconusPrice.OriginDomesticPriceAreaCode)
		}

		destinationRateArea := stageConusToOconusPrice.DestinationIntlPriceAreaID
		destinationRateAreaID, found := gre.internationalRateAreaToIDMap[destinationRateArea]
		if !found {
			return fmt.Errorf("could not find service [%s] in map", stageConusToOconusPrice.DestinationIntlPriceAreaID)
		}

		var perUnitCentsHHG int
		perUnitCentsHHG, err = priceToCents(stageConusToOconusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageConusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		var perUnitCentsUB int
		perUnitCentsUB, err = priceToCents(stageConusToOconusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageConusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		intlPricingModelICOLH := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceICOLH.ID,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelICOLH)

		intlPricingModelICOUB := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceICOUB.ID,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelICOUB)

		for _, model := range intlPricingModels {
			verrs, dbErr := dbTx.ValidateAndSave(&model)
			if dbErr != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, dbErr)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	//tab 3c OCONUS to CONUS data
	var oconusToConusPrices []models.StageOconusToConusPrice
	err = dbTx.All(&oconusToConusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged OCONUS to CONUS prices: %w", err)
	}

	//Int'l O->C Shipping & LH
	var serviceIOCLH models.ReService
	err = dbTx.Where("code = 'IOCLH'").First(&serviceIOCLH)
	if err != nil {
		return fmt.Errorf("failed importing re_intl_prices from StageOconousToConus with code IOCLH: %w", err)
	}

	//Int'l O->C UB
	var serviceIOCUB models.ReService
	err = dbTx.Where("code = 'IOCUB'").First(&serviceIOCUB)
	if err != nil {
		return fmt.Errorf("failed importing re_intl_prices from StageOconousToConus with code IOCUB: %w", err)
	}

	//loop through the OCONUS to CONUS data and store in db
	for _, stageOconusToConusPrice := range oconusToConusPrices {
		var intlPricingModels models.ReIntlPrices
		isPeakPeriod, err := isPeakPeriod(stageOconusToConusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageOconusToConusPrice.Season, err)
		}

		originRateArea := stageOconusToConusPrice.OriginIntlPriceAreaID
		originRateAreaID, found := gre.internationalRateAreaToIDMap[originRateArea]
		if !found {
			return fmt.Errorf("could not find service [%s] in map", stageOconusToConusPrice.OriginIntlPriceAreaID)
		}

		destinationRateArea := stageOconusToConusPrice.DestinationDomesticPriceAreaCode
		destinationRateAreaID, found := gre.domesticRateAreaToIDMap[destinationRateArea]
		if !found {
			return fmt.Errorf("could not find service [%s] in map", stageOconusToConusPrice.DestinationDomesticPriceAreaCode)
		}

		perUnitCentsHHG, err := priceToCents(stageOconusToConusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToConusPrice.HHGShippingLinehaulPrice, err)
		}

		perUnitCentsUB, err := priceToCents(stageOconusToConusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToConusPrice.HHGShippingLinehaulPrice, err)
		}

		intlPricingModelIOCLH := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOCLH.ID,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelIOCLH)

		intlPricingModelIOCUB := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOCUB.ID,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelIOCUB)

		for _, model := range intlPricingModels {
			verrs, err := dbTx.ValidateAndSave(&model)
			if err != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, err)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	return nil
}
