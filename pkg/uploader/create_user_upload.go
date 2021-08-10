package uploader

import (
	"io"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// CreateUserUploadForDocumentWrapper wrapper/helper function to create a user upload
func CreateUserUploadForDocumentWrapper(appCfg appconfig.AppConfig, userID uuid.UUID, storer storage.FileStorer, file io.ReadCloser, filename string, fileSizeLimit ByteSize, docID *uuid.UUID) (*models.UserUpload, string, *validate.Errors, error) {
	userUploader, err := NewUserUploader(storer, fileSizeLimit)
	if err != nil {
		appCfg.Logger().Fatal("could not instantiate uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
	}

	aFile, err := userUploader.PrepareFileForUpload(appCfg, file, filename)
	if err != nil {
		appCfg.Logger().Fatal("could not prepare file for uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
	}

	newUserUpload, verrs, err := userUploader.CreateUserUploadForDocument(appCfg, docID, userID, File{File: aFile}, AllowedTypesServiceMember)
	if verrs.HasAny() || err != nil {
		return nil, "", verrs, err
	}

	url, err := userUploader.PresignedURL(appCfg, newUserUpload)
	if err != nil {
		appCfg.Logger().Error("failed to get presigned url", zap.Error(err))
		return nil, "", &validate.Errors{}, err
	}

	return newUserUpload, url, &validate.Errors{}, err
}
