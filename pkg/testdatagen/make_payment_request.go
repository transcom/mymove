package testdatagen

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakePaymentRequest creates a single PaymentRequest and associated set relationships
func MakePaymentRequest(db *pop.Connection, assertions Assertions) (models.PaymentRequest, error) {
	// Create new PaymentRequest if not provided
	// ID is required because it must be populated for Eager saving to work.
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		var err error
		moveTaskOrder, err = makeMove(db, assertions)
		if err != nil {
			return models.PaymentRequest{}, err
		}
	}

	paymentRequestNumber := assertions.PaymentRequest.PaymentRequestNumber
	sequenceNumber := assertions.PaymentRequest.SequenceNumber
	if paymentRequestNumber == "" {
		if sequenceNumber == 0 {
			sequenceNumber = 1
		}
		paymentRequestNumber = fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, sequenceNumber)
	}

	paymentRequest := models.PaymentRequest{
		CreatedAt:            assertions.PaymentRequest.CreatedAt,
		MoveTaskOrder:        moveTaskOrder,
		MoveTaskOrderID:      moveTaskOrder.ID,
		IsFinal:              false,
		RejectionReason:      nil,
		Status:               models.PaymentRequestStatusPending,
		PaymentRequestNumber: paymentRequestNumber,
		SequenceNumber:       sequenceNumber,
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentRequest, assertions.PaymentRequest)

	mustCreate(db, &paymentRequest, assertions.Stub)

	return paymentRequest, nil
}
