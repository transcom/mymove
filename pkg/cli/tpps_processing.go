package cli

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	// ProcessTPPSCustomDateFile is the env var for the date of a file that can be customized if we want to process a payment file other than the daily run of the task
	ProcessTPPSCustomDateFile string = "process_tpps_custom_date_file"
	// TPPSS3Bucket is the env var for the S3 bucket for TPPS payment files that we import from US bank
	TPPSS3Bucket string = "tpps_s3_bucket"
	// TPPSS3Folder is the env var for the S3 folder inside the tpps_s3_bucket for TPPS payment files that we import from US bank
	TPPSS3Folder string = "tpps_s3_folder"
)

// InitTPPSFlags initializes TPPS SFTP command line flags
func InitTPPSFlags(flag *pflag.FlagSet) {
	flag.String(ProcessTPPSCustomDateFile, "", "Custom date for TPPS filename to process, format of MILMOVE-enYYYYMMDD.csv")
	flag.String(TPPSS3Bucket, "", "S3 bucket for TPPS payment files that we import from US bank")
	flag.String(TPPSS3Folder, "", "S3 folder inside the TPPSS3Bucket for TPPS payment files that we import from US bank")
}

// CheckTPPSFlags validates the TPPS processing command line flags
func CheckTPPSFlags(v *viper.Viper) error {
	ProcessTPPSCustomDateFile := v.GetString(ProcessTPPSCustomDateFile)
	if ProcessTPPSCustomDateFile == "" {
		return fmt.Errorf("invalid ProcessTPPSCustomDateFile %s, expecting the format of MILMOVE-enYYYYMMDD.csv", ProcessTPPSCustomDateFile)
	}

	TPPSS3Bucket := v.GetString(TPPSS3Bucket)
	if TPPSS3Bucket == "" {
		return fmt.Errorf("no value for TPPSS3Bucket found")
	}

	TPPSS3Folder := v.GetString(TPPSS3Folder)
	if TPPSS3Folder == "" {
		return fmt.Errorf("no value for TPPSS3Folder found")
	}

	return nil
}
