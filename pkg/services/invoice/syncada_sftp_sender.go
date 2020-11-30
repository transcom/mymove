package invoice

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
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
}

// NewSyncadaSFTPSession creates a new SyncadaSFTPSession service object
func NewSyncadaSFTPSession(port string, userID string, remote string, password string, syncadaInboundDirectory string) services.SyncadaSFTPSender {
	return &SyncadaSenderSFTPSession{
		port,
		userID,
		remote,
		password,
		syncadaInboundDirectory,
	}
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

	return NewSyncadaSFTPSession(port, userID, ipAddress, password, inboundDir), nil
}

// SendToSyncadaViaSFTP copies specified local content to Syncada's SFTP server
func (s *SyncadaSenderSFTPSession) SendToSyncadaViaSFTP(localDataReader io.Reader, syncadaFileName string) (int64, error) {
	config := &ssh.ClientConfig{
		User: s.userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.password),
		},

		//RA Summary: gosec - G106 - Audit the use of ssh.InsecureIgnoreHostKey
		//RA: The linter is flagging this line of code because we are setting insecure ignore host key.
		//RA: The hostKey was removed because authentication is performed using a user ID and password
		//RA: If hostKey configuration is needed, please see PR #5039: https://github.com/transcom/mymove/pull/5039
		//RA Developer Status: {RA Request, RA Accepted, POA&M Request, POA&M Accepted, Mitigated, Need Developer Fix, False Positive, Bad Practice}
		//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
		//RA Validator: jneuner@mitre.org
		//RA Modified Severity:
		// #nosec G106
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	connection, err := ssh.Dial("tcp", s.remote+":"+s.port, config)
	if err != nil {
		return 0, err
	}
	defer connection.Close()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		return 0, err
	}
	defer client.Close()

	// create destination file
	syncadaFilePath := fmt.Sprintf("/%s/%s/%s", s.userID, s.syncadaInboundDirectory, syncadaFileName)
	syncadaFile, err := client.Create(syncadaFilePath)
	if err != nil {
		return 0, err
	}
	defer syncadaFile.Close()

	// copy source file to destination file
	bytes, err := io.Copy(syncadaFile, localDataReader)
	if err != nil {
		return 0, err
	}

	return bytes, err
}
