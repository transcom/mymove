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
func InitNewSyncadaSFTPSession() services.SyncadaSFTPSender {
	return NewSyncadaSFTPSession(os.Getenv("SYNCADA_SFTP_PORT"), os.Getenv("SYNCADA_SFTP_USER_ID"), os.Getenv("SYNCADA_SFTP_IP_ADDRESS"), os.Getenv("SYNCADA_SFTP_PASSWORD"), os.Getenv("SYNCADA_SFTP_INBOUND_DIRECTORY"))
}

// SendToSyncadaViaSFTP copies specified local content to Syncada's SFTP server
func (s *SyncadaSenderSFTPSession) SendToSyncadaViaSFTP(localDataReader io.Reader, syncadaFileName string) (int64, error) {
	config := &ssh.ClientConfig{
		User: s.userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.password),
		},
		/* nosec */
		// The hostKey was removed because authentication is performed using a user ID and password
		// If hostKey configuration is needed, please see PR #5039: https://github.com/transcom/mymove/pull/5039
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
