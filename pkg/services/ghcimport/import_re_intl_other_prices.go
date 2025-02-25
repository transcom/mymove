package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREInternationalOtherPrices(appCtx appcontext.AppContext) error {
	// Tab 3d: Other International Prices
	var stageOtherIntlPrices []models.StageOtherIntlPrice
	err := appCtx.DB().All(&stageOtherIntlPrices)
	if err != nil {
		return fmt.Errorf("could not read staged other international prices: %w", err)
	}

	for _, stagePrice := range stageOtherIntlPrices {
		isPeakPeriod, err := isPeakPeriod(stagePrice.Season)
		if err != nil {
			return fmt.Errorf("could not process season [%s]: %w", stagePrice.Season, err)
		}

		// First look in the international rate area map; if not found, try the domestic rate
		// area map.  Both international and domestic are possibilities for this column.
		rateAreaCode := stagePrice.RateAreaCode
		rateAreaID, found := gre.internationalRateAreaToIDMap[rateAreaCode]
		if !found {
			rateAreaID, found = gre.domesticRateAreaToIDMap[rateAreaCode]
			if !found {
				return fmt.Errorf("could not find rate area [%s] in map", rateAreaCode)
			}
		}

		servicesToInsert := []struct {
			service models.ReServiceCode
			price   string
		}{
			{models.ReServiceCodeIHPK, stagePrice.HHGOriginPackPrice},
			{models.ReServiceCodeIHUPK, stagePrice.HHGDestinationUnPackPrice},
			{models.ReServiceCodeIUBPK, stagePrice.UBOriginPackPrice},
			{models.ReServiceCodeIUBUPK, stagePrice.UBDestinationUnPackPrice},
			{models.ReServiceCodeIOFSIT, stagePrice.OriginDestinationSITFirstDayWarehouse},
			{models.ReServiceCodeIDFSIT, stagePrice.OriginDestinationSITFirstDayWarehouse},
			{models.ReServiceCodeIOASIT, stagePrice.OriginDestinationSITAddlDays},
			{models.ReServiceCodeIDASIT, stagePrice.OriginDestinationSITAddlDays},
			{models.ReServiceCodeIOPSIT, stagePrice.SITLte50Miles},
			{models.ReServiceCodeIDDSIT, stagePrice.SITGt50Miles},
		}

		for _, serviceToInsert := range servicesToInsert {
			service := serviceToInsert.service
			price := serviceToInsert.price

			priceCents, err := priceToCents(price)
			if err != nil {
				return fmt.Errorf("could not process price [%s] for service [%s]: %w", price, service, err)
			}
			serviceID, found := gre.serviceToIDMap[service]
			if !found {
				return fmt.Errorf("missing service [%s] in map of services", service)
			}

			intlOtherPrice := models.ReIntlOtherPrice{
				ContractID:   gre.ContractID,
				ServiceID:    serviceID,
				RateAreaID:   rateAreaID,
				IsPeakPeriod: isPeakPeriod,
				PerUnitCents: unit.Cents(priceCents),
			}

			verrs, err := appCtx.DB().ValidateAndSave(&intlOtherPrice)
			if verrs.HasAny() {
				return fmt.Errorf("validation errors when saving other international price [%+v]: %w", intlOtherPrice, verrs)
			}
			if err != nil {
				return fmt.Errorf("could not save other international price [%+v]: %w", intlOtherPrice, err)
			}
		}
	}

	return nil
}
