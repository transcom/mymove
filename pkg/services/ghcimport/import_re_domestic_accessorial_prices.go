package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticAccessorialPrices(dbTx *pop.Connection) error {
	//tab 5a) Access. and Add. Prices
	var domesticAccessorialPrices []models.StageDomesticMoveAccessorialPrice
	err := dbTx.All(&domesticAccessorialPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic accessorial prices: %w", err)
	}

	services := []struct {
		serviceCode     string
		serviceProvided string
	}{
		{"DCRT", "Crating (per cubic ft.)"},
		{"DCRTSA", "Crating (per cubic ft.)"},
		{"DUCRT", "Uncrating (per cubic ft.)"},
		{"DDSHUT", "Shuttle Service (per cwt)"},
		{"DOSHUT", "Shuttle Service (per cwt)"},
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
					return fmt.Errorf("missing service [%s] in map of services", service)
				}

				domesticAccessorial := models.ReDomesticAccessorialPrice{
					ContractID:       gre.ContractID,
					ServicesSchedule: servicesSchedule,
					ServiceID:        serviceID,
					PerUnitCents:     unit.Cents(perUnitCentsService),
				}

				verrs, dbErr := dbTx.ValidateAndSave(&domesticAccessorial)
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
