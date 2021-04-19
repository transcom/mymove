package cli

import (
	"fmt"

	"github.com/pkg/sftp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// TODO: Figure out what to do in terms of static analysis on the false positive for G101
const (
	// GEXSFTPPortFlag is the ENV var for the GEX SFTP port
	GEXSFTPPortFlag string = "gex-sftp-port"
	// GEXSFTPUserIDFlag is the ENV var for the GEX SFTP user ID
	GEXSFTPUserIDFlag string = "gex-sftp-user-id"
	// GEXSFTPIPAddressFlag is the ENV var for the GEX SFTP IP address
	GEXSFTPIPAddressFlag string = "gex-sftp-ip-address"
	// GEXSFTPPasswordFlag is the ENV var for the GEX SFTP password
	GEXSFTPPasswordFlag string = "gex-sftp-password" // #nosec G101
	// GEXSFTPHostKeyFlag is the ENV var for the GEX SFTP host key
	GEXSFTPHostKeyFlag string = "gex-sftp-host-key"
	// GEXSFTP997PickupDirectory is the ENV var for the directory where GEX delivers responses
	GEXSFTP997PickupDirectory string = "gex-sftp-997-pickup-directory"
	// GEXSFTP824PickupDirectory is the ENV var for the directory where GEX delivers responses
	GEXSFTP824PickupDirectory string = "gex-sftp-824-pickup-directory"
)

// InitGEXSFTPFlags initializes GEX SFTP command line flags
func InitGEXSFTPFlags(flag *pflag.FlagSet) {
	flag.Int(GEXSFTPPortFlag, 22, "GEX SFTP Port")
	flag.String(GEXSFTPUserIDFlag, "", "GEX SFTP User ID")
	flag.String(GEXSFTPIPAddressFlag, "localhost", "GEX SFTP IP Address")
	flag.String(GEXSFTPPasswordFlag, "", "GEX SFTP Password")
	flag.String(GEXSFTPHostKeyFlag, "", "GEX SFTP Host Key")
	flag.String(GEXSFTP997PickupDirectory, "", "GEX 997 SFTP Pickup Directory")
	flag.String(GEXSFTP824PickupDirectory, "", "GEX 834 SFTP Pickup Directory")
}

// CheckGEXSFTP validates GEX SFTP command line flags
func CheckGEXSFTP(v *viper.Viper) error {
	if err := ValidatePort(v, GEXSFTPPortFlag); err != nil {
		return err
	}

	if err := ValidateHost(v, GEXSFTPIPAddressFlag); err != nil {
		return err
	}

	return nil
}

// InitGEXSSH initializes a GEX SSH client from command line flags.
func InitGEXSSH(v *viper.Viper, logger Logger) (*ssh.Client, error) {
	userID := v.GetString(GEXSFTPUserIDFlag)
	password := v.GetString(GEXSFTPPasswordFlag)
	hostKeyString := v.GetString(GEXSFTPHostKeyFlag)
	remote := v.GetString(GEXSFTPIPAddressFlag)
	port := v.GetString(GEXSFTPPortFlag)

	logger.Info("Parsing GEX SFTP host key...")
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		logger.Error("Failed to parse GEX SFTP host key", zap.Error(err))
		return nil, fmt.Errorf("failed to parse host key %w", err)
	}
	logger.Info("...Parsing GEX SFTP host key successful")

	config := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// Connect to SSH client
	address := remote + ":" + port
	logger.Info("Connecting to GEX SSH...", zap.String("address", address))
	sshClient, err := ssh.Dial("tcp", address, config)
	if err != nil {
		logger.Error("Failed to connect to GEX SSH", zap.Error(err))
		return nil, err
	}
	logger.Info("...GEX SSH connection successful")

	return sshClient, nil
}

// InitGEXSFTP initializes a GEX SFTP client from command line flags.
func InitGEXSFTP(sshClient *ssh.Client, logger Logger) (*sftp.Client, error) {
	// Create new SFTP client
	logger.Info("Connecting to GEX SFTP...")
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		logger.Error("Failed to connect to GEX SFTP", zap.Error(err))
		return nil, err
	}
	logger.Info("...GEX SFTP connection successful")

	return client, nil
}
