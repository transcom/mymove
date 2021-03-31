package ghcimport

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREInternationalPrices(dbTx *pop.Connection) error {
	if err := gre.importOconusToOconusPrices(dbTx); err != nil {
		return fmt.Errorf("could not import OCONUS to OCONUS prices: %w", err)
	}

	if err := gre.importConusToOconusPrices(dbTx); err != nil {
		return fmt.Errorf("could not import CONUS to OCONUS prices: %w", err)
	}

	if err := gre.importOconusToConusPrices(dbTx); err != nil {
		return fmt.Errorf("could not import OCONUS to CONUS prices: %w", err)
	}

	if err := gre.importNonStandardLocationPrices(dbTx); err != nil {
		return fmt.Errorf("could not import non-standard location prices: %w", err)
	}

	return nil
}

func (gre *GHCRateEngineImporter) importOconusToOconusPrices(dbTx *pop.Connection) error {
	// tab 3a) OCONUS to OCONUS data
	var oconusToOconusPrices []models.StageOconusToOconusPrice
	err := dbTx.All(&oconusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged OCONUS to OCONUS prices: %w", err)
	}

	// Int'l O->O Shipping & LH
	serviceIOOLH, foundService := gre.serviceToIDMap["IOOLH"]
	if !foundService {
		return fmt.Errorf("missing service IOOLH in map of services")
	}

	// Int'l O->O UB
	serviceIOOUB, foundService := gre.serviceToIDMap["IOOUB"]
	if !foundService {
		return fmt.Errorf("missing service IOOUB in map of services")
	}

	// loop through the OCONUS to OCONUS data and store in db
	for _, stageOconusToOconusPrice := range oconusToOconusPrices {
		var intlPricingModels models.ReIntlPrices
		peakPeriod, err := isPeakPeriod(stageOconusToOconusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageOconusToOconusPrice.Season, err)
		}

		originRateAreaID, found := gre.internationalRateAreaToIDMap[stageOconusToOconusPrice.OriginIntlPriceAreaID]
		if !found {
			return fmt.Errorf("could not find origin rate area [%s] in map", stageOconusToOconusPrice.OriginIntlPriceAreaID)
		}

		destinationRateAreaID, found := gre.internationalRateAreaToIDMap[stageOconusToOconusPrice.DestinationIntlPriceAreaID]
		if !found {
			return fmt.Errorf("could not find destination rate area [%s] in map", stageOconusToOconusPrice.DestinationIntlPriceAreaID)
		}

		perUnitCentsHHG, err := priceToCents(stageOconusToOconusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		perUnitCentsUB, err := priceToCents(stageOconusToOconusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageOconusToOconusPrice.UBPrice, err)
		}

		intlPricingModelIOOLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceIOOLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelIOOLH)

		intlPricingModelIOOUB := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceIOOUB,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelIOOUB)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := dbTx.ValidateAndSave(&copyOfModel)
			if dbErr != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, dbErr)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) importConusToOconusPrices(dbTx *pop.Connection) error {
	// tab 3b CONUS to OCONUS data
	var conusToOconusPrices []models.StageConusToOconusPrice
	err := dbTx.All(&conusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged CONUS to OCONUS prices: %w", err)
	}

	// Int'l C->O Shipping & LH
	serviceICOLH, foundService := gre.serviceToIDMap["ICOLH"]
	if !foundService {
		return fmt.Errorf("missing service ICOLH in map of services")
	}

	// Int'l C->O UB
	serviceICOUB, foundService := gre.serviceToIDMap["ICOUB"]
	if !foundService {
		return fmt.Errorf("missing service ICOUB in map of services")
	}

	// loop through the CONUS to OCONUS data and store in db
	for _, stageConusToOconusPrice := range conusToOconusPrices {
		var intlPricingModels models.ReIntlPrices

		peakPeriod, err := isPeakPeriod(stageConusToOconusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageConusToOconusPrice.Season, err)
		}

		originRateAreaID, found := gre.domesticRateAreaToIDMap[stageConusToOconusPrice.OriginDomesticPriceAreaCode]
		if !found {
			return fmt.Errorf("could not find domestic rate area [%s] in map", stageConusToOconusPrice.OriginDomesticPriceAreaCode)
		}

		destinationRateAreaID, found := gre.internationalRateAreaToIDMap[stageConusToOconusPrice.DestinationIntlPriceAreaID]
		if !found {
			return fmt.Errorf("could not find international rate area [%s] in map", stageConusToOconusPrice.DestinationIntlPriceAreaID)
		}

		perUnitCentsHHG, err := priceToCents(stageConusToOconusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageConusToOconusPrice.HHGShippingLinehaulPrice, err)
		}

		perUnitCentsUB, err := priceToCents(stageConusToOconusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageConusToOconusPrice.UBPrice, err)
		}

		intlPricingModelICOLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceICOLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelICOLH)

		intlPricingModelICOUB := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceICOUB,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelICOUB)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := dbTx.ValidateAndSave(&copyOfModel)
			if dbErr != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, dbErr)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) importOconusToConusPrices(dbTx *pop.Connection) error {
	// tab 3c OCONUS to CONUS data
	var oconusToConusPrices []models.StageOconusToConusPrice
	err := dbTx.All(&oconusToConusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged OCONUS to CONUS prices: %w", err)
	}

	// Int'l O->C Shipping & LH
	serviceIOCLH, foundService := gre.serviceToIDMap["IOCLH"]
	if !foundService {
		return fmt.Errorf("missing service IOCLH in map of services")
	}

	// Int'l O->C UB
	serviceIOCUB, foundService := gre.serviceToIDMap["IOCUB"]
	if !foundService {
		return fmt.Errorf("missing service IOCUB in map of services")
	}

	// loop through the OCONUS to CONUS data and store in db
	for _, stageOconusToConusPrice := range oconusToConusPrices {
		var intlPricingModels models.ReIntlPrices

		isPeakPeriod, err := isPeakPeriod(stageOconusToConusPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageOconusToConusPrice.Season, err)
		}

		originRateAreaID, found := gre.internationalRateAreaToIDMap[stageOconusToConusPrice.OriginIntlPriceAreaID]
		if !found {
			return fmt.Errorf("could not find international rate area [%s] in map", stageOconusToConusPrice.OriginIntlPriceAreaID)
		}

		destinationRateAreaID, found := gre.domesticRateAreaToIDMap[stageOconusToConusPrice.DestinationDomesticPriceAreaCode]
		if !found {
			return fmt.Errorf("could not find domestic rate area [%s] in map", stageOconusToConusPrice.DestinationDomesticPriceAreaCode)
		}

		perUnitCentsHHG, err := priceToCents(stageOconusToConusPrice.HHGShippingLinehaulPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageOconusToConusPrice.HHGShippingLinehaulPrice, err)
		}

		perUnitCentsUB, err := priceToCents(stageOconusToConusPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageOconusToConusPrice.UBPrice, err)
		}

		intlPricingModelIOCLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceIOCLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelIOCLH)

		intlPricingModelIOCUB := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceIOCUB,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelIOCUB)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := dbTx.ValidateAndSave(&copyOfModel)
			if dbErr != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, dbErr)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) importNonStandardLocationPrices(dbTx *pop.Connection) error {
	// tab 3e) Non-standard location prices
	var nonStandardLocnPrices []models.StageNonStandardLocnPrice
	err := dbTx.All(&nonStandardLocnPrices)
	if err != nil {
		return fmt.Errorf("could not read staged non-standard location prices: %w", err)
	}

	// Int'l non-standard HHG
	serviceNSTH, foundService := gre.serviceToIDMap["NSTH"]
	if !foundService {
		return fmt.Errorf("missing service NSTH in map of services")
	}

	// Int'l non-standard UB
	serviceNSTUB, foundService := gre.serviceToIDMap["NSTUB"]
	if !foundService {
		return fmt.Errorf("missing service NSTUB in map of services")
	}

	// loop through the non-standard location data and store in db
	for _, stageNonStandardLocnPrice := range nonStandardLocnPrices {
		var intlPricingModels models.ReIntlPrices

		peakPeriod, err := isPeakPeriod(stageNonStandardLocnPrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stageNonStandardLocnPrice.Season, err)
		}

		moveToAndFromKind := strings.Split(stageNonStandardLocnPrice.MoveType, " to ")
		if len(moveToAndFromKind) != 2 {
			return fmt.Errorf("could not parse move type [%s]", stageNonStandardLocnPrice.MoveType)
		}

		originRateAreaID, err := gre.getRateAreaIDForKind(stageNonStandardLocnPrice.OriginID, moveToAndFromKind[0])
		if err != nil {
			return err
		}

		destinationRateAreaID, err := gre.getRateAreaIDForKind(stageNonStandardLocnPrice.DestinationID, moveToAndFromKind[1])
		if err != nil {
			return err
		}

		perUnitCentsHHG, err := priceToCents(stageNonStandardLocnPrice.HHGPrice)
		if err != nil {
			return fmt.Errorf("could not process linehaul price [%s]: %w", stageNonStandardLocnPrice.HHGPrice, err)
		}

		perUnitCentsUB, err := priceToCents(stageNonStandardLocnPrice.UBPrice)
		if err != nil {
			return fmt.Errorf("could not process UB price [%s]: %w", stageNonStandardLocnPrice.UBPrice, err)
		}

		intlPricingModelNSTH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceNSTH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelNSTH)

		intlPricingModelNSTUB := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceNSTUB,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelNSTUB)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := dbTx.ValidateAndSave(&copyOfModel)
			if dbErr != nil {
				return fmt.Errorf("error saving ReIntlPrices: %+v with error: %w", model, dbErr)
			}
			if verrs.HasAny() {
				return fmt.Errorf("error saving ReIntlPrices: %+v with validation errors: %w", model, verrs)
			}
		}
	}

	return nil
}

func (gre *GHCRateEngineImporter) getRateAreaIDForKind(rateArea string, kind string) (uuid.UUID, error) {
	switch kind {
	case "NSRA", "OCONUS":
		intlRateAreaID, found := gre.internationalRateAreaToIDMap[rateArea]
		if !found {
			return uuid.Nil, fmt.Errorf("could not find rate area [%s] in international rate area map", rateArea)
		}
		return intlRateAreaID, nil
	case "CONUS":
		domesticRateAreaID, found := gre.domesticRateAreaToIDMap[rateArea]
		if !found {
			return uuid.Nil, fmt.Errorf("could not find rate area [%s] in domestic rate area map", rateArea)
		}
		return domesticRateAreaID, nil
	}

	return uuid.Nil, fmt.Errorf("unexpected rate area kind [%s]", kind)
}
