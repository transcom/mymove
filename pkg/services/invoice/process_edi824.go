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

	edi824 := ediResponse824.EDI{}
	err := edi824.Parse(stringEDI824)
	if err != nil {
		e.logger.Error("unable to parse EDI824", zap.Error(err))
		return fmt.Errorf("unable to parse EDI824")
	}

	transactionSet := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0]
	icn := edi824.InterchangeControlEnvelope.ISA.InterchangeControlNumber
	otiGCN := transactionSet.OTIs[0].GroupControlNumber
	bgn := transactionSet.BGN

	var paymentRequest models.PaymentRequest
	err = e.db.Q().
		Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
		Where("payment_request_to_interchange_control_numbers.interchange_control_number = ?", int(otiGCN)).
		First(&paymentRequest)
	if err != nil {
		e.logger.Error("unable to find PaymentRequest with GCN", zap.Error(err))
		return fmt.Errorf("unable to find PaymentRequest with GCN: %s, %d", err.Error(), int(otiGCN))
	}

	var move models.Move
	err = e.db.Q().
		Find(&move, paymentRequest.MoveTaskOrderID)
	if err != nil {
		e.logger.Error("unable to find move with associated payment request", zap.Error(err))
		return fmt.Errorf("unable to find move with associated payment request: %w", err)
	}

	bgnRefIdentification := bgn.ReferenceIdentification
	mtoRefID := move.ReferenceID
	if mtoRefID == nil {
		e.logger.Error("unable to find PaymentRequest with GCN", zap.Error(err))
		return fmt.Errorf("The BGN02 Reference Identification field: %s doesn't match the Reference ID %s of the associated move", bgnRefIdentification, *mtoRefID)
	}
	if bgnRefIdentification != *mtoRefID {
		e.logger.Error("unable to find PaymentRequest with GCN", zap.Error(err))
		return fmt.Errorf("The BGN02 Reference Identification field: %s doesn't match the Reference ID %v of the associated move", bgnRefIdentification, *mtoRefID)
	}

	err = edi824.Validate()
	if err != nil {
		e.logger.Error("Validation error(s) detected with the EDI824", zap.Error(err))
		return fmt.Errorf("Validation error(s) detected with the EDI824: %w", err)
	}

	teds := fetchTEDSegments(edi824)

	var transactionError error
	transactionError = e.db.Transaction(func(tx *pop.Connection) error {
		prToICN := models.PaymentRequestToInterchangeControlNumber{
			InterchangeControlNumber: int(icn),
			PaymentRequestID:         paymentRequest.ID,
		}
		err = tx.Save(&prToICN)
		if err != nil {
			return fmt.Errorf("failure saving payment request to interchange control number: %w", err)
		}
		for _, ted := range teds {
			code := ted.ApplicationErrorConditionCode
			desc := ted.FreeFormMessage
			ediError := models.EdiError{
				Code:                       &code,
				Description:                &desc,
				PaymentRequestID:           paymentRequest.ID,
				InterchangeControlNumberID: prToICN.ID,
				EDIType:                    models.EDI824,
			}
			err = tx.Save(&ediError)
			if err != nil {
				return fmt.Errorf("failure saving edi technical error description: %w", err)
			}
		}

		paymentRequest.Status = models.PaymentRequestStatusEDIError
		err = tx.Update(&paymentRequest)
		if err != nil {
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		return nil
	})

	if transactionError != nil {
		e.logger.Error(transactionError.Error())
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
