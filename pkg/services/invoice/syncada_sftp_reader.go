package invoice

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/sftp"

	"github.com/transcom/mymove/pkg/services"
)

// SyncadaReaderSFTPSession contains information to create a new Syncada SFTP session
// TODO should i just modify the SyncadaSenderSFTPSession to have a more generic name for the directory and reuse it?
// TODO Is it worth keeping the client in this struct or should i just pass as an arg?
type SyncadaReaderSFTPSession struct {
	client services.SFTPClient
}

// InitNewSyncadaSFTPReaderSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPReader
func InitNewSyncadaSFTPReaderSession(client services.SFTPClient) services.SyncadaSFTPReader {
	return &SyncadaReaderSFTPSession{
		client: client,
	}
}

// ReadFromSyncadaViaSFTP fetches contents of files on SFTP server modified after lastRead
// TODO should i be returning the connection? or taking a connection as input?
// TODO we've also got both a connection and a client and they both want to be closed, should i be passing both?
func (s *SyncadaReaderSFTPSession) ReadFromSyncadaViaSFTP(syncadaPath string, lastRead time.Time) ([]services.RawSyncadaFile, time.Time, error) {
	fileList, err := s.client.ReadDir(syncadaPath)

	if err != nil {
		return []services.RawSyncadaFile{}, time.Time{}, err
	}

	readFiles := make([]services.RawSyncadaFile, 0, len(fileList))

	mostRecentFileModTime := lastRead

	for _, f := range fileList {
		if f.ModTime().After(lastRead) {
			if f.ModTime().After(mostRecentFileModTime) {
				mostRecentFileModTime = f.ModTime()
			}
			syncadaFilePath := sftp.Join(syncadaPath, f.Name())

			syncadaFile, err := s.client.Open(syncadaFilePath)
			if err != nil {
				fmt.Printf("Failed to open %s: %s\n", syncadaFilePath, err.Error())
				continue
			}

			buf := new(bytes.Buffer)
			_, err = syncadaFile.WriteTo(buf)
			if err != nil {
				fmt.Printf("Failed to read %s: %s\n", syncadaFilePath, err.Error())
				syncadaFile.Close()
				continue
			}
			// TODO should this be done with defer? Not sure since we're in a loop
			syncadaFile.Close()

			fd := services.RawSyncadaFile{Path: syncadaFilePath, Text: buf.String()}
			readFiles = append(readFiles, fd)
		}
	}

	return readFiles, mostRecentFileModTime, nil
}

// TODO what should we do with errors from this function?
// TODO I think i'd want to unify them into one error for ease of checking
// TODO I think it makes sense for this to try each file even if some of them fail.
// TODO seems like the most likely causes for errors would be if the files were already deleted or if the connection was dropped
// TODO Also, this function could be somewhere else but it's so small

// RemoveFromSyncadaViaSFTP attempts to remove every file path passed to it from an SFTP server
func (s *SyncadaReaderSFTPSession) RemoveFromSyncadaViaSFTP(filePaths []string) []error {
	fileDeleteErrors := make([]error, 0, len(filePaths))
	for _, f := range filePaths {
		err := s.client.Remove(f)
		if err != nil {
			fileDeleteErrors = append(fileDeleteErrors, err)
		}
	}
	return fileDeleteErrors
}
