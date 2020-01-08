package ghcimport

import (
	"fmt"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (gre *GHCRateEngineImporter) importREIntlAccessorialPrices(dbTx *pop.Connection) error {
	//tab 5a) Access. and Add. Prices
	var intlAccessorialPrices []models.StageInternationalMoveAccessorialPrices
	err := dbTx.All(&intlAccessorialPrices)
	if err != nil {
		return fmt.Errorf("could not read staged intl accessorial prices: %w", err)
	}

	//loop through the intl accessorial price data and store in db
	for _, stageIntlAccessorialPrice := range intlAccessorialPrices {
		serviceICRT, foundService := gre.serviceToIDMap["ICRT"]
		if !foundService {
			return fmt.Errorf("missing service DCRT in map of services")
		}

		serviceIUCRT, foundService := gre.serviceToIDMap["IUCRT"]
		if !foundService {
			return fmt.Errorf("missing service DUCRT in map of services")
		}

		serviceIDSHUT, foundService := gre.serviceToIDMap["IDSHUT"]
		if !foundService {
			return fmt.Errorf("missing service DDSHUT in map of services")
		}

		var perUnitCentsService int
		perUnitCentsService, err = priceToCents(stageIntlAccessorialPrice.PricePerUnit)
		if err != nil {
			return fmt.Errorf("could not process per unit price [%s]: %w", stageIntlAccessorialPrice.PricePerUnit, err)
		}

		intlAccessorial := models.ReIntlAccessorialPrice{
			ContractID:   gre.contractID,
			PerUnitCents: unit.Cents(perUnitCentsService),
		}

		if stageIntlAccessorialPrice.ServiceProvided == "Crating (per cubic ft.)" {
			intlAccessorial.ServiceID = serviceICRT
		} else if stageIntlAccessorialPrice.ServiceProvided == "Uncrating (per cubic ft.)" {
			intlAccessorial.ServiceID = serviceIUCRT
		} else if stageIntlAccessorialPrice.ServiceProvided == "Shuttle Service (per cwt)" {
			intlAccessorial.ServiceID = serviceIDSHUT
		}

		if stageIntlAccessorialPrice.Market == "CONUS" {
			intlAccessorial.Market = models.MarketConus
		} else if stageIntlAccessorialPrice.Market == "OCONUS" {
			intlAccessorial.Market = models.MarketOconus
		}

		verrs, dbErr := dbTx.ValidateAndSave(&intlAccessorial)
		if dbErr != nil {
			return fmt.Errorf("error saving ReIntlAccessorialPrices: %+v with error: %w", intlAccessorial, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReIntlAccessorialPrices: %+v with validation errors: %w", intlAccessorial, verrs)
		}
	}

	return nil
}
