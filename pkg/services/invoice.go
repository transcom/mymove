package services

import (
	"net/http"

	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
)

// GexSender is an interface for sending and receiving a request
type GexSender interface {
	SendToGex(edi string, transactionName string) (resp *http.Response, err error)
}

// GHCPaymentRequestInvoiceGenerator is the exported interface for generating an invoice
//go:generate mockery -name GHCPaymentRequestInvoiceGenerator
type GHCPaymentRequestInvoiceGenerator interface {
	Generate(paymentRequest models.PaymentRequest, sendProductionInvoice bool) (ediinvoice.Invoice858C, error)
}

// GHCJobRunner is the exported interface for finding payment request ready to be sent to GEX
//go:generate mockery -name GHCJobRunner
type GHCJobRunner interface {
	ApprovedPaymentRequestFetcher() (models.PaymentRequests, error)
}
