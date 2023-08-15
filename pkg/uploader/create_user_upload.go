package uploader

import (
	"io"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
)

// CreateUserUploadForDocumentWrapper wrapper/helper function to create a user upload
func CreateUserUploadForDocumentWrapper(
	appCtx appcontext.AppContext, userID uuid.UUID,
	storer storage.FileStorer, file io.ReadCloser,
	filename string,
	fileSizeLimit ByteSize,
	allowedFileTypes AllowedFileTypes,
	docID *uuid.UUID,
) (*models.UserUpload, string, *validate.Errors, error) {

	userUploader, err := NewUserUploader(storer, fileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFailedToInitUploader{message: err.Error()}
	}

	aFile, err := userUploader.PrepareFileForUpload(appCtx, file, filename)
	if err != nil {
		appCtx.Logger().Fatal("could not prepare file for uploader", zap.Error(err))
		return nil, "", &validate.Errors{}, ErrFile{message: err.Error()}
	}

	newUserUpload, verrs, err := userUploader.CreateUserUploadForDocument(appCtx, docID, userID, File{File: aFile}, allowedFileTypes)
	if verrs.HasAny() || err != nil {
		return nil, "", verrs, err
	}

	url, err := userUploader.PresignedURL(appCtx, newUserUpload)
	if err != nil {
		appCtx.Logger().Error("failed to get presigned url", zap.Error(err))
		return nil, "", &validate.Errors{}, err
	}

	return newUserUpload, url, &validate.Errors{}, err
}
