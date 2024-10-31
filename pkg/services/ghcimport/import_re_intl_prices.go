package ghcimport

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREInternationalPrices(appCtx appcontext.AppContext) error {
	if err := gre.importInternationalPrices(appCtx); err != nil {
		return fmt.Errorf("could not import INTL prices: %w", err)
	}

	return nil
}

func (gre *GHCRateEngineImporter) importInternationalPrices(appCtx appcontext.AppContext) error {

	var oconusToOconusPrices []models.StageOconusToOconusPrice
	err := appCtx.DB().All(&oconusToOconusPrices)
	if err != nil {
		return fmt.Errorf("could not read staged INTL prices: %w", err)
	}

	serviceUBP, foundService := gre.serviceToIDMap[models.ReServiceCodeUBP]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeUBP)
	}

	serviceISLH, foundService := gre.serviceToIDMap[models.ReServiceCodeISLH]
	if !foundService {
		return fmt.Errorf("missing service %s in map of services", models.ReServiceCodeISLH)
	}

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

		intlPricingModelUBP := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceUBP,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsUB),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelUBP)

		intlPricingModelISLH := models.ReIntlPrice{
			ContractID:            gre.ContractID,
			ServiceID:             serviceISLH,
			OriginRateAreaID:      originRateAreaID,
			DestinationRateAreaID: destinationRateAreaID,
			IsPeakPeriod:          peakPeriod,
			PerUnitCents:          unit.Cents(perUnitCentsHHG),
		}
		intlPricingModels = append(intlPricingModels, intlPricingModelISLH)

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
