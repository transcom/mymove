package uploader

import (
	"io"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// CreateUserUploadForDocumentWrapper wrapper/helper function to create a user upload
func CreateUserUploadForDocumentWrapper(
	appCtx appcontext.AppContext, userID uuid.UUID,
	storer storage.FileStorer, file io.ReadCloser,
	filename string,
	fileSizeLimit ByteSize,
	allowedFileTypes AllowedFileTypes,
	docID *uuid.UUID,
	uploadType models.UploadType,
) (*models.UserUpload, string, *validate.Errors, error) {

	var userUploader *UserUploader

	if uploadType == models.UploadTypeUSER {
		var err error
		userUploader, err = NewUserUploader(storer, fileSizeLimit)
		if err != nil {
			appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
		}
	} else if uploadType == models.UploadTypeOFFICE {
		var err error
		userUploader, err = NewOfficeUploader(storer, fileSizeLimit)
		if err != nil {
			appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
			return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
		}
	} else {
		appCtx.Logger().Fatal("could not instantiate uploader")
		return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: "could not instantiate uploader"}
	}

	aFile, err := userUploader.PrepareFileForUpload(appCtx, file, filename)
	if err != nil {
		appCtx.Logger().Fatal("could not prepare file for uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
	}

	newUserUpload, verrs, err := userUploader.CreateUserUploadForDocument(appCtx, docID, userID, File{File: aFile}, allowedFileTypes)
	if verrs.HasAny() || err != nil {
		return nil, "", verrs, err
	}

	if storer.StorageType() == "S3" {
		// If the file storer is S3 then wait for AV to complete
		s3Key := newUserUpload.Upload.StorageKey
		pollErr := waitForAVScanToComplete(appCtx, storer, s3Key)
		if pollErr != nil {
			return nil, "", &validate.Errors{}, pollErr
		}
	}

	url, err := userUploader.PresignedURL(appCtx, newUserUpload)
	if err != nil {
		appCtx.Logger().Error("failed to get presigned url", zap.Error(err))
		return nil, "", &validate.Errors{}, err
	}

	return newUserUpload, url, &validate.Errors{}, err
}

// This is a blocking poller that will not consider the upload as "complete" until the anti virus has scanned it.
// This function should only be called in AWS environments, and is not a great permanent solution.
func waitForAVScanToComplete(
	appCtx appcontext.AppContext,
	storer storage.FileStorer,
	s3Key string,
) error {
	maxWait := 2 * time.Minute
	pollInterval := 5 * time.Second
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	timer := time.NewTimer(maxWait)
	defer timer.Stop()

	// Immediate AV scan check outside the routine
	done, err := checkAVScanStatus(appCtx, storer, s3Key)
	if err != nil {
		return err
	}
	if done {
		// Scan has completed with no errors
		return nil
	}

	for {
		select {
		case <-ticker.C:
			// Execute ticker
			done, err := checkAVScanStatus(appCtx, storer, s3Key)
			if err != nil {
				return err
			}
			if done {
				// Scan has completed with no errors
				return nil
			}
		case <-timer.C:
			// Timer has finished before the AV scan could
			errMsg := "timed out waiting for AV scan"
			appCtx.Logger().Error(errMsg, zap.String("s3Key", s3Key))
			return errors.New(errMsg)
		}
	}
}

// Returns a bool indicating the scan is done and the err if any
func checkAVScanStatus(appCtx appcontext.AppContext, storer storage.FileStorer, s3Key string) (bool, error) {
	tags, err := storer.Tags(s3Key)
	if err != nil {
		appCtx.Logger().Error("Failed to get S3 object tags", zap.Error(err))
		return false, err
	}

	status, ok := tags["av-status"]
	if !ok {
		// Bleh, keep looping AV isn't done yet
		status = "SCANNING"
	}

	switch status {
	case "CLEAN":
		// Pack it up, we're done here
		return true, nil

	case "INFECTED":
		err := errors.New("S3 object is infected")
		appCtx.Logger().Error("Uploaded S3 object is infected",
			zap.String("path", s3Key),
			zap.Error(err),
		)
		return true, err
	}
	return false, nil
}
