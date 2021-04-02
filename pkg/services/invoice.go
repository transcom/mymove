package services

import (
	"io"
	"net/http"
	"os"
	"time"

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

// SFTPFiler is the exported interface for an SFTP client file
//go:generate mockery --name SFTPFiler
type SFTPFiler interface {
	Close() error
	WriteTo(w io.Writer) (int64, error)
}

// SFTPClient is the exported interface for an SFTP client created for reading from Syncada
//go:generate mockery --name SFTPClient
type SFTPClient interface {
	ReadDir(p string) ([]os.FileInfo, error)
	Open(path string) (SFTPFiler, error)
	Remove(path string) error
}

// SyncadaSFTPReader is the exported interface for reading files from Syncada
//go:generate mockery -name SyncadaSFTPReader
type SyncadaSFTPReader interface {
	FetchAndProcessSyncadaFiles(syncadaPath string, lastRead time.Time, processor SyncadaFileProcessor) (time.Time, error)
}

// SyncadaFileProcessor is the exported interface for processing EDI files from Syncada
//go:generate mockery -name SyncadaFileProcessor
type SyncadaFileProcessor interface {
	ProcessFile(syncadaPath string, text string) error
	EDIType() models.EDIType
}
