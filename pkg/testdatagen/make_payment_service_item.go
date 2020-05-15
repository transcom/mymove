package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
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

	paymentServiceItem := models.PaymentServiceItem{
		PaymentRequest:   paymentRequest,
		PaymentRequestID: paymentRequest.ID,
		MTOServiceItem:   mtoServiceItem,
		MTOServiceItemID: mtoServiceItem.ID,
		Status:           models.PaymentServiceItemStatusRequested,
		RequestedAt:      time.Now(),
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentServiceItem, assertions.PaymentServiceItem)

	mustCreate(db, &paymentServiceItem)

	return paymentServiceItem
}

// MakeDefaultPaymentServiceItem makes a PaymentServiceItem with default values
func MakeDefaultPaymentServiceItem(db *pop.Connection) models.PaymentServiceItem {
	return MakePaymentServiceItem(db, Assertions{})
}
