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
	case models.ReServiceCodeIHPK, models.ReServiceCodeINPK:
		if p.ServiceItem.ReService.Code == models.ReServiceCodeINPK {
			// If this is an iNTS iHHG packing scenario, we need to make sure to
			// use the IHPK packing for reIntlOtherPrice fetching because INPK pricing doesn't exist
			ihpkService, err := models.FetchReServiceByCode(appCtx.DB(), models.ReServiceCodeIHPK)
			if err != nil {
				return "", fmt.Errorf("error fetching rate area id for shipment ID: %s and service ID %s: %s", shipmentID, serviceID, err)
			}
			serviceID = ihpkService.ID
		}
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

	default:
		return "", fmt.Errorf("unsupported service code to retrieve service item param PerUnitCents")
	}
}
