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

// ServiceRequestUploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type ServiceRequestUploader struct {
	uploader *Uploader
}

// NewServiceRequestUploader creates and returns a new uploader
func NewServiceRequestUploader(storer storage.FileStorer, fileSizeLimit ByteSize) (*ServiceRequestUploader, error) {
	uploader, err := NewUploader(storer, fileSizeLimit, models.UploadTypePRIME)
	if err != nil {
		if err == ErrFileSizeLimitExceedsMax {
			return nil, err
		}
		return nil, fmt.Errorf("could not create uploader.ServiceRequestUploader for ServiceRequestUpload: %w", err)
	}
	return &ServiceRequestUploader{
		uploader: uploader,
	}, nil
}

// PrepareFileForUpload called Uploader.PrepareFileForUpload
func (u *ServiceRequestUploader) PrepareFileForUpload(appCtx appcontext.AppContext, file io.ReadCloser, filename string) (afero.File, error) {
	// Read the incoming data into a temporary afero.File for consumption
	return u.uploader.PrepareFileForUpload(appCtx, file, filename)
}

func (u *ServiceRequestUploader) createAndStore(appCtx appcontext.AppContext, serviceItemDocID *uuid.UUID, contractorID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.ServiceRequestDocumentUpload, *validate.Errors, error) {
	// If storage key is not set assign a default
	if u.GetUploadStorageKey() == "" {
		u.uploader.DefaultStorageKey = path.Join("prime", contractorID.String())
	}

	newUpload, verrs, err := u.uploader.CreateUpload(appCtx, File{File: file}, allowedTypes)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("error creating and storing new upload for prime", zap.Error(err))
		return nil, verrs, err
	}

	id := uuid.Must(uuid.NewV4())

	newUploadForPrime := &models.ServiceRequestDocumentUpload{
		ID:                       id,
		ServiceRequestDocumentID: *serviceItemDocID,
		ContractorID:             contractorID,
		UploadID:                 newUpload.ID,
		Upload:                   *newUpload,
	}

	verrs, err = appCtx.DB().ValidateAndCreate(newUploadForPrime)
	if err != nil || verrs.HasAny() {
		appCtx.Logger().Error("error creating new prime upload", zap.Error(err))
		return nil, verrs, err
	}

	return newUploadForPrime, &validate.Errors{}, nil
}

// CreateServiceRequestUploadForDocument creates a new ServiceRequestUpload by performing validations, storing the specified
// file using the supplied storer, and saving an ServiceRequestUpload object to the database containing
// the file's metadata.
func (u *ServiceRequestUploader) CreateServiceRequestUploadForDocument(appCtx appcontext.AppContext, posID *uuid.UUID, contractorID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.ServiceRequestDocumentUpload, *validate.Errors, error) {

	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewServiceRequestUploader before calling this function")
	}

	var serviceRequestUpload *models.ServiceRequestDocumentUpload
	var verrs *validate.Errors
	var uploadError error

	txError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		transactionError := errors.New("Rollback The transaction")
		serviceRequestUpload, verrs, uploadError = u.createAndStore(txnAppCtx, posID, contractorID, file, allowedTypes)
		if verrs.HasAny() || uploadError != nil {
			txnAppCtx.Logger().Error("error creating new prime upload", zap.Error(uploadError))
			return transactionError
		}

		txnAppCtx.Logger().Info("created a prime upload with id and key ", zap.Any("new_prime_upload_id", serviceRequestUpload.ID), zap.String("key", serviceRequestUpload.Upload.StorageKey))
		return nil
	})
	if txError != nil {
		return nil, verrs, uploadError
	}

	return serviceRequestUpload, &validate.Errors{}, nil
}

// DeleteServiceRequestUpload removes an ServiceRequestUpload from the database and deletes its file from the
// storer.
func (u *ServiceRequestUploader) DeleteServiceRequestUpload(appCtx appcontext.AppContext, serviceRequestUpload *models.ServiceRequestDocumentUpload) error {
	return appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := u.uploader.DeleteUpload(txnAppCtx, &serviceRequestUpload.Upload); err != nil {
			return err
		}
		return models.DeleteServiceRequestDocumentUpload(txnAppCtx.DB(), serviceRequestUpload)
	})
}

// PresignedURL returns a URL that can be used to access an ServiceRequestUpload's file.
func (u *ServiceRequestUploader) PresignedURL(appCtx appcontext.AppContext, serviceRequestUpload *models.ServiceRequestDocumentUpload) (string, error) {
	if serviceRequestUpload == nil {
		appCtx.Logger().Error("failed to get ServiceRequestUploader presigned url")
		return "", errors.New("failed to get ServiceRequestUploader presigned url")
	}
	url, err := u.uploader.PresignedURL(appCtx, &serviceRequestUpload.Upload)
	if err != nil {
		appCtx.Logger().Error("failed to get ServiceRequestUploader presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// FileSystem return Uploader file system
func (u *ServiceRequestUploader) FileSystem() *afero.Afero {
	return u.uploader.Storer.FileSystem()
}

// Uploader return Uploader
func (u *ServiceRequestUploader) Uploader() *Uploader {
	return u.uploader
}

// SetUploadStorageKey set the ServiceRequestUpload.Upload.StorageKey member
func (u *ServiceRequestUploader) SetUploadStorageKey(key string) {
	if u.uploader != nil {
		u.uploader.SetUploadStorageKey(key)
	}
}

// GetUploadStorageKey returns the ServiceRequestUpload.Upload.StorageKey member
func (u *ServiceRequestUploader) GetUploadStorageKey() string {
	if u.uploader == nil {
		return ""
	}
	return u.uploader.UploadStorageKey
}

// Download fetches an Upload's file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (u *ServiceRequestUploader) Download(appCtx appcontext.AppContext, serviceRequestUpload *models.ServiceRequestDocumentUpload) (io.ReadCloser, error) {
	return u.uploader.Download(appCtx, &serviceRequestUpload.Upload)
}
