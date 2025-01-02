package serviceparamvaluelookups

import (
	"database/sql"
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
)

// PerUnitCents does lookup on the per unit cents value associated with a service item
type PerUnitCentsLookup struct {
	ServiceItem models.MTOServiceItem
	MTOShipment models.MTOShipment
}

func (p PerUnitCentsLookup) lookup(appCtx appcontext.AppContext, s *ServiceItemParamKeyData) (string, error) {
	var isPeakPeriod bool
	serviceID := p.ServiceItem.ReServiceID
	contractID := s.ContractID
	if p.ServiceItem.ReService.Code == models.ReServiceCodeIHPK {
		// IHPK we need the rate area id for the pickup address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.PickupAddressID, p.ServiceItem.ReServiceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, p.ServiceItem.ReServiceID, err)
		}
		isPeakPeriod = ghcrateengine.IsPeakPeriod(*p.MTOShipment.RequestedPickupDate)
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", rateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IHPK per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, and rateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, rateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil
	}
	if p.ServiceItem.ReService.Code == models.ReServiceCodeIHUPK {
		// IHUPK we need the rate area id for the destination address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.PickupAddressID, p.ServiceItem.ReServiceID, contractID)
		if err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, p.ServiceItem.ReServiceID, err)
		}
		isPeakPeriod = ghcrateengine.IsPeakPeriod(*p.MTOShipment.RequestedPickupDate)
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", rateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IHUPK per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, and rateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, rateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil
	} else if p.ServiceItem.ReService.Code == models.ReServiceCodeISLH {
		// IHUPK we need the rate area id for the destination address
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.PickupAddressID, p.ServiceItem.ReServiceID, contractID)
		if err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("error fetching rate area id for origina address for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, p.ServiceItem.ReServiceID, err)
		}
		destRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.DestinationAddressID, p.ServiceItem.ReServiceID, contractID)
		if err != nil && err != sql.ErrNoRows {
			return "", fmt.Errorf("error fetching rate area id for destination address for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, p.ServiceItem.ReServiceID, err)
		}
		isPeakPeriod = ghcrateengine.IsPeakPeriod(*p.MTOShipment.RequestedPickupDate)
		var reIntlPrice models.ReIntlPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("origin_rate_area_id = ?", originRateAreaID).
			Where("destination_rate_area_id = ?", destRateAreaID).
			First(&reIntlPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching ISLH per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, originRateAreaID: %s, and destRateAreaid: %s: %s", contractID, serviceID, isPeakPeriod, originRateAreaID, destRateAreaID, err)
		}
		return reIntlPrice.PerUnitCents.ToMillicents().ToCents().String(), nil
	} else {
		return "", fmt.Errorf("unsupported service code to retrieve service item param PerUnitCents")
	}
}
