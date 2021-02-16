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

const hoursInADay float64 = 24
const fullBillingPeriod int = 29

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

	billableMTOServiceItemSITDays, err := calculateBillableMTOServiceItemSITDays(mtoShipmentSITPaymentServiceItems, keyData.MTOServiceItem)
	if err != nil {
		return "", err
	}

	if remainingMoveTaskOrderSITDays <= 0 {
		return "", fmt.Errorf("Move Task Order %v has 0 remaining SIT Days", s.MTOShipment.MoveTaskOrderID)
	} else if notEnoughRemainingMoveTaskOrderSITDays(remainingMoveTaskOrderSITDays, billableMTOServiceItemSITDays) {
		return strconv.Itoa(remainingMoveTaskOrderSITDays), nil
	}

	if billableMTOServiceItemSITDays <= 0 {
		return "", fmt.Errorf("MTO Service Item %v has 0 billable SIT Days", keyData.MTOServiceItemID)
	}

	return strconv.Itoa(billableMTOServiceItemSITDays), nil
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
		} else if isAdditionalDaysSIT(moveTaskOrderSITPaymentServiceItem.MTOServiceItem) {
			paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(moveTaskOrderSITPaymentServiceItem)
			if err != nil {
				return 0, err
			}
			remainingMoveTaskOrderSITDays = remainingMoveTaskOrderSITDays - paymentServiceItemSITDays
		}
	}

	return remainingMoveTaskOrderSITDays, nil
}

func calculateBillableMTOServiceItemSITDays(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems, mtoServiceItem models.MTOServiceItem) (int, error) {
	billableMTOServiceItemSITDays := 0
	originMTOShipmentEntryDate, destinationMTOShipmentEntryDate, err := fetchAndVerifyMTOShipmentSITDates(mtoShipmentSITPaymentServiceItems, mtoServiceItem)
	if err != nil {
		return 0, err
	}

	originSubmittedMTOShipmentSITDays, destinationSubmittedMTOShipmentSITDays, err := calculateSubmittedMTOShipmentSITDays(mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return 0, err
	}

	// submittedMTOShipmentSITDays

	if isDOASIT(mtoServiceItem) {
		originMTOShipmentDepartureDate := mtoServiceItem.SITDepartureDate
		isNotFullBillingPeriod, sitDaysAvailableForBilling := isNotFullBillingPeriod(originMTOShipmentEntryDate, originSubmittedMTOShipmentSITDays)

		if originMTOShipmentDepartureDate != nil {
			mtoShipmentSITDuration := int(originMTOShipmentDepartureDate.Sub(originMTOShipmentEntryDate).Hours() / hoursInADay)
			billableMTOServiceItemSITDays = mtoShipmentSITDuration - originSubmittedMTOShipmentSITDays
			if billableMTOServiceItemSITDays < 0 {
				billableMTOServiceItemSITDays = 0
			}
			return billableMTOServiceItemSITDays, nil
		} else if originMTOShipmentDepartureDate == nil && isNotFullBillingPeriod {
			billableMTOServiceItemSITDays = sitDaysAvailableForBilling
			if billableMTOServiceItemSITDays < 0 {
				billableMTOServiceItemSITDays = 0
			}
			return 0, fmt.Errorf("MTO Shipment %v has no departure date and only %v billable SIT day(s)", mtoServiceItem.MTOShipmentID, billableMTOServiceItemSITDays)
		} else {
			return fullBillingPeriod, nil
		}
	} else if isDDASIT(mtoServiceItem) {
		destinationMTOShipmentDepartureDate := mtoServiceItem.SITDepartureDate
		isNotFullBillingPeriod, sitDaysAvailableForBilling := isNotFullBillingPeriod(destinationMTOShipmentEntryDate, destinationSubmittedMTOShipmentSITDays)

		if destinationMTOShipmentDepartureDate != nil {
			mtoShipmentSITDuration := int(destinationMTOShipmentDepartureDate.Sub(destinationMTOShipmentEntryDate).Hours() / hoursInADay)
			billableMTOServiceItemSITDays = mtoShipmentSITDuration - destinationSubmittedMTOShipmentSITDays
			if billableMTOServiceItemSITDays < 0 {
				billableMTOServiceItemSITDays = 0
			}
			return billableMTOServiceItemSITDays, nil
		} else if destinationMTOShipmentDepartureDate == nil && isNotFullBillingPeriod {
			billableMTOServiceItemSITDays = sitDaysAvailableForBilling
			if billableMTOServiceItemSITDays < 0 {
				billableMTOServiceItemSITDays = 0
			}
			return 0, fmt.Errorf("MTO Shipment %v has no departure date and only %v billable SIT day(s)", mtoServiceItem.MTOShipmentID, billableMTOServiceItemSITDays)
		} else {
			return fullBillingPeriod, nil
		}
	}

	return billableMTOServiceItemSITDays, nil
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

	for _, mtoShipmentSITPaymentServiceItem := range mtoShipmentSITPaymentServiceItems {
		if isDomesticOrigin(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
			if isDomesticOrigin(mtoServiceItem) {
				if mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil && mtoShipmentSITPaymentServiceItem.MTOServiceItem.ID != mtoServiceItem.ID {
					return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has an Origin MTO Service Item %v with a SIT Departure Date of %v", mtoShipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID, mtoShipmentSITPaymentServiceItem.ID, mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate)
				} else if mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate != nil {
					sitEntryDate := mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate
					if originSITEntryDate.IsZero() {
						originSITEntryDate = *sitEntryDate
					} else if originSITEntryDate != *sitEntryDate {
						return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v has multiple Origin MTO Service Items with different SIT Entry Dates", mtoShipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID)
					}
				}
			}
		} else if isDomesticDestination(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
			if isDomesticDestination(mtoServiceItem) {
				if mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate != nil && mtoShipmentSITPaymentServiceItem.MTOServiceItem.ID != mtoServiceItem.ID {
					return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has a Destination MTO Service Item %v with a SIT Departure Date of %v", mtoShipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID, mtoShipmentSITPaymentServiceItem.ID, mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITDepartureDate)
				} else if mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate != nil {
					sitEntryDate := mtoShipmentSITPaymentServiceItem.MTOServiceItem.SITEntryDate
					if destinationSITEntryDate.IsZero() {
						destinationSITEntryDate = *sitEntryDate
					} else if destinationSITEntryDate != *sitEntryDate {
						return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v has multiple Destination MTO Service Items with different SIT Entry Dates", mtoShipmentSITPaymentServiceItem.MTOServiceItem.MTOShipment.ID)
					}
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

func calculateSubmittedMTOShipmentSITDays(mtoShipmentSITPaymentServiceItems models.PaymentServiceItems) (int, int, error) {
	originSubmittedMTOShipmentSITDays := 0
	destinationSubmittedMTOShipmentSITDays := 0

	for _, mtoShipmentSITPaymentServiceItem := range mtoShipmentSITPaymentServiceItems {
		if isDomesticOrigin(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
			if isFirstDaySIT(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
				originSubmittedMTOShipmentSITDays++
			} else if isAdditionalDaysSIT(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
				paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(mtoShipmentSITPaymentServiceItem)
				if err != nil {
					return 0, 0, err
				}
				originSubmittedMTOShipmentSITDays = originSubmittedMTOShipmentSITDays + paymentServiceItemSITDays
			}
		} else if isDomesticDestination(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
			if isFirstDaySIT(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
				destinationSubmittedMTOShipmentSITDays++
			} else if isAdditionalDaysSIT(mtoShipmentSITPaymentServiceItem.MTOServiceItem) {
				paymentServiceItemSITDays, err := fetchNumberDaysSITParamValue(mtoShipmentSITPaymentServiceItem)
				if err != nil {
					return 0, 0, err
				}
				destinationSubmittedMTOShipmentSITDays = destinationSubmittedMTOShipmentSITDays + paymentServiceItemSITDays
			}
		}
	}

	return originSubmittedMTOShipmentSITDays, destinationSubmittedMTOShipmentSITDays, nil
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

func isNotFullBillingPeriod(mtoShipmentEntryDate time.Time, submittedMTOShipmentSITDays int) (bool, int) {
	todaysDate := time.Now()
	mtoShipmentSITDuration := int(todaysDate.Sub(mtoShipmentEntryDate).Hours() / hoursInADay)
	sitDaysAvailableForBilling := mtoShipmentSITDuration - submittedMTOShipmentSITDays

	if sitDaysAvailableForBilling < fullBillingPeriod {
		return true, sitDaysAvailableForBilling
	}

	return false, 0
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

func isAdditionalDaysSIT(mtoServiceItem models.MTOServiceItem) bool {
	if isDOASIT(mtoServiceItem) || isDDASIT(mtoServiceItem) {
		return true
	}

	return false
}
