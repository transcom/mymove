package upload

import (
	"io"
	"strings"
	"time"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/uploader"
)

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
				verrs := validate.NewErrors()
				verrs.Add("file", err.Error())
				return apperror.NewInvalidCreateInputError(verrs, "File cannot be uploaded.")
			}
			return err
		}

		// Prefix the filename with a timestamp for uniqueness
		fileName := assembleUploadFilePathName(uploadFilename)
		aFile, err := newUploader.PrepareFileForUpload(txnAppCtx, file, fileName)
		if err != nil {
			return err
		}

		newUploader.SetUploadStorageKey(fileName)

		newUpload, verrs, err := newUploader.CreateUpload(txnAppCtx, uploader.File{File: aFile}, u.allowedTypes)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidCreateInputError(verrs, "Validation errors found while uploading file.")
		} else if err != nil {
			return apperror.NewQueryError("Upload", err, "Failed to upload file")
		}

		upload = newUpload
		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return upload, nil
}

// filenameTimeFormat is the format for the timestamp we use in the filename of the upload.
// Go needs an example string when reformatting time.Time objects.
const filenameTimeFormat string = "20060102150405"

// assembleUploadFilePathName puts a timestamp prefix on the file name while preserving the rest of the path
func assembleUploadFilePathName(filePathName string) string {
	splitPath := strings.Split(filePathName, "/")

	// The last element in the slice will be the actual file name
	fileName := splitPath[len(splitPath)-1]

	// Replace the actual file name with a timestamped version, to ensure uniqueness
	splitPath[len(splitPath)-1] = time.Now().Format(filenameTimeFormat) + "-" + fileName

	// Reconnect the file path name and return the whole string
	return strings.Join(splitPath, "/")
}
