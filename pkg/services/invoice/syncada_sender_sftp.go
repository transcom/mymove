package invoice

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SyncadaSenderSFTPSession contains information to create a new Syncada SFTP session
type SyncadaSenderSFTPSession struct {
	port                     string
	userID                   string
	remote                   string
	password                 string
	destinationFileDirectory string
}

// NewSyncadaSFTPSession creates a new SyncadaSFTPSession service object
func NewSyncadaSFTPSession(port string, userID string, remote string, password string, destinationFileDirectory string) SyncadaSenderSFTPSession {
	return SyncadaSenderSFTPSession{
		port,
		userID,
		remote,
		password,
		destinationFileDirectory,
	}
}

// SendToSyncada converts a speicified file to a string and copies it to Syncada's SFTP server
func (s *SyncadaSenderSFTPSession) SendToSyncada(localFilePath string, destinationFileName string) (resp string, err error) {
	config := &ssh.ClientConfig{
		User: s.userID,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.password),
		},
		/* #nosec */
		// The hostKey was removed because authentication is performed using a user ID and password
		// If hostKey configuration is needed, please see PR #5039: https://github.com/transcom/mymove/pull/5039
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	connection, err := ssh.Dial("tcp", s.remote+":"+s.port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	// create new SFTP client
	client, err := sftp.NewClient(connection)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// open local file
	localFile, err := os.Open(filepath.Clean(localFilePath))
	if err != nil {
		log.Fatal(err)
	}

	// create destination file
	destinationFilePath := fmt.Sprintf("/%s/%s/%s", s.userID, s.destinationFileDirectory, destinationFileName)
	destinationFile, err := client.Create(destinationFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer destinationFile.Close()

	// copy source file to destination file
	bytes, err := io.Copy(destinationFile, localFile)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%d bytes copied\n", bytes), err
}
