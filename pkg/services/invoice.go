package services

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/pkg/sftp"

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

// SFTPFile is the exported interface for a file retrieved from Syncada via SFTP
//go:generate mockery --name SFTPFile --outpkg ghcmocks --output ./ghcmocks
type SFTPFile interface {
	WriteTo(w io.Writer) (written int64, err error)
	Close() error
}

// SFTPClient is the exported interface for an SFTP client created for reading from Syncada
//go:generate mockery --name SFTPClient --outpkg ghcmocks --output ./ghcmocks
type SFTPClient interface {
	ReadDir(p string) ([]os.FileInfo, error)
	Open(path string) (*sftp.File, error)
	Remove(path string) error
	Close() error
}

// SyncadaSFTPReader is the exported interface for reading files from Syncada
//go:generate mockery -name SyncadaSFTPReader
type SyncadaSFTPReader interface {
	FetchAndProcessSyncadaFiles(syncadaPath string, lastRead time.Time, processor SyncadaFileProcessor) error
}

// SyncadaFileProcessor is the exported interface for processing EDI files from Syncada
//go:generate mockery -name SyncadaFileProcessor
type SyncadaFileProcessor interface {
	ProcessFile(syncadaPath string, text string) error
}
