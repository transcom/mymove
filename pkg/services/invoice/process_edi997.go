package invoice

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v5"

	ediResponse997 "github.com/transcom/mymove/pkg/edi/edi997"
	edi "github.com/transcom/mymove/pkg/edi/invoice"
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
		return edi997, err
	}
	var paymentRequestNumber = "2"
	errString := ""
	ediHeader := edi.InvoiceResponseHeader{}
	ediHeader.ISA = edi997.InterchangeControlEnvelope.ISA
	ediHeader.GS = edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS
	// validate header fields
	err = ValidateEDIHeader(ediHeader, paymentRequestNumber, e.logger)
	if err != nil {
		errString += err.Error()
	}

	// validate 997 specific fields
	for _, functionalGroup := range edi997.InterchangeControlEnvelope.FunctionalGroups {
		for _, transactionSet := range functionalGroup.TransactionSets {

			functionalGroupResponse := transactionSet.FunctionalGroupResponse
			transactionSetResponses := functionalGroupResponse.TransactionSetResponses
			ak1 := functionalGroupResponse.AK1
			functionalIDCode := strings.TrimSpace(ak1.FunctionalIdentifierCode)
			if functionalIDCode != "SI" {
				msg := fmt.Sprintf("Invalid FunctionalIdentifierCode in AK1 segment: %s  for PaymentRequestNumber: %s", functionalIDCode, paymentRequestNumber)
				e.logger.Error(msg)
				errString += msg + "\n"
			}
			for _, transactionSetResponse := range transactionSetResponses {
				ak2 := transactionSetResponse.AK2
				transactionSetIdentifierCode := strings.TrimSpace(ak2.TransactionSetIdentifierCode)
				if transactionSetIdentifierCode != "858" {
					msg := fmt.Sprintf("Invalid TransactionSetIdentifierCode in AK2 segment: %s for PaymentRequestNumber: %s", transactionSetIdentifierCode, paymentRequestNumber)
					e.logger.Error(msg)
					errString += msg + "\n"
				}

				ak5 := transactionSetResponse.AK5
				transactionSetAcknowledgmentCode := strings.TrimSpace(ak5.TransactionSetAcknowledgmentCode)
				if transactionSetAcknowledgmentCode != "A" {
					msg := fmt.Sprintf("Invalid TransactionSetAcknowledgmentCode in AK5 segment: %s for PaymentRequestNumber: %s", transactionSetAcknowledgmentCode, paymentRequestNumber)
					e.logger.Error(msg)
					errString += msg + "\n"
				}
			}
		}
	}
	if errString != "" {
		return edi997, fmt.Errorf(errString)
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

// ValidateEDIHeader validates an EDI header for a 997 or 824
func ValidateEDIHeader(ediHeader edi.InvoiceResponseHeader, paymentRequestNumber string, logger Logger) error {
	returnErr := ""
	isa := ediHeader.ISA
	icn := isa.InterchangeControlNumber

	if icn < 1 || icn > 1000000000 {
		msg := fmt.Sprintf("Invalid InterchangeControlNumber in ISA segment: %d for PaymentRequestNumber: %s", icn, paymentRequestNumber)
		logger.Error(msg)
		returnErr += msg + "\n"
	}
	ackRequested := isa.AcknowledgementRequested
	if ackRequested != 0 && ackRequested != 1 {
		msg := fmt.Sprintf("Invalid AcknowledgementRequested in ISA segment: %d for PaymentRequestNumber: %s", ackRequested, paymentRequestNumber)
		logger.Error(msg)
		returnErr += msg + "\n"
	}
	usageIndicator := strings.TrimSpace(isa.UsageIndicator)
	if usageIndicator != "T" && usageIndicator != "P" {
		msg := fmt.Sprintf("Invalid UsageIndicator in ISA segment %s for PaymentRequestNumber: %s", usageIndicator, paymentRequestNumber)
		logger.Error(msg)
		returnErr += msg + "\n"
	}
	gs := ediHeader.GS
	functionalIDCode := strings.TrimSpace(gs.FunctionalIdentifierCode)
	if functionalIDCode != "SI" {
		msg := fmt.Sprintf("Invalid FunctionalIdentifierCode in GS segment: %s for PaymentRequestNumber: %s", functionalIDCode, paymentRequestNumber)
		logger.Error(msg)
		returnErr += msg + "\n"
	}
	gcn := gs.GroupControlNumber
	if gcn < 1 || gcn > 1000000000 {
		msg := fmt.Sprintf("Invalid GroupControlNumber in GS segment: %d  for PaymentRequestNumber: %s", gcn, paymentRequestNumber)
		logger.Error(msg)
		returnErr += msg + "\n"
	}
	if returnErr != "" {
		return fmt.Errorf(returnErr)
	}
	return nil
}
