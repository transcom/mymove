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

// ZeroLengthFile represents an error caused by a file with no content
type ZeroLengthFile struct {
	message string
}

func (z ZeroLengthFile) Error() string {
	return z.message
}

// NewLocalFile creates a *runtime.File from a file on the local filesystem
func NewLocalFile(filePath string) (*runtime.File, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not stat file")
	}
	header := multipart.FileHeader{
		Filename: info.Name(),
		Size:     info.Size(),
	}

	/*
		#nosec - this path should never come from user input
	*/
	data, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "could not open file")
	}
	return &runtime.File{
		Header: &header,
		Data:   data,
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
func (u *Uploader) CreateUpload(documentID uuid.UUID, userID uuid.UUID, file *runtime.File) (*models.Upload, *validate.Errors, error) {
	if file.Header.Size == 0 {
		return nil, nil, errors.WithStack(&ZeroLengthFile{"File has length of 0"})
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

	// validate upload before pushing file to S3
	verrs, err := newUpload.Validate(u.db)
	if err != nil {
		u.logger.Error("Failed to validate", zap.Error(err))
		return nil, nil, err
	} else if verrs.HasAny() {
		return nil, verrs, nil
	}

	// Push file to S3
	key := u.uploadKey(newUpload)
	if _, err = u.storer.Store(key, file.Data, checksum); err != nil {
		u.logger.Error("failed to store object", zap.Error(err))
		return nil, nil, err
	}

	// Already validated upload, so just save
	err = u.db.Create(newUpload)
	if err != nil {
		u.logger.Error("DB Insertion", zap.Error(err))
		return nil, nil, err
	}

	u.logger.Info("created an upload with id and key ", zap.Any("new_upload_id", newUpload.ID), zap.String("key", key))

	return newUpload, nil, nil
}

// PresignedURL returns a URL that can be used to access an Upload's file.
func (u *Uploader) PresignedURL(upload *models.Upload) (string, error) {
	key := u.uploadKey(upload)
	url, err := u.storer.PresignedURL(key, upload.ContentType)
	if err != nil {
		u.logger.Error("failed to get presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// DeleteUpload removes an Upload from the database and deletes its file from the
// storer.
func (u *Uploader) DeleteUpload(upload *models.Upload) error {
	key := u.uploadKey(upload)

	if err := u.storer.Delete(key); err != nil {
		return err
	}

	if err := models.DeleteUpload(u.db, upload); err != nil {
		return err
	}
	return nil
}

func (u *Uploader) uploadKey(upload *models.Upload) string {
	return u.storer.Key("documents", upload.DocumentID.String(), "uploads", upload.ID.String())
}
