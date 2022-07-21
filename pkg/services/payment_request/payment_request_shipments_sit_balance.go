package paymentrequest

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
				start, err = time.Parse(sitParamDateFormat, paymentServiceItemParam.Value)
			}
		} else if paymentServiceItemParam.ServiceItemParamKey.Key == models.ServiceItemParamNameSITPaymentRequestEnd {
			// remove once the pricer work is done so a 500 server error isn't returned for an unparseable date
			if paymentServiceItemParam.Value != "NOT IMPLEMENTED" {
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

func calculateReviewedSITBalance(paymentServiceItems []models.PaymentServiceItem, shipmentsSITBalances map[string]services.ShipmentPaymentSITBalance) error {
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
				shipmentSITBalance.TotalSITDaysRemaining -= *shipmentSITBalance.PreviouslyBilledDays

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

				if shipment.SITDaysAllowance != nil {
					shipmentSITBalance.TotalSITDaysAuthorized = *shipment.SITDaysAllowance
					shipmentSITBalance.TotalSITDaysRemaining = shipmentSITBalance.TotalSITDaysAuthorized - daysInSIT
				}

				shipmentsSITBalances[shipment.ID.String()] = shipmentSITBalance
			}
		}
	}

	return nil
}

func calculatePendingSITBalance(paymentServiceItems []models.PaymentServiceItem, shipmentsSITBalances map[string]services.ShipmentPaymentSITBalance) error {
	for _, paymentServiceItem := range paymentServiceItems {
		if !isAdditionalDaySIT(paymentServiceItem.MTOServiceItem.ReService.Code) {
			continue
		}

		shipment := paymentServiceItem.MTOServiceItem.MTOShipment

		_, end, err := getStartAndEndParams(paymentServiceItem.PaymentServiceItemParams)
		if err != nil {
			return err
		}

		daysInSIT, err := lookupDaysInSIT(paymentServiceItem.PaymentServiceItemParams)
		if err != nil {
			return err
		}

		if shipmentSITBalance, ok := shipmentsSITBalances[shipment.ID.String()]; ok {
			shipmentSITBalance.PendingSITDaysInvoiced = daysInSIT

			if shipment.SITDaysAllowance != nil {
				shipmentSITBalance.TotalSITDaysRemaining -= shipmentSITBalance.PendingSITDaysInvoiced
				// start counting from the day after the last day in the SIT payment range
				shipmentSITBalance.TotalSITEndDate = end.AddDate(0, 0, shipmentSITBalance.TotalSITDaysRemaining+1)
			}

			// I think this would be accurate for the scenario there were 2 pending payment requests, they would see
			// dates reflective of only their SIT items. I think we would need to do something different if we wanted
			// to show different values for origin and dest SIT service items on the same payment request and shipment
			shipmentSITBalance.PendingBilledEndDate = end
			shipmentsSITBalances[shipment.ID.String()] = shipmentSITBalance
		} else {
			shipmentSITBalance := services.ShipmentPaymentSITBalance{
				ShipmentID:             shipment.ID,
				PendingSITDaysInvoiced: daysInSIT,
				PendingBilledEndDate:   end,
			}

			if shipment.SITDaysAllowance != nil {
				shipmentSITBalance.TotalSITDaysAuthorized = *shipment.SITDaysAllowance
				shipmentSITBalance.TotalSITDaysRemaining = shipmentSITBalance.TotalSITDaysAuthorized - daysInSIT
				// start counting from the day after the last day in the SIT payment range
				shipmentSITBalance.TotalSITEndDate = end.AddDate(0, 0, shipmentSITBalance.TotalSITDaysRemaining+1)
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
	err := appCtx.DB().Eager("PaymentServiceItems.MTOServiceItem.ReService", "PaymentServiceItems.MTOServiceItem.MTOShipment", "PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey").Find(&paymentRequest, paymentRequestID)
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
	var paymentServiceItems []models.PaymentServiceItem
	err = appCtx.DB().Q().EagerPreload("MTOServiceItem", "MTOServiceItem.ReService", "MTOServiceItem.MTOShipment", "PaymentRequest", "PaymentServiceItemParams", "PaymentServiceItemParams.ServiceItemParamKey").
		InnerJoin("payment_requests", "payment_requests.id = payment_service_items.payment_request_id").
		InnerJoin("mto_service_items", "mto_service_items.id = payment_service_items.mto_service_item_id").
		InnerJoin("re_services", "re_services.id = mto_service_items.re_service_id").
		Where("payment_requests.status = ?", models.PaymentRequestStatusReviewed).
		Where("payment_requests.move_id = ?", paymentRequest.MoveTaskOrderID).
		Where("re_services.code IN (?)", string(models.ReServiceCodeDOASIT), string(models.ReServiceCodeDDASIT)).
		Order("payment_requests.created_at asc, mto_service_items.sit_entry_date asc").
		All(&paymentServiceItems)

	if err != nil {
		return nil, err
	}

	shipmentsSITBalances := map[string]services.ShipmentPaymentSITBalance{}

	// first go through the previously billed SIT service items
	err = calculateReviewedSITBalance(paymentServiceItems, shipmentsSITBalances)
	if err != nil {
		return nil, err
	}

	// review the pending SIT service items on the open payment request
	// we may need to change how this works once a pending payment request is reviewed by the TIO because the numbers
	// will look different when looking at the reviewed payment request again.
	err = calculatePendingSITBalance(paymentRequest.PaymentServiceItems, shipmentsSITBalances)
	if err != nil {
		return nil, err
	}

	var sitBalances []services.ShipmentPaymentSITBalance

	for i := range shipmentsSITBalances {
		sitBalances = append(sitBalances, shipmentsSITBalances[i])
	}

	return sitBalances, nil
}
