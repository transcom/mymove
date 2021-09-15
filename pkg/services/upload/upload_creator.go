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
	fileStorer   storage.FileStorer
	allowedTypes uploader.AllowedFileTypes
}

// NewUploadCreator returns a new uploadCreator
func NewUploadCreator(fileStorer storage.FileStorer) services.UploadCreator {
	return &uploadCreator{fileStorer, uploader.AllowedTypesPDFImages}
}

// CreateUpload uploads a new document to an AWS S3 bucket
func (u *uploadCreator) CreateUpload(
	appCtx appcontext.AppContext,
	file io.ReadCloser,
	uploadFilename string,
	uploadType models.UploadType,
) (*models.Upload, error) {
	var upload *models.Upload

	txErr := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		newUploader, err := uploader.NewUploader(u.fileStorer, uploader.MaxFileSizeLimit, uploadType)
		if err != nil {
			if err == uploader.ErrFileSizeLimitExceedsMax {
				return services.NewBadDataError(err.Error()) // preserves the error message from the uploader err
			}
			return err
		}

		fileName := time.Now().Format(filenameTimeFormat) + "-" + uploadFilename
		aFile, err := newUploader.PrepareFileForUpload(txnAppCtx, file, fileName)
		if err != nil {
			return err
		}

		newUploader.SetUploadStorageKey(fileName)

		newUpload, verrs, err := newUploader.CreateUpload(txnAppCtx, uploader.File{File: aFile}, u.allowedTypes)
		if verrs != nil && verrs.HasAny() {
			return services.NewInvalidCreateInputError(verrs, "Validation errors found while uploading file.")
		} else if err != nil {
			return fmt.Errorf("Failure to upload file: %v", err)
		}

		upload = newUpload
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return upload, nil
}
