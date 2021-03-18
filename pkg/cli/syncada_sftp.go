package cli

import (
	"fmt"

	"github.com/pkg/sftp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// TODO: Figure out what to do in terms of static analysis on the false positive for G101
const (
	// SyncadaSFTPPortFlag is the ENV var for the Syncada SFTP port
	SyncadaSFTPPortFlag string = "syncada-sftp-port"
	// SyncadaSFTPUserIDFlag is the ENV var for the Syncada SFTP user ID
	SyncadaSFTPUserIDFlag string = "syncada-sftp-user-id"
	// SyncadaSFTPIPAddressFlag is the ENV var for the Syncada SFTP IP address
	SyncadaSFTPIPAddressFlag string = "syncada-sftp-ip-address"
	// SyncadaSFTPPasswordFlag is the ENV var for the Syncada SFTP password
	SyncadaSFTPPasswordFlag string = "syncada-sftp-password" // #nosec G101
	// SyncadaSFTPHostKeyFlag is the ENV var for the Syncada SFTP host key
	SyncadaSFTPHostKeyFlag string = "syncada-sftp-host-key"
)

// InitSyncadaSFTPFlags initializes Syncada SFTP command line flags
func InitSyncadaSFTPFlags(flag *pflag.FlagSet) {
	flag.Int(SyncadaSFTPPortFlag, 22, "Syncada SFTP Port")
	flag.String(SyncadaSFTPUserIDFlag, "", "Syncada SFTP User ID")
	flag.String(SyncadaSFTPIPAddressFlag, "", "Syncada SFTP IP Address")
	flag.String(SyncadaSFTPPasswordFlag, "", "Syncada SFTP Password")
	flag.String(SyncadaSFTPHostKeyFlag, "", "Syncada SFTP Host Key")
}

// CheckSyncadaSFTP validates Syncada SFTP command line flags
func CheckSyncadaSFTP(v *viper.Viper) error {
	if err := ValidatePort(v, SyncadaSFTPPortFlag); err != nil {
		return err
	}

	if err := ValidateHost(v, SyncadaSFTPIPAddressFlag); err != nil {
		return err
	}

	return nil
}

// InitSyncadaSFTP initializes a Syncada SFTP client from command line flags.
func InitSyncadaSFTP(v *viper.Viper, logger Logger) (*sftp.Client, error) {
	userID := v.GetString(SyncadaSFTPUserIDFlag)
	password := v.GetString(SyncadaSFTPPasswordFlag)
	hostKeyString := v.GetString(SyncadaSFTPHostKeyFlag)
	remote := v.GetString(SyncadaSFTPIPAddressFlag)
	port := v.GetString(SyncadaSFTPPortFlag)

	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		return nil, fmt.Errorf("failed to parse host key %w", err)
	}

	config := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	connection, err := ssh.Dial("tcp", remote+":"+port, config)
	if err != nil {
		return nil, err
	}

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return nil, err
	}

	return client, nil
}
