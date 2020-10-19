package cli

import "github.com/spf13/pflag"

const (
	// SyncadaSFTPUserIDFlag is the Syncada SFTP User ID Flag
	SyncadaSFTPUserIDFlag string = "syncada-sftp-user-id"
	// SyncadaSFTPPsswrdFlag is the Syncada SFTP Password Flag
	SyncadaSFTPPsswrdFlag string = "syncada-sftp-password"
	// SyncadaSFTPIPAddressFlag is the Syncada SFTP IP Address Flag
	SyncadaSFTPIPAddressFlag string = "syncada-sftp-ip-address"
	// SyncadaSFTPInboundDirectoryFlag is the Syncada SFTP Inbound Directory Flag
	SyncadaSFTPInboundDirectoryFlag string = "syncada-sftp-inbound-directory"
	// SyncadaSFTPOutboundDirectoryFlag is the Syncada SFTP Outbound Directory Flag
	SyncadaSFTPOutboundDirectoryFlag string = "syncada-sftp-outbound-directory"
	// SyncadaSFTPPortFlag is the Syncada SFTP Port Flag
	SyncadaSFTPPortFlag string = "syncada-sftp-port"
	// SyncadaSFTPTransportURLFlag is the Syncada SFTP Transport URL Flag
	SyncadaSFTPTransportURLFlag string = "syncada-sftp-transport-url"
	// SyncadaSFTPRoutingIDFlag is the Syncada SFTP Routing ID Flag
	SyncadaSFTPRoutingIDFlag string = "syncada-sftp-routing-id"
)

// InitSyncadaFlags initializes the Syncada SFTP command line flags
func InitSyncadaFlags(flag *pflag.FlagSet) {
	flag.String(SyncadaSFTPUserIDFlag, "", "Syncada SFTP User ID")
	flag.String(SyncadaSFTPPsswrdFlag, "", "Syncada SFTP Password")
	flag.String(SyncadaSFTPIPAddressFlag, "", "Syncada SFTP IP Address")
	flag.String(SyncadaSFTPInboundDirectoryFlag, "", "Syncada SFTP Inbound Directory")
	flag.String(SyncadaSFTPOutboundDirectoryFlag, "", "Syncada SFTP Outbound Directory")
	flag.String(SyncadaSFTPPortFlag, "", "Syncada SFTP Port")
	flag.String(SyncadaSFTPTransportURLFlag, "", "Syncada SFTP Transport URL")
	flag.String(SyncadaSFTPRoutingIDFlag, "", "Syncada SFTP Rotuing ID")
}
