package uploader

import (
	"io"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// CreateUserUploadForDocumentWrapper wrapper/helper function to create a user upload
func CreateUserUploadForDocumentWrapper(db *pop.Connection, logger Logger, userID uuid.UUID, storer storage.FileStorer, file io.ReadCloser, filename string, fileSizeLimit ByteSize, docID *uuid.UUID) (*models.UserUpload, string, *validate.Errors, error) {
	userUploader, err := NewUserUploader(db, logger, storer, fileSizeLimit)
	if err != nil {
		logger.Fatal("could not instantiate uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
	}

	aFile, err := userUploader.PrepareFileForUpload(file, filename)
	if err != nil {
		logger.Fatal("could not prepare file for uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
	}

	newUserUpload, verrs, err := userUploader.CreateUserUploadForDocument(docID, userID, File{File: aFile}, AllowedTypesServiceMember)
	if verrs.HasAny() || err != nil {
		return nil, "", verrs, err
	}

	url, err := userUploader.PresignedURL(newUserUpload)
	if err != nil {
		logger.Error("failed to get presigned url", zap.Error(err))
		return nil, "", &validate.Errors{}, err
	}

	return newUserUpload, url, &validate.Errors{}, err
}
