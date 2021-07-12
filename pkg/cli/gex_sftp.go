package cli

import (
	"fmt"
	"net"

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

// GetLocalIP returns the non loopback local IP of the host
// This might be source address, however, we can also try to use the
// aws sdk to get the IP of the load balancer, since I think that's where the connection request
// would be coming from. https://docs.aws.amazon.com/sdk-for-go/api/service/elbv2/#ELBV2.DescribeLoadBalancerAttributes
// The only way I could possibly see which is which is to implement both methods, and compare the IPs, and see which matches
// the IP needed for troubleshooting?
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
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

	logger.Info("Connecting to GEX SSH...", zap.String("destination_address", address), zap.String("source_address", GetLocalIP()))

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
