package ghcimport

import (
	"fmt"
	"strconv"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREShipmentTypePrices(dbTx *pop.Connection) error {
	//tab 5a) Access. and Add. Prices
	var domesticIntlAddlPrices []models.StageDomesticInternationalAdditionalPrice
	err := dbTx.All(&domesticIntlAddlPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic international additional prices: %w", err)
	}

	var serviceToCodeMap = map[string]string{
		//concatenating market with shipment type so that keys in  map are unique
		"CONUS:Mobile Homes":            "DMHF",
		"CONUS:Tow Away Boat Service":   "DBTF",
		"OCONUS:Tow Away Boat Service":  "IBTF",
		"CONUS:Haul Away Boat Service":  "DBHF",
		"OCONUS:Haul Away Boat Service": "IBHF",
		"CONUS:NTS Packing Factor":      "DNPKF",
		"OCONUS:NTS Packing Factor":     "INPKF",
	}

	//loop through the domestic international additional prices data and store in db
	for _, stageDomesticIntlAddlPrices := range domesticIntlAddlPrices {
		//shipment codes stored in the re_services table
		factor, err := strconv.ParseFloat(stageDomesticIntlAddlPrices.Factor, 64)
		if err != nil {
			return fmt.Errorf("could not process factor [%s]: %w", stageDomesticIntlAddlPrices.Factor, err)
		}

		market, err := getMarket(stageDomesticIntlAddlPrices.Market)
		if err != nil {
			return fmt.Errorf("could not process market [%s]: %w", stageDomesticIntlAddlPrices.Market, err)
		}

		shipmentTypeFound := false
		for shipmentType, serviceCode := range serviceToCodeMap {
			if shipmentType == stageDomesticIntlAddlPrices.Market+":"+stageDomesticIntlAddlPrices.ShipmentType {
				shipmentTypeFound = true
				serviceID, found := gre.serviceToIDMap[serviceCode]
				if !found {
					return fmt.Errorf("missing service [%s] in map of services", serviceCode)
				}

				shipmentTypePrice := models.ReShipmentTypePrice{
					ContractID: gre.ContractID,
					Market:     market,
					ServiceID:  serviceID,
					Factor:     factor,
				}

				verrs, dbErr := dbTx.ValidateAndSave(&shipmentTypePrice)
				if dbErr != nil {
					return fmt.Errorf("error saving ReShipmentTypePrices: %+v with error: %w", shipmentTypePrice, dbErr)
				}
				if verrs.HasAny() {
					return fmt.Errorf("error saving ReShipmentTypePrices: %+v with validation errors: %w", shipmentTypePrice, verrs)
				}
			}
		}
		if !shipmentTypeFound {
			return fmt.Errorf("shipment type [%s] not found", stageDomesticIntlAddlPrices.ShipmentType)
		}
	}

	return nil
}
