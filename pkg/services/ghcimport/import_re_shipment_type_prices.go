package ghcimport

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREShipmentTypePrices(dbTx *pop.Connection) error {
	//tab 5a) Access. and Add. Prices
	var domesticIntlAddlPrices []models.StageDomesticInternationalAdditionalPrices
	err := dbTx.All(&domesticIntlAddlPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic international additional prices: %w", err)
	}

	shipmentCodePositionInSlice := 0

	//loop through the domestic international additional prices data and store in db
	for _, stageDomesticIntlAddlPrices := range domesticIntlAddlPrices {
		//shipment codes stored in the re_services table
		shipmentCodes := []string{"DMHF", "DBTF", "DBHF", "IBTF", "IBHF", "DNPKF", "INPKF"}

		//pass each shipmentcode  to the serviceToIDMap to get its serviceID
		service := shipmentCodes[shipmentCodePositionInSlice]
		serviceID, found := gre.serviceToIDMap[service]
		if !found {
			return fmt.Errorf("missing service [%s] in map of services", service)
		}

		factor, err := strconv.ParseFloat(stageDomesticIntlAddlPrices.Factor, 64)
		if err != nil {
			return fmt.Errorf("could not process factor [%s]: %w", stageDomesticIntlAddlPrices.Factor, err)
		}
		shipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: gre.contractID,
			ServiceID:  serviceID,
			Factor:     factor,
		}
		shipmentCodePositionInSlice++

		if stageDomesticIntlAddlPrices.Market == "CONUS" {
			shipmentTypePrice.Market = models.MarketConus
		} else if stageDomesticIntlAddlPrices.Market == "OCONUS" {
			shipmentTypePrice.Market = models.MarketOconus
		} else {
			return fmt.Errorf("market [%s] is not a valid market", stageDomesticIntlAddlPrices.Market)
		}

		verrs, dbErr := dbTx.ValidateAndSave(&shipmentTypePrice)
		if dbErr != nil {
			return fmt.Errorf("error saving ReShipmentTypePrices: %+v with error: %w", shipmentTypePrice, dbErr)
		}
		if verrs.HasAny() {
			return fmt.Errorf("error saving ReShipmentTypePrices: %+v with validation errors: %w", shipmentTypePrice, verrs)
		}
	}

	return nil
}
