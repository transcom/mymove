package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// NumberDaysSITLookup does lookup of the number of SIT days a Move Task Orders MTO Shipment can bill for
type NumberDaysSITLookup struct {
	MTOShipment models.MTOShipment
}

func (s NumberDaysSITLookup) lookup(keyData *ServiceItemParamKeyData) (string, error) {
	moveTaskOrderSITPaymentServiceItems, err := fetchMoveTaskOrderSITPaymentServiceItems(keyData.db, s.MTOShipment)
	if err != nil {
		return "", err
	}

	mtoShipmentSITPaymentServiceItems, err := fetchMTOShipmentSITPaymentServiceItems(keyData.db, s.MTOShipment)
	if err != nil {
		return "", err
	}

	remainingMoveTaskOrderSITDays, err := calculateRemainingMoveTaskOrderSITDays(moveTaskOrderSITPaymentServiceItems)
	if err != nil {
		return "", err
	}

	mtoServiceItemSITDays, err := calculateMTOServiceItemSITDays(mtoShipmentSITPaymentServiceItems, keyData.MTOServiceItem)
	if err != nil {
		return "", err
	}

	if notEnoughRemainingMoveTaskOrderSITDays(remainingMoveTaskOrderSITDays, mtoServiceItemSITDays) {
		return strconv.Itoa(remainingMoveTaskOrderSITDays), nil
	}

	return strconv.Itoa(mtoServiceItemSITDays), nil
}

func fetchMoveTaskOrderSITPaymentServiceItems(db *pop.Connection, mtoShipment models.MTOShipment) (models.PaymentServiceItems, error) {
	moveTaskOrderSITPaymentServiceItems := models.PaymentServiceItems{}

	err := db.Q().
		Join("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Join("re_services", "re_services.id = mto_service_items.re_service_id").
		Eager("MTOServiceItem.ReService", "PaymentServiceItemParams.ServiceItemParamKey").
		Where("payment_service_items.status IN ($1, $2, $3, $4) AND mto_service_items.move_id = ($5) AND re_services.code IN ($6, $7, $8, $9)", models.PaymentServiceItemStatusRequested, models.PaymentServiceItemStatusApproved, models.PaymentServiceItemStatusSentToGex, models.PaymentServiceItemStatusPaid, mtoShipment.MoveTaskOrderID, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT).
		All(&moveTaskOrderSITPaymentServiceItems)
	if err != nil {
		return models.PaymentServiceItems{}, err
	}

	return moveTaskOrderSITPaymentServiceItems, nil
}

func fetchMTOShipmentSITPaymentServiceItems(db *pop.Connection, mtoShipment models.MTOShipment) (models.PaymentServiceItems, error) {
	mtoShipmentSITPaymentServiceItems := models.PaymentServiceItems{}

	err := db.Q().
		Join("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Join("re_services", "re_services.id = mto_service_items.re_service_id").
		Eager("MTOServiceItem.ReService", "PaymentServiceItemParams.ServiceItemParamKey").
		Where("payment_service_items.status IN ($1, $2, $3, $4) AND mto_service_items.mto_shipment_id = ($5) AND re_services.code IN ($6, $7, $8, $9)", models.PaymentServiceItemStatusRequested, models.PaymentServiceItemStatusApproved, models.PaymentServiceItemStatusSentToGex, models.PaymentServiceItemStatusPaid, mtoShipment.ID, models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT).
		All(&mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return models.PaymentServiceItems{}, err
	}

	return mtoShipmentSITPaymentServiceItems, nil
}

func calculateRemainingMoveTaskOrderSITDays(moveTaskOrderSITPaymentServiceItems models.PaymentServiceItems) (int, error) {
	remainingMoveTaskOrderSITDays := 90

	for _, moveTaskOrderSITPaymentServiceItem := range moveTaskOrderSITPaymentServiceItems {
		if isFirstDaySIT(moveTaskOrderSITPaymentServiceItem.MTOServiceItem) {
			remainingMoveTaskOrderSITDays--
		} else {
			paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(moveTaskOrderSITPaymentServiceItem)
			if err != nil {
				return 0, err
			}
			remainingMoveTaskOrderSITDays = remainingMoveTaskOrderSITDays - paymentServiceItemSITDays
		}
	}

	return remainingMoveTaskOrderSITDays, nil
}

func calculateMTOServiceItemSITDays(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems, mtoServiceItem models.MTOServiceItem) (int, error) {
	mtoServiceItemSITDays := 0

	originMTOShipmentEntryDate, destinationMTOShipmentEntryDate, err := fetchAndVerifyMTOShipmentSITDates(mtoShipmentSITPaymentServiceItems, mtoServiceItem)
	if err != nil {
		return 0, err
	}

	submittedMTOShipmentSITDays, err := calculateSubmittedMTOShipmentSITDays(mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return 0, err
	}

	if isDOASIT(mtoServiceItem) {
		originMTOShipmentDepartureDate := mtoServiceItem.SITDepartureDate
		isNotFullPaymentPeriod, daysSITAvailableForPayment := isNotFullPaymentPeriod(originMTOShipmentEntryDate)

		if originMTOShipmentDepartureDate != nil {
			mtoShipmentSITDuration := int(originMTOShipmentDepartureDate.Sub(originMTOShipmentEntryDate).Hours() / 24)
			mtoServiceItemSITDays = mtoShipmentSITDuration - submittedMTOShipmentSITDays
			return mtoServiceItemSITDays, nil
		} else if originMTOShipmentDepartureDate == nil && isNotFullPaymentPeriod {
			mtoServiceItemSITDays = daysSITAvailableForPayment
			return mtoServiceItemSITDays, nil
		} else {
			mtoServiceItemSITDays = 29
			return mtoServiceItemSITDays, nil
		}
	} else if isDDASIT(mtoServiceItem) {
		destinationMTOShipmentDepartureDate := mtoServiceItem.SITDepartureDate
		isNotFullPaymentPeriod, daysSITAvailableForPayment := isNotFullPaymentPeriod(originMTOShipmentEntryDate)

		if destinationMTOShipmentDepartureDate != nil {
			mtoShipmentSITDuration := int(destinationMTOShipmentDepartureDate.Sub(destinationMTOShipmentEntryDate).Hours() / 24)
			mtoServiceItemSITDays = mtoShipmentSITDuration - submittedMTOShipmentSITDays
			return mtoServiceItemSITDays, nil
		} else if destinationMTOShipmentDepartureDate == nil && isNotFullPaymentPeriod {
			mtoServiceItemSITDays = daysSITAvailableForPayment
			return mtoServiceItemSITDays, nil
		} else {
			mtoServiceItemSITDays = 29
			return mtoServiceItemSITDays, nil
		}
	}

	return mtoServiceItemSITDays, nil
}

func notEnoughRemainingMoveTaskOrderSITDays(remainingMoveTaskOrderSITDays int, mtoServiceItemSITDays int) bool {
	if remainingMoveTaskOrderSITDays < mtoServiceItemSITDays {
		return true
	}

	return false
}

func fetchAndVerifyMTOShipmentSITDates(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems, mtoServiceItem models.MTOServiceItem) (time.Time, time.Time, error) {
	var originSITEntryDate time.Time
	var destinationSITEntryDate time.Time

	for _, shipmentSITPaymentServiceItem := range mtoShipmentSITPaymentServiceItems {
		if isDomesticOrigin(shipmentSITPaymentServiceItem.MTOServiceItem) {

			if isDomesticOrigin(mtoServiceItem) && shipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
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

			if isDomesticDestination(mtoServiceItem) && shipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil {
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

	if isDomesticOrigin(mtoServiceItem) && !originSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate != nil && originSITEntryDate != *mtoServiceItem.SITEntryDate {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has an Origin MTO Service Item with a different SIT Entry Date of %v", mtoServiceItem.MTOShipment.ID, originSITEntryDate)
	} else if isDomesticDestination(mtoServiceItem) && !destinationSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate != nil && destinationSITEntryDate != *mtoServiceItem.SITEntryDate {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has a Destination MTO Service Item with a different SIT Entry Date of %v", mtoServiceItem.MTOShipment.ID, originSITEntryDate)
	}

	if isDomesticOrigin(mtoServiceItem) && originSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v does not have an Origin MTO Service Item with a SIT Entry Date", mtoServiceItem.MTOShipment.ID)
	} else if isDomesticOrigin(mtoServiceItem) && originSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate != nil {
		sitEntryDate := mtoServiceItem.SITEntryDate
		originSITEntryDate = *sitEntryDate
	} else if isDomesticDestination(mtoServiceItem) && destinationSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v does not have a Destination MTO Service Item with a SIT Entry Date", mtoServiceItem.MTOShipment.ID)
	} else if isDomesticDestination(mtoServiceItem) && destinationSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate != nil {
		sitEntryDate := mtoServiceItem.SITEntryDate
		destinationSITEntryDate = *sitEntryDate
	}

	return originSITEntryDate, destinationSITEntryDate, nil
}

func calculateSubmittedMTOShipmentSITDays(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems) (int, error) {
	submittedMTOShipmentSITDays := 0

	for _, mtoShipmentSITPaymentServiceItem := range mtoShipmentSITPaymentServiceItems {
		if isFirstDaySIT(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
			submittedMTOShipmentSITDays++
		} else {
			paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(mtoShipmentSITPaymentServiceItem)
			if err != nil {
				return 0, err
			}
			submittedMTOShipmentSITDays = submittedMTOShipmentSITDays + paymentServiceItemSITDays
		}
	}

	return submittedMTOShipmentSITDays, nil
}

func fetchNumberDaysSITParamValue(paymentServiceItem models.PaymentServiceItem) (int, error) {
	for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
		if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameNumberDaysSIT {

			if paymentServiceItemParam.ServiceItemParamKey.Type != models.ServiceItemParamTypeInteger {
				return 0, fmt.Errorf("trying to convert %s to an int, but param is of type %s", models.ServiceItemParamNameNumberDaysSIT, paymentServiceItemParam.ServiceItemParamKey.Type)
			}

			numberDaysSITParamValue, err := strconv.Atoi(paymentServiceItemParam.Value)
			if err != nil {
				return 0, fmt.Errorf("could not convert value %s to an int: %w", paymentServiceItemParam.Value, err)
			}

			return numberDaysSITParamValue, err
		}
	}

	return 0, nil
}

func isNotFullPaymentPeriod(mtoShipmentEntryDate time.Time) (bool, int) {
	todaysDate := time.Now()
	daysSITAvailableForPayment := int(todaysDate.Sub(mtoShipmentEntryDate).Hours() / 24)

	if daysSITAvailableForPayment < 29 {
		return true, daysSITAvailableForPayment
	}

	return false, 29
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

func isDomesticOrigin(mtoServiceItem models.MTOServiceItem) bool {
	if isDOFSIT(mtoServiceItem) || isDOASIT(mtoServiceItem) {
		return true
	}

	return false
}

func isDomesticDestination(mtoServiceItem models.MTOServiceItem) bool {
	if isDDFSIT(mtoServiceItem) || isDDASIT(mtoServiceItem) {
		return true
	}

	return false
}

func isFirstDaySIT(mtoServiceItem models.MTOServiceItem) bool {
	if isDOFSIT(mtoServiceItem) || isDDFSIT(mtoServiceItem) {
		return true
	}

	return false
}
