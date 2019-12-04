package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakePaymentRequest creates a single PaymentRequest and associated set relationships
func MakePaymentRequest(db *pop.Connection, assertions Assertions) models.PaymentRequest {

	// Create new PaymentRequest if not provided
	// ID is required because it must be populated for Eager saving to work.
	moveTaskOrder := assertions.PaymentRequest.MoveTaskOrder
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMoveTaskOrder(db, assertions)
	}

	paymentRequest := models.PaymentRequest{
		MoveTaskOrder:   moveTaskOrder,
		MoveTaskOrderID: moveTaskOrder.ID,
		IsFinal:         false,
		RejectionReason: "Not good enough",
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentRequest, assertions.PaymentRequest)

	mustCreate(db, &paymentRequest)

	return paymentRequest
}

// MakeDefaultPaymentRequest makes an PaymentRequest with default values
func MakeDefaultPaymentRequest(db *pop.Connection) models.PaymentRequest {
	return MakePaymentRequest(db, Assertions{})
}
