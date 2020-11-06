package uploader

import (
	"fmt"
	"io"
	"path"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// UserUploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type UserUploader struct {
	db       *pop.Connection
	logger   Logger
	uploader *Uploader
}

// NewUserUploader creates and returns a new uploader
func NewUserUploader(db *pop.Connection, logger Logger, storer storage.FileStorer, fileSizeLimit ByteSize) (*UserUploader, error) {
	uploader, err := NewUploader(db, logger, storer, fileSizeLimit, models.UploadTypeUSER)
	if err != nil {
		return nil, fmt.Errorf("could not create uploader.UserUploader for UserUpload: %w", err)
	}
	return &UserUploader{
		db:       db,
		logger:   logger,
		uploader: uploader,
	}, nil
}

// PrepareFileForUpload called Uploader.PrepareFileForUpload
func (u *UserUploader) PrepareFileForUpload(file io.ReadCloser, filename string) (afero.File, error) {
	// Read the incoming data into a temporary afero.File for consumption
	return u.uploader.PrepareFileForUpload(file, filename)
}

func (u *UserUploader) createAndStore(documentID *uuid.UUID, userID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.UserUpload, *validate.Errors, error) {
	// If storage key is not set assign a default
	if u.GetUploadStorageKey() == "" {
		u.uploader.DefaultStorageKey = path.Join("user", userID.String())
	}

	newUpload, verrs, err := u.uploader.CreateUpload(File{File: file}, allowedTypes)
	if verrs.HasAny() || err != nil {
		u.logger.Error("error creating and storing new upload for user", zap.Error(err))
		return nil, verrs, err
	}

	id := uuid.Must(uuid.NewV4())

	newUploadForUser := &models.UserUpload{
		ID:         id,
		DocumentID: documentID,
		UploaderID: userID,
		UploadID:   newUpload.ID,
		Upload:     *newUpload,
	}

	verrs, err = u.db.ValidateAndCreate(newUploadForUser)
	if err != nil || verrs.HasAny() {
		u.logger.Error("error creating new user upload", zap.Error(err))
		return nil, verrs, err
	}

	return newUploadForUser, &validate.Errors{}, nil
}

// CreateUserUploadForDocument creates a new UserUpload by performing validations, storing the specified
// file using the supplied storer, and saving an UserUpload object to the database containing
// the file's metadata.
func (u *UserUploader) CreateUserUploadForDocument(documentID *uuid.UUID, userID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.UserUpload, *validate.Errors, error) {

	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewUserUploader before calling this function")
	}

	var userUpload *models.UserUpload
	var verrs *validate.Errors
	var uploadError error

	userUpload, verrs, uploadError = u.createAndStore(documentID, userID, file, allowedTypes)
	if verrs.HasAny() || uploadError != nil {
		u.logger.Error("error creating new user upload (existing TX)", zap.Error(uploadError))
	} else {
		u.logger.Info("created a user upload with id and key (existing TX)", zap.Any("new_user_upload_id", userUpload.ID), zap.String("key", userUpload.Upload.StorageKey))
	}

	return userUpload, verrs, uploadError
}

// CreateUserUpload stores UserUpload but does not assign a Document
func (u *UserUploader) CreateUserUpload(userID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.UserUpload, *validate.Errors, error) {
	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewUserUploader before calling this function")
	}
	return u.CreateUserUploadForDocument(nil, userID, file, allowedTypes)
}

// DeleteUserUpload removes an UserUpload from the database and deletes its file from the
// storer.
func (u *UserUploader) DeleteUserUpload(userUpload *models.UserUpload) error {

	if u.db.TX != nil {
		if err := u.uploader.DeleteUpload(&userUpload.Upload); err != nil {
			return err
		}
		return models.DeleteUserUpload(u.db, userUpload)
	}
	return u.db.Transaction(func(db *pop.Connection) error {
		if err := u.uploader.DeleteUpload(&userUpload.Upload); err != nil {
			return err
		}
		return models.DeleteUserUpload(db, userUpload)
	})
}

// PresignedURL returns a URL that can be used to access an UserUpload's file.
func (u *UserUploader) PresignedURL(userUpload *models.UserUpload) (string, error) {
	if userUpload == nil {
		u.logger.Error("failed to get UserUploader presigned url")
		return "", errors.New("failed to get UserUploader presigned url")
	}
	url, err := u.uploader.PresignedURL(&userUpload.Upload)
	if err != nil {
		u.logger.Error("failed to get UserUploader presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// FileSystem return file system from Uploader file storer
func (u *UserUploader) FileSystem() *afero.Afero {
	return u.uploader.Storer.FileSystem()
}

// Uploader return the Uploader for UserUploader
func (u *UserUploader) Uploader() *Uploader {
	return u.uploader
}

// SetUploadStorageKey set the UserUpload.Upload.StorageKey member
func (u *UserUploader) SetUploadStorageKey(key string) {
	if u.uploader != nil {
		u.uploader.SetUploadStorageKey(key)
	}
}

// GetUploadStorageKey returns the UserUpload.Upload.StorageKey member
func (u *UserUploader) GetUploadStorageKey() string {
	if u.uploader == nil {
		return ""
	}
	return u.uploader.UploadStorageKey
}

// Download fetches an Upload's file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (u *UserUploader) Download(userUpload *models.UserUpload) (io.ReadCloser, error) {
	return u.uploader.Download(&userUpload.Upload)
}
