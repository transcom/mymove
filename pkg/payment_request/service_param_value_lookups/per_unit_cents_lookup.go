package serviceparamvaluelookups

import (
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
	serviceID := p.ServiceItem.ReServiceID
	contractID := s.ContractID

	switch p.ServiceItem.ReService.Code {
	case models.ReServiceCodeIHPK:
		// IHPK: Need rate area ID for the pickup address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.PickupAddressID, serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(*p.MTOShipment.RequestedPickupDate)
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

	case models.ReServiceCodeIHUPK:
		// IHUPK: Need rate area ID for the destination address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.PickupAddressID, serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(*p.MTOShipment.RequestedPickupDate)
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

	case models.ReServiceCodeISLH:
		// ISLH: Need rate area IDs for origin and destination
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.PickupAddressID, serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for origin address for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, serviceID, err)
		}
		destRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.MTOShipment.DestinationAddressID, serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for destination address for shipment ID: %s and service ID %s: %s", p.MTOShipment.ID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(*p.MTOShipment.RequestedPickupDate)
		var reIntlPrice models.ReIntlPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("origin_rate_area_id = ?", originRateAreaID).
			Where("destination_rate_area_id = ?", destRateAreaID).
			First(&reIntlPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching ISLH per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, originRateAreaID: %s, and destRateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, originRateAreaID, destRateAreaID, err)
		}
		return reIntlPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	default:
		return "", fmt.Errorf("unsupported service code to retrieve service item param PerUnitCents")
	}
}
