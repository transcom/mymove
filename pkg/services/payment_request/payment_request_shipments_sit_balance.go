package paymentrequest

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	sitstatus "github.com/transcom/mymove/pkg/services/sit_status"
)

var sitParamDateFormat = "2006-01-02"

type paymentRequestShipmentsSITBalance struct {
}

// NewPaymentRequestShipmentsSITBalance constructs a new service for the SIT balances of a payment request's shipments
func NewPaymentRequestShipmentsSITBalance() services.ShipmentsPaymentSITBalance {
	return &paymentRequestShipmentsSITBalance{}
}

func getStartAndEndParams(params models.PaymentServiceItemParams) (start time.Time, end time.Time, err error) {
	for _, paymentServiceItemParam := range params {
		if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameSITPaymentRequestStart {
			// remove once the pricer work is done so a 500 server error isn't returned for an unparseable date
			if paymentServiceItemParam.Value != "NOT IMPLEMENTED" {
				fmt.Println(paymentServiceItemParam.PaymentServiceItemID)
				start, err = time.Parse(sitParamDateFormat, paymentServiceItemParam.Value)
			}
		} else if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameSITPaymentRequestEnd {
			// remove once the pricer work is done so a 500 server error isn't returned for an unparseable date
			if paymentServiceItemParam.Value != "NOT IMPLEMENTED" {
				fmt.Println(paymentServiceItemParam.PaymentServiceItemID)
				end, err = time.Parse(sitParamDateFormat, paymentServiceItemParam.Value)
			}
		}
		if err != nil {
			return start, end, err
		}
	}

	return start, end, err
}

func lookupDaysInSIT(params models.PaymentServiceItemParams) (int, error) {
	for _, paymentServiceItemParam := range params {
		if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameNumberDaysSIT {
			daysInSIT, err := strconv.Atoi(paymentServiceItemParam.Value)
			if err != nil {
				return 0, err
			}
			return daysInSIT, nil
		}
	}

	return 0, nil
}

func isAdditionalDaySIT(code models.ReServiceCode) bool {
	return code == models.ReServiceCodeDOASIT || code == models.ReServiceCodeDDASIT
}

func hasSITServiceItem(paymentServiceItems models.PaymentServiceItems) bool {
	for _, paymentServiceItem := range paymentServiceItems {
		if code := paymentServiceItem.MTOServiceItem.ReService.Code; isAdditionalDaySIT(code) {
			return true
		}
	}

	return false
}

func calculateReviewedSITBalance(appCtx appcontext.AppContext, paymentServiceItems []models.PaymentServiceItem, shipmentsSITBalances map[string]services.ShipmentPaymentSITBalance) error {
	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	for _, paymentServiceItem := range paymentServiceItems {
		// Ignoring potentially rejected SIT service items here
		if paymentServiceItem.Status == models.PaymentServiceItemStatusApproved {
			_, end, err := getStartAndEndParams(paymentServiceItem.PaymentServiceItemParams)
			if err != nil {
				return err
			}

			daysInSIT, err := lookupDaysInSIT(paymentServiceItem.PaymentServiceItemParams)
			if err != nil {
				return err
			}

			shipment := paymentServiceItem.MTOServiceItem.MTOShipment

			if shipmentSITBalance, ok := shipmentsSITBalances[shipment.ID.String()]; ok {
				totalPreviouslyBilledDays := daysInSIT + *shipmentSITBalance.PreviouslyBilledDays
				shipmentSITBalance.PreviouslyBilledDays = &totalPreviouslyBilledDays

				// try to use most recent SIT billed end date
				if shipmentSITBalance.PreviouslyBilledEndDate.Before(end) {
					// If the DaysInSIT is different than the start and end rage should we change this to be the cutoff
					// date?
					shipmentSITBalance.PreviouslyBilledEndDate = &end
				}

				shipmentsSITBalances[shipment.ID.String()] = shipmentSITBalance
			} else {
				shipmentSITBalance := services.ShipmentPaymentSITBalance{
					ShipmentID:              shipment.ID,
					PreviouslyBilledDays:    &daysInSIT,
					PreviouslyBilledEndDate: &end,
				}

				// sort the SIT service items into past, current and future to aid in the upcoming calculations
				shipmentSIT, err := sitstatus.NewShipmentSITStatus().RetrieveShipmentSIT(appCtx, shipment)
				if err != nil {
					return err
				}
				sortedShipmentSIT := sitstatus.SortShipmentSITs(shipmentSIT, today)

				totalSITDaysAuthorized, err := sitstatus.NewShipmentSITStatus().CalculateShipmentSITAllowance(appCtx, shipment)
				if err != nil {
					return err
				}
				totalSITDaysUsed := sitstatus.CalculateTotalDaysInSIT(sortedShipmentSIT, today)
				totalSITDaysRemaining := totalSITDaysAuthorized - totalSITDaysUsed

				shipmentSITBalance.TotalSITDaysAuthorized = totalSITDaysAuthorized
				shipmentSITBalance.TotalSITDaysRemaining = totalSITDaysRemaining

				shipmentsSITBalances[shipment.ID.String()] = shipmentSITBalance
			}
		}
	}

	return nil
}

func calculatePendingSITBalance(appCtx appcontext.AppContext, paymentServiceItems []models.PaymentServiceItem, shipmentsSITBalances map[string]services.ShipmentPaymentSITBalance) error {
	year, month, day := time.Now().Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	for _, paymentServiceItem := range paymentServiceItems {
		if !isAdditionalDaySIT(paymentServiceItem.MTOServiceItem.ReService.Code) {
			continue
		}

		shipment := paymentServiceItem.MTOServiceItem.MTOShipment

		start, end, err := getStartAndEndParams(paymentServiceItem.PaymentServiceItemParams)
		if err != nil {
			return err
		}

		daysInSIT, err := lookupDaysInSIT(paymentServiceItem.PaymentServiceItemParams)
		if err != nil {
			return err
		}
		// sort the SIT service items into past, current and future to aid in the upcoming calculations
		shipmentSIT, err := sitstatus.NewShipmentSITStatus().RetrieveShipmentSIT(appCtx, shipment)
		if err != nil {
			return err
		}
		sortedShipmentSIT := sitstatus.SortShipmentSITs(shipmentSIT, today)

		if shipmentSITBalance, ok := shipmentsSITBalances[shipment.ID.String()]; ok {
			shipmentSITBalance.PendingSITDaysInvoiced = daysInSIT
			shipmentSITBalance.PendingBilledStartDate = start
			// I think this would be accurate for the scenario there were 2 pending payment requests, they would see
			// dates reflective of only their SIT items. I think we would need to do something different if we wanted
			// to show different values for origin and dest SIT service items on the same payment request and shipment
			shipmentSITBalance.PendingBilledEndDate = end

			// Even though these have been set before, we should do these calculations again in order to recalculate the
			// totalSITEndDate using this service item's entry date.
			// Additionally retrieve the latest SIT Departure date from the current SIT if it exists. The first current SIT is chosen as there is not currently support for more than one SIT
			// Per AC under B-20899
			shipmentSIT, err = sitstatus.NewShipmentSITStatus().RetrieveShipmentSIT(appCtx, shipment)
			if err != nil {
				return err
			}
			sortedShipmentSIT = sitstatus.SortShipmentSITs(shipmentSIT, today)
			// Get the latest authorized end date
			if len(sortedShipmentSIT.CurrentSITs) == 0 {
				// No current SIT, get the most recent authorized end date
				shipmentSITBalance.TotalSITEndDate = sortedShipmentSIT.PastSITs[len(sortedShipmentSIT.PastSITs)-1].Summary.SITAuthorizedEndDate
			} else {
				shipmentSITBalance.TotalSITEndDate = sortedShipmentSIT.CurrentSITs[0].Summary.SITAuthorizedEndDate
			}
			shipmentsSITBalances[shipment.ID.String()] = shipmentSITBalance
		} else {
			shipmentSITBalance := services.ShipmentPaymentSITBalance{
				ShipmentID:             shipment.ID,
				PendingSITDaysInvoiced: daysInSIT,
				PendingBilledStartDate: start,
				PendingBilledEndDate:   end,
			}
			shipmentSIT, err = sitstatus.NewShipmentSITStatus().RetrieveShipmentSIT(appCtx, shipment)
			if err != nil {
				return err
			}
			sortedShipmentSIT = sitstatus.SortShipmentSITs(shipmentSIT, today)

			totalSITDaysAuthorized, err := sitstatus.NewShipmentSITStatus().CalculateShipmentSITAllowance(appCtx, shipment)
			if err != nil {
				return err
			}
			totalSITDaysUsed := sitstatus.CalculateTotalDaysInSIT(sortedShipmentSIT, today)
			totalSITDaysRemaining := totalSITDaysAuthorized - totalSITDaysUsed

			// Retrieve the latest SIT Departure date from the current SIT if it exists. The first current SIT is chosen as there is not currently support for more than one SIT
			// Per AC under B-20899

			shipmentSITBalance.TotalSITDaysAuthorized = totalSITDaysAuthorized
			shipmentSITBalance.TotalSITDaysRemaining = totalSITDaysRemaining
			// Get the latest authorized end date
			if len(sortedShipmentSIT.CurrentSITs) == 0 {
				// No current SIT, get the most recent authorized end date
				shipmentSITBalance.TotalSITEndDate = sortedShipmentSIT.PastSITs[len(sortedShipmentSIT.PastSITs)-1].Summary.SITAuthorizedEndDate
			} else {
				shipmentSITBalance.TotalSITEndDate = sortedShipmentSIT.CurrentSITs[0].Summary.SITAuthorizedEndDate
			}

			shipmentsSITBalances[shipment.ID.String()] = shipmentSITBalance
		}
	}

	return nil
}

func (m paymentRequestShipmentsSITBalance) ListShipmentPaymentSITBalance(appCtx appcontext.AppContext, paymentRequestID uuid.UUID) ([]services.ShipmentPaymentSITBalance, error) {
	var paymentRequest models.PaymentRequest
	// Keeping this query simple in case the payment request is not found as opposed to existing but with no SIT service
	// items
	err := appCtx.DB().Eager(
		"PaymentServiceItems.MTOServiceItem.ReService",
		"PaymentServiceItems.MTOServiceItem.MTOShipment",
		"PaymentServiceItems.MTOServiceItem.MTOShipment.MTOServiceItems",
		"PaymentServiceItems.MTOServiceItem.MTOShipment.MTOServiceItems.ReService",
		"PaymentServiceItems.MTOServiceItem.MTOShipment.SITDurationUpdates",
		"PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey").Find(&paymentRequest, paymentRequestID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(paymentRequestID, "no payment request exists with that id")
		default:
			return nil, apperror.NewQueryError("PaymentRequest", err, "")
		}
	}

	// Check for SIT payment service items, if there are none we don't need to return the SIT balance
	if !hasSITServiceItem(paymentRequest.PaymentServiceItems) {
		return nil, nil
	}

	// We already have the current payment request service items, find all previously reviewed payment requests for this
	// move with billed SIT days
	var reviewedPaymentServiceItems []models.PaymentServiceItem
	err = appCtx.DB().Q().Eager("MTOServiceItem",
		"MTOServiceItem.ReService",
		"MTOServiceItem.MTOShipment",
		"MTOServiceItem.MTOShipment.MTOServiceItems",
		"MTOServiceItem.MTOShipment.MTOServiceItems.ReService",
		"MTOServiceItem.MTOShipment.SITDurationUpdates",
		"PaymentRequest",
		"PaymentServiceItemParams",
		"PaymentServiceItemParams.ServiceItemParamKey").
		InnerJoin("payment_requests", "payment_requests.id = payment_service_items.payment_request_id").
		InnerJoin("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		InnerJoin("re_services", "re_services.id = mto_service_items.re_service_id").
		Where("payment_requests.status = ?", models.PaymentRequestStatusReviewed).
		Where("payment_requests.move_id = ?", paymentRequest.MoveTaskOrderID).
		Where("re_services.code IN (?)", string(models.ReServiceCodeDOASIT), string(models.ReServiceCodeDDASIT)).
		Order("payment_requests.created_at asc, mto_service_items.sit_entry_date asc").
		All(&reviewedPaymentServiceItems)

	if err != nil {
		return nil, err
	}

	shipmentsSITBalances := map[string]services.ShipmentPaymentSITBalance{}

	// first go through the previously billed SIT service items
	err = calculateReviewedSITBalance(appCtx, reviewedPaymentServiceItems, shipmentsSITBalances)
	if err != nil {
		return nil, err
	}

	// review the pending SIT service items on the open payment request
	// we may need to change how this works once a pending payment request is reviewed by the TIO because the numbers
	// will look different when looking at the reviewed payment request again.
	err = calculatePendingSITBalance(appCtx, paymentRequest.PaymentServiceItems, shipmentsSITBalances)
	if err != nil {
		return nil, err
	}

	var sitBalances []services.ShipmentPaymentSITBalance

	for i := range shipmentsSITBalances {
		sitBalances = append(sitBalances, shipmentsSITBalances[i])
	}

	return sitBalances, nil
}
