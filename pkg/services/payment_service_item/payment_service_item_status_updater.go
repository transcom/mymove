package paymentserviceitem

import (
	"database/sql"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type paymentServiceItemUpdater struct {
}

// NewPaymentServiceItemStatusUpdater returns a new updater for payment service item status
func NewPaymentServiceItemStatusUpdater() services.PaymentServiceItemStatusUpdater {
	return &paymentServiceItemUpdater{}
}

func (p *paymentServiceItemUpdater) UpdatePaymentServiceItemStatus(appCtx appcontext.AppContext, paymentServiceItemID uuid.UUID,
	desiredStatus models.PaymentServiceItemStatus, rejectionReason *string, eTag string) (models.PaymentServiceItem, *validate.Errors, error) {

	// Fetch the existing record
	paymentServiceItem, verrs, err := p.fetchPaymentServiceItem(appCtx, paymentServiceItemID)
	if err != nil || verrs != nil && verrs.HasAny() {
		return models.PaymentServiceItem{}, verrs, err
	}

	// Update the record
	updatedPaymentServiceItem, verrs, err := p.updatePaymentServiceItem(appCtx, paymentServiceItem, desiredStatus,
		rejectionReason, eTag, checkETag(), rejectionRequiresRejectionReason())
	if err != nil || verrs != nil && verrs.HasAny() {
		return models.PaymentServiceItem{}, verrs, err
	}

	// Return the updated object
	return updatedPaymentServiceItem, nil, nil
}

// Fetch the existing service item
func (p *paymentServiceItemUpdater) fetchPaymentServiceItem(appCtx appcontext.AppContext, paymentServiceItemID uuid.UUID) (models.PaymentServiceItem,
	*validate.Errors, error) {
	var paymentServiceItem models.PaymentServiceItem
	err := appCtx.DB().EagerPreload("PaymentRequest").Find(&paymentServiceItem, paymentServiceItemID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			// If we don't find a record let's return something that will cause a 404
			return models.PaymentServiceItem{}, nil, apperror.NewNotFoundError(paymentServiceItemID,
				"while looking for payment service item")
		default:
			return models.PaymentServiceItem{}, nil, apperror.NewQueryError("PaymentServiceItem", err, "")
		}
	}
	return paymentServiceItem, nil, nil
}

// Update the service item based on the requested new status
func (p *paymentServiceItemUpdater) updatePaymentServiceItem(appCtx appcontext.AppContext,
	paymentServiceItem models.PaymentServiceItem, desiredStatus models.PaymentServiceItemStatus,
	rejectionReason *string, eTag string, checks ...validator) (models.PaymentServiceItem, *validate.Errors, error) {

	// Validate the change we're trying to make to the payment service item
	if verr := validatePaymentServiceItem(appCtx, &paymentServiceItem, desiredStatus,
		rejectionReason, eTag, checks...); verr != nil {
		return models.PaymentServiceItem{}, nil, verr
	}

	// If we're denying this thing we want to make sure to update the DeniedAt field and nil out ApprovedAt.
	if desiredStatus == models.PaymentServiceItemStatusDenied {
		paymentServiceItem.RejectionReason = rejectionReason
		paymentServiceItem.DeniedAt = swag.Time(time.Now())
		paymentServiceItem.ApprovedAt = nil
		paymentServiceItem.Status = desiredStatus
	}
	// If we're approving this thing then we don't want there to be a rejection reason
	// We also will want to update the ApprovedAt field and nil out the DeniedAt field.
	if desiredStatus == models.PaymentServiceItemStatusApproved {
		paymentServiceItem.RejectionReason = nil
		paymentServiceItem.DeniedAt = nil
		paymentServiceItem.ApprovedAt = swag.Time(time.Now())
		paymentServiceItem.Status = desiredStatus
	}

	// Save the record
	verrs, err := appCtx.DB().ValidateAndSave(&paymentServiceItem)
	if err != nil || verrs != nil && verrs.HasAny() {
		return models.PaymentServiceItem{}, verrs, err
	}

	return paymentServiceItem, nil, nil
}
