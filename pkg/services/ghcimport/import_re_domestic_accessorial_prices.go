package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREDomesticAccessorialPrices(dbTx *pop.Connection) error {
	//tab 5a) Access. and Add. Prices
	var domesticAccessorialPrices []models.StageDomesticMoveAccessorialPrices
	err := dbTx.All(&domesticAccessorialPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic accessorial prices: %w", err)
	}

	//loop through the domestic accessorial price data and store in db
	for _, stageDomesticAccessorialPrice := range domesticAccessorialPrices {
		serviceDCRT, foundService := gre.serviceToIDMap["DCRT"]
		if !foundService {
			return fmt.Errorf("missing service DCRT in map of services")
		}

		serviceDUCRT, foundService := gre.serviceToIDMap["DUCRT"]
		if !foundService {
			return fmt.Errorf("missing service DUCRT in map of services")
		}

		serviceDDSHUT, foundService := gre.serviceToIDMap["DDSHUT"]
		if !foundService {
			return fmt.Errorf("missing service DDSHUT in map of services")
		}

		servicesSchedule, err := stringToInteger(stageDomesticAccessorialPrice.ServicesSchedule)
		if err != nil {
			return fmt.Errorf("could not process services schedule [%s]: %w", stageDomesticAccessorialPrice.ServicesSchedule, err)
		}

		var perUnitCentsService int
		perUnitCentsService, err = priceToCents(stageDomesticAccessorialPrice.PricePerUnit)
		if err != nil {
			return fmt.Errorf("could not process price per unit [%s]: %w", stageDomesticAccessorialPrice.PricePerUnit, err)
		}

		domesticAccessorial := models.ReDomesticAccessorialPrice{
			ContractID:       gre.contractID,
			ServicesSchedule: servicesSchedule,
			PerUnitCents:     unit.Cents(perUnitCentsService),
		}

		if stageDomesticAccessorialPrice.ServiceProvided == "Crating (per cubic ft.)" {
			domesticAccessorial.ServiceID = serviceDCRT
		} else if stageDomesticAccessorialPrice.ServiceProvided == "Uncrating (per cubic ft.)" {
			domesticAccessorial.ServiceID = serviceDUCRT
		} else if stageDomesticAccessorialPrice.ServiceProvided == "Shuttle Service (per cwt)" {
			domesticAccessorial.ServiceID = serviceDDSHUT
		}

		verrs, dbErr := dbTx.ValidateAndSave(&domesticAccessorial)
		if dbErr != nil {
			return fmt.Errorf("error saving ReDomesticAccessorialPrices: %+v with error: %w", domesticAccessorial, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReDomesticAccessorialPrices: %+v with validation errors: %w", domesticAccessorial, verrs)
		}
	}

	return nil
}
