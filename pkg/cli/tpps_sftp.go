package cli

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Set of flags used for SFTPTPPSPaid
const (
	// SFTPTPPSPaidInvoiceReportPickupDirectory is the ENV var for the directory where TPPS delivers the TPPS paid invoice report

	// maria evaluated whether you should actually keep this in here
	SFTPTPPSPaidInvoiceReportPickupDirectory string = "S3 BUCKET HERE"
)

// maria i don't know if you want to even keep this function if we don't need it for
// tpps processing

// InitTPPSSFTPFlags initializes TPPS SFTP command line flags
func InitTPPSSFTPFlags(flag *pflag.FlagSet) {
	// flag.Int(GEXSFTPPortFlag, 22, "GEX SFTP Port")
	// flag.String(GEXSFTPUserIDFlag, "", "GEX SFTP User ID")
	// flag.String(GEXSFTPIPAddressFlag, "localhost", "GEX SFTP IP Address")
	// flag.String(GEXSFTPPasswordFlag, "", "GEX SFTP Password")
	// flag.String(GEXPrivateKeyFlag, "", "GEX Private Key")
	// flag.String(GEXSFTPHostKeyFlag, "", "GEX SFTP Host Key")
	// flag.String(GEXSFTP997PickupDirectory, "", "GEX 997 SFTP Pickup Directory")
	// flag.String(GEXSFTP824PickupDirectory, "", "GEX 834 SFTP Pickup Directory")
	flag.String(SFTPTPPSPaidInvoiceReportPickupDirectory, "", "TPPS Paid Invoice SFTP Pickup Directory")
}

// CheckTPPSSFTP validates TPPS SFTP command line flags
func CheckTPPSSFTP(v *viper.Viper) error {
	return nil
}
