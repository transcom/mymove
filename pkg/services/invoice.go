package services

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	tppsResponse "github.com/transcom/mymove/pkg/edi/tpps_paid_invoice_report"
	"github.com/transcom/mymove/pkg/models"
)

// GEXChannel type to define constants for valid GEX Channels
type GEXChannel string

const (
	// GEXChannelInvoice is the URL query parameter that we use when sending EDI invoices to US Bank via GEX
	GEXChannelInvoice GEXChannel = "TRANSCOM-DPS-MILMOVE-CPS-IN-USBANK-RCOM"
	// GEXChannelDataWarehouse is the URL query parameter that we use when sending data to the IGC data warehouse
	GEXChannelDataWarehouse GEXChannel = "TRANSCOM-DPS-MILMOVE-GHG-IN-IGC-RCOM"
)

// GexSender is an interface for sending and receiving a request
//
//go:generate mockery --name GexSender
type GexSender interface {
	SendToGex(channel GEXChannel, edi string, filename string) (resp *http.Response, err error)
}

// GHCPaymentRequestInvoiceGenerator is the exported interface for generating an invoice
//
//go:generate mockery --name GHCPaymentRequestInvoiceGenerator
type GHCPaymentRequestInvoiceGenerator interface {
	Generate(appCtx appcontext.AppContext, paymentRequest models.PaymentRequest, sendProductionInvoice bool) (ediinvoice.Invoice858C, error)
}

// SyncadaSFTPSender is the exported interface for sending an EDI to Syncada
//
//go:generate mockery --name SyncadaSFTPSender
type SyncadaSFTPSender interface {
	SendToSyncadaViaSFTP(appCtx appcontext.AppContext, localDataReader io.Reader, syncadaFileName string) (int64, error)
}

// SFTPFiler is the exported interface for an SFTP client file
//
//go:generate mockery --name SFTPFiler
type SFTPFiler interface {
	Close() error
	WriteTo(w io.Writer) (int64, error)
}

// SFTPClient is the exported interface for an SFTP client created for reading from and SFTP connection
//
//go:generate mockery --name SFTPClient
type SFTPClient interface {
	ReadDir(p string) ([]os.FileInfo, error)
	Open(path string) (SFTPFiler, error)
	Remove(path string) error
}

// SyncadaSFTPReader is the exported interface for reading files from an SFTP connection
//
//go:generate mockery --name SyncadaSFTPReader
type SyncadaSFTPReader interface {
	FetchAndProcessSyncadaFiles(appCtx appcontext.AppContext, pickupPath string, lastRead time.Time, processor SyncadaFileProcessor) (time.Time, error)
}

// SyncadaFileProcessor is the exported interface for processing EDI files from Syncada
//
//go:generate mockery --name SyncadaFileProcessor
type SyncadaFileProcessor interface {
	ProcessFile(appCtx appcontext.AppContext, syncadaPath string, text string) error
	EDIType() models.EDIType
}

// TPPSPaidInvoiceReportProcessor defines an interface for storing TPPS payment files in the database
//
//go:generate mockery --name TPPSPaidInvoiceReportProcessor
type TPPSPaidInvoiceReportProcessor interface {
	ProcessFile(appCtx appcontext.AppContext, syncadaPath string, text string) error
	StoreTPPSPaidInvoiceReportInDatabase(appCtx appcontext.AppContext, tppsData []tppsResponse.TPPSData) (*validate.Errors, int, int, error)
}
