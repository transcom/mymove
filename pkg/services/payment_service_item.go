package services

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PaymentServiceItemStatusUpdater is the exported interface for updating a payment service item
//
//go:generate mockery --name PaymentServiceItemStatusUpdater --disable-version-string
type PaymentServiceItemStatusUpdater interface {
	UpdatePaymentServiceItemStatus(appCtx appcontext.AppContext, paymentServiceItemID uuid.UUID,
		status models.PaymentServiceItemStatus, rejectionReason *string, eTag string) (models.PaymentServiceItem, *validate.Errors, error)
}
