package invoice

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v5"

	ediResponse997 "github.com/transcom/mymove/pkg/edi/edi997"
	"github.com/transcom/mymove/pkg/services"
)

type edi997Processor struct {
	db     *pop.Connection
	logger Logger
}

// NewEDI997Processor returns a new EDI997 processor
func NewEDI997Processor(db *pop.Connection,
	logger Logger) services.EDI997Processor {

	return &edi997Processor{
		db:     db,
		logger: logger,
	}
}

//ProcessEDI997 parses an EDI 997 response and updates the payment request status
func (e *edi997Processor) ProcessEDI997(stringEDI997 string) (ediResponse997.EDI, error) {
	edi997 := ediResponse997.EDI{}
	err := edi997.Parse(stringEDI997)
	if err != nil {
		// e.logger.Error()
		return edi997, err
	}
	var paymentRequestNumber = "2"

	// validate header messages
	isa := edi997.InterchangeControlEnvelope.ISA
	icn := isa.InterchangeControlNumber
	if icn < 1 || icn > 1000000000 {
		return edi997, fmt.Errorf("Invalid InterchangeControlNumber: %d for PaymentRequestNumber: %s", icn, paymentRequestNumber)
	}
	ackRequested := isa.AcknowledgementRequested
	if ackRequested != 0 && ackRequested != 1 {
		return edi997, fmt.Errorf("Invalid AcknowledgementRequested field: %d for PaymentRequestNumber: %s", ackRequested, paymentRequestNumber)
	}
	usageIndicator := strings.TrimSpace(isa.UsageIndicator)
	if usageIndicator != "T" && usageIndicator != "P" {
		return edi997, fmt.Errorf("Invalid UsageIndicator field: %s for PaymentRequestNumber: %s", usageIndicator, paymentRequestNumber)
	}
	gs := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS
	functionalIDCode := strings.TrimSpace(gs.FunctionalIdentifierCode)
	if functionalIDCode != "SI" {
		return edi997, fmt.Errorf("Invalid FunctionalIdentifierCode field: %s for PaymentRequestNumber: %s", functionalIDCode, paymentRequestNumber)
	}
	gcn := gs.GroupControlNumber
	if gcn < 1 || gcn > 1000000000 {
		return edi997, fmt.Errorf("Invalid GroupControlNumber: %d for PaymentRequestNumber: %s", gcn, paymentRequestNumber)
	}
	return edi997, nil

	// var ediPaymentRequest models.PaymentRequest
	// err := e.db.Q().
	// 	Where("status = ?", models.PaymentRequestStatusReviewed).
	// 	All(&reviewedPaymentRequests)
	// if err != nil {
	// 	return reviewedPaymentRequests, services.NewQueryError("PaymentRequests", err, fmt.Sprintf("Could not find reviewed payment requests: %s", err))
	// }
	// return reviewedPaymentRequests, err
}
