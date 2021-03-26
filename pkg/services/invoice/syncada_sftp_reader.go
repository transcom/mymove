package invoice

import (
	"bytes"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/pkg/sftp"

	"github.com/transcom/mymove/pkg/services"
)

// syncadaReaderSFTPSession contains information to create a new Syncada SFTP session
type syncadaReaderSFTPSession struct {
	client                     services.SFTPClient
	logger                     Logger
	deleteFilesAfterProcessing bool
}

// NewSyncadaSFTPReaderSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPReader
func NewSyncadaSFTPReaderSession(client services.SFTPClient, logger Logger, deleteFilesAfterProcessing bool) services.SyncadaSFTPReader {
	return &syncadaReaderSFTPSession{
		client,
		logger,
		deleteFilesAfterProcessing,
	}
}

// FetchAndProcessSyncadaFiles downloads Syncada files with SFTP, processes them using the provided processor, and deletes them from the SFTP server if they were successfully processed
func (s *syncadaReaderSFTPSession) FetchAndProcessSyncadaFiles(syncadaPath string, lastRead time.Time, processor services.SyncadaFileProcessor) (time.Time, error) {
	fileList, err := s.client.ReadDir(syncadaPath)
	if err != nil {
		s.logger.Error("Error reading SFTP directory", zap.String("directory", syncadaPath))
		return time.Time{}, err
	}

	mostRecentFileModTime := lastRead

	for _, fileInfo := range fileList {
		if fileInfo.ModTime().After(lastRead) {
			if fileInfo.ModTime().After(mostRecentFileModTime) {
				mostRecentFileModTime = fileInfo.ModTime()
			}
			filePath := sftp.Join(syncadaPath, fileInfo.Name())

			fileText, err := s.downloadFile(filePath)
			if err != nil {
				s.logger.Info("Error while downloading Syncada file", zap.String("path", filePath), zap.Error(err))
				continue
			}

			err = processor.ProcessFile(filePath, fileText)
			if err != nil {
				s.logger.Error("Error while processing Syncada file", zap.String("path", filePath), zap.String("file contents", fileText), zap.Error(err))
				continue
			}

			if s.deleteFilesAfterProcessing {
				err = s.client.Remove(filePath)
				if err != nil {
					s.logger.Error("Error while deleting Syncada file", zap.String("path", filePath))
				} else {
					s.logger.Info("Deleted Syncada file", zap.String("path", filePath))
				}
			}
		}
	}

	return mostRecentFileModTime, nil
}

func (s *syncadaReaderSFTPSession) downloadFile(path string) (string, error) {
	file, err := s.client.Open(path)
	if err != nil {
		// This is expected (at least in the US Bank testing environment) because they
		// upload some files that we don't have permission to read, and we don't know
		// how to tell the different files apart yet.
		return "", fmt.Errorf("failed to open file over SFTP: %w", err)
	}

	// Note: Avoiding a defer on Close here because we want a Close to actually cause
	// an error to be returned.

	buf := new(bytes.Buffer)
	_, err = file.WriteTo(buf)
	if err != nil {
		// If close fails, just log it as we're already in an error situation.
		if closeErr := file.Close(); closeErr != nil {
			s.logger.Error("could not close file", zap.Error(closeErr))
		}
		return "", fmt.Errorf("failed to read file over SFTP: %w", err)
	}

	return buf.String(), file.Close()
}
