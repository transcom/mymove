package paymentrequest

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop"

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
}

// NewPaymentRequestReviewedProcessor returns a new payment request reviewed processor
func NewPaymentRequestReviewedProcessor(db *pop.Connection, logger Logger, fetcher services.PaymentRequestReviewedFetcher, generator services.GHCPaymentRequestInvoiceGenerator) services.PaymentRequestReviewedProcessor {
	return &paymentRequestReviewedProcessor{db, logger, fetcher, generator}
}

func (p *paymentRequestReviewedProcessor) ProcessReviewedPaymentRequest() error {

	// Fetch all payment request that have been reviewed
	reviewedPaymentRequests, err := p.reviewedPaymentRequestFetcher.FetchReviewedPaymentRequest()
	if err != nil {
		return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to FetchReviewedPaymentRequest: %w", err)
	}

	// records for successfully sent PRs
	var sentToGexStatuses []string
	// records for PRs that failed to send
	var failed []string

	// Send all reviewed payment request to Syncada
	paymentHelper := paymentrequesthelper.RequestPaymentHelper{DB: p.db, Logger: p.logger}
	for _, pr := range reviewedPaymentRequests {

		// generate EDI file
		var edi858c ediinvoice.Invoice858C
		edi858c, err = p.ediGenerator.Generate(pr, false)
		if err != nil {
			return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to generator.Generate: %w", err)

		}
		var edi858cString string
		edi858cString, err = edi858c.EDIString()
		if err != nil {
			return fmt.Errorf("function ProcessReviewedPaymentRequest failed call to edi858c.EDIString: %w", err)

		}

		// TODO: USING STUB FUNCTION FOR SEND TO SYNCADA

		// Send EDI string to Syncada
		// If sent successfully to GEX, update payment request status to SENT_TO_GEX.
		err = paymentHelper.SendToSyncada(edi858cString)
		if err != nil {
			// save payment request ID and error
			// TODO: if there is an error, no way to flag it with a status.
			// (ID, error string) to be returned in an error message
			f := []string{pr.ID.String(), err.Error()}
			value := "('" + strings.Join(f, "','") + "')"
			failed = append(failed, value)
		} else {
			// (ID, status) to be used in update query
			// ('a2c34dba-015f-4f96-a38b-0c0b9272e208'::uuid,'SENT_TO_GEX'::payment_request_status)
			status := []string{"'" + pr.ID.String() + "'::uuid", "'" + models.PaymentRequestStatusSentToGex.String() + "'::payment_request_status"}
			value := "(" + strings.Join(status, ",") + ")"
			sentToGexStatuses = append(sentToGexStatuses, value)
		}
	}

	// save error messages from failed sends
	var errFailedToSendString string
	if len(failed) > 0 {
		errFailedToSendString = "error sending the following EDIs (PaymentRequest.ID, error string) to Syncada:\n\t"
		for _, e := range failed {
			errFailedToSendString += "\t" + e + "\n"
		}
	}

	transactionError := p.db.Transaction(func(tx *pop.Connection) error {
		// Save PRs with successful sent to GEX

		/* Use `update...from` postgres syntax described here https://stackoverflow.com/a/18799497
		   To have one call to the DB if we have multiple updates.
		UPDATE payment_requests AS pr SET
		    status = c.status
		FROM (VALUES
		    ('id1', 'status1'),
		    ('id2', 'status2')
		    ) AS c(id, status)
		WHERE c.id = pr.id;
		*/

		values := strings.Join(sentToGexStatuses, ",")
		q := `
UPDATE payment_requests AS pr SET
    status = c.status
FROM (VALUES
	%s
    ) AS c(id, status)
WHERE c.id = pr.id;`
		qq := fmt.Sprintf(q, values)
		err = tx.RawQuery(qq).Exec()
		if err != nil {
			return fmt.Errorf("failure updating payment request status: %w", err)
		}

		return nil
	})

	// Build up error string
	returnError := ""
	if errFailedToSendString != "" {
		returnError += errFailedToSendString
	}
	if transactionError != nil {
		if returnError != "" {
			returnError += "\n"
		}
		returnError += transactionError.Error()
	}
	if returnError != "" {
		return fmt.Errorf("function ProcessReviewedPaymentRequest has failure(s): %s", returnError)
	}

	return nil
}
