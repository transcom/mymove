package invoice

import (
	"os"

	"github.com/pkg/sftp"

	"github.com/transcom/mymove/pkg/services"
)

// sftpClientWrapper wraps an SFTP client to facilitate testing/mocking
type sftpClientWrapper struct {
	client *sftp.Client
}

// NewSFTPClientWrapper initializes a new SFTPClientWrapper and returns a testable SFTPClient
func NewSFTPClientWrapper(client *sftp.Client) services.SFTPClient {
	return &sftpClientWrapper{
		client,
	}
}

func (s sftpClientWrapper) ReadDir(p string) ([]os.FileInfo, error) {
	return s.client.ReadDir(p)
}

func (s sftpClientWrapper) Open(path string) (services.SFTPFiler, error) {
	// Note that this method returns a generic SFTPFiler interface as opposed to an explicit sftp.File
	// so we can inject a mock version of the file for testing.
	return s.client.Open(path)
}

func (s sftpClientWrapper) Remove(path string) error {
	return s.client.Remove(path)
}
