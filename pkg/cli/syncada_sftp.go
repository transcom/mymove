package cli

import (
	"fmt"

	"github.com/pkg/sftp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

const (
	// SyncadaSFTPPortFlag is the ENV var for the Syncada SFTP port
	SyncadaSFTPPortFlag string = "syncada-sftp-port"
	// SyncadaSFTPUserIDFlag is the ENV var for the Syncada SFTP user ID
	SyncadaSFTPUserIDFlag string = "syncada-sftp-user-id"
	// SyncadaSFTPIPAddressFlag is the ENV var for the Syncada SFTP IP address
	SyncadaSFTPIPAddressFlag string = "syncada-sftp-ip-address"

	//RA Summary: gosec - G101 - Password Management: Hardcoded Password
	//RA: This line was flagged because of use of the word "password"
	//RA: This line is used to identify the name of the flag. SyncadaSFTPPasswordFlag is the Syncada SFTP Password Flag.
	//RA: See MB-7727 and MB-7728 for tracking future work to resolve this issue
	//RA: App should implement public-key authentication; issue remains open while interface control is negotiated for this connection.
	//RA Developer Status: Mitigated
	//RA Validator Status: Known Issue
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity: CAT III

	// SyncadaSFTPPasswordFlag is the ENV var for the Syncada SFTP password
	SyncadaSFTPPasswordFlag string = "syncada-sftp-password" // #nosec G101
	// SyncadaSFTPHostKeyFlag is the ENV var for the Syncada SFTP host key
	SyncadaSFTPHostKeyFlag string = "syncada-sftp-host-key"
	// SyncadaSFTPOutboundDirectory is the ENV var for the directory where Syncada uploads responses
	SyncadaSFTPOutboundDirectory string = "syncada-sftp-outbound-directory"
)

// InitSyncadaSFTPFlags initializes Syncada SFTP command line flags
func InitSyncadaSFTPFlags(flag *pflag.FlagSet) {
	flag.Int(SyncadaSFTPPortFlag, 22, "Syncada SFTP Port")
	flag.String(SyncadaSFTPUserIDFlag, "", "Syncada SFTP User ID")
	flag.String(SyncadaSFTPIPAddressFlag, "localhost", "Syncada SFTP IP Address")
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

// InitSyncadaSSH initializes a Syncada SSH client from command line flags.
func InitSyncadaSSH(v *viper.Viper, logger Logger) (*ssh.Client, error) {
	userID := v.GetString(SyncadaSFTPUserIDFlag)
	password := v.GetString(SyncadaSFTPPasswordFlag)
	hostKeyString := v.GetString(SyncadaSFTPHostKeyFlag)
	remote := v.GetString(SyncadaSFTPIPAddressFlag)
	port := v.GetString(SyncadaSFTPPortFlag)

	logger.Info("Parsing Syncada SFTP host key...")
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		logger.Error("Failed to parse Syncada SFTP host key", zap.Error(err))
		return nil, fmt.Errorf("failed to parse host key %w", err)
	}
	logger.Info("...Parsing Syncada SFTP host key successful")

	config := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Connect to SSH client
	address := remote + ":" + port
	logger.Info("Connecting to Syncada SSH...", zap.String("address", address))
	sshClient, err := ssh.Dial("tcp", address, config)
	if err != nil {
		logger.Error("Failed to connect to Syncada SSH", zap.Error(err))
		return nil, err
	}
	logger.Info("...Syncada SSH connection successful")

	return sshClient, nil
}

// InitSyncadaSFTP initializes a Syncada SFTP client from command line flags.
func InitSyncadaSFTP(sshClient *ssh.Client, logger Logger) (*sftp.Client, error) {
	// Create new SFTP client
	logger.Info("Connecting to Syncada SFTP...")
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		logger.Error("Failed to connect to Syncada SFTP", zap.Error(err))
		return nil, err
	}
	logger.Info("...Syncada SFTP connection successful")

	return client, nil
}
