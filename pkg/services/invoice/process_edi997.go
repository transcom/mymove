package invoice

import (
	"fmt"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appconfig"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"

	ediResponse997 "github.com/transcom/mymove/pkg/edi/edi997"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type edi997Processor struct {
}

// NewEDI997Processor returns a new EDI997 processor
func NewEDI997Processor() services.SyncadaFileProcessor {

	return &edi997Processor{}
}

//ProcessFile parses an EDI 997 response and updates the payment request status
func (e *edi997Processor) ProcessFile(appCfg appconfig.AppConfig, path string, stringEDI997 string) error {
	fmt.Print(path)

	edi997 := ediResponse997.EDI{}
	err := edi997.Parse(stringEDI997)
	if err != nil {
		appCfg.Logger().Error("unable to parse EDI997", zap.Error(err))
		return fmt.Errorf("unable to parse EDI997")
	}
	appCfg.Logger().Info("RECEIVED: 997 Processor received a 997")
	e.logEDI(appCfg, edi997)

	// Find the PaymentRequestID that matches the GCN
	var gcn int64
	var ediTypeFromAK2 string
	if edi997.InterchangeControlEnvelope.FunctionalGroups != nil {
		if edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets != nil {
			ak1 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
			gcn = ak1.GroupControlNumber

			ediTypeFromAK2 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK2.TransactionSetIdentifierCode
		} else {
			appCfg.Logger().Error("Validation error(s) detected with the EDI997. EDI Errors could not be saved", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI997. EDI Errors could not be saved: %w", err)
		}
	} else {
		appCfg.Logger().Error("Validation error(s) detected with the EDI997. EDI Errors could not be saved", zap.Error(err))
		return fmt.Errorf("Validation error(s) detected with the EDI997. EDI Errors could not be saved: %w", err)
	}

	// In the 858, the EDI only has 1 group, and the ICN and the GCN are the same. Therefore, we'll query the PR to ICN table
	// to find the associated payment request using the reported GCN from the 997.
	var paymentRequest models.PaymentRequest
	err = appCfg.DB().Q().
		Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
		Where("payment_request_to_interchange_control_numbers.interchange_control_number = ? and payment_request_to_interchange_control_numbers.edi_type = ?", int(gcn), ediTypeFromAK2).
		First(&paymentRequest)
	if err != nil {
		appCfg.Logger().Error("unable to find PaymentRequest with GCN", zap.Error(err))
		return fmt.Errorf("unable to find PaymentRequest with GCN: %s, %d", err.Error(), int(gcn))
	}

	icn := edi997.InterchangeControlEnvelope.ISA.InterchangeControlNumber
	prToICN := models.PaymentRequestToInterchangeControlNumber{
		InterchangeControlNumber: int(icn),
		PaymentRequestID:         paymentRequest.ID,
		EDIType:                  models.EDIType997,
	}

	transactionError := appCfg.NewTransaction(func(txnAppCfg appconfig.AppConfig) error {
		lookupErr := txnAppCfg.DB().Where("payment_request_id = ? and interchange_control_number = ? and edi_type = ?", prToICN.PaymentRequestID, prToICN.InterchangeControlNumber, prToICN.EDIType).First(&prToICN)
		if lookupErr != nil {
			txnAppCfg.Logger().Error("failure looking up payment request to interchange control number", zap.Error(err))
		}
		if prToICN.ID == uuid.Nil {
			err = txnAppCfg.DB().Save(&prToICN)
			if err != nil {
				txnAppCfg.Logger().Error("failure saving payment request to interchange control number", zap.Error(err))
				return fmt.Errorf("failure saving payment request to interchange control number: %w", err)
			}
		} else {
			txnAppCfg.Logger().Info(fmt.Sprintf("duplicate EDI %s processed for payment request: %s with ICN: %d", prToICN.EDIType, prToICN.PaymentRequestID, prToICN.InterchangeControlNumber))
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
			err = txnAppCfg.DB().Save(&ediError)
			if err != nil {
				txnAppCfg.Logger().Error("failure saving edi validation errors", zap.Error(err))
				return fmt.Errorf("failure saving edi validation errors: %w", err)
			}
			txnAppCfg.Logger().Error("Validation error(s) detected with the EDI997", zap.Error(err))
			return fmt.Errorf("Validation error(s) detected with the EDI997: %w, %v", err, desc)
		}

		paymentRequest.Status = models.PaymentRequestStatusReceivedByGex
		err = txnAppCfg.DB().Update(&paymentRequest)
		if err != nil {
			txnAppCfg.Logger().Error("failure updating payment request", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		txnAppCfg.Logger().Info("SUCCESS: 997 Processor updated Payment Request to new status")
		e.logEDIWithPaymentRequest(txnAppCfg, edi997, paymentRequest)
		return nil
	})

	if transactionError != nil {
		appCfg.Logger().Error(transactionError.Error())
		return transactionError
	}

	return nil
}

func (e *edi997Processor) EDIType() models.EDIType {
	return models.EDIType997
}

func (e *edi997Processor) logEDI(appCfg appconfig.AppConfig, edi ediResponse997.EDI) {
	var ak1 edisegment.AK1
	if len(edi.InterchangeControlEnvelope.FunctionalGroups) > 0 && len(edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets) > 0 {
		ak1 = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
	} else {
		appCfg.Logger().Warn("unable to log EDI 997, failed functional group or transaction set index check")
		return
	}

	appCfg.Logger().Info("EDI 997 log",
		zap.Int64("997 ICN", edi.InterchangeControlEnvelope.ISA.InterchangeControlNumber),
		zap.Int64("858 GCN/ICN", ak1.GroupControlNumber),
		zap.String("UsageIndicator (ISA-15)", edi.InterchangeControlEnvelope.ISA.UsageIndicator),
	)
}

func (e *edi997Processor) logEDIWithPaymentRequest(appCfg appconfig.AppConfig, edi ediResponse997.EDI, paymentRequest models.PaymentRequest) {
	var ak1 edisegment.AK1
	if len(edi.InterchangeControlEnvelope.FunctionalGroups) > 0 && len(edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets) > 0 {
		ak1 = edi.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
	} else {
		appCfg.Logger().Warn("unable to log EDI 997, failed functional group or transaction set index check")
		return
	}

	appCfg.Logger().Info("EDI 997 log",
		zap.Int64("997 ICN", edi.InterchangeControlEnvelope.ISA.InterchangeControlNumber),
		zap.Int64("858 GCN/ICN", ak1.GroupControlNumber),
		zap.String("PaymentRequestNumber", paymentRequest.PaymentRequestNumber),
		zap.String("PaymentRequest.Status", string(paymentRequest.Status)),
		zap.String("UsageIndicator (ISA-15)", edi.InterchangeControlEnvelope.ISA.UsageIndicator),
	)
}
