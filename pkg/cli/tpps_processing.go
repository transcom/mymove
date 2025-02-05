package cli

import "github.com/spf13/pflag"

const (
	// ProcessTPPSInvoiceReportPickupDirectory is the ENV var for the directory where TPPS paid invoice files are stored to be processed
	ProcessTPPSInvoiceReportPickupDirectory string = "process_tpps_invoice_report_pickup_directory"
	// ProcessTPPSCustomDateFile is the env var for the date of a file that can be customized if we want to process a payment file other than the daily run of the task
	ProcessTPPSCustomDateFile string = "process_tpps_custom_date_file"
	// TPPSS3Bucket is the env var for the S3 bucket for TPPS payment files that we import from US bank
	TPPSS3Bucket string = "tpps_s3_bucket"
	// TPPSS3Folder is the env var for the S3 folder inside the tpps_s3_bucket for TPPS payment files that we import from US bank
	TPPSS3Folder string = "tpps_s3_folder"
)

// InitTPPSFlags initializes TPPS SFTP command line flags
func InitTPPSFlags(flag *pflag.FlagSet) {
	flag.String(ProcessTPPSInvoiceReportPickupDirectory, "", "TPPS Paid Invoice SFTP Pickup Directory")
	flag.String(ProcessTPPSCustomDateFile, "", "Custom date for TPPS filename to process, format of MILMOVE-enYYYYMMDD.csv")
	flag.String(TPPSS3Bucket, "", "S3 bucket for TPPS payment files that we import from US bank")
	flag.String(TPPSS3Folder, "", "S3 folder inside the TPPSS3Bucket for TPPS payment files that we import from US bank")
}
