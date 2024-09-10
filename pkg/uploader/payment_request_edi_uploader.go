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

type PaymentRequestEDIUploader struct {
	uploader *Uploader
}

// NewPaymentRequestEdiUploader creates and returns a new uploader
func NewPaymentRequestEDIUploader(storer storage.FileStorer, fileSizeLimit ByteSize) (*PaymentRequestEDIUploader, error) {
	uploader, err := NewUploader(storer, fileSizeLimit, models.UploadTypeAPP)
	if err != nil {
		if err == ErrFileSizeLimitExceedsMax {
			return nil, err
		}
		return nil, fmt.Errorf("could not create uploader.PaymentRequestEdiUploader for PaymentRequestEdiUpload: %w", err)
	}
	return &PaymentRequestEDIUploader{
		uploader: uploader,
	}, nil
}

// PrepareFileForUpload called Uploader.PrepareFileForUpload
func (u *PaymentRequestEDIUploader) PrepareFileForUpload(appCtx appcontext.AppContext, file io.ReadCloser, filename string) (afero.File, error) {
	// Read the incoming data into a temporary afero.File for consumption

	// Convert io.ReadCloser to io.ReadSeeker
	seeker, ok := file.(io.ReadSeeker)
	if !ok {
		appCtx.Logger().Error("file is not seekable")
		return nil, errors.New("file is not seekable")
	}

	fileSize, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		appCtx.Logger().Error("error getting file size", zap.Error(err))
		return nil, err
	}
	_, err = seeker.Seek(0, io.SeekStart)
	if err != nil {
		appCtx.Logger().Error("error resetting file position", zap.Error(err))
		return nil, err
	}

	appCtx.Logger().Info("File size", zap.Int64("size", fileSize))

	return u.uploader.PrepareFileForUpload(appCtx, file, filename)
}
func (u *PaymentRequestEDIUploader) createAndStore(appCtx appcontext.AppContext, file *File, allowedTypes AllowedFileTypes) (*models.PaymentRequestEdiUpload, *validate.Errors, error) {
	// If storage key is not set assign a default
	id := uuid.Must(uuid.NewV4())
	if u.GetUploadStorageKey() == "" {
		u.uploader.DefaultStorageKey = path.Join("app", id.String())
	}
	aFile, err := u.PrepareFileForUpload(appCtx, file.File, file.File.Name())
	if err != nil {
		appCtx.Logger().Error("error preparing file for upload", zap.Error(err))
		return nil, nil, err
	}
	preppedFile := &File{
		File: aFile,
	}

	newUpload, verrs, err := u.uploader.CreateUpload(appCtx, *preppedFile, allowedTypes)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("error creating and storing new upload", zap.Error(err))
		return nil, verrs, err
	}

	newUploadForApp := &models.PaymentRequestEdiUpload{
		ID:       id,
		UploadID: newUpload.ID,
		Upload:   *newUpload,
	}

	verrs, err = appCtx.DB().ValidateAndCreate(newUploadForApp)
	if err != nil || verrs.HasAny() {
		appCtx.Logger().Error("error creating new app upload", zap.Error(err))
		return nil, verrs, err
	}

	return newUploadForApp, &validate.Errors{}, nil
}
func (u *PaymentRequestEDIUploader) CreatePaymentRequestEDIUploadForDocument(appCtx appcontext.AppContext, file *File, allowedTypes AllowedFileTypes) (*models.PaymentRequestEdiUpload, *validate.Errors, error) {

	if u.uploader == nil {
		return nil, &validate.Errors{}, errors.New("Did not call NewPaymentRequestEDIUploader before calling this function")
	}

	if file != nil && file.File != nil {
		tags := "858c"
		file.Tags = &tags
	}

	paymentRequestEDIUpload, verrs, err := u.createAndStore(appCtx, file, allowedTypes)
	return paymentRequestEDIUpload, verrs, err

}

func (u *PaymentRequestEDIUploader) PresignedURL(appCtx appcontext.AppContext, paymentRequestEDIUpload *models.PaymentRequestEdiUpload) (string, error) {
	if paymentRequestEDIUpload == nil {
		appCtx.Logger().Error("failed to get PaymentRequestEDIUploader presigned url")
		return "", errors.New("failed to get PaymentRequestEDIUploader presigned url")
	}
	url, err := u.uploader.PresignedURL(appCtx, &paymentRequestEDIUpload.Upload)
	if err != nil {
		appCtx.Logger().Error("failed to get PaymentRequestEDIUploader presigned url", zap.Error(err))
		return "", err
	}
	return url, nil
}

// FileSystem return Uploader file system
func (u *PaymentRequestEDIUploader) FileSystem() *afero.Afero {
	return u.uploader.Storer.FileSystem()
}

// Uploader return Uploader
func (u *PaymentRequestEDIUploader) Uploader() *Uploader {
	return u.uploader
}

// SetUploadStorageKey set the PaymentRequestEDIUpload.Upload.StorageKey member
func (u *PaymentRequestEDIUploader) SetUploadStorageKey(key string) {
	if u.uploader != nil {
		u.uploader.SetUploadStorageKey(key)
	}
}

// GetUploadStorageKey returns the PaymentRequestEDIUpload.Upload.StorageKey member
func (u *PaymentRequestEDIUploader) GetUploadStorageKey() string {
	if u.uploader == nil {
		return ""
	}
	return u.uploader.UploadStorageKey
}

// Download fetches an Upload's file and stores it in a tempfile. The path to this
// file is returned.
//
// It is the caller's responsibility to delete the tempfile.
func (u *PaymentRequestEDIUploader) Download(appCtx appcontext.AppContext, paymentRequestEDIUpload *models.PaymentRequestEdiUpload) (io.ReadCloser, error) {
	return u.uploader.Download(appCtx, &paymentRequestEDIUpload.Upload)
}
