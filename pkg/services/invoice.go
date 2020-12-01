package services

import (
	"io"
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

// SyncadaSFTPSender is the exported interface for sending an EDI to Syncada
//go:generate mockery -name SyncadaSFTPSender
type SyncadaSFTPSender interface {
	SendToSyncadaViaSFTP(localDataReader io.Reader, syncadaFileName string) (int64, error)
}
