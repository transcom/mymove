package factory

import (
	"fmt"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildPaymentRequest(db *pop.Connection, customs []Customization, traits []Trait) models.PaymentRequest {
	customs = setupCustomizations(customs, traits)

	// Find address customization and extract the custom address
	var cPaymentRequest models.PaymentRequest
	if result := findValidCustomization(customs, PaymentRequest); result != nil {
		cPaymentRequest = result.Model.(models.PaymentRequest)
		if result.LinkOnly {
			return cPaymentRequest
		}
	}

	move := BuildMove(db, customs, traits)

	sequenceNumber := 1
	if cPaymentRequest.SequenceNumber != 0 {
		sequenceNumber = cPaymentRequest.SequenceNumber
	}
	paymentRequestNumber := fmt.Sprintf("%s-%d", *move.ReferenceID, sequenceNumber)

	// Create default PaymentRequest
	paymentRequest := models.PaymentRequest{
		MoveTaskOrder:        move,
		MoveTaskOrderID:      move.ID,
		IsFinal:              false,
		RejectionReason:      nil,
		Status:               models.PaymentRequestStatusPending,
		PaymentRequestNumber: paymentRequestNumber,
		SequenceNumber:       sequenceNumber,
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&paymentRequest, cPaymentRequest)

	// If db is false, it's a stub. No need to create in database.
	if db != nil {
		mustCreate(db, &paymentRequest)
	}

	return paymentRequest
}
