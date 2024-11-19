package cli

import "github.com/spf13/pflag"

// Set of flags used for SFTPTPPSPaid
const (
	// SFTPTPPSPaidInvoiceReportPickupDirectory is the ENV var for the directory where TPPS delivers the TPPS paid invoice report
	// TODO: Create a parameter called /{environment_name}/s3_filepath to test getting files from the S3 path in the experiemental and follow on environments
	SFTPTPPSPaidInvoiceReportPickupDirectory string = "s3-filepath"
)

// InitTPPSSFTPFlags initializes TPPS SFTP command line flags
func InitTPPSSFTPFlags(flag *pflag.FlagSet) {
	flag.String(SFTPTPPSPaidInvoiceReportPickupDirectory, "", "TPPS Paid Invoice SFTP Pickup Directory")
}
