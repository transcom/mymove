package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
)

// NumberDaysSITLookup does lookup of the number of SIT days a Move Task Orders MTO Shipment can bill for
type NumberDaysSITLookup struct {
	MTOShipment models.MTOShipment
}

const hoursInADay float64 = 24

func (s NumberDaysSITLookup) lookup(appCtx appcontext.AppContext, keyData *ServiceItemParamKeyData) (string, error) {
	mtoShipmentSITPaymentServiceItems, err := fetchMTOShipmentSITPaymentServiceItems(appCtx, s.MTOShipment)
	if err != nil {
		return "", err
	}

	_, _, err = fetchAndVerifyMTOShipmentSITDates(mtoShipmentSITPaymentServiceItems, keyData.MTOServiceItem)
	if err != nil {
		return "", err
	}

	currentPaymentServiceItem, priorPaymentServiceItems, err := findCurrentPaymentServiceItem(mtoShipmentSITPaymentServiceItems, keyData.PaymentRequestID, keyData.MTOServiceItemID)
	if err != nil {
		return "", err
	}

	start, end, err := fetchSITStartAndEndDateParamValues(currentPaymentServiceItem)
	if err != nil {
		return "", fmt.Errorf("failed to parse params for PaymentServiceItem %v: %w", currentPaymentServiceItem.ID, err)
	}

	hasOverlappingDate := hasOverlappingSITDates(priorPaymentServiceItems, keyData.MTOServiceItem, start, end)
	if hasOverlappingDate {
		return "", errors.New("new requested SIT dates overlap previously requested dates")
	}

	shipmentSITStatus := mtoshipment.NewShipmentSITStatus()
	totalSITAllowance, err := shipmentSITStatus.CalculateShipmentSITAllowance(appCtx, s.MTOShipment)
	if err != nil {
		return "", err
	}

	remainingShipmentSITDays, err := calculateRemainingSITDays(priorPaymentServiceItems, totalSITAllowance)
	if err != nil {
		return "", err
	}

	if remainingShipmentSITDays <= 0 {
		return "", fmt.Errorf("MTOShipment %v has 0 remaining SIT Days", s.MTOShipment.ID)
	}

	billableShipmentSITDays, err := calculateNumberSITAdditionalDays(currentPaymentServiceItem)
	if err != nil {
		return "", err
	}

	if remainingShipmentSITDays < billableShipmentSITDays {
		return "", fmt.Errorf("only %d additional days in SIT can be billed for MTOShipment %v", remainingShipmentSITDays, s.MTOShipment.ID)
	}
	return strconv.Itoa(billableShipmentSITDays), nil
}

func hasOverlappingSITDates(shipmentSITPaymentServiceItems models.PaymentServiceItems, mtoServiceItem models.MTOServiceItem, sitStart time.Time, sitEnd time.Time) bool {
	for _, paymentServiceItem := range shipmentSITPaymentServiceItems {
		// Check for overlapping requested SIT dates with previously billed additional days SIT service items at the same origin or destination
		if isAdditionalDaysSIT(paymentServiceItem.MTOServiceItem) && paymentServiceItem.MTOServiceItem.ReService.Code == mtoServiceItem.ReService.Code {
			// Get the payment request service item param start and end dates
			start, end, err := fetchSITStartAndEndDateParamValues(paymentServiceItem)
			if err != nil {
				return false // TODO do we want to say non overlapping if we can't parse?
			}
			// Check if the start or end date has already be used for billing.
			// dateInRange() checks inclusively.
			if dateInRange(sitStart, start, end) || dateInRange(sitEnd, start, end) {
				return true
			}
		}
	}

	return false
}

// Check dates inclusively
func dateInRange(check time.Time, start time.Time, end time.Time) bool {
	checkDateOnly := check.Truncate(24 * time.Hour)
	startDateOnly := start.Truncate(24 * time.Hour)
	endDateOnly := end.Truncate(24 * time.Hour)

	// If the check date equals the start or end date, return true
	if checkDateOnly.Equal(startDateOnly) || checkDateOnly.Equal(endDateOnly) {
		return true
	}

	// If the check date is between the start and end date range, return true
	if checkDateOnly.After(startDateOnly) && checkDateOnly.Before(endDateOnly) {
		return true
	}

	return false
}

func fetchMTOShipmentSITPaymentServiceItems(appCtx appcontext.AppContext, mtoShipment models.MTOShipment) (models.PaymentServiceItems, error) {
	mtoShipmentSITPaymentServiceItems := models.PaymentServiceItems{}

	err := appCtx.DB().Q().
		Join("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		Join("re_services", "re_services.id = mto_service_items.re_service_id").
		Join("payment_requests", "payment_requests.id = payment_service_items.payment_request_id").
		Eager("MTOServiceItem.ReService", "PaymentServiceItemParams.ServiceItemParamKey").
		Where("mto_service_items.mto_shipment_id = ($1)", mtoShipment.ID).
		Where("payment_requests.status != $2", models.PaymentRequestStatusDeprecated).
		Where("payment_service_items.status IN ($3, $4, $5, $6)", models.PaymentServiceItemStatusRequested, models.PaymentServiceItemStatusApproved, models.PaymentServiceItemStatusSentToGex, models.PaymentServiceItemStatusPaid).
		Where("re_services.code IN ($7, $8, $9, $10)", models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT, models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT).
		All(&mtoShipmentSITPaymentServiceItems)
	if err != nil {
		return models.PaymentServiceItems{}, err
	}

	return mtoShipmentSITPaymentServiceItems, nil
}

func calculateRemainingSITDays(sitPaymentServiceItems models.PaymentServiceItems, sitDaysAllowance int) (int, error) {
	remainingSITDays := sitDaysAllowance

	for _, sitPaymentServiceItem := range sitPaymentServiceItems {
		if isFirstDaySIT(sitPaymentServiceItem.MTOServiceItem) {
			remainingSITDays--
		} else if isAdditionalDaysSIT(sitPaymentServiceItem.MTOServiceItem) {
			paymentServiceItemSITDays, err := calculateNumberSITAdditionalDays(sitPaymentServiceItem)
			if err != nil {
				return 0, err
			}
			remainingSITDays -= paymentServiceItemSITDays
		}
	}

	return remainingSITDays, nil
}

func calculateNumberSITAdditionalDays(paymentServiceItem models.PaymentServiceItem) (int, error) {
	startDate, endDate, err := fetchSITStartAndEndDateParamValues(paymentServiceItem)
	if err != nil {
		return 0, err
	}

	days := 1 + endDate.Sub(startDate).Hours()/hoursInADay

	return int(days), nil
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

	if isDomesticOrigin(mtoServiceItem) && !originSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate != nil && !originSITEntryDate.Equal(*mtoServiceItem.SITEntryDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("MTO Shipment %v already has an Origin MTO Service Item with a different SIT Entry Date of %v", mtoServiceItem.MTOShipment.ID, originSITEntryDate)
	} else if isDomesticDestination(mtoServiceItem) && !destinationSITEntryDate.IsZero() && mtoServiceItem.SITEntryDate != nil && !destinationSITEntryDate.Equal(*mtoServiceItem.SITEntryDate) {
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

func findCurrentPaymentServiceItem(paymentServiceItems models.PaymentServiceItems, paymentRequestID uuid.UUID, mtoServiceItemID uuid.UUID) (models.PaymentServiceItem, models.PaymentServiceItems, error) {
	currentPaymentServiceItem := models.PaymentServiceItem{}
	priorPaymentServiceItems := models.PaymentServiceItems{}
	found := false
	for _, psi := range paymentServiceItems {
		if psi.PaymentRequestID == paymentRequestID && psi.MTOServiceItemID == mtoServiceItemID {
			if found {
				return models.PaymentServiceItem{}, models.PaymentServiceItems{}, fmt.Errorf("multiple PaymentServiceItems for MTOServiceItem %v found within the same PaymentRequest", mtoServiceItemID)
			}
			currentPaymentServiceItem = psi
			found = true
		} else {
			priorPaymentServiceItems = append(priorPaymentServiceItems, psi)
		}
	}
	if !found {
		return models.PaymentServiceItem{}, models.PaymentServiceItems{}, fmt.Errorf("failed to find a PaymentServiceItem for MTOServiceItem %v in PaymentRequest %v", mtoServiceItemID, paymentRequestID)
	}

	return currentPaymentServiceItem, priorPaymentServiceItems, nil
}

func fetchSITStartAndEndDateParamValues(paymentServiceItem models.PaymentServiceItem) (time.Time, time.Time, error) {
	start := time.Time{}
	end := time.Time{}

	for _, paymentServiceItemParam := range paymentServiceItem.PaymentServiceItemParams {
		var err error
		if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameSITPaymentRequestStart {
			start, err = time.Parse(ghcrateengine.DateParamFormat, paymentServiceItemParam.Value)
			if err != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("failed to parse SITPaymentRequestStart as a date: %w", err)
			}
		}
		if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameSITPaymentRequestEnd {
			end, err = time.Parse(ghcrateengine.DateParamFormat, paymentServiceItemParam.Value)
			if err != nil {
				return time.Time{}, time.Time{}, fmt.Errorf("failed to parse SITPaymentRequestEnd as a date: %w", err)
			}
		}
	}

	return start, end, nil
}

func isDOFSIT(mtoServiceItem models.MTOServiceItem) bool {
	return mtoServiceItem.ReService.Code == models.ReServiceCodeDOFSIT
}

func isDOASIT(mtoServiceItem models.MTOServiceItem) bool {
	return mtoServiceItem.ReService.Code == models.ReServiceCodeDOASIT
}

func isDDFSIT(mtoServiceItem models.MTOServiceItem) bool {
	return mtoServiceItem.ReService.Code == models.ReServiceCodeDDFSIT
}

func isDDASIT(mtoServiceItem models.MTOServiceItem) bool {
	return mtoServiceItem.ReService.Code == models.ReServiceCodeDDASIT
}

func isDomesticOrigin(mtoServiceItem models.MTOServiceItem) bool {
	return isDOFSIT(mtoServiceItem) || isDOASIT(mtoServiceItem)
}

func isDomesticDestination(mtoServiceItem models.MTOServiceItem) bool {
	return isDDFSIT(mtoServiceItem) || isDDASIT(mtoServiceItem)
}

func isFirstDaySIT(mtoServiceItem models.MTOServiceItem) bool {
	return isDOFSIT(mtoServiceItem) || isDDFSIT(mtoServiceItem)
}

func isAdditionalDaysSIT(mtoServiceItem models.MTOServiceItem) bool {
	return isDOASIT(mtoServiceItem) || isDDASIT(mtoServiceItem)
}
