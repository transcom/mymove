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

// PrimeUploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type PrimeUploader struct {
	uploader *Uploader
}

// NewPrimeUploader creates and returns a new uploader
func NewPrimeUploader(storer storage.FileStorer, fileSizeLimit ByteSize) (*PrimeUploader, error) {
	uploader, err := NewUploader(storer, fileSizeLimit, models.UploadTypePRIME)
	if err != nil {
		if err == ErrFileSizeLimitExceedsMax {
			return nil, err
		}
		return nil, fmt.Errorf("could not create uploader.PrimeUploader for PrimeUpload: %w", err)
	}
	return &PrimeUploader{
		uploader: uploader,
	}, nil
}

// PrepareFileForUpload called Uploader.PrepareFileForUpload
func (u *PrimeUploader) PrepareFileForUpload(appCtx appcontext.AppContext, file io.ReadCloser, filename string) (afero.File, error) {
	// Read the incoming data into a temporary afero.File for consumption
	return u.uploader.PrepareFileForUpload(appCtx, file, filename)
}

func (u *PrimeUploader) createAndStore(appCtx appcontext.AppContext, posID *uuid.UUID, contractorID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.PrimeUpload, *validate.Errors, error) {
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

	newUploadForUser := &models.PrimeUpload{
		ID:                  id,
		ProofOfServiceDocID: *posID,
		ContractorID:        contractorID,
		UploadID:            newUpload.ID,
		Upload:              *newUpload,
	}

	verrs, err = appCtx.DB().ValidateAndCreate(newUploadForUser)
	if err != nil || verrs.HasAny() {
		appCtx.Logger().Error("error creating new prime upload", zap.Error(err))
		return nil, verrs, err
	}

	return newUploadForUser, &validate.Errors{}, nil
}

// CreatePrimeUploadForDocument creates a new PrimeUpload by performing validations, storing the specified
// file using the supplied storer, and saving an PrimeUpload object to the database containing
// the file's metadata.
func (u *PrimeUploader) CreatePrimeUploadForDocument(appCtx appcontext.AppContext, posID *uuid.UUID, contractorID uuid.UUID, file File, allowedTypes AllowedFileTypes) (*models.PrimeUpload, *validate.Errors, error) {

	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewPrimeUploader before calling this function")
	}

	var primeUpload *models.PrimeUpload
	var verrs *validate.Errors
	var uploadError error

	// If we are already in a transaction, don't start one
	if appCtx.DB().TX != nil {
		primeUpload, verrs, uploadError = u.createAndStore(appCtx, posID, contractorID, file, allowedTypes)
		if verrs.HasAny() || uploadError != nil {
			appCtx.Logger().Error("error creating new prime upload (existing TX)", zap.Error(uploadError))
		} else {
			appCtx.Logger().Info("created a prime upload with id and key (existing TX)", zap.Any("new_prime_upload_id", primeUpload.ID), zap.String("key", primeUpload.Upload.StorageKey))
		}

		return primeUpload, verrs, uploadError
	}

	txError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		transactionError := errors.New("Rollback The transaction")
		primeUpload, verrs, uploadError = u.createAndStore(txnAppCtx, posID, contractorID, file, allowedTypes)
		if verrs.HasAny() || uploadError != nil {
			txnAppCtx.Logger().Error("error creating new prime upload", zap.Error(uploadError))
			return transactionError
		}

		txnAppCtx.Logger().Info("created a prime upload with id and key ", zap.Any("new_prime_upload_id", primeUpload.ID), zap.String("key", primeUpload.Upload.StorageKey))
		return nil
	})
	if txError != nil {
		return nil, verrs, uploadError
	}

	return primeUpload, &validate.Errors{}, nil
}

// DeletePrimeUpload removes an PrimeUpload from the database and deletes its file from the
// storer.
func (u *PrimeUploader) DeletePrimeUpload(appCtx appcontext.AppContext, primeUpload *models.PrimeUpload) error {
	if appCtx.DB().TX != nil {
		if err := u.uploader.DeleteUpload(appCtx, &primeUpload.Upload); err != nil {
			return err
		}
		return models.DeletePrimeUpload(appCtx.DB(), primeUpload)

	}
	return appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := u.uploader.DeleteUpload(txnAppCtx, &primeUpload.Upload); err != nil {
			return err
		}
		return models.DeletePrimeUpload(txnAppCtx.DB(), primeUpload)
	})
}

// PresignedURL returns a URL that can be used to access an PrimeUpload's file.
func (u *PrimeUploader) PresignedURL(appCtx appcontext.AppContext, primeUpload *models.PrimeUpload) (string, error) {
	if primeUpload == nil {
		appCtx.Logger().Error("failed to get PrimeUploader presigned url")
		return "", errors.New("failed to get PrimeUploader presigned url")
	}
	url, err := u.uploader.PresignedURL(appCtx, &primeUpload.Upload)
	if err != nil {
		appCtx.Logger().Error("failed to get PrimeUploader presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// FileSystem return Uploader file system
func (u *PrimeUploader) FileSystem() *afero.Afero {
	return u.uploader.Storer.FileSystem()
}

// Uploader return Uploader
func (u *PrimeUploader) Uploader() *Uploader {
	return u.uploader
}

// SetUploadStorageKey set the PrimeUpload.Upload.StorageKey member
func (u *PrimeUploader) SetUploadStorageKey(key string) {
	if u.uploader != nil {
		u.uploader.SetUploadStorageKey(key)
	}
}

// GetUploadStorageKey returns the PrimeUpload.Upload.StorageKey member
func (u *PrimeUploader) GetUploadStorageKey() string {
	if u.uploader == nil {
		return ""
	}
	return u.uploader.UploadStorageKey
}

// Download fetches an Upload's file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (u *PrimeUploader) Download(appCtx appcontext.AppContext, primeUpload *models.PrimeUpload) (io.ReadCloser, error) {
	return u.uploader.Download(appCtx, &primeUpload.Upload)
}
