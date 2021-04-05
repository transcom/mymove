package invoice

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/gobuffalo/pop/v5"
	"go.uber.org/zap"

	"github.com/pkg/sftp"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// syncadaReaderSFTPSession contains information to create a new Syncada SFTP session
type syncadaReaderSFTPSession struct {
	client                     services.SFTPClient
	db                         *pop.Connection
	logger                     Logger
	deleteFilesAfterProcessing bool
}

// NewSyncadaSFTPReaderSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPReader
func NewSyncadaSFTPReaderSession(client services.SFTPClient, db *pop.Connection, logger Logger, deleteFilesAfterProcessing bool) services.SyncadaSFTPReader {
	return &syncadaReaderSFTPSession{
		client,
		db,
		logger,
		deleteFilesAfterProcessing,
	}
}

// FetchAndProcessSyncadaFiles downloads Syncada files with SFTP, processes them using the provided processor, and deletes them from the SFTP server if they were successfully processed
func (s *syncadaReaderSFTPSession) FetchAndProcessSyncadaFiles(syncadaPath string, lastRead time.Time, processor services.SyncadaFileProcessor) (time.Time, error) {
	// Store/log metrics about EDI processing upon exiting this method.
	numProcessed := 0
	start := time.Now()
	defer func() {
		ediProcessing := models.EDIProcessing{
			EDIType:          processor.EDIType(),
			ProcessStartedAt: start,
			ProcessEndedAt:   time.Now(),
			NumEDIsProcessed: numProcessed,
		}
		s.logger.Info("EDIs processed", zap.Object("EDIs processed", &ediProcessing))

		verrs, err := s.db.ValidateAndCreate(&ediProcessing)
		if err != nil {
			s.logger.Error("failed to create EDIProcessing record", zap.Error(err))
		}
		if verrs.HasAny() {
			s.logger.Error("failed to validate EDIProcessing record", zap.Error(err))
		}
	}()

	fileList, err := s.client.ReadDir(syncadaPath)
	if err != nil {
		s.logger.Error("Error reading SFTP directory", zap.String("directory", syncadaPath))
		return time.Time{}, err
	}

	// This check is intended to be temporary. We are currently getting all EDI
	// responses in the same directory. Once we start getting them in separate
	// dirs, then we won't need to inspect the file to determine the type.
	var checkEDIType *regexp.Regexp
	if processor.EDIType() == models.EDIType997 {
		checkEDIType, err = regexp.Compile(`ST\*997\*\d*`)
	} else if processor.EDIType() == models.EDIType824 {
		checkEDIType, err = regexp.Compile(`ST\*824\*\d*`)
	}
	if err != nil {
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

			if checkEDIType != nil {
				if checkEDIType.FindStringIndex(fileText) == nil {
					s.logger.Info("Skipping file that does not match current response type", zap.String("path", filePath))
					continue
				}
			}

			err = processor.ProcessFile(filePath, fileText)
			if err != nil {
				s.logger.Error("Error while processing Syncada file", zap.String("path", filePath), zap.String("file contents", fileText), zap.Error(err))
				continue
			}

			numProcessed++

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
