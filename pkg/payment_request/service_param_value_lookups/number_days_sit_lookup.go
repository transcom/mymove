package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// TODO: Write NumberDaysSITLookup description
type NumberDaysSITLookup struct {
	MTOShipment models.MTOShipment
}

func (s NumberDaysSITLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	moveTaskOrderSITPaymentServiceItems, err := fetchMoveTaskOrderSITPaymentServiceItems(keyData.db, s.MTOShipment)
	if err != nil {
		return "", err
	}

	remainingMoveTaskOrderSITDays, err := calculateRemainingMoveTaskOrderSITDays(moveTaskOrderSITPaymentServiceItems)
	if err != nil {
		return "", err
	}

	mtoShipmentSITPaymentServiceItems, err := fetchMTOShipmentSITPaymentServiceItems(keyData.db, s.MTOShipment)
	if err != nil {
		return "", err
	}

	mtoServiceItemSITDays, err := calculateMTOServiceItemSITDays(mtoShipmentSITPaymentServiceItems, keyData.MTOServiceItem)
	if err != nil {
		return "", err
	}

	if mtoServiceItemSITDays > remainingMoveTaskOrderSITDays {
		return strconv.Itoa(remainingMoveTaskOrderSITDays), nil
	}

	return strconv.Itoa(mtoServiceItemSITDays), nil
}

func fetchMoveTaskOrderSITPaymentServiceItems(db *pop.Connection, mtoShipment models.MTOShipment) (models.PaymentServiceItems, error) {
	moveTaskOrderSITPaymentServiceItems := models.PaymentServiceItems{}

	err := db.Q().
		Join("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Join("re_services", "re_services.id = mto_service_items.re_service_id").
		Eager("MTOServiceItem.ReService").
		Eager("PaymentServiceItemParams.ServiceItemParamKey").
		Where("payment_service_items.status IN ($1, $2, $3, $4) AND mto_service_items.move_id = ($5) AND re_services.code IN ($6, $7, $8, $9)", models.PaymentServiceItemStatusRequested, models.PaymentServiceItemStatusApproved, models.PaymentServiceItemStatusSentToGex, models.PaymentServiceItemStatusPaid, mtoShipment.MoveTaskOrderID, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT).
		All(&moveTaskOrderSITPaymentServiceItems)
	if err != nil {
		return models.PaymentServiceItems{}, err
	}

	return moveTaskOrderSITPaymentServiceItems, nil
}

func calculateRemainingMoveTaskOrderSITDays(moveTaskOrderSITPaymentServiceItems models.PaymentServiceItems) (int, error) {
	remainingMoveTaskOrderSITDays := 90

	for _, moveTaskOrderSITPaymentServiceItem := range moveTaskOrderSITPaymentServiceItems {
		paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(moveTaskOrderSITPaymentServiceItem)
		remainingMoveTaskOrderSITDays = remainingMoveTaskOrderSITDays - paymentServiceItemSITDays
		if err != nil {
			return 0, err
		}
	}

	return remainingMoveTaskOrderSITDays, nil
}

func fetchMTOShipmentSITPaymentServiceItems(db *pop.Connection, mtoShipment models.MTOShipment) (models.PaymentServiceItems, error) {
	mtoShipmentSITPaymentServiceItems := models.PaymentServiceItems{}

	err := db.Q().
		Join("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Join("re_services", "re_services.id = mto_service_items.re_service_id").
		Eager("MTOServiceItem.ReService").
		Eager("PaymentServiceItemParams.ServiceItemParamKey").
		Where("payment_service_items.status IN ($1, $2, $3, $4) AND mto_service_items.mto_shipment_id = ($5) AND re_services.code IN ($6, $7, $8, $9)", models.PaymentServiceItemStatusRequested, models.PaymentServiceItemStatusApproved, models.PaymentServiceItemStatusSentToGex, models.PaymentServiceItemStatusPaid, mtoShipment.ID, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT).
		All(&mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return models.PaymentServiceItems{}, err
	}

	return mtoShipmentSITPaymentServiceItems, nil
}

func calculateMTOServiceItemSITDays(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems, currentMTOServiceItem models.MTOServiceItem) (int, error) {
	originMTOShipmentEntryDate, destinationMTOShipmentEntryDate, err := fetchAndVerifyMTOShipmentSITDates(mtoShipmentSITPaymentServiceItems, currentMTOServiceItem)
	if err != nil {
		return 0, err
	}

	currentMTOShipmentSITDays, err := calculateCurrentMTOShipmentSITDays(mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return 0, err
	}

	if isDomesticOrigin(currentMTOServiceItem) {
		originMTOShipmentDepartureDate := currentMTOServiceItem.SITDepartureDate

		if isDOFSIT(currentMTOServiceItem) {
			return 1, nil
		} else if isDOASIT(currentMTOServiceItem) && originMTOShipmentDepartureDate != nil {
			totalCurrentMTOShipmentSITDays := int(originMTOShipmentDepartureDate.Sub(originMTOShipmentEntryDate).Hours() / 24)
			return totalCurrentMTOShipmentSITDays - currentMTOShipmentSITDays, nil
		} else if isDOASIT(currentMTOServiceItem) && originMTOShipmentDepartureDate == nil {
			return 29, nil
		}
	} else if isDomesticDestination(currentMTOServiceItem) {
		destinationMTOShipmentDepartureDate := currentMTOServiceItem.SITDepartureDate
		if isDDFSIT(currentMTOServiceItem) {
			return 1, nil
		} else if isDDASIT(currentMTOServiceItem) && destinationMTOShipmentDepartureDate != nil {
			totalCurrentMTOShipmentSITDays := int(destinationMTOShipmentDepartureDate.Sub(destinationMTOShipmentEntryDate).Hours() / 24)
			return totalCurrentMTOShipmentSITDays - currentMTOShipmentSITDays, nil
		} else if isDDASIT(currentMTOServiceItem) && destinationMTOShipmentDepartureDate == nil {
			return 29, nil
		}
	}

	return 0, nil
}

func fetchAndVerifyMTOShipmentSITDates(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems, lookupMTOServiceItem models.MTOServiceItem) (time.Time, time.Time, error) {
	var originSITEntryDate time.Time
	var destinationSITEntryDate time.Time

	for _, shipmentSITPaymentServiceItem := range mtoShipmentSITPaymentServiceItems {
		if isDomesticOrigin(shipmentSITPaymentServiceItem.MTOServiceItem) {

			if isDomesticOrigin(lookupMTOServiceItem) && shipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has an Origin MTO Service Item %v with a SIT Departure Date of %v", shipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID, shipmentSITPaymentServiceItem.ID, shipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate)
			} else if shipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate != nil {
				sitEntryDate := shipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate

				if originSITEntryDate.IsZero() || originSITEntryDate == *sitEntryDate {
					originSITEntryDate = *sitEntryDate
				} else {
					return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v has multiple Origin MTO Service Items with different SIT Entry Dates", shipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID)
				}
			}
		} else if isDomesticDestination(shipmentSITPaymentServiceItem.MTOServiceItem) {

			if isDomesticDestination(lookupMTOServiceItem) && shipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has a Destination MTO Service Item %v with a SIT Departure Date of %v", shipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID, shipmentSITPaymentServiceItem.ID, shipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate)
			} else if shipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate != nil {
				sitEntryDate := shipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate

				if destinationSITEntryDate.IsZero() || destinationSITEntryDate == *sitEntryDate {
					destinationSITEntryDate = *sitEntryDate
				} else {
					return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v has multiple Destination MTO Service Items with different SIT Entry Dates", shipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID)
				}
			}
		}
	}

	if isDomesticOrigin(lookupMTOServiceItem) && !originSITEntryDate.IsZero() && lookupMTOServiceItem.SITEntryDate != nil && originSITEntryDate != *lookupMTOServiceItem.SITEntryDate {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has an Origin MTO Service Item with a different SIT Entry Date of %v", lookupMTOServiceItem.MTOShipment.ID, originSITEntryDate)
	} else if isDomesticDestination(lookupMTOServiceItem) && !destinationSITEntryDate.IsZero() && lookupMTOServiceItem.SITEntryDate != nil && destinationSITEntryDate != *lookupMTOServiceItem.SITEntryDate {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has a Destination MTO Service Item with a different SIT Entry Date of %v", lookupMTOServiceItem.MTOShipment.ID, originSITEntryDate)
	}

	if isDomesticOrigin(lookupMTOServiceItem) && originSITEntryDate.IsZero() && lookupMTOServiceItem.SITEntryDate == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v does not have an Origin MTO Service Item with a SIT Entry Date", lookupMTOServiceItem.MTOShipment.ID)
	} else if isDomesticOrigin(lookupMTOServiceItem) && originSITEntryDate.IsZero() && lookupMTOServiceItem.SITEntryDate != nil {
		sitEntryDate := lookupMTOServiceItem.SITEntryDate
		originSITEntryDate = *sitEntryDate
	} else if isDomesticDestination(lookupMTOServiceItem) && destinationSITEntryDate.IsZero() && lookupMTOServiceItem.SITEntryDate == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v does not have a Destination MTO Service Item with a SIT Entry Date", lookupMTOServiceItem.MTOShipment.ID)
	} else if isDomesticDestination(lookupMTOServiceItem) && destinationSITEntryDate.IsZero() && lookupMTOServiceItem.SITEntryDate != nil {
		sitEntryDate := lookupMTOServiceItem.SITEntryDate
		destinationSITEntryDate = *sitEntryDate
	}

	return originSITEntryDate, destinationSITEntryDate, nil
}

func calculateCurrentMTOShipmentSITDays(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems) (int, error) {
	currentMTOShipmentSITDays := 0

	for _, moveTaskOrderSITPaymentServiceItem := range mtoShipmentSITPaymentServiceItems {
		paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(moveTaskOrderSITPaymentServiceItem)
		if err != nil {
			return 0, err
		}

		currentMTOShipmentSITDays = currentMTOShipmentSITDays + paymentServiceItemSITDays
	}

	return currentMTOShipmentSITDays, nil
}

func fetchNumberDaysSITParamValue(paymentServiceItem models.PaymentServiceItem) (int, error) {
	for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
		if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameNumberDaysSIT {

			if paymentServiceItemParam.ServiceItemParamKey.Type != models.ServiceItemParamTypeInteger {
				return 0, fmt.Errorf("trying to convert %s to an int, but param is of type %s", models.ServiceItemParamNameNumberDaysSIT, paymentServiceItemParam.ServiceItemParamKey.Type)
			}

			paymentServiceItemParamValue, err := strconv.Atoi(paymentServiceItemParam.Value)
			if err != nil {
				return 0, fmt.Errorf("could not convert value %s to an int: %w", paymentServiceItemParam.Value, err)
			}

			return paymentServiceItemParamValue, err
		}
	}

	return 0, nil
}

func isDomesticOrigin(mtoServiceItem models.MTOServiceItem) bool {
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDOFSIT || mtoServiceItem.ReService.Code == models.ReServiceCodeDOASIT {
		return true
	}

	return false
}

func isDOFSIT(mtoServiceItem models.MTOServiceItem) bool {
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDOFSIT {
		return true
	}

	return false
}

func isDOASIT(mtoServiceItem models.MTOServiceItem) bool {
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDOASIT {
		return true
	}

	return false
}

func isDomesticDestination(mtoServiceItem models.MTOServiceItem) bool {
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDDFSIT || mtoServiceItem.ReService.Code == models.ReServiceCodeDDASIT {
		return true
	}

	return false
}

func isDDFSIT(mtoServiceItem models.MTOServiceItem) bool {
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDDFSIT {
		return true
	}

	return false
}

func isDDASIT(mtoServiceItem models.MTOServiceItem) bool {
	if mtoServiceItem.ReService.Code == models.ReServiceCodeDDASIT {
		return true
	}

	return false
}
