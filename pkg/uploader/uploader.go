package uploader

import (
	"io"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// ErrZeroLengthFile represents an error caused by a file with no content
var ErrZeroLengthFile = errors.New("File has length of 0")

// Uploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type Uploader struct {
	db     *pop.Connection
	logger *zap.Logger
	Storer storage.FileStorer
}

// NewUploader creates and returns a new uploader
func NewUploader(db *pop.Connection, logger *zap.Logger, storer storage.FileStorer) *Uploader {
	return &Uploader{
		db:     db,
		logger: logger,
		Storer: storer,
	}
}

// CreateUpload creates a new Upload by performing validations, storing the specified
// file using the supplied storer, and saving an Upload object to the database containing
// the file's metadata.
func (u *Uploader) CreateUpload(documentID *uuid.UUID, userID uuid.UUID, file afero.File) (*models.Upload, *validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	info, err := file.Stat()
	if err != nil {
		u.logger.Error("Could not get file info", zap.Error(err))
	}

	if info.Size() == 0 {
		return nil, responseVErrors, ErrZeroLengthFile
	}

	contentType, err := storage.DetectContentType(file)
	if err != nil {
		u.logger.Error("Could not detect content type", zap.Error(err))
		return nil, responseVErrors, err
	}

	checksum, err := storage.ComputeChecksum(file)
	if err != nil {
		u.logger.Error("Could not compute checksum", zap.Error(err))
		return nil, responseVErrors, err
	}

	id := uuid.Must(uuid.NewV4())

	newUpload := &models.Upload{
		ID:          id,
		DocumentID:  documentID,
		UploaderID:  userID,
		Filename:    file.Name(),
		Bytes:       info.Size(),
		ContentType: contentType,
		Checksum:    checksum,
	}

	u.db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")

		verrs, err := db.ValidateAndCreate(newUpload)
		if err != nil || verrs.HasAny() {
			u.logger.Error("Error creating new upload", zap.Error(err))
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error creating new upload")
			return transactionError
		}

		// Push file to S3
		if _, err := u.Storer.Store(newUpload.StorageKey, file, checksum); err != nil {
			u.logger.Error("failed to store object", zap.Error(err))
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "failed to store object")
			return transactionError
		}

		u.logger.Info("created an upload with id and key ", zap.Any("new_upload_id", newUpload.ID), zap.String("key", newUpload.StorageKey))
		return nil

	})

	return newUpload, responseVErrors, responseError
}

// CreateUploadNoDocument stores Upload but does not create a Document
func (u *Uploader) CreateUploadNoDocument(userID uuid.UUID, aFile *afero.File) (*models.Upload, *validate.Errors, error) {
	return u.CreateUpload(nil, userID, *aFile)
}

// PresignedURL returns a URL that can be used to access an Upload's file.
func (u *Uploader) PresignedURL(upload *models.Upload) (string, error) {
	url, err := u.Storer.PresignedURL(upload.StorageKey, upload.ContentType)
	if err != nil {
		u.logger.Error("failed to get presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// DeleteUpload removes an Upload from the database and deletes its file from the
// storer.
func (u *Uploader) DeleteUpload(upload *models.Upload) error {
	if err := u.Storer.Delete(upload.StorageKey); err != nil {
		return err
	}

	return models.DeleteUpload(u.db, upload)
}

// Download fetches an Upload's file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (u *Uploader) Download(upload *models.Upload) (io.ReadCloser, error) {
	return u.Storer.Fetch(upload.StorageKey)
}
