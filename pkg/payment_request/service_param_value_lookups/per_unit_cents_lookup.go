package serviceparamvaluelookups

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

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
	if serviceID == uuid.Nil {
		reService, err := models.FetchReServiceByCode(appCtx.DB(), p.ServiceItem.ReService.Code)
		if err != nil {
			return "", fmt.Errorf("error fetching ReService Code %s: %w", p.ServiceItem.ReService.Code, err)
		}
		serviceID = reService.ID
	}
	contractID := s.ContractID
	var shipmentID uuid.UUID
	var pickupAddressID uuid.UUID
	var destinationAddressID uuid.UUID
	var moveDate time.Time
	// HHG shipment
	if p.MTOShipment.ShipmentType != models.MTOShipmentTypePPM {
		shipmentID = p.MTOShipment.ID
		if p.MTOShipment.RequestedPickupDate != nil {
			moveDate = *p.MTOShipment.RequestedPickupDate
		} else {
			return "", fmt.Errorf("requested pickup date is required for shipment with id: %s", shipmentID)
		}
		if p.MTOShipment.PickupAddressID != nil {
			pickupAddressID = *p.MTOShipment.PickupAddressID
		} else {
			return "", fmt.Errorf("pickup address is required for shipment with id: %s", shipmentID)
		}
		if p.MTOShipment.DestinationAddressID != nil {
			destinationAddressID = *p.MTOShipment.DestinationAddressID
		} else {
			return "", fmt.Errorf("destination address is required for shipment with id: %s", shipmentID)
		}
	} else { // PPM shipment
		shipmentID = p.MTOShipment.PPMShipment.ID
		if p.MTOShipment.ActualPickupDate != nil {
			moveDate = *p.MTOShipment.ActualPickupDate
		} else if p.MTOShipment.RequestedPickupDate != nil {
			moveDate = *p.MTOShipment.RequestedPickupDate
		} else {
			return "", fmt.Errorf("actual move date is required for PPM shipment with id: %s", shipmentID)
		}

		if p.MTOShipment.PPMShipment.PickupAddressID != nil {
			pickupAddressID = *p.MTOShipment.PPMShipment.PickupAddressID
		} else {
			return "", fmt.Errorf("pickup address is required for PPM shipment with id: %s", shipmentID)
		}

		if p.MTOShipment.PPMShipment.DestinationAddressID != nil {
			destinationAddressID = *p.MTOShipment.PPMShipment.DestinationAddressID
		} else {
			return "", fmt.Errorf("destination address is required for PPM shipment with id: %s", shipmentID)
		}
	}

	switch p.ServiceItem.ReService.Code {
	case models.ReServiceCodeIHPK:
		// IHPK: Need rate area ID for the pickup address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
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
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
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
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for origin address for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		destRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for destination address for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
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

	case models.ReServiceCodeIOFSIT:
		// IOFSIT: Need rate area ID for origin
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for origin address for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", originRateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IOFSIT per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, originRateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, originRateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeIOASIT:
		// IOASIT: Need rate area ID for origin
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for origin address for shipment ID: %s, service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", originRateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IOASIT per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, originRateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, originRateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeIDFSIT:
		// IDFSIT: Need rate area ID for destination
		destRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for destination address for shipment ID: %s, service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", destRateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IDFSIT per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, destRateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, destRateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeIDASIT:
		// IDASIT: Need rate area ID for destination
		destRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for destination address for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", destRateAreaID).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IDASIT per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, destRateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, destRateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeIUBPK:
		// IUBPK: Need rate area ID for the pickup address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
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
			return "", fmt.Errorf("error fetching IUBPK per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, and rateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, rateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeIUBUPK:
		// IUBUPK: Need rate area ID for the destination address
		rateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
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
			return "", fmt.Errorf("error fetching IUBUPK per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, and rateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, rateAreaID, err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeUBP:
		// UBP: Need rate area IDs for origin and destination
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), pickupAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for origin address for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
		}
		destRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), destinationAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for destination address for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
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
			return "", fmt.Errorf("error fetching UBP per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, originRateAreaID: %s, and destRateAreaID: %s: %s", contractID, serviceID, isPeakPeriod, originRateAreaID, destRateAreaID, err)
		}
		return reIntlPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	case models.ReServiceCodeIOPSIT:
		// IOPSIT: Need rate area ID for origin
		if p.ServiceItem.SITOriginHHGActualAddressID == nil {
			return "", fmt.Errorf("ServiceItem.SITOriginHHGActualAddressID is not set. serviceID %s", serviceID)
		}
		originRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.ServiceItem.SITOriginHHGActualAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for SIT origin address for shipment ID: %s, service ID %s, addressID: %s: %s", shipmentID, serviceID, p.ServiceItem.SITOriginHHGActualAddressID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
		if p.ServiceItem.SITDeliveryMiles == nil {
			return "", fmt.Errorf("ServiceItem.SITDeliveryMiles is not set. serviceID %s", serviceID)
		}
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", originRateAreaID).
			Where("is_less_50_miles = ?", (*p.ServiceItem.SITDeliveryMiles <= 50)).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IOPSIT per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, originRateAreaID: %s, SITDeliveryMiles: %d: %s", contractID, serviceID, isPeakPeriod, originRateAreaID, int(*p.ServiceItem.SITDeliveryMiles), err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil
	case models.ReServiceCodeIDDSIT:
		// IDDSIT: Need rate area ID for destination
		if p.ServiceItem.SITDestinationFinalAddressID == nil {
			return "", fmt.Errorf("ServiceItem.SITDestinationFinalAddressID is not set. serviceID %s", serviceID)
		}
		destinationRateAreaID, err := models.FetchRateAreaID(appCtx.DB(), *p.ServiceItem.SITDestinationFinalAddressID, &serviceID, contractID)
		if err != nil {
			return "", fmt.Errorf("error fetching rate area id for SIT destination address for shipment ID: %s, service ID %s, addressID: %s: %s", shipmentID, serviceID, p.ServiceItem.SITDestinationFinalAddressID, err)
		}
		isPeakPeriod := ghcrateengine.IsPeakPeriod(moveDate)
		if p.ServiceItem.SITDeliveryMiles == nil {
			return "", fmt.Errorf("ServiceItem.SITDeliveryMiles is not set. serviceID %s", serviceID)
		}
		var reIntlOtherPrice models.ReIntlOtherPrice
		err = appCtx.DB().Q().
			Where("contract_id = ?", contractID).
			Where("service_id = ?", serviceID).
			Where("is_peak_period = ?", isPeakPeriod).
			Where("rate_area_id = ?", destinationRateAreaID).
			Where("is_less_50_miles = ?", (*p.ServiceItem.SITDeliveryMiles <= 50)).
			First(&reIntlOtherPrice)
		if err != nil {
			return "", fmt.Errorf("error fetching IDDSIT per unit cents for contractID: %s, serviceID %s, isPeakPeriod: %t, destRateAreaID: %s, SITDeliveryMiles: %d: %s", contractID, serviceID, isPeakPeriod, destinationRateAreaID, int(*p.ServiceItem.SITDeliveryMiles), err)
		}
		return reIntlOtherPrice.PerUnitCents.ToMillicents().ToCents().String(), nil

	default:
		return "", fmt.Errorf("unsupported service code to retrieve service item param PerUnitCents")
	}
}
