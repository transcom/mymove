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

// ErrZeroLengthFile represents an error caused by a file with no content
var ErrZeroLengthFile = errors.New("File has length of 0")

// ErrTooLarge is an error where the file size exceeds the limit
type ErrTooLarge struct {
	FileSize      int64
	FileSizeLimit ByteSize
}

// ErrFileSizeLimitExceedsMax is an error where file size exceeds max size
var ErrFileSizeLimitExceedsMax = errors.Errorf("FileSizeLimit exceeds max of %d bytes", MaxFileSizeLimit)

// MaxFileSizeLimit sets the maximum file size limit
// Anti-Virus scanning won't be able to scan files larger than 250MB
// Any unscanned files will not be available for download so while we can upload a larger
// file of any size the file will be locked from downloading forever.
const MaxFileSizeLimit = 250 * MB

// ErrTooLarge is the string representation of an error
func (e ErrTooLarge) Error() string {
	return fmt.Sprintf("file is too large: %d > %d filesize limit", e.FileSize, e.FileSizeLimit)
}

// ByteSize is a snack
type ByteSize int64

const (
	// B Byte
	B ByteSize = 1
	// KB KiloByte
	KB = 1000
	// MB MegaByte
	MB = 1000 * 1000
)

// Int64 returns an integer of the byte size
func (b ByteSize) Int64() int64 {
	return int64(b)
}

// File type to be used by Uploader. A wrapper around afero.File that allows attaching
// some additional metadata
type File struct {
	afero.File
	Tags *string
}

// Uploader encapsulates a few common processes: creating Uploads for a Document,
// generating pre-signed URLs for file access, and deleting Uploads.
type Uploader struct {
	db                *pop.Connection
	logger            Logger
	Storer            storage.FileStorer
	UploadStorageKey  string
	DefaultStorageKey string
	FileSizeLimit     ByteSize
	UploadType        models.UploadType
}

// NewUploader creates and returns a new uploader
func NewUploader(db *pop.Connection, logger Logger, storer storage.FileStorer, fileSizeLimit ByteSize, uploadType models.UploadType) (*Uploader, error) {
	if fileSizeLimit > MaxFileSizeLimit {
		return nil, ErrFileSizeLimitExceedsMax
	}
	return &Uploader{
		db:               db,
		logger:           logger,
		Storer:           storer,
		UploadStorageKey: "",
		FileSizeLimit:    fileSizeLimit,
		UploadType:       uploadType,
	}, nil
}

// SetUploadStorageKey set the Upload.StorageKey member
func (u *Uploader) SetUploadStorageKey(key string) {
	u.UploadStorageKey = key
}

// PrepareFileForUpload copy file buffer into Afero file, return Afero file
func (u *Uploader) PrepareFileForUpload(file io.ReadCloser, filename string) (afero.File, error) {
	// Read the incoming data into a temporary afero.File for consumption
	aFile, err := u.Storer.TempFileSystem().Create(filename)
	if err != nil {
		errorString := "Error opening afero file"
		u.logger.Error(errorString, zap.Error(err))
		return aFile, fmt.Errorf("%s %v", errorString, zap.Error(err))
	}

	_, err = io.Copy(aFile, file)
	if err != nil {
		errorString := "Error copying incoming data into afero file."
		u.logger.Error(errorString, zap.Error(err))
		return aFile, fmt.Errorf("%s %v", errorString, zap.Error(err))
	}

	return aFile, nil
}

func (u *Uploader) createAndPushUploadToS3(txConn *pop.Connection, file File, upload *models.Upload) (*models.Upload, *validate.Errors, error) {

	verrs, err := txConn.ValidateAndCreate(upload)
	if err != nil || verrs.HasAny() {
		u.logger.Error("Error creating new upload", zap.Error(err))
		return nil, verrs, fmt.Errorf("error creating upload %w", err)
	}

	// Push file to S3
	if _, err := u.Storer.Store(upload.StorageKey, file.File, upload.Checksum, file.Tags); err != nil {
		responseVErrors := validate.NewErrors()
		u.logger.Error("failed to store object", zap.Error(err))
		responseVErrors.Append(verrs)
		return nil, responseVErrors, fmt.Errorf("failed to store object %w", err)
	}

	u.logger.Info("created an upload with id and key ", zap.Any("new_upload_id", upload.ID), zap.String("key", upload.StorageKey))
	return upload, verrs, nil
}

// CreateUpload creates a new Upload by performing validations, storing the specified
// file using the supplied storer, and saving an Upload object to the database containing
// the file's metadata.
func (u *Uploader) CreateUpload(file File, allowedTypes AllowedFileTypes) (*models.Upload, *validate.Errors, error) {
	responseVErrors := validate.NewErrors()

	info, fileStatErr := file.Stat()
	if fileStatErr != nil {
		u.logger.Error("Could not get file info", zap.Error(fileStatErr))
	}

	if info.Size() == 0 {
		return nil, responseVErrors, ErrZeroLengthFile
	}

	if info.Size() > u.FileSizeLimit.Int64() {
		u.logger.Error("upload exceeds file size limit",
			zap.Int64("FileSize", info.Size()),
			zap.Int64("FileSizeLimit", u.FileSizeLimit.Int64()),
		)
		return nil, responseVErrors, ErrTooLarge{info.Size(), u.FileSizeLimit}
	}

	contentType, detectContentTypeErr := storage.DetectContentType(file)
	if detectContentTypeErr != nil {
		u.logger.Error("Could not detect content type", zap.Error(detectContentTypeErr))
		return nil, responseVErrors, detectContentTypeErr
	}

	validator := models.NewStringInList(contentType, "ContentType", allowedTypes)
	validator.IsValid(responseVErrors)
	if responseVErrors.HasAny() {
		u.logger.Error("Invalid content type for upload",
			zap.String("ContentType", contentType),
		)
		return nil, responseVErrors, nil
	}

	checksum, computeChecksumErr := storage.ComputeChecksum(file)
	if computeChecksumErr != nil {
		u.logger.Error("Could not compute checksum", zap.Error(computeChecksumErr))
		return nil, responseVErrors, computeChecksumErr
	}

	id := uuid.Must(uuid.NewV4())

	newUpload := &models.Upload{
		ID:          id,
		Filename:    file.Name(),
		Bytes:       info.Size(),
		ContentType: contentType,
		Checksum:    checksum,
		UploadType:  u.UploadType,
	}

	// Set the Upload.StorageKey if set
	if u.UploadStorageKey != "" {
		newUpload.StorageKey = u.UploadStorageKey
	} else if u.DefaultStorageKey != "" {
		newUpload.StorageKey = path.Join(u.DefaultStorageKey, "uploads", id.String())
	}

	var uploadError error
	// If we are already in a transaction, don't start one
	if u.db.TX != nil {
		responseCreateAndPushVerrs := validate.NewErrors()
		var responseCreateAndPushErr error
		newUpload, responseCreateAndPushVerrs, responseCreateAndPushErr = u.createAndPushUploadToS3(u.db, file, newUpload)
		if responseCreateAndPushErr != nil || responseCreateAndPushVerrs.HasAny() {
			responseVErrors.Append(responseCreateAndPushVerrs)
			uploadError = errors.Wrap(responseCreateAndPushErr, "failed to create and store upload object")
			return nil, responseVErrors, uploadError
		}
		u.logger.Info("created an upload with id and key ", zap.Any("new_upload_id", newUpload.ID), zap.String("key", newUpload.StorageKey))
		return newUpload, responseVErrors, nil
	}

	err := u.db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("Rollback The transaction")
		responseCreateAndPushVerrs := validate.NewErrors()
		var responseCreateAndPushErr error
		newUpload, responseCreateAndPushVerrs, responseCreateAndPushErr = u.createAndPushUploadToS3(db, file, newUpload)
		if responseCreateAndPushErr != nil || responseCreateAndPushVerrs.HasAny() {
			responseVErrors.Append(responseCreateAndPushVerrs)
			uploadError = errors.Wrap(responseCreateAndPushErr, "failed to create and store upload object")
			return transactionError
		}
		return nil
	})
	if err != nil {
		return nil, responseVErrors, errors.Wrap(uploadError, "could not create upload")
	}

	return newUpload, responseVErrors, nil
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
