package ghcimport

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

func (gre *GHCRateEngineImporter) importREShipmentTypePrices(appCtx appcontext.AppContext) error {
	//tab 5a) Access. and Add. Prices
	var domesticIntlAddlPrices []models.StageDomesticInternationalAdditionalPrice
	err := appCtx.DB().All(&domesticIntlAddlPrices)
	if err != nil {
		return fmt.Errorf("could not read staged domestic international additional prices: %w", err)
	}

	var serviceToCodeMap = map[string]models.ReServiceCode{
		//concatenating market with shipment type so that keys in  map are unique
		"CONUS:Mobile Homes":            models.ReServiceCodeDMHF,
		"CONUS:Tow Away Boat Service":   models.ReServiceCodeDBTF,
		"OCONUS:Tow Away Boat Service":  models.ReServiceCodeIBTF,
		"CONUS:Haul Away Boat Service":  models.ReServiceCodeDBHF,
		"OCONUS:Haul Away Boat Service": models.ReServiceCodeIBHF,
		"CONUS:NTS Packing Factor":      models.ReServiceCodeDNPK,
		"OCONUS:NTS Packing Factor":     models.ReServiceCodeINPK,
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

				verrs, dbErr := appCtx.DB().ValidateAndSave(&shipmentTypePrice)
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
