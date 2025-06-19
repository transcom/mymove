package uploader

import (
	"fmt"
	"io"
	"path"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// UserUploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type UserUploader struct {
	uploader *Uploader
}

// NewUserUploader creates and returns a new uploader
func NewUserUploader(storer storage.FileStorer, fileSizeLimit ByteSize) (*UserUploader, error) {
	uploader, err := NewUploader(storer, fileSizeLimit, models.UploadTypeUSER)
	if err != nil {
		return nil, fmt.Errorf("could not create uploader.UserUploader for UserUpload: %w", err)
	}
	return &UserUploader{
		uploader: uploader,
	}, nil
}

// NewOfficeUploader creates and returns a new uploader
func NewOfficeUploader(storer storage.FileStorer, fileSizeLimit ByteSize) (*UserUploader, error) {
	uploader, err := NewUploader(storer, fileSizeLimit, models.UploadTypeOFFICE)
	if err != nil {
		return nil, fmt.Errorf("could not create uploader.UserUploader for UserUpload: %w", err)
	}
	return &UserUploader{
		uploader: uploader,
	}, nil
}

// PrepareFileForUpload calls Uploader.PrepareFileForUpload
func (u *UserUploader) PrepareFileForUpload(appCtx appcontext.AppContext, file io.ReadCloser, filename string) (afero.File, error) {
	// Read the incoming data into a temporary afero.File for consumption
	return u.uploader.PrepareFileForUpload(appCtx, file, filename)
}

func (u *UserUploader) createAndStore(appCtx appcontext.AppContext, documentID *uuid.UUID, userID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.UserUpload, *validate.Errors, error) {
	// If storage key is not set assign a default
	if u.GetUploadStorageKey() == "" {
		u.uploader.DefaultStorageKey = path.Join("user", userID.String())
	}

	newUpload, verrs, err := u.uploader.CreateUpload(appCtx, File{File: file}, allowedTypes)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("error creating and storing new upload for user", zap.Error(err))
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

	verrs, err = appCtx.DB().ValidateAndCreate(newUploadForUser)
	if err != nil || verrs.HasAny() {
		appCtx.Logger().Error("error creating new user upload", zap.Error(err))
		return nil, verrs, err
	}

	return newUploadForUser, &validate.Errors{}, nil
}

// CreateUserUploadForDocument creates a new UserUpload by performing validations, storing the specified
// file using the supplied storer, and saving an UserUpload object to the database containing
// the file's metadata.
func (u *UserUploader) CreateUserUploadForDocument(appCtx appcontext.AppContext, documentID *uuid.UUID, userID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.UserUpload, *validate.Errors, error) {

	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewUserUploader before calling this function")
	}

	var userUpload *models.UserUpload
	var verrs *validate.Errors
	var uploadError error

	userUpload, verrs, uploadError = u.createAndStore(appCtx, documentID, userID, file, allowedTypes)
	if verrs.HasAny() || uploadError != nil {
		appCtx.Logger().Error("error creating new user upload", zap.Error(uploadError))
	} else {
		appCtx.Logger().Info("created a user upload with id and key", zap.Any("new_user_upload_id", userUpload.ID), zap.String("key", userUpload.Upload.StorageKey))
	}

	defer func() {
		if file.File != nil {
			file.File.Close()
		}
		err := u.uploader.Storer.TempFileSystem().Remove(file.File.Name())

		if err != nil {
			appCtx.Logger().Error("error removing file from memory", zap.Error(err))
		}
	}()

	return userUpload, verrs, uploadError
}

// CreateUserUpload stores UserUpload but does not assign a Document
func (u *UserUploader) CreateUserUpload(appCtx appcontext.AppContext, userID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.UserUpload, *validate.Errors, error) {
	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewUserUploader before calling this function")
	}
	return u.CreateUserUploadForDocument(appCtx, nil, userID, file, allowedTypes)
}

func (u *UserUploader) UpdateUserXlsxUploadFilename(appCtx appcontext.AppContext, userUpload *models.UserUpload, newFilename string) (*models.UserUpload, *validate.Errors, error) {
	// 1) Mutate the in-memory struct
	userUpload.Upload.Filename = newFilename

	// 2) Persist only the Upload table change.
	//    ValidateAndUpdate will run any model validations and issue an UPDATE.
	verrs, err := appCtx.DB().ValidateAndUpdate(&userUpload.Upload)
	if err != nil || verrs.HasAny() {
		appCtx.Logger().Error("failed to update upload filename",
			zap.Error(err),
			zap.Any("validation_errors", verrs),
		)
		return nil, verrs, errors.Wrap(err, "could not update upload filename")
	}

	// Return the updated UserUpload
	return userUpload, &validate.Errors{}, nil
}

// DeleteUserUpload removes an UserUpload from the database and deletes its file from the
// storer.
func (u *UserUploader) DeleteUserUpload(appCtx appcontext.AppContext, userUpload *models.UserUpload) error {

	if appCtx.DB().TX != nil {
		if err := u.uploader.DeleteUpload(appCtx, &userUpload.Upload); err != nil {
			return err
		}
		return models.DeleteUserUpload(appCtx.DB(), userUpload)
	}
	return appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := u.uploader.DeleteUpload(txnAppCtx, &userUpload.Upload); err != nil {
			return err
		}
		return models.DeleteUserUpload(txnAppCtx.DB(), userUpload)
	})
}

// PresignedURL returns a URL that can be used to access an UserUpload's file.
func (u *UserUploader) PresignedURL(appCtx appcontext.AppContext, userUpload *models.UserUpload) (string, error) {
	if userUpload == nil {
		appCtx.Logger().Error("failed to get UserUploader presigned url")
		return "", errors.New("failed to get UserUploader presigned url")
	}
	url, err := u.uploader.PresignedURL(appCtx, &userUpload.Upload)
	if err != nil {
		appCtx.Logger().Error("failed to get UserUploader presigned url", zap.Error(err))
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
func (u *UserUploader) Download(appCtx appcontext.AppContext, userUpload *models.UserUpload) (io.ReadCloser, error) {
	return u.uploader.Download(appCtx, &userUpload.Upload)
}
