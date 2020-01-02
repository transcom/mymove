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
	serviceIOOLH, foundService := gre.serviceToIDMap["IOOLH"]
	if !foundService {
		return fmt.Errorf("missing service IOOLH in map of services")
	}

	//Int'l O->O UB
	serviceIOOUB, foundService := gre.serviceToIDMap["IOOUB"]
	if !foundService {
		return fmt.Errorf("missing service IOOUB in map of services")
	}

	//loop through the OCONUS to OCONUS data and store in db
	for _, stageOconusToOconusPrice := range oconusToOconusPrices {
		var intlPricingModels models.ReIntlPrices
		var peakPeriod bool
		peakPeriod, err = isPeakPeriod(stageOconusToOconusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageOconusToOconusPrice.Season, err)
		}

		originRateArea := stageOconusToOconusPrice.OriginIntlPriceAreaID
		originRateAreaID, found := gre.internationalRateAreaToIDMap[originRateArea]
		if !found {
			return fmt.Errorf("could not find origin rate area [%s] in map", stageOconusToOconusPrice.OriginIntlPriceAreaID)
		}

		destinationRateArea := stageOconusToOconusPrice.DestinationIntlPriceAreaID
		destinationRateAreaID, found := gre.internationalRateAreaToIDMap[destinationRateArea]
		if !found {
			return fmt.Errorf("could not find destination rate area [%s] in map", stageOconusToOconusPrice.DestinationIntlPriceAreaID)
		}

		var perUnitCentsHHG int
		perUnitCentsHHG, err = priceToCents(stageOconusToOconusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		var perUnitCentsUB int
		perUnitCentsUB, err = priceToCents(stageOconusToOconusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageOconusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		intlPricingModelIOOLH := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOOLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelIOOLH)

		intlPricingModelIOOUB := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOOUB,
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
	serviceICOLH, foundService := gre.serviceToIDMap["ICOLH"]
	if !foundService {
		return fmt.Errorf("missing service ICOLH in map of services")
	}

	//Int'l C->O UB
	serviceICOUB, foundService := gre.serviceToIDMap["ICOUB"]
	if !foundService {
		return fmt.Errorf("missing service ICOUB in map of services")
	}

	//loop through the CONUS to OCONUS data and store in db
	for _, stageConusToOconusPrice := range conusToOconusPrices {
		var intlPricingModels models.ReIntlPrices
		var peakPeriod bool
		peakPeriod, err = isPeakPeriod(stageConusToOconusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageConusToOconusPrice.Season, err)
		}

		originRateArea := stageConusToOconusPrice.OriginDomesticPriceAreaCode
		originRateAreaID, found := gre.domesticRateAreaToIDMap[originRateArea]
		if !found {
			return fmt.Errorf("could not find domestic rate area [%s] in map", stageConusToOconusPrice.OriginDomesticPriceAreaCode)
		}

		destinationRateArea := stageConusToOconusPrice.DestinationIntlPriceAreaID
		destinationRateAreaID, found := gre.internationalRateAreaToIDMap[destinationRateArea]
		if !found {
			return fmt.Errorf("could not find international rate area [%s] in map", stageConusToOconusPrice.DestinationIntlPriceAreaID)
		}

		var perUnitCentsHHG int
		perUnitCentsHHG, err = priceToCents(stageConusToOconusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageConusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		var perUnitCentsUB int
		perUnitCentsUB, err = priceToCents(stageConusToOconusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageConusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		intlPricingModelICOLH := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceICOLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelICOLH)

		intlPricingModelICOUB := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceICOUB,
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
	serviceIOCLH, foundService := gre.serviceToIDMap["IOCLH"]
	if !foundService {
		return fmt.Errorf("missing service IOCLH in map of services")
	}

	//Int'l O->C UB
	serviceIOCUB, foundService := gre.serviceToIDMap["IOCUB"]
	if !foundService {
		return fmt.Errorf("missing service IOCUB in map of services")
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
			return fmt.Errorf("could not find international rate area [%s] in map", stageOconusToConusPrice.OriginIntlPriceAreaID)
		}

		destinationRateArea := stageOconusToConusPrice.DestinationDomesticPriceAreaCode
		destinationRateAreaID, found := gre.domesticRateAreaToIDMap[destinationRateArea]
		if !found {
			return fmt.Errorf("could not find domestic rate area [%s] in map", stageOconusToConusPrice.DestinationDomesticPriceAreaCode)
		}

		perUnitCentsHHG, err := priceToCents(stageOconusToConusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToConusPrice.HHGShippingLinehaulPrice, err)
		}

		perUnitCentsUB, err := priceToCents(stageOconusToConusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageOconusToConusPrice.HHGShippingLinehaulPrice, err)
		}

		intlPricingModelIOCLH := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOCLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}

		intlPricingModels = append(intlPricingModels, intlPricingModelIOCLH)

		intlPricingModelIOCUB := models.ReIntlPrice{
			ContractID:            gre.contractID,
			ServiceID:             serviceIOCUB,
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
