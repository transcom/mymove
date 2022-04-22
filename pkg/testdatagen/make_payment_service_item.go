package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakePaymentServiceItem creates a single PaymentServiceItem and associated relationships
func MakePaymentServiceItem(db *pop.Connection, assertions Assertions) models.PaymentServiceItem {
	paymentRequest := assertions.PaymentRequest
	if isZeroUUID(paymentRequest.ID) {
		paymentRequest = MakePaymentRequest(db, assertions)
	}

	mtoServiceItem := assertions.MTOServiceItem
	if isZeroUUID(mtoServiceItem.ID) {
		mtoServiceItem = MakeMTOServiceItem(db, assertions)
	}

	var cents = unit.Cents(888)
	paymentServiceItem := models.PaymentServiceItem{
		PaymentRequest:   paymentRequest,
		PaymentRequestID: paymentRequest.ID,
		MTOServiceItem:   mtoServiceItem,
		MTOServiceItemID: mtoServiceItem.ID,
		PriceCents:       &cents,
		Status:           models.PaymentServiceItemStatusRequested,
		RequestedAt:      time.Now(),
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentServiceItem, assertions.PaymentServiceItem)

	mustCreate(db, &paymentServiceItem, assertions.Stub)

	return paymentServiceItem
}

// MakeDefaultPaymentServiceItem makes a PaymentServiceItem with default values
func MakeDefaultPaymentServiceItem(db *pop.Connection) models.PaymentServiceItem {
	return MakePaymentServiceItem(db, Assertions{})
}
