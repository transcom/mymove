package paymentrequest

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/invoice"

	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	paymentrequesthelper "github.com/transcom/mymove/pkg/payment_request"
	"github.com/transcom/mymove/pkg/services"
)

type paymentRequestReviewedProcessor struct {
	db                            *pop.Connection
	logger                        Logger
	reviewedPaymentRequestFetcher services.PaymentRequestReviewedFetcher
	ediGenerator                  services.GHCPaymentRequestInvoiceGenerator
	runSendToSyncada              bool // if false, do not send to Syncada, e.g. UT shouldn't send to Syncada
	gexSender                     services.GexSender
	sftpSender                    services.SyncadaSFTPSender
}

// NewPaymentRequestReviewedProcessor returns a new payment request reviewed processor
func NewPaymentRequestReviewedProcessor(db *pop.Connection,
	logger Logger,
	fetcher services.PaymentRequestReviewedFetcher,
	generator services.GHCPaymentRequestInvoiceGenerator,
	runSendToSyncada bool,
	gexSender services.GexSender,
	sftpSender services.SyncadaSFTPSender) services.PaymentRequestReviewedProcessor {

	return &paymentRequestReviewedProcessor{
		db:                            db,
		logger:                        logger,
		reviewedPaymentRequestFetcher: fetcher,
		ediGenerator:                  generator,
		gexSender:                     gexSender,
		sftpSender:                    sftpSender,
		runSendToSyncada:              runSendToSyncada}
}

// InitNewPaymentRequestReviewedProcessor initialize NewPaymentRequestReviewedProcessor for production use
func InitNewPaymentRequestReviewedProcessor(db *pop.Connection, logger Logger, sendToSyncada bool, icnSequencer sequence.Sequencer) (services.PaymentRequestReviewedProcessor, error) {
	reviewedPaymentRequestFetcher := NewPaymentRequestReviewedFetcher(db)
	generator := invoice.NewGHCPaymentRequestInvoiceGenerator(db, icnSequencer, clock.New())
	var sftpSession services.SyncadaSFTPSender
	sftpSession, err := invoice.InitNewSyncadaSFTPSession()
	if err != nil {
		// just log the error, sftpSession is set to nil if there is an error
		logger.Error(fmt.Errorf("configuration of SyncadaSFTPSession failed: %w", err).Error())
		return nil, err
	}
	var gexSender services.GexSender
	gexSender = nil

	return NewPaymentRequestReviewedProcessor(
		db,
		logger,
		reviewedPaymentRequestFetcher,
		generator,
		sendToSyncada,
		gexSender,
		sftpSession), nil
}

func (p *paymentRequestReviewedProcessor) ProcessAndLockReviewedPR(pr models.PaymentRequest) error {
	var transactionError error

	transactionError = p.db.Transaction(func(tx *pop.Connection) error {
		var lockedPR models.PaymentRequest

		query := `
			SELECT * FROM payment_requests
			WHERE id = $1 FOR UPDATE SKIP LOCKED;
		`
		err := p.db.RawQuery(query, pr.ID).First(&lockedPR)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return fmt.Errorf("failure retrieving payment request with ID: %s. Err: %w", pr.ID, err)
		}

		// generate EDI file
		var edi858c ediinvoice.Invoice858C
		edi858c, err = p.ediGenerator.Generate(lockedPR, false)
		if err != nil {
			return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to generator.Generate: %w", err)
		}
		var edi858cString string
		edi858cString, err = edi858c.EDIString(p.logger)
		if err != nil {
			return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to edi858c.EDIString: %w", err)
		}

		p.logger.Info("858 Processor calling SendToSyncada...",
			zap.Int64("858 ICN", edi858c.ISA.InterchangeControlNumber),
			zap.String("ShipmentIdentificationNumber/PaymentRequestNumber", edi858c.Header.ShipmentInformation.ShipmentIdentificationNumber),
			zap.String("ReferenceIdentification/PaymentRequestNumber", edi858c.Header.PaymentRequestNumber.ReferenceIdentification),
			zap.String("Date", edi858c.ISA.InterchangeDate),
			zap.String("Time", edi858c.ISA.InterchangeTime),
		)
		// Send EDI string to Syncada
		// If sent successfully to GEX, update payment request status to SENT_TO_GEX.
		err = paymentrequesthelper.SendToSyncada(edi858cString, p.gexSender, p.sftpSender, p.runSendToSyncada, p.logger)
		if err != nil {
			return fmt.Errorf("error sending the following EDI (PaymentRequest.ID: %s, error string) to Syncada: %s", lockedPR.ID, err)
		}
		sentToGexAt := time.Now()
		lockedPR.SentToGexAt = &sentToGexAt
		lockedPR.Status = models.PaymentRequestStatusSentToGex
		err = p.db.Update(&lockedPR)

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
		verrs, err := p.db.ValidateAndCreate(&errToSave)

		// We are just logging these errors instead of returning them to avoid obscuring the original error
		if err != nil {
			p.logger.Error(
				"failed to save EDI 858 error",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(err),
			)
		} else if verrs != nil && verrs.HasAny() {
			p.logger.Error(
				"failed to save EDI 858 error due to validation errors",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(verrs),
			)
		}

		pr.Status = models.PaymentRequestStatusEDIError
		verrs, err = p.db.ValidateAndUpdate(&pr)
		if err != nil {
			p.logger.Error(
				"error while updating payment request status",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(err),
			)
		} else if verrs != nil && verrs.HasAny() {
			p.logger.Error(
				"failed to update payment request status due to validation errors",
				zap.String("PaymentRequestID", pr.ID.String()),
				zap.Error(verrs),
			)
		}

		return transactionError
	}
	return nil
}

func (p *paymentRequestReviewedProcessor) ProcessReviewedPaymentRequest() error {
	// Store/log metrics about EDI processing upon exiting this method.
	numProcessed := 0
	start := time.Now()
	defer func() {
		ediProcessing := models.EDIProcessing{
			EDIType:          models.EDIType858,
			ProcessStartedAt: start,
			ProcessEndedAt:   time.Now(),
			NumEDIsProcessed: numProcessed,
		}
		p.logger.Info("EDIs processed", zap.Object("EDIs processed", &ediProcessing))

		verrs, err := p.db.ValidateAndCreate(&ediProcessing)
		if err != nil {
			p.logger.Error("failed to create EDIProcessing record", zap.Error(err))
		}
		if verrs.HasAny() {
			p.logger.Error("failed to validate EDIProcessing record", zap.Error(err))
		}
	}()

	// Fetch all payment request that have been reviewed
	reviewedPaymentRequests, err := p.reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest()
	if err != nil {
		return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to FetchReviewedPaymentRequest: %w", err)
	}

	if len(reviewedPaymentRequests) == 0 {
		// No reviewed payment requests to process
		return nil
	}

	// Send all reviewed payment request to Syncada
	for _, pr := range reviewedPaymentRequests {
		err := p.ProcessAndLockReviewedPR(pr)
		if err != nil {
			return err
		}
		numProcessed++
	}

	return nil
}
