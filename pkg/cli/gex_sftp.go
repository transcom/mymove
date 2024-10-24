package cli

import (
	"fmt"

	"github.com/pkg/sftp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

// Set of flags used for GEXSFTP
const (
	// GEXSFTPPortFlag is the ENV var for the GEX SFTP port
	GEXSFTPPortFlag string = "gex-sftp-port"
	// GEXSFTPUserIDFlag is the ENV var for the GEX SFTP user ID
	GEXSFTPUserIDFlag string = "gex-sftp-user-id"
	// GEXSFTPIPAddressFlag is the ENV var for the GEX SFTP IP address
	GEXSFTPIPAddressFlag string = "gex-sftp-ip-address"
	//RA Summary: gosec - G101 - Password Management: Hardcoded Password
	//RA: This line was flagged because of use of the word "password"
	//RA: This line is used to identify the name of the flag. GEXSFTPPasswordFlag is the GEX SFTP Password Flag.
	//RA: See MB-7727 and MB-7728 for tracking future work to resolve this issue
	//RA: App should implement public-key authentication; issue remains open while interface control is negotiated for this connection.
	//RA Developer Status: Mitigated
	//RA Validator Status: Known Issue
	//RA Validator: leodis.f.scott.civ@mail.mil
	//RA Modified Severity: CAT III
	// #nosec G101
	// GEXSFTPPasswordFlag is the ENV var for the GEX SFTP password
	GEXSFTPPasswordFlag string = "gex-sftp-password"
	// GEXXPrivateKeyFlag is the ENV var for the private key which is used in establishing an
	// ssh connection to the GEX server. The GEX server has the public key.
	GEXPrivateKeyFlag string = "gex-private-key"
	// GEXSFTPHostKeyFlag is the ENV var for the GEX SFTP host key
	GEXSFTPHostKeyFlag string = "gex-sftp-host-key"
	// GEXSFTP997PickupDirectory is the ENV var for the directory where GEX delivers responses
	GEXSFTP997PickupDirectory string = "gex-sftp-997-pickup-directory"
	// GEXSFTP824PickupDirectory is the ENV var for the directory where GEX delivers responses
	GEXSFTP824PickupDirectory string = "gex-sftp-824-pickup-directory"
)

// Pending completion of B-20560, uncomment the code below
/*
// Set of flags used for SFTPTPPSPaid
const (
	// SFTPTPPSPaidInvoiceReportPickupDirectory is the ENV var for the directory where TPPS delivers the TPPS paid invoice report
	SFTPTPPSPaidInvoiceReportPickupDirectory string = "pending" // pending completion of B-20560
)
*/

// InitGEXSFTPFlags initializes GEX SFTP command line flags
func InitGEXSFTPFlags(flag *pflag.FlagSet) {
	flag.Int(GEXSFTPPortFlag, 22, "GEX SFTP Port")
	flag.String(GEXSFTPUserIDFlag, "", "GEX SFTP User ID")
	flag.String(GEXSFTPIPAddressFlag, "localhost", "GEX SFTP IP Address")
	flag.String(GEXSFTPPasswordFlag, "", "GEX SFTP Password")
	flag.String(GEXPrivateKeyFlag, "", "GEX Private Key")
	flag.String(GEXSFTPHostKeyFlag, "", "GEX SFTP Host Key")
	flag.String(GEXSFTP997PickupDirectory, "", "GEX 997 SFTP Pickup Directory")
	flag.String(GEXSFTP824PickupDirectory, "", "GEX 834 SFTP Pickup Directory")
	// flag.String(SFTPTPPSPaidInvoiceReportPickupDirectory, "", "TPPS Paid Invoice SFTP Pickup Directory") // pending completion of B-20560
}

// CheckGEXSFTP validates GEX SFTP command line flags
func CheckGEXSFTP(v *viper.Viper) error {

	port := v.GetString(GEXSFTPPortFlag)
	if port == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_PORT")
	}

	if err := ValidatePort(v, GEXSFTPPortFlag); err != nil {
		return err
	}

	hostKeyString := v.GetString(GEXSFTPHostKeyFlag)
	if hostKeyString == "" {
		return fmt.Errorf("Invalid credentials sftp missing GEX_SFTP_HOST_KEY")
	}
	_, _, _, _, hostKeyErr := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if hostKeyErr != nil {
		return hostKeyErr
	}

	userID := v.GetString(GEXSFTPUserIDFlag)
	if userID == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_USER_ID")
	}

	remote := v.GetString(GEXSFTPIPAddressFlag)
	if remote == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_IP_ADDRESS")
	}

	password := v.GetString(GEXSFTPPasswordFlag)
	if password == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_SFTP_PASSWORD")
	}

	privateKeyString := v.GetString(GEXPrivateKeyFlag)
	if privateKeyString == "" {
		return fmt.Errorf("Invalid credentials SFTP missing GEX_PRIVATE_KEY")
	}

	_, signerErr := ssh.ParsePrivateKey([]byte(privateKeyString))
	if signerErr != nil {
		return signerErr
	}

	return ValidateHost(v, GEXSFTPIPAddressFlag)
}

// InitGEXSSH initializes a GEX SSH client from command line flags.
func InitGEXSSH(logger *zap.Logger, v *viper.Viper) (*ssh.Client, error) {
	userID := v.GetString(GEXSFTPUserIDFlag)
	password := v.GetString(GEXSFTPPasswordFlag)
	hostKeyString := v.GetString(GEXSFTPHostKeyFlag)
	remote := v.GetString(GEXSFTPIPAddressFlag)
	port := v.GetString(GEXSFTPPortFlag)
	privateKeyString := v.GetString(GEXPrivateKeyFlag)

	CheckOutboundIP(logger)

	logger.Info("Parsing GEX SFTP host key...")
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		logger.Error("Failed to parse GEX SFTP host key", zap.Error(err))
		return nil, fmt.Errorf("failed to parse host key %w", err)
	}
	logger.Info("...Parsing GEX SFTP host key successful")

	logger.Info("Parsing GEX SFTP private key...")

	signer, err := ssh.ParsePrivateKey([]byte(privateKeyString))
	if err != nil {
		logger.Error("Failed to parse GEX SFTP private key", zap.Error(err))
		return nil, fmt.Errorf("failed to parse private key %w", err)
	}

	logger.Info("...Parsing GEX SFTP private key successful")

	config := &ssh.ClientConfig{
		User: userID,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
			// Fall back to the password if the private key doesn't work.
			ssh.Password(password),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	// Connect to SSH client
	address := remote + ":" + port

	logger.Info("Connecting to GEX SSH...", zap.String("destination_address", address))

	sshClient, err := ssh.Dial("tcp", address, config)
	if err != nil {
		logger.Error("Failed to connect to GEX SSH", zap.Error(err))
		return nil, err
	}
	logger.Info("...GEX SSH connection successful")

	return sshClient, nil
}

// InitGEXSFTP initializes a GEX SFTP client from command line flags.
func InitGEXSFTP(logger *zap.Logger, sshClient *ssh.Client) (*sftp.Client, error) {
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
