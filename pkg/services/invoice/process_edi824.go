package invoice

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	ediResponse824 "github.com/transcom/mymove/pkg/edi/edi824"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type edi824Processor struct {
	db     *pop.Connection
	logger Logger
}

// NewEDI824Processor returns a new EDI824 processor
func NewEDI824Processor(db *pop.Connection,
	logger Logger) services.EDI824Processor {

	return &edi824Processor{
		db:     db,
		logger: logger,
	}
}

//ProcessFile parses an EDI 824 response and updates the payment request status
func (e *edi824Processor) ProcessFile(path string, stringEDI824 string) error {
	fmt.Printf(path)
	errString := ""

	edi824 := ediResponse824.EDI{}
	err := edi824.Parse(stringEDI824)
	if err != nil {
		errString += err.Error()
	}

	// Find the PaymentRequestID that matches the ICN
	icn := edi824.InterchangeControlEnvelope.ISA.InterchangeControlNumber

	var prToICN models.PaymentRequestToInterchangeControlNumber
	err = e.db.Q().
		Where("interchange_control_number = ?", int(icn)).
		First(&prToICN)
	if err != nil {
		errString += fmt.Sprintf("unable to find PaymentRequestTOInterchangeControlNumber with ICN: %s, %d", err.Error(), int(icn)) + "\n"
	}
	err = edi824.Validate()
	if err != nil {
		errString += err.Error()
	}

	teds := fetchTEDSegments(edi824)

	if errString != "" {
		e.logger.Error(errString)
		return fmt.Errorf(errString)
	}

	var transactionError error
	transactionError = e.db.Transaction(func(tx *pop.Connection) error {
		for _, ted := range teds {
			code := ted.ApplicationErrorConditionCode
			desc := ted.FreeFormMessage
			ediError := models.EdiError{
				Code:                       &code,
				Description:                &desc,
				PaymentRequestID:           prToICN.PaymentRequestID,
				InterchangeControlNumberID: prToICN.ID,
				InterchangeControlNumber:   prToICN,
				EDIType:                    models.EDI824,
			}
			err = tx.Save(&ediError)
			if err != nil {
				e.logger.Error("failure saving edi technical error description", zap.Error(err))
				return fmt.Errorf("failure saving edi technical error description: %w", err)
			}
		}
		var paymentRequest models.PaymentRequest
		err = e.db.Q().
			Where("id = ?", prToICN.PaymentRequestID).
			First(&paymentRequest)
		if err != nil {
			errString += fmt.Sprintf("unable to find payment request with ID: %s, %d", err.Error(), int(icn)) + "\n"
		}
		paymentRequest.Status = models.PaymentRequestStatusEDIError
		err = tx.Update(&paymentRequest)
		if err != nil {
			e.logger.Error("failure updating payment request", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		return nil
	})

	if transactionError != nil {
		return transactionError
	}

	return nil
}

func fetchTEDSegments(edi ediResponse824.EDI) []edisegment.TED {
	var teds []edisegment.TED
	for _, functionalGroup := range edi.InterchangeControlEnvelope.FunctionalGroups {
		for _, transactionSet := range functionalGroup.TransactionSets {
			for _, ted := range transactionSet.TEDs {
				teds = append(teds, ted)
			}
		}
	}
	return teds
}
