package invoice

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	ediResponse997 "github.com/transcom/mymove/pkg/edi/edi997"
	edi "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
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
	errString := ""
	ediHeader := edi.InvoiceResponseHeader{}
	ediHeader.ISA = edi997.InterchangeControlEnvelope.ISA
	ediHeader.GS = edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS

	// Find the PaymentRequestID that matches the ICN
	icn := ediHeader.ISA.InterchangeControlNumber
	var paymentRequest models.PaymentRequest
	err = e.db.Q().
		Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
		Where("payment_request_to_interchange_control_numbers.interchange_control_number = ?", int(icn)).
		First(&paymentRequest)

	// validate header fields
	err = ValidateEDIHeader(ediHeader, paymentRequest.ID)
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
				msg := fmt.Sprintf("Invalid FunctionalIdentifierCode in AK1 segment: %s  for PaymentRequestID: %s", functionalIDCode, paymentRequest.ID)
				errString += msg + "\n"
			}
			for _, transactionSetResponse := range transactionSetResponses {
				ak2 := transactionSetResponse.AK2
				transactionSetIdentifierCode := strings.TrimSpace(ak2.TransactionSetIdentifierCode)
				if transactionSetIdentifierCode != "858" {
					msg := fmt.Sprintf("Invalid TransactionSetIdentifierCode in AK2 segment: %s for PaymentRequestID: %s", transactionSetIdentifierCode, paymentRequest.ID)
					errString += msg + "\n"
				}

				ak5 := transactionSetResponse.AK5
				transactionSetAcknowledgmentCode := strings.TrimSpace(ak5.TransactionSetAcknowledgmentCode)
				if transactionSetAcknowledgmentCode != "A" {
					msg := fmt.Sprintf("Invalid TransactionSetAcknowledgmentCode in AK5 segment: %s for PaymentRequestID: %s", transactionSetAcknowledgmentCode, paymentRequest.ID)
					errString += msg + "\n"
				}
			}
		}
	}
	if errString != "" {
		e.logger.Error(errString)
		return edi997, fmt.Errorf(errString)
	}

	var transactionError error
	transactionError = e.db.Transaction(func(tx *pop.Connection) error {
		paymentRequest.Status = models.PaymentRequestStatusReceivedByGex
		err = e.db.Update(&paymentRequest)
		if err != nil {
			e.logger.Error("failure updating payment request", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		return nil
	})

	if transactionError != nil {
		return edi997, transactionError
	}

	return edi997, nil
}

// ValidateEDIHeader validates an EDI header for a 997 or 824
func ValidateEDIHeader(ediHeader edi.InvoiceResponseHeader, paymentRequestID uuid.UUID) error {
	returnErr := ""
	isa := ediHeader.ISA
	icn := isa.InterchangeControlNumber

	if icn < 1 || icn > 1000000000 {
		msg := fmt.Sprintf("Invalid InterchangeControlNumber in ISA segment: %d for PaymentRequestID: %s", icn, paymentRequestID)
		returnErr += msg + "\n"
	}
	ackRequested := isa.AcknowledgementRequested
	if ackRequested != 0 && ackRequested != 1 {
		msg := fmt.Sprintf("Invalid AcknowledgementRequested in ISA segment: %d for PaymentRequestID: %s", ackRequested, paymentRequestID)
		returnErr += msg + "\n"
	}
	usageIndicator := strings.TrimSpace(isa.UsageIndicator)
	if usageIndicator != "T" && usageIndicator != "P" {
		msg := fmt.Sprintf("Invalid UsageIndicator in ISA segment %s for PaymentRequestID: %s", usageIndicator, paymentRequestID)
		returnErr += msg + "\n"
	}
	gs := ediHeader.GS
	functionalIDCode := strings.TrimSpace(gs.FunctionalIdentifierCode)
	if functionalIDCode != "SI" {
		msg := fmt.Sprintf("Invalid FunctionalIdentifierCode in GS segment: %s for PaymentRequestID: %s", functionalIDCode, paymentRequestID)
		returnErr += msg + "\n"
	}
	gcn := gs.GroupControlNumber
	if gcn < 1 || gcn > 1000000000 {
		msg := fmt.Sprintf("Invalid GroupControlNumber in GS segment: %d  for PaymentRequestID: %s", gcn, paymentRequestID)
		returnErr += msg + "\n"
	}
	if returnErr != "" {
		return fmt.Errorf(returnErr)
	}
	return nil
}
