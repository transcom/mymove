package invoice

import (
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	ediResponse824 "github.com/transcom/mymove/pkg/edi/edi824"
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
		// TODO: save error to the db
		errString += err.Error()
	}

	// Find the PaymentRequestID that matches the ICN
	icn := edi824.InterchangeControlEnvelope.ISA.InterchangeControlNumber
	var paymentRequest models.PaymentRequest
	err = e.db.Q().
		Join("payment_request_to_interchange_control_numbers", "payment_request_to_interchange_control_numbers.payment_request_id = payment_requests.id").
		Where("payment_request_to_interchange_control_numbers.interchange_control_number = ?", int(icn)).
		First(&paymentRequest)
	if err != nil {
		// TODO: save error to the db
		errString += fmt.Sprintf("unable to find payment request with ID: %s, %d", err.Error(), int(icn)) + "\n"
	}

	err = edi824.Validate()
	if err != nil {
		// TODO: save error to the db
		errString += err.Error()
	}

	if errString != "" {
		e.logger.Error(errString)
		return fmt.Errorf(errString)
	}

	var transactionError error
	transactionError = e.db.Transaction(func(tx *pop.Connection) error {
		paymentRequest.Status = models.PaymentRequestStatusReceivedByGex
		err = e.db.Update(&paymentRequest)
		if err != nil {
			// TODO: save error to the db
			e.logger.Error("failure updating payment request", zap.Error(err))
			return fmt.Errorf("failure updating payment request status: %w", err)
		}
		return nil
	})

	if transactionError != nil {
		// TODO: save error to the db
		return transactionError
	}

	return nil
}
