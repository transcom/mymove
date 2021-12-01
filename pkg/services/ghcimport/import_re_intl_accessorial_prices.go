package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREIntlAccessorialPrices(appCtx appcontext.AppContext) error {
	//tab 5a) Access. and Add. Prices
	var intlAccessorialPrices []models.StageInternationalMoveAccessorialPrice
	err := appCtx.DB().All(&intlAccessorialPrices)
	if err != nil {
		return fmt.Errorf("could not read staged intl accessorial prices: %w", err)
	}

	services := []struct {
		serviceCode     models.ReServiceCode
		serviceProvided string
	}{
		{models.ReServiceCodeICRT, "Crating (per cubic ft.)"},
		{models.ReServiceCodeICRTSA, "Crating (per cubic ft.)"},
		{models.ReServiceCodeIUCRT, "Uncrating (per cubic ft.)"},
		{models.ReServiceCodeIDSHUT, "Shuttle Service (per cwt)"},
		{models.ReServiceCodeIOSHUT, "Shuttle Service (per cwt)"},
	}

	//loop through the intl accessorial price data and store in db
	for _, stageIntlAccessorialPrice := range intlAccessorialPrices {
		var perUnitCentsService int
		perUnitCentsService, err = priceToCents(stageIntlAccessorialPrice.PricePerUnit)
		if err != nil {
			return fmt.Errorf("could not process per unit price [%s]: %w", stageIntlAccessorialPrice.PricePerUnit, err)
		}

		market, err := getMarket(stageIntlAccessorialPrice.Market)
		if err != nil {
			return fmt.Errorf("could not process market [%s]: %w", stageIntlAccessorialPrice.Market, err)
		}

		serviceProvidedFound := false
		for _, service := range services {
			serviceCode := service.serviceCode
			serviceProvided := service.serviceProvided

			if stageIntlAccessorialPrice.ServiceProvided == serviceProvided {
				serviceProvidedFound = true
				serviceID, found := gre.serviceToIDMap[serviceCode]
				if !found {
					return fmt.Errorf("missing service [%s] in map of services", serviceCode)
				}

				intlAccessorial := models.ReIntlAccessorialPrice{
					ContractID:   gre.ContractID,
					Market:       market,
					ServiceID:    serviceID,
					PerUnitCents: unit.Cents(perUnitCentsService),
				}

				verrs, dbErr := appCtx.DB().ValidateAndSave(&intlAccessorial)
				if dbErr != nil {
					return fmt.Errorf("error saving ReIntlAccessorialPrices: %+v with error: %w", intlAccessorial, dbErr)
				}
				if verrs.HasAny() {
					return fmt.Errorf("error saving ReIntlAccessorialPrices: %+v with validation errors: %w", intlAccessorial, verrs)
				}
			}
		}
		if !serviceProvidedFound {
			return fmt.Errorf("service [%s] not found", stageIntlAccessorialPrice.ServiceProvided)
		}
	}

	return nil
}
