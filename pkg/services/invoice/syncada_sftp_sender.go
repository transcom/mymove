//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to close an asynchronous connection
//RA: Given the functions causing the lint errors are used close an asynchronous connection in order to prevent it
//RA: from running indefinitely, it is not deemed a risk
//RA Developer Status: Mitigated
//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
//RA Validator: jneuner@mitre.org
//RA Modified Severity:
// nolint:errcheck
package invoice

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/transcom/mymove/pkg/services"
)

// SyncadaSenderSFTPSession contains information to create a new Syncada SFTP session
type SyncadaSenderSFTPSession struct {
	port                    string
	userID                  string
	remote                  string
	password                string
	syncadaInboundDirectory string
	hostKey                 ssh.PublicKey
}

// InitNewSyncadaSFTPSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPSender
func InitNewSyncadaSFTPSession() (services.SyncadaSFTPSender, error) {
	port := os.Getenv("SYNCADA_SFTP_PORT")
	if port == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_PORT")
	}

	userID := os.Getenv("SYNCADA_SFTP_USER_ID")
	if userID == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_USER_ID")
	}

	ipAddress := os.Getenv("SYNCADA_SFTP_IP_ADDRESS")
	if ipAddress == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_IP_ADDRESS")
	}

	password := os.Getenv("SYNCADA_SFTP_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_PASSWORD")
	}

	inboundDir := os.Getenv("SYNCADA_SFTP_INBOUND_DIRECTORY")
	if inboundDir == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_INBOUND_DIRECTORY")
	}

	hostKeyString := os.Getenv("SYNCADA_SFTP_HOST_KEY")
	if hostKeyString == "" {
		return nil, fmt.Errorf("Invalid credentials sftp missing SYNCADA_SFTP_HOST_KEY")
	}
	hostKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(hostKeyString))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse host key %w", err)
	}

	return &SyncadaSenderSFTPSession{
		port,
		userID,
		ipAddress,
		password,
		inboundDir,
		hostKey,
	}, nil
}

// SendToSyncadaViaSFTP copies specified local content to Syncada's SFTP server
func (s *SyncadaSenderSFTPSession) SendToSyncadaViaSFTP(localDataReader io.Reader, syncadaFileName string) (int64, error) {
	config := &ssh.ClientConfig{
		User: s.userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.password),
		},
		HostKeyCallback: ssh.FixedHostKey(s.hostKey),
	}

	// connect
	connection, err := ssh.Dial("tcp", s.remote+":"+s.port, config)
	if err != nil {
		return 0, err
	}

	defer func() {
		if closeErr := connection.Close(); closeErr != nil {
			fmt.Println("Failed to close connection", zap.Error(closeErr))
		}
	}()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return 0, err
	}

	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			fmt.Println("Failed to close SFTP client", zap.Error(closeErr))
		}
	}()

	// create destination file
	syncadaFilePath := fmt.Sprintf("/%s/%s/%s", s.userID, s.syncadaInboundDirectory, syncadaFileName)
	syncadaFile, err := client.Create(syncadaFilePath)
	if err != nil {
		return 0, err
	}

	defer func() {
		if closeErr := syncadaFile.Close(); closeErr != nil {
			fmt.Println("Failed to close Syncada destination file", zap.Error(closeErr))
		}
	}()

	// copy source file to destination file
	bytes, err := io.Copy(syncadaFile, localDataReader)
	if err != nil {
		return 0, err
	}

	return bytes, err
}
