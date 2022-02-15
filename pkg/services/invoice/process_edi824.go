package invoice

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	ediResponse824 "github.com/transcom/mymove/pkg/edi/edi824"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type edi824Processor struct {
}

// NewEDI824Processor returns a new EDI824 processor
func NewEDI824Processor() services.SyncadaFileProcessor {

	return &edi824Processor{}
}

//ProcessFile parses an EDI 824 response and updates the payment request status
func (e *edi824Processor) ProcessFile(appCtx appcontext.AppContext, path string, stringEDI824 string) error {
	edi824 := ediResponse824.EDI{}
	err := edi824.Parse(stringEDI824)
	if err != nil {
		appCtx.Logger().Error("unable to parse EDI824", zap.Error(err))
		return fmt.Errorf("unable to parse EDI824")
	}

	appCtx.Logger().Info("RECEIVED: 824 Processor received a 824")
	e.logEDI(appCtx, edi824)

	var transactionError error
	var otiGCN int64
	var bgn edisegment.BGN
	transactionError = appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		icn := edi824.InterchangeControlEnvelope.ISA.InterchangeControlNumber
		if edi824.InterchangeControlEnvelope.FunctionalGroups != nil {
			if edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets != nil {
				transactionSet := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0]
				if len(transactionSet.OTIs) > 0 {
					otiGCN = transactionSet.OTIs[0].GroupControlNumber
				} else {
					txnAppCtx.Logger().Error("Validation error(s) detected with the EDI824. EDI Errors could not be saved", zap.Error(err))
					return fmt.Errorf("Validation error(s) detected with the EDI824. EDI Errors could not be saved: %w", err)
				}
				bgn = transactionSet.BGN
			} else {
				txnAppCtx.Logger().Error("Validation error(s) detected with the EDI824. EDI Errors could not be saved", zap.Error(err))
				return fmt.Errorf("Validation error(s) detected with the EDI824. EDI Errors could not be saved: %w", err)
			}
		} else {
			txnAppCtx.Logger().Error("Validation error(s) detected with the EDI824. EDI Errors could not be saved", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI824. EDI Errors could not be saved: %w", err)
		}

		// In the 858, the EDI only has 1 group, and the ICN and the GCN are the same. Therefore, we'll query the PR to ICN table
		// to find the associated payment request using the reported GCN from the 824.
		// we are only processing 824s in response to 858s
		var paymentRequest models.PaymentRequest
		err = txnAppCtx.DB().Q().
			Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
			Where("payment_request_to_interchange_control_numbers.interchange_control_number = ? and payment_request_to_interchange_control_numbers.edi_type = ?", int(otiGCN), models.EDIType858).
			First(&paymentRequest)
		if err != nil {
			txnAppCtx.Logger().Error("unable to find PaymentRequest with GCN", zap.Error(err))
			return fmt.Errorf("unable to find PaymentRequest with GCN: %s, %d", err.Error(), int(otiGCN))
		}

		prToICN := models.PaymentRequestToInterchangeControlNumber{
			InterchangeControlNumber: int(icn),
			PaymentRequestID:         paymentRequest.ID,
			EDIType:                  models.EDIType824,
		}
		err = txnAppCtx.DB().Save(&prToICN)
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
			err = txnAppCtx.DB().Save(&ediError)
			if err != nil {
				return fmt.Errorf("failure saving edi validation errors: %w", err)
			}
			txnAppCtx.Logger().Error("Validation error(s) detected with the EDI824", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI824: %w, %v", err, desc)
		}

		var move models.Move
		err = txnAppCtx.DB().Q().
			Find(&move, paymentRequest.MoveTaskOrderID)
		if err != nil {
			txnAppCtx.Logger().Error("unable to find move with associated payment request", zap.Error(err))
			switch err {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(paymentRequest.MoveTaskOrderID, "looking for Move")
			default:
				return apperror.NewQueryError("Move", err, "")
			}
		}

		// The BGN02 Reference Identification field from the 824 stores the payment request number used in the 858.
		// For MilMove we use the Payment Request Number in the 858
		bgnRefIdentification := bgn.ReferenceIdentification
		paymentRequestNumber := paymentRequest.PaymentRequestNumber
		if bgnRefIdentification != paymentRequestNumber {
			txnAppCtx.Logger().Error(fmt.Sprintf("The BGN02 Reference Identification field: %s doesn't match the PaymentRequestNumber %s of the associated payment request", bgnRefIdentification, paymentRequestNumber), zap.Error(err))
			return fmt.Errorf("The BGN02 Reference Identification field: %s doesn't match the PaymentRequestNumber %v of the associated payment request", bgnRefIdentification, paymentRequestNumber)
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
			err = txnAppCtx.DB().Save(&ediError)
			if err != nil {
				txnAppCtx.Logger().Error("failure saving edi technical error description", zap.Error(err))
				return fmt.Errorf("failure saving edi technical error description: %w", err)
			}
		}

		paymentRequest.Status = models.PaymentRequestStatusEDIError
		err = txnAppCtx.DB().Update(&paymentRequest)
		if err != nil {
			txnAppCtx.Logger().Error("failure updating payment request status:", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		txnAppCtx.Logger().Info("SUCCESS: 824 Processor updated Payment Request to new status")
		e.logEDIWithPaymentRequest(txnAppCtx, edi824, paymentRequest)
		return nil
	})

	if transactionError != nil {
		appCtx.Logger().Error(transactionError.Error())
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

func (e *edi824Processor) logEDI(appCtx appcontext.AppContext, edi ediResponse824.EDI) {
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
			appCtx.Logger().Warn("unable to log EDI 824, failed OTI index check")
			return
		}
	} else {
		appCtx.Logger().Warn("unable to log EDI 824, failed functional group or transaction sets index check")
		return
	}

	appCtx.Logger().Info("EDI 824 log",
		zap.Int64("824 ICN", icn),
		zap.String("BGN.ReferenceIdentification", bgn.ReferenceIdentification),
		zap.Int64("858 GCN", otiGCN),
		zap.String("UsageIndicator (ISA-15)", edi.InterchangeControlEnvelope.ISA.UsageIndicator),
	)
}

func (e *edi824Processor) logEDIWithPaymentRequest(appCtx appcontext.AppContext, edi ediResponse824.EDI, paymentRequest models.PaymentRequest) {
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
			appCtx.Logger().Warn("unable to log EDI 824, failed OTI index check")
			return
		}
	} else {
		appCtx.Logger().Warn("unable to log EDI 824, failed functional group or transaction sets index check")
		return
	}

	appCtx.Logger().Info("EDI 824 log",
		zap.Int64("824 ICN", icn),
		zap.String("BGN.ReferenceIdentification", bgn.ReferenceIdentification),
		zap.Int64("858 GCN", otiGCN),
		zap.String("PaymentRequestNumber", paymentRequest.PaymentRequestNumber),
		zap.String("PaymentRequest.Status", string(paymentRequest.Status)),
		zap.String("UsageIndicator (ISA-15)", edi.InterchangeControlEnvelope.ISA.UsageIndicator),
	)
}
