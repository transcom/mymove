package ghcimport

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticAccessorialPrices(appCtx appcontext.AppContext) error {
	//tab 5a) Access. and Add. Prices
	var domesticAccessorialPrices []models.StageDomesticMoveAccessorialPrice
	err := appCtx.DB().All(&domesticAccessorialPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic accessorial prices: %w", err)
	}

	services := []struct {
		serviceCode     models.ReServiceCode
		serviceProvided string
	}{
		{models.ReServiceCodeDCRT, "Crating (per cubic ft.)"},
		{models.ReServiceCodeDCRTSA, "Crating (per cubic ft.)"},
		{models.ReServiceCodeDUCRT, "Uncrating (per cubic ft.)"},
		{models.ReServiceCodeDDSHUT, "Shuttle Service (per cwt)"},
		{models.ReServiceCodeDOSHUT, "Shuttle Service (per cwt)"},
	}

	//loop through the domestic accessorial price data and store in db
	for _, stageDomesticAccessorialPrice := range domesticAccessorialPrices {
		servicesSchedule, err := stringToInteger(stageDomesticAccessorialPrice.ServicesSchedule)
		if err != nil {
			return fmt.Errorf("could not process services schedule [%s]: %w", stageDomesticAccessorialPrice.ServicesSchedule, err)
		}

		var perUnitCentsService int
		perUnitCentsService, err = priceToCents(stageDomesticAccessorialPrice.PricePerUnit)
		if err != nil {
			return fmt.Errorf("could not process price per unit [%s]: %w", stageDomesticAccessorialPrice.PricePerUnit, err)
		}

		serviceProvidedFound := false
		for _, service := range services {
			serviceCode := service.serviceCode
			serviceProvided := service.serviceProvided

			if stageDomesticAccessorialPrice.ServiceProvided == serviceProvided {
				serviceProvidedFound = true
				serviceID, found := gre.serviceToIDMap[serviceCode]
				if !found {
					return fmt.Errorf("missing service [%s] in map of services", serviceCode)
				}

				domesticAccessorial := models.ReDomesticAccessorialPrice{
					ContractID:       gre.ContractID,
					ServicesSchedule: servicesSchedule,
					ServiceID:        serviceID,
					PerUnitCents:     unit.Cents(perUnitCentsService),
				}

				verrs, dbErr := appCtx.DB().ValidateAndSave(&domesticAccessorial)
				if dbErr != nil {
					return fmt.Errorf("error saving ReDomesticAccessorialPrices: %+v with error: %w", domesticAccessorial, dbErr)
				}
				if verrs.HasAny() {
					return fmt.Errorf("error saving ReDomesticAccessorialPrices: %+v with validation errors: %w", domesticAccessorial, verrs)
				}
			}
		}
		if !serviceProvidedFound {
			return fmt.Errorf("service [%s] not found", stageDomesticAccessorialPrice.ServiceProvided)
		}
	}

	return nil
}
