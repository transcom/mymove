package invoice

import (
	"bytes"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/sftp"

	"github.com/transcom/mymove/pkg/services"
)

// SyncadaReaderSFTPSession contains information to create a new Syncada SFTP session
type SyncadaReaderSFTPSession struct {
	client services.SFTPClient
	logger Logger
}

// InitNewSyncadaSFTPReaderSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPReader
func InitNewSyncadaSFTPReaderSession(client services.SFTPClient, logger Logger) services.SyncadaSFTPReader {
	return &SyncadaReaderSFTPSession{
		client,
		logger,
	}
}

// FetchAndProcessSyncadaFiles downloads Syncada files with SFTP, processes them using the provided processor, and deletes them from the SFTP server if they were successfully processed
func (s *SyncadaReaderSFTPSession) FetchAndProcessSyncadaFiles(syncadaPath string, lastRead time.Time, processor services.SyncadaFileProcessor) error {
	files, _, err := s.ReadFromSyncadaViaSFTP(syncadaPath, lastRead)
	if err != nil {
		return err
	}

	pathsToDelete := make([]string, 0, len(files))

	for _, file := range files {
		err := processor.ProcessFile(file.Path, file.Text)
		if err != nil {
			s.logger.Error("Error while processing Syncada file", zap.String("path", file.Path), zap.String("file contents", file.Text), zap.Error(err))
		} else {
			pathsToDelete = append(pathsToDelete, file.Path)
		}
	}

	s.RemoveFromSyncadaViaSFTP(pathsToDelete)

	// TODO What would be useful for us to return here? Error information? Number of files processed?
	return nil
}

// ReadFromSyncadaViaSFTP fetches contents of files on SFTP server modified after lastRead
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
				// TODO this will happen all the time and is not a big deal
				s.logger.Warn("Failed to open Syncada file over SFTP", zap.String("path", syncadaFilePath), zap.Error(err))
				continue
			}

			buf := new(bytes.Buffer)
			_, err = syncadaFile.WriteTo(buf)
			if err != nil {
				s.logger.Error("Failed to read Syncada file over SFTP", zap.String("path", syncadaFilePath), zap.Error(err))
				syncadaFile.Close()
				continue
			}
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
			s.logger.Error("Error while deleting Syncada file", zap.String("path", f))
			fileDeleteErrors = append(fileDeleteErrors, err)
		}
	}
	return fileDeleteErrors
}
