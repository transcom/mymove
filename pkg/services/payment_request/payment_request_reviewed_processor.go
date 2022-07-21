package paymentrequest

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/invoice"
)

//GexSendError is returned when there is an error sending an EDI to GEX
type GexSendError struct {
	paymentRequestID uuid.UUID
	err              error
}

func (e GexSendError) Error() string {
	return fmt.Sprintf("error sending the following EDI (PaymentRequest.ID: %s) to GEX: %s", e.paymentRequestID, e.err.Error())
}

type paymentRequestReviewedProcessor struct {
	reviewedPaymentRequestFetcher services.PaymentRequestReviewedFetcher
	ediGenerator                  services.GHCPaymentRequestInvoiceGenerator
	runSendToSyncada              bool // if false, do not send to Syncada, e.g. UT shouldn't send to Syncada
	gexSender                     services.GexSender
	sftpSender                    services.SyncadaSFTPSender
}

// NewPaymentRequestReviewedProcessor returns a new payment request reviewed processor
func NewPaymentRequestReviewedProcessor(
	fetcher services.PaymentRequestReviewedFetcher,
	generator services.GHCPaymentRequestInvoiceGenerator,
	runSendToSyncada bool,
	gexSender services.GexSender,
	sftpSender services.SyncadaSFTPSender) services.PaymentRequestReviewedProcessor {

	return &paymentRequestReviewedProcessor{
		reviewedPaymentRequestFetcher: fetcher,
		ediGenerator:                  generator,
		gexSender:                     gexSender,
		sftpSender:                    sftpSender,
		runSendToSyncada:              runSendToSyncada}
}

// InitNewPaymentRequestReviewedProcessor initialize NewPaymentRequestReviewedProcessor for production use
func InitNewPaymentRequestReviewedProcessor(appCtx appcontext.AppContext, sendToSyncada bool, icnSequencer sequence.Sequencer, gexSender services.GexSender) (services.PaymentRequestReviewedProcessor, error) {
	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher()
	generator := invoice.NewGHCPaymentRequestInvoiceGenerator(icnSequencer, clock.New())
	var sftpSession services.SyncadaSFTPSender
	if gexSender == nil {
		var err error
		sftpSession, err = invoice.InitNewSyncadaSFTPSession()
		if err != nil {
			// just log the error, sftpSession is set to nil if there is an error
			appCtx.Logger().Error(fmt.Errorf("configuration of SyncadaSFTPSession failed: %w", err).Error())
			return nil, err
		}
	}

	return NewPaymentRequestReviewedProcessor(
		reviewedPaymentRequestFetcher,
		generator,
		sendToSyncada,
		gexSender,
		sftpSession), nil
}

func (p *paymentRequestReviewedProcessor) ProcessAndLockReviewedPR(appCtx appcontext.AppContext, pr models.PaymentRequest) error {
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var lockedPR models.PaymentRequest

		query := `
			SELECT * FROM payment_requests
			WHERE id = $1 FOR NO KEY UPDATE SKIP LOCKED;
		`
		err := txnAppCtx.DB().RawQuery(query, pr.ID).First(&lockedPR)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return fmt.Errorf("failure retrieving payment request with ID: %s. Err: %w", pr.ID, err)
		}

		// generate EDI file
		var edi858c ediinvoice.Invoice858C
		edi858c, err = p.ediGenerator.Generate(txnAppCtx, lockedPR, false)
		icn := edi858c.ISA.InterchangeControlNumber
		if err != nil {
			return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to generator.Generate: %w", err)
		}
		var edi858cString string
		edi858cString, err = edi858c.EDIString(txnAppCtx.Logger())
		if err != nil {
			return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to edi858c.EDIString: %w", err)
		}

		txnAppCtx.Logger().Info("858 Processor calling SendToSyncada...",
			zap.Int64("858 ICN", edi858c.ISA.InterchangeControlNumber),
			zap.String("ShipmentIdentificationNumber/PaymentRequestNumber", edi858c.Header.ShipmentInformation.ShipmentIdentificationNumber),
			zap.String("ReferenceIdentification/PaymentRequestNumber", edi858c.Header.PaymentRequestNumber.ReferenceIdentification),
			zap.String("Date", edi858c.ISA.InterchangeDate),
			zap.String("Time", edi858c.ISA.InterchangeTime),
			zap.String("UsageIndicator (ISA-15)", edi858c.ISA.UsageIndicator),
		)
		// Send EDI string to Syncada
		// If sent successfully to GEX, update payment request status to SENT_TO_GEX.
		err = paymentrequesthelper.SendToSyncada(txnAppCtx, edi858cString, icn, p.gexSender, p.sftpSender, p.runSendToSyncada)
		if err != nil {
			return GexSendError{paymentRequestID: lockedPR.ID, err: err}
		}
		sentToGexAt := time.Now()
		lockedPR.SentToGexAt = &sentToGexAt
		lockedPR.Status = models.PaymentRequestStatusSentToGex
		err = txnAppCtx.DB().Update(&lockedPR)

		if err != nil {
			return fmt.Errorf("failure updating payment request status: %w", err)
		}

		return nil
	})
	if transactionError != nil {
		errDescription := transactionError.Error()

		errToSave := models.EdiError{
			PaymentRequestID:           pr.ID,
			InterchangeControlNumberID: nil,
			Code:                       nil,
			Description:                &errDescription,
			EDIType:                    models.EDIType858,
		}
		verrs, err := appCtx.DB().ValidateAndCreate(&errToSave)

		// We are just logging these errors instead of returning them to avoid obscuring the original error
		if err != nil {
			appCtx.Logger().Error(
				"failed to save EDI 858 error",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(err),
			)
		} else if verrs != nil && verrs.HasAny() {
			appCtx.Logger().Error(
				"failed to save EDI 858 error due to validation errors",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(verrs),
			)
		}

		switch transactionError.(type) {
		case GexSendError:
			// if we failed in sending there is nothing to do here but retry later so keep the status the same
		default:
			pr.Status = models.PaymentRequestStatusEDIError
		}
		verrs, err = appCtx.DB().ValidateAndUpdate(&pr)
		if err != nil {
			appCtx.Logger().Error(
				"error while updating payment request status",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(err),
			)
		} else if verrs != nil && verrs.HasAny() {
			appCtx.Logger().Error(
				"failed to update payment request status due to validation errors",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(verrs),
			)
		}

		return transactionError
	}
	return nil
}

func (p *paymentRequestReviewedProcessor) ProcessReviewedPaymentRequest(appCtx appcontext.AppContext) {
	// Store/log metrics about EDI processing upon exiting this method.
	numProcessed := 0
	start := time.Now()
	logger := appCtx.Logger()
	defer func() {
		ediProcessing := models.EDIProcessing{
			EDIType:          models.EDIType858,
			ProcessStartedAt: start,
			ProcessEndedAt:   time.Now(),
			NumEDIsProcessed: numProcessed,
		}
		logger.Info("EDIs processed", zap.Object("edisProcessed", &ediProcessing))

		verrs, err := appCtx.DB().ValidateAndCreate(&ediProcessing)
		if err != nil {
			appCtx.Logger().Error("failed to create EDIProcessing record", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("failed to validate EDIProcessing record", zap.Error(err))
		}
	}()

	// Fetch all payment request that have been reviewed
	reviewedPaymentRequests, err := p.reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest(appCtx)
	if err != nil {
		appCtx.Logger().Error("function ProcessReviewedPaymentRequest failed call to FetchReviewedPaymentRequest", zap.Error(err))
		return
	}

	if len(reviewedPaymentRequests) == 0 {
		// No reviewed payment requests to process
		logger.Info("no payment requests to process found")
		return
	}

	// Send all reviewed payment request to Syncada
	for _, pr := range reviewedPaymentRequests {
		err := p.ProcessAndLockReviewedPR(appCtx, pr)
		if err != nil {
			// only log the error and keep working, one failure shouldn't stop the processing of others
			appCtx.Logger().Error(fmt.Sprintf("failed to process payment request id: %s", pr.ID), zap.Error(err))
		} else {
			numProcessed++
		}
	}
}
