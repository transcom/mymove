package invoice

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pkg/sftp"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// syncadaReaderSFTPSession contains information to create a new SFTP session
type syncadaReaderSFTPSession struct {
	client                     services.SFTPClient
	deleteFilesAfterProcessing bool
}

// NewSyncadaSFTPReaderSession initialize a NewSyncadaSFTPSession and return services.SyncadaSFTPReader
func NewSyncadaSFTPReaderSession(client services.SFTPClient, deleteFilesAfterProcessing bool) services.SyncadaSFTPReader {
	return &syncadaReaderSFTPSession{
		client,
		deleteFilesAfterProcessing,
	}
}

// FetchAndProcessSyncadaFiles downloads Syncada files with SFTP, processes them using the provided processor, and deletes them from the SFTP server if they were successfully processed
func (s *syncadaReaderSFTPSession) FetchAndProcessSyncadaFiles(appCtx appcontext.AppContext, pickupPath string, lastRead time.Time, processor services.SyncadaFileProcessor) (time.Time, error) {
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
		appCtx.Logger().Info("EDIs processed", zap.Object("edisProcessed", &ediProcessing))

		verrs, err := appCtx.DB().ValidateAndCreate(&ediProcessing)
		if err != nil {
			appCtx.Logger().Error("failed to create EDIProcessing record", zap.Error(err))
		}
		if verrs.HasAny() {
			appCtx.Logger().Error("failed to validate EDIProcessing record", zap.Error(err))
		}
	}()

	fileList, err := s.client.ReadDir(pickupPath)
	if err != nil {
		appCtx.Logger().Error("Error reading SFTP directory", zap.String("directory", pickupPath))
		return time.Time{}, err
	}

	mostRecentFileModTime := lastRead

	for _, fileInfo := range fileList {
		if fileInfo.ModTime().After(lastRead) {
			if fileInfo.ModTime().After(mostRecentFileModTime) {
				mostRecentFileModTime = fileInfo.ModTime()
			}
			filePath := sftp.Join(pickupPath, fileInfo.Name())

			fileText, err := s.downloadFile(appCtx, filePath)
			if err != nil {
				appCtx.Logger().Info("Error while downloading Syncada file", zap.String("path", filePath), zap.Error(err))
				continue
			}

			err = processor.ProcessFile(appCtx, filePath, fileText)
			if err != nil {
				appCtx.Logger().Error("Error while processing Syncada file", zap.String("path", filePath), zap.String("file contents", fileText), zap.Error(err))
				continue
			}

			numProcessed++

			if s.deleteFilesAfterProcessing {
				err = s.client.Remove(filePath)
				if err != nil {
					appCtx.Logger().Error("Error while deleting Syncada file", zap.String("path", filePath))
				} else {
					appCtx.Logger().Info("Deleted Syncada file", zap.String("path", filePath))
				}
			} else {
				appCtx.Logger().Info("Delete sftp files: false", zap.String("path", filePath))
			}
		}
	}

	return mostRecentFileModTime, nil
}

func (s *syncadaReaderSFTPSession) downloadFile(appCtx appcontext.AppContext, path string) (string, error) {
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
			appCtx.Logger().Error("could not close file", zap.Error(closeErr))
		}
		return "", fmt.Errorf("failed to read file over SFTP: %w", err)
	}

	return buf.String(), file.Close()
}
