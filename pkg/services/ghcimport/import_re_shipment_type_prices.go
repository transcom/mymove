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

	//loop through the domestic international additional prices data and store in db
	for _, stageDomesticIntlAddlPrices := range domesticIntlAddlPrices {
		shipmentType := stageDomesticIntlAddlPrices.ShipmentType
		shipmentTypeID, found := gre.shipmentTypeToIDMap[shipmentType]
		if !found {
			return fmt.Errorf("could not find shipment type %s in map", stageDomesticIntlAddlPrices.ShipmentType)
		}

		factorHundredths, err := strconv.ParseFloat(stageDomesticIntlAddlPrices.Factor, 64)
		if err != nil {
			return fmt.Errorf("could not process factor [%s]: %w", stageDomesticIntlAddlPrices.Factor, err)
		}
		shipmentTypePrice := models.ReShipmentTypePrice{
			ContractID:       gre.contractID,
			ShipmentTypeID:   shipmentTypeID,
			FactorHundredths: factorHundredths,
		}

		if stageDomesticIntlAddlPrices.Market == "CONUS" {
			shipmentTypePrice.Market = models.MarketConus
		} else if stageDomesticIntlAddlPrices.Market == "OCONUS" {
			shipmentTypePrice.Market = models.MarketOconus
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
