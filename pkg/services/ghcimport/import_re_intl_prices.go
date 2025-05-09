package ghcimport

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREInternationalPrices(appCtx appcontext.AppContext) error {
	if err := gre.importOconusToOconusPrices(appCtx); err != nil {
		return fmt.Errorf("could not import OCONUS to OCONUS prices: %w", err)
	}

	if err := gre.importConusToOconusPrices(appCtx); err != nil {
		return fmt.Errorf("could not import CONUS to OCONUS prices: %w", err)
	}

	if err := gre.importOconusToConusPrices(appCtx); err != nil {
		return fmt.Errorf("could not import OCONUS to CONUS prices: %w", err)
	}

	if err := gre.importNonStandardLocationPrices(appCtx); err != nil {
		return fmt.Errorf("could not import non-standard location prices: %w", err)
	}
	return nil
}

func (gre *GHCRateEngineImporter) importOconusToOconusPrices(appCtx appcontext.AppContext) error {
	// tab 3a) OCONUS to OCONUS data
	var oconusToOconusPrices []models.StageOconusToOconusPrice
	err := appCtx.DB().All(&oconusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged OCONUS to OCONUS prices: %w", err)
	}

	// Int'l O->O Shipping & LH
	serviceISLH, foundService := gre.serviceToIDMap[models.ReServiceCodeISLH]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeISLH)
	}

	// Int'l O->O UB
	serviceUBP, foundService := gre.serviceToIDMap[models.ReServiceCodeUBP]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeUBP)
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

		intlPricingModelISLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceISLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelISLH)

		intlPricingModelUBP := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceUBP,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelUBP)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := appCtx.DB().ValidateAndSave(&copyOfModel)
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

func (gre *GHCRateEngineImporter) importConusToOconusPrices(appCtx appcontext.AppContext) error {
	// tab 3b CONUS to OCONUS data
	var conusToOconusPrices []models.StageConusToOconusPrice
	err := appCtx.DB().All(&conusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged CONUS to OCONUS prices: %w", err)
	}

	// Int'l C->O Shipping & LH
	serviceISLH, foundService := gre.serviceToIDMap[models.ReServiceCodeISLH]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeISLH)
	}

	// Int'l C->O UB
	serviceUBP, foundService := gre.serviceToIDMap[models.ReServiceCodeUBP]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeUBP)
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

		intlPricingModelISLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceISLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelISLH)

		intlPricingModelUBP := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceUBP,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelUBP)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := appCtx.DB().ValidateAndSave(&copyOfModel)
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

func (gre *GHCRateEngineImporter) importOconusToConusPrices(appCtx appcontext.AppContext) error {
	// tab 3c OCONUS to CONUS data
	var oconusToConusPrices []models.StageOconusToConusPrice
	err := appCtx.DB().All(&oconusToConusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged OCONUS to CONUS prices: %w", err)
	}

	// Int'l O->C Shipping & LH
	serviceISLH, foundService := gre.serviceToIDMap[models.ReServiceCodeISLH]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeISLH)
	}

	// Int'l O->C UB
	serviceUBP, foundService := gre.serviceToIDMap[models.ReServiceCodeUBP]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeUBP)
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

		intlPricingModelISLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceISLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelISLH)

		intlPricingModelUBP := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceUBP,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          isPeakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelUBP)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := appCtx.DB().ValidateAndSave(&copyOfModel)
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

func (gre *GHCRateEngineImporter) importNonStandardLocationPrices(appCtx appcontext.AppContext) error {
	// tab 3e) Non-standard location prices
	var nonStandardLocnPrices []models.StageNonStandardLocnPrice
	err := appCtx.DB().All(&nonStandardLocnPrices)
	if err != nil {
		return fmt.Errorf("could not read staged non-standard location prices: %w", err)
	}

	// Int'l non-standard HHG
	serviceISLH, foundService := gre.serviceToIDMap[models.ReServiceCodeISLH]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeISLH)
	}

	// Int'l non-standard UB
	serviceUBP, foundService := gre.serviceToIDMap[models.ReServiceCodeUBP]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeUBP)
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

		intlPricingModelISLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceISLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelISLH)

		intlPricingModelUBP := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceUBP,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelUBP)

		for _, model := range intlPricingModels {
			copyOfModel := model // Make copy to avoid implicit memory aliasing of items from a range statement.
			verrs, dbErr := appCtx.DB().ValidateAndSave(&copyOfModel)
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
