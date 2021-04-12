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
	logger Logger) services.SyncadaFileProcessor {

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

	e.logger.Info("RECEIVED: 824 Processor received a 824")
	e.logEDI(edi824)

	var transactionError error
	var otiGCN int64
	var bgn edisegment.BGN
	transactionError = e.db.Transaction(func(tx *pop.Connection) error {
		icn := edi824.InterchangeControlEnvelope.ISA.InterchangeControlNumber
		if edi824.InterchangeControlEnvelope.FunctionalGroups != nil {
			if edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets != nil {
				transactionSet := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0]
				if len(transactionSet.OTIs) > 0 {
					otiGCN = transactionSet.OTIs[0].GroupControlNumber
				} else {
					e.logger.Error("Validation error(s) detected with the EDI824. EDI Errors could not be saved", zap.Error(err))
					return fmt.Errorf("Validation error(s) detected with the EDI824. EDI Errors could not be saved: %w", err)
				}
				bgn = transactionSet.BGN
			} else {
				e.logger.Error("Validation error(s) detected with the EDI824. EDI Errors could not be saved", zap.Error(err))
				return fmt.Errorf("Validation error(s) detected with the EDI824. EDI Errors could not be saved: %w", err)
			}
		} else {
			e.logger.Error("Validation error(s) detected with the EDI824. EDI Errors could not be saved", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI824. EDI Errors could not be saved: %w", err)
		}

		// In the 858, the EDI only has 1 group, and the ICN and the GCN are the same. Therefore, we'll query the PR to ICN table
		// to find the associated payment request using the reported GCN from the 824.
		var paymentRequest models.PaymentRequest
		err = e.db.Q().
			Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
			Where("payment_request_to_interchange_control_numbers.interchange_control_number = ?", int(otiGCN)).
			First(&paymentRequest)
		if err != nil {
			e.logger.Error("unable to find PaymentRequest with GCN", zap.Error(err))
			return fmt.Errorf("unable to find PaymentRequest with GCN: %s, %d", err.Error(), int(otiGCN))
		}

		prToICN := models.PaymentRequestToInterchangeControlNumber{
			InterchangeControlNumber: int(icn),
			PaymentRequestID:         paymentRequest.ID,
		}
		err = tx.Save(&prToICN)
		if err != nil {
			return fmt.Errorf("failure saving payment request to interchange control number: %w", err)
		}

		err = edi824.Validate()
		if err != nil {
			code := "MilMove"
			desc := err.Error()
			ediError := models.EdiError{
				Code:                       &code,
				Description:                &desc,
				PaymentRequestID:           paymentRequest.ID,
				InterchangeControlNumberID: &prToICN.ID,
				EDIType:                    models.EDIType824,
			}
			err = tx.Save(&ediError)
			if err != nil {
				return fmt.Errorf("failure saving edi validation errors: %w", err)
			}
			e.logger.Error("Validation error(s) detected with the EDI824", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI824: %w, %v", err, desc)
		}

		var move models.Move
		err = e.db.Q().
			Find(&move, paymentRequest.MoveTaskOrderID)
		if err != nil {
			e.logger.Error("unable to find move with associated payment request", zap.Error(err))
			return fmt.Errorf("unable to find move with associated payment request: %w", err)
		}

		// The BGN02 Reference Identification field from the 824 stores the reference identification used in the 858.
		// For MilMove we use the MTO Reference ID in the 858 (which used to the field for the GBLOC, but is not relevant for GHC MilMove).
		bgnRefIdentification := bgn.ReferenceIdentification
		mtoRefID := move.ReferenceID
		if mtoRefID == nil {
			e.logger.Error(fmt.Sprintf("An associated move with mto.ReferenceID: %s was not found", *mtoRefID), zap.Error(err))
			return fmt.Errorf("An associated move with mto.ReferenceID: %s was not found", *mtoRefID)
		}
		if bgnRefIdentification != *mtoRefID {
			e.logger.Error(fmt.Sprintf("The BGN02 Reference Identification field: %s doesn't match the Reference ID %s of the associated move", bgnRefIdentification, *mtoRefID), zap.Error(err))
			return fmt.Errorf("The BGN02 Reference Identification field: %s doesn't match the Reference ID %v of the associated move", bgnRefIdentification, *mtoRefID)
		}

		teds := fetchTEDSegments(edi824)

		for _, ted := range teds {
			code := ted.ApplicationErrorConditionCode
			desc := ted.FreeFormMessage
			ediError := models.EdiError{
				Code:                       &code,
				Description:                &desc,
				PaymentRequestID:           paymentRequest.ID,
				InterchangeControlNumberID: &prToICN.ID,
				EDIType:                    models.EDIType824,
			}
			err = tx.Save(&ediError)
			if err != nil {
				e.logger.Error("failure saving edi technical error description", zap.Error(err))
				return fmt.Errorf("failure saving edi technical error description: %w", err)
			}
		}

		paymentRequest.Status = models.PaymentRequestStatusEDIError
		err = tx.Update(&paymentRequest)
		if err != nil {
			e.logger.Error("failure updating payment request status:", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		e.logger.Info("SUCCESS: 824 Processor updated Payment Request to new status")
		e.logEDIWithPaymentRequest(edi824, paymentRequest)
		return nil
	})

	if transactionError != nil {
		e.logger.Error(transactionError.Error())
		return transactionError
	}

	return nil
}

func fetchTEDSegments(edi ediResponse824.EDI) []edisegment.TED {
	var teds []edisegment.TED
	for _, functionalGroup := range edi.InterchangeControlEnvelope.FunctionalGroups {
		for _, transactionSet := range functionalGroup.TransactionSets {
			teds = append(teds, transactionSet.TEDs...)
		}
	}
	return teds
}

func (e *edi824Processor) EDIType() models.EDIType {
	return models.EDIType824
}

func (e *edi824Processor) logEDI(edi ediResponse824.EDI) {
	var transactionSet0 ediResponse824.TransactionSet
	var bgn edisegment.BGN
	var otiGCN int64
	icn := edi.InterchangeControlEnvelope.ISA.InterchangeControlNumber
	if len(edi.InterchangeControlEnvelope.FunctionalGroups) > 0 && len(edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets) > 0 {
		transactionSet0 = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0]
		bgn = transactionSet0.BGN
		if len(transactionSet0.OTIs) > 0 {
			otiGCN = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].OTIs[0].GroupControlNumber
		} else {
			e.logger.Warn("unable to log EDI 824, failed OTI index check")
			return
		}
	} else {
		e.logger.Warn("unable to log EDI 824, failed functional group or transaction sets index check")
		return
	}

	e.logger.Info("EDI 824 log",
		zap.Int64("824 ICN", icn),
		zap.String("BGN.ReferenceIdentification", bgn.ReferenceIdentification),
		zap.Int64("858 GCN", otiGCN),
	)
}

func (e *edi824Processor) logEDIWithPaymentRequest(edi ediResponse824.EDI, paymentRequest models.PaymentRequest) {
	var transactionSet0 ediResponse824.TransactionSet
	var bgn edisegment.BGN
	var otiGCN int64
	icn := edi.InterchangeControlEnvelope.ISA.InterchangeControlNumber
	if len(edi.InterchangeControlEnvelope.FunctionalGroups) > 0 && len(edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets) > 0 {
		transactionSet0 = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0]
		bgn = transactionSet0.BGN
		if len(transactionSet0.OTIs) > 0 {
			otiGCN = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].OTIs[0].GroupControlNumber
		} else {
			e.logger.Warn("unable to log EDI 824, failed OTI index check")
			return
		}
	} else {
		e.logger.Warn("unable to log EDI 824, failed functional group or transaction sets index check")
		return
	}

	e.logger.Info("EDI 824 log",
		zap.Int64("824 ICN", icn),
		zap.String("BGN.ReferenceIdentification", bgn.ReferenceIdentification),
		zap.Int64("858 GCN", otiGCN),
		zap.String("PaymentRequestNumber", paymentRequest.PaymentRequestNumber),
		zap.String("PaymentRequest.Status", string(paymentRequest.Status)),
	)
}
