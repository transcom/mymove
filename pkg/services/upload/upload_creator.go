package upload

import (
	"fmt"
	"io"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

// filenameTimeFormat is the format for the timestamp we use in the filename of the upload.
// Go needs an example string when reformatting time.Time objects.
const filenameTimeFormat string = "20060102150405"

type uploadCreator struct {
	fileStorer storage.FileStorer
}

// NewUploadCreator returns a new uploadCreator
func NewUploadCreator(fileStorer storage.FileStorer) services.UploadCreator {
	return &uploadCreator{fileStorer}
}

// CreateUpload uploads a new document to an AWS S3 bucket
func (u *uploadCreator) CreateUpload(
	appCtx appcontext.AppContext,
	file io.ReadCloser,
	uploadFilename string,
	uploadType models.UploadType,
) (*models.Upload, error) {
	var upload *models.Upload
	var uploadErr error

	// If we are already in a transaction, don't start one
	if appCtx.DB().TX != nil {
		upload, uploadErr = u.createUploadTxn(appCtx, file, uploadFilename, uploadType)
	} else {
		// This error is ignored because the value is saved directly to the variable defined outside of the
		// transaction function, uploadError
		_ = appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			upload, uploadErr = u.createUploadTxn(txnAppCtx, file, uploadFilename, uploadType)
			if uploadErr != nil {
				return uploadErr
			}
			return nil
		})
	}
	if uploadErr != nil {
		return nil, uploadErr
	}

	return upload, nil
}

// createUploadTxn contains the bare code to create an models.Upload record from within or without a transaction
func (u *uploadCreator) createUploadTxn(
	appCtx appcontext.AppContext,
	file io.ReadCloser,
	uploadFilename string,
	uploadType models.UploadType,
) (*models.Upload, error) {

	newUploader, err := uploader.NewUploader(u.fileStorer, uploader.MaxFileSizeLimit, uploadType)
	if err != nil {
		if err == uploader.ErrFileSizeLimitExceedsMax {
			return nil, services.NewBadDataError(err.Error()) //todo - improve this messaging
		}
		return nil, err
	}

	fileName := time.Now().Format(filenameTimeFormat) + "-" + uploadFilename

	aFile, err := newUploader.PrepareFileForUpload(appCtx, file, fileName)
	if err != nil {
		return nil, err
	}

	newUploader.SetUploadStorageKey(fileName)

	upload, verrs, err := newUploader.CreateUpload(appCtx, uploader.File{File: aFile}, uploader.AllowedTypesAny)
	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidCreateInputError(verrs, "Validation errors found while uploading file.")
	} else if err != nil {
		return nil, fmt.Errorf("Failure to upload file: %v", err)
	}

	return upload, nil
}
