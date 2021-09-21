package paymentserviceitem

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type paymentServiceItemUpdater struct {
}

// NewPaymentServiceItemStatusUpdater returns a new updater for payment service item status
func NewPaymentServiceItemStatusUpdater() services.PaymentServiceItemStatusUpdater {
	return &paymentServiceItemUpdater{}
}

func (p *paymentServiceItemUpdater) UpdatePaymentServiceItemStatus(appCtx appcontext.AppContext, paymentServiceItemID uuid.UUID,
	status models.PaymentServiceItemStatus, rejectionReason *string, eTag string) (models.PaymentServiceItem, *validate.Errors, error) {
	// Fetch the existing record
	var paymentServiceItem models.PaymentServiceItem
	err := appCtx.DB().Find(&paymentServiceItem, paymentServiceItemID)
	if err != nil {
		// If we don't find a record let's return something that will cause a 404
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return models.PaymentServiceItem{}, nil, services.NewNotFoundError(paymentServiceItemID,
				"while looking for payment service item")
		}
		// If it's something else let's still return the error variable (err)
		return models.PaymentServiceItem{}, nil, err
	}

	// Check the etag to make sure we're good to update (not stale)
	currentEtag := etag.GenerateEtag(paymentServiceItem.UpdatedAt)
	if currentEtag != eTag {
		return models.PaymentServiceItem{}, nil, services.NewPreconditionFailedError(paymentServiceItem.ID,
			query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	// Update the record
	// If we're denying this thing we want to make sure to update the DeniedAt field and nil out ApprovedAt.
	if status == models.PaymentServiceItemStatusDenied {
		paymentServiceItem.RejectionReason = rejectionReason
		paymentServiceItem.DeniedAt = swag.Time(time.Now())
		paymentServiceItem.ApprovedAt = nil
		paymentServiceItem.Status = status
	}
	// If we're approving this thing then we don't want there to be a rejection reason
	// We also will want to update the ApprovedAt field and nil out the DeniedAt field.
	if status == models.PaymentServiceItemStatusApproved {
		paymentServiceItem.RejectionReason = nil
		paymentServiceItem.DeniedAt = nil
		paymentServiceItem.ApprovedAt = swag.Time(time.Now())
		paymentServiceItem.Status = status
	}

	// Save the record
	verrs, err := appCtx.DB().ValidateAndSave(&paymentServiceItem)
	if err != nil || verrs != nil && verrs.HasAny() {
		return models.PaymentServiceItem{}, verrs, err
	}
	// Return the updated object
	return paymentServiceItem, nil, nil
}
