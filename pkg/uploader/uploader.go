package uploader

import (
	"mime/multipart"
	"os"

	"github.com/go-openapi/runtime"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// ErrZeroLengthFile represents an error caused by a file with no content
var ErrZeroLengthFile = errors.New("File has length of 0")

// NewLocalFile creates a *runtime.File from a local filepath
func NewLocalFile(filePath string) (*runtime.File, error) {
	// #nosec never comes from user input
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}

	return RuntimeFile(file)
}

// RuntimeFile creates a *runtime.File from an os.File
func RuntimeFile(file *os.File) (*runtime.File, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "could not get file stats")
	}

	header := multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	return &runtime.File{
		Header: &header,
		Data:   file,
	}, nil
}

// Uploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type Uploader struct {
	db     *pop.Connection
	logger *zap.Logger
	storer storage.FileStorer
}

// NewUploader creates and returns a new uploader
func NewUploader(db *pop.Connection, logger *zap.Logger, storer storage.FileStorer) *Uploader {
	return &Uploader{
		db:     db,
		logger: logger,
		storer: storer,
	}
}

// CreateUpload creates a new Upload by performing validations, storing the specified
// file using the supplied storer, and saving an Upload object to the database containing
// the file's metadata.
func (u *Uploader) CreateUpload(documentID *uuid.UUID, userID uuid.UUID, file *runtime.File) (*models.Upload, *validate.Errors, error) {
	if file.Header.Size == 0 {
		return nil, nil, ErrZeroLengthFile
	}

	contentType, err := storage.DetectContentType(file.Data)
	if err != nil {
		u.logger.Error("Could not detect content type", zap.Error(err))
		return nil, nil, err
	}

	checksum, err := storage.ComputeChecksum(file.Data)
	if err != nil {
		u.logger.Error("Could not compute checksum", zap.Error(err))
		return nil, nil, err
	}

	id := uuid.Must(uuid.NewV4())

	newUpload := &models.Upload{
		ID:          id,
		DocumentID:  documentID,
		UploaderID:  userID,
		Filename:    file.Header.Filename,
		Bytes:       int64(file.Header.Size),
		ContentType: contentType,
		Checksum:    checksum,
	}

	responseVErrors := validate.NewErrors()
	var responseError error

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
		if _, err := u.storer.Store(newUpload.StorageKey, file.Data, checksum); err != nil {
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

// PresignedURL returns a URL that can be used to access an Upload's file.
func (u *Uploader) PresignedURL(upload *models.Upload) (string, error) {
	url, err := u.storer.PresignedURL(upload.StorageKey, upload.ContentType)
	if err != nil {
		u.logger.Error("failed to get presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// DeleteUpload removes an Upload from the database and deletes its file from the
// storer.
func (u *Uploader) DeleteUpload(upload *models.Upload) error {
	if err := u.storer.Delete(upload.StorageKey); err != nil {
		return err
	}

	if err := models.DeleteUpload(u.db, upload); err != nil {
		return err
	}
	return nil
}

// Download fetches an Upload's file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (u *Uploader) Download(upload *models.Upload) (string, error) {
	return u.storer.Fetch(upload.StorageKey)
}
