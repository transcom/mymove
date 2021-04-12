package invoice

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"

	ediResponse997 "github.com/transcom/mymove/pkg/edi/edi997"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type edi997Processor struct {
	db     *pop.Connection
	logger Logger
}

// NewEDI997Processor returns a new EDI997 processor
func NewEDI997Processor(db *pop.Connection,
	logger Logger) services.SyncadaFileProcessor {

	return &edi997Processor{
		db:     db,
		logger: logger,
	}
}

//ProcessFile parses an EDI 997 response and updates the payment request status
func (e *edi997Processor) ProcessFile(path string, stringEDI997 string) error {
	fmt.Printf(path)

	edi997 := ediResponse997.EDI{}
	err := edi997.Parse(stringEDI997)
	if err != nil {
		e.logger.Error("unable to parse EDI997", zap.Error(err))
		return fmt.Errorf("unable to parse EDI997")
	}
	e.logger.Info("RECEIVED: 997 Processor received a 997")
	e.logEDI(edi997)

	// Find the PaymentRequestID that matches the GCN
	icn := edi997.InterchangeControlEnvelope.ISA.InterchangeControlNumber
	var gcn int64
	if edi997.InterchangeControlEnvelope.FunctionalGroups != nil {
		if edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets != nil {
			ak1 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
			gcn = ak1.GroupControlNumber
		} else {
			e.logger.Error("Validation error(s) detected with the EDI997. EDI Errors could not be saved", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI997. EDI Errors could not be saved: %w", err)
		}
	} else {
		e.logger.Error("Validation error(s) detected with the EDI997. EDI Errors could not be saved", zap.Error(err))
		return fmt.Errorf("Validation error(s) detected with the EDI997. EDI Errors could not be saved: %w", err)
	}

	// In the 858, the EDI only has 1 group, and the ICN and the GCN are the same. Therefore, we'll query the PR to ICN table
	// to find the associated payment request using the reported GCN from the 997.
	var paymentRequest models.PaymentRequest
	err = e.db.Q().
		Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
		Where("payment_request_to_interchange_control_numbers.interchange_control_number = ?", int(gcn)).
		First(&paymentRequest)
	if err != nil {
		e.logger.Error("unable to find PaymentRequest with GCN", zap.Error(err))
		return fmt.Errorf("unable to find PaymentRequest with GCN: %s, %d", err.Error(), int(gcn))
	}

	prToICN := models.PaymentRequestToInterchangeControlNumber{
		InterchangeControlNumber: int(icn),
		PaymentRequestID:         paymentRequest.ID,
	}

	var transactionError error
	transactionError = e.db.Transaction(func(tx *pop.Connection) error {
		err = tx.Save(&prToICN)
		if err != nil {
			e.logger.Error("failure saving payment request to interchange control number", zap.Error(err))
			return fmt.Errorf("failure saving payment request to interchange control number: %w", err)
		}
		err = edi997.Validate()
		if err != nil {
			code := "MilMove"
			desc := err.Error()
			ediError := models.EdiError{
				Code:                       &code,
				Description:                &desc,
				PaymentRequestID:           paymentRequest.ID,
				InterchangeControlNumberID: &prToICN.ID,
				EDIType:                    models.EDIType997,
			}
			err = tx.Save(&ediError)
			if err != nil {
				e.logger.Error("failure saving edi validation errors", zap.Error(err))
				return fmt.Errorf("failure saving edi validation errors: %w", err)
			}
			e.logger.Error("Validation error(s) detected with the EDI997", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI997: %w, %v", err, desc)
		}

		paymentRequest.Status = models.PaymentRequestStatusReceivedByGex
		err = tx.Update(&paymentRequest)
		if err != nil {
			e.logger.Error("failure updating payment request", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		e.logger.Info("SUCCESS: 997 Processor updated Payment Request to new status")
		e.logEDIWithPaymentRequest(edi997, paymentRequest)
		return nil
	})

	if transactionError != nil {
		e.logger.Error(transactionError.Error())
		return transactionError
	}

	return nil
}

func (e *edi997Processor) EDIType() models.EDIType {
	return models.EDIType997
}

func (e *edi997Processor) logEDI(edi ediResponse997.EDI) {
	var ak1 edisegment.AK1
	if len(edi.InterchangeControlEnvelope.FunctionalGroups) > 0 && len(edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets) > 0 {
		ak1 = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
	} else {
		e.logger.Warn("unable to log EDI 997, failed functional group or transaction set index check")
		return
	}

	e.logger.Info("EDI 997 log",
		zap.Int64("997 ICN", edi.InterchangeControlEnvelope.ISA.InterchangeControlNumber),
		zap.Int64("858 GCN/ICN", ak1.GroupControlNumber),
	)
}

func (e *edi997Processor) logEDIWithPaymentRequest(edi ediResponse997.EDI, paymentRequest models.PaymentRequest) {
	var ak1 edisegment.AK1
	if len(edi.InterchangeControlEnvelope.FunctionalGroups) > 0 && len(edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets) > 0 {
		ak1 = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
	} else {
		e.logger.Warn("unable to log EDI 997, failed functional group or transaction set index check")
		return
	}

	e.logger.Info("EDI 997 log",
		zap.Int64("997 ICN", edi.InterchangeControlEnvelope.ISA.InterchangeControlNumber),
		zap.Int64("858 GCN/ICN", ak1.GroupControlNumber),
		zap.String("PaymentRequestNumber", paymentRequest.PaymentRequestNumber),
		zap.String("PaymentRequest.Status", string(paymentRequest.Status)),
	)
}
