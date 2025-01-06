package cli

import "github.com/spf13/pflag"

const (
	// ProcessTPPSInvoiceReportPickupDirectory is the ENV var for the directory where TPPS paid invoice files are stored to be processed
	ProcessTPPSInvoiceReportPickupDirectory string = "process_tpps_invoice_report_pickup_directory"
)

// InitTPPSFlags initializes TPPS SFTP command line flags
func InitTPPSFlags(flag *pflag.FlagSet) {
	flag.String(ProcessTPPSInvoiceReportPickupDirectory, "", "TPPS Paid Invoice SFTP Pickup Directory")
}
