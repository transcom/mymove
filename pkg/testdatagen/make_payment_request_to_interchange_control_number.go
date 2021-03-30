package testdatagen

import (
	"fmt"
	"log"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
)

// MakePaymentRequestToInterchangeControlNumber creates a single PaymentRequest and PaymentRequestToInterchangeControlNumber
func MakePaymentRequestToInterchangeControlNumber(db *pop.Connection, assertions Assertions) models.PaymentRequestToInterchangeControlNumber {
	paymentRequestID := assertions.PaymentRequestToInterchangeControlNumber.PaymentRequestID
	if isZeroUUID(paymentRequestID) {
		paymentRequest := MakePaymentRequest(db, assertions)
		paymentRequestID = paymentRequest.ID
	}

	icnSequencer, err := sequence.NewRandomSequencer(ediinvoice.ICNRandomMin, ediinvoice.ICNRandomMax)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered creating random sequencer: %v", err))
	}

	icn, err := icnSequencer.NextVal()
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered getting random interchange control number: %v", err))
	}

	pr2icn := models.PaymentRequestToInterchangeControlNumber{
		PaymentRequestID:         paymentRequestID,
		InterchangeControlNumber: int(icn),
	}

	mergeModels(&pr2icn, assertions.PaymentRequestToInterchangeControlNumber)

	mustCreate(db, &pr2icn, assertions.Stub)

	return pr2icn
}
