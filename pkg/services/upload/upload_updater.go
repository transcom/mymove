package upload

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

// uploadUpdater is a service object to update Upload
type uploadUpdater struct {
	*models.Upload
}

// NewUploadUpdater returns a new UploadUpdater
func NewUploadUpdater() *uploadUpdater {
	return &uploadUpdater{}
}

func (f *uploadUpdater) UpdateUploadForRotation(appCtx appcontext.AppContext, uploadID uuid.UUID, newRotation *int64) (*models.Upload, error) {
	upload, err := models.FetchUpload(appCtx.DB(), uploadID)
	if err != nil {
		return &models.Upload{}, apperror.NewNotFoundError(uploadID, "no upload found")
	}

	if newRotation == nil {
		return &models.Upload{}, apperror.NewInvalidInputError(uploadID, nil, nil, "rotation is required")
	}

	upload.Rotation = newRotation

	err = appCtx.DB().Save(upload)
	if err != nil {
		return &models.Upload{}, apperror.NewQueryError("upload", err, "")
	}

	return upload, nil
}
