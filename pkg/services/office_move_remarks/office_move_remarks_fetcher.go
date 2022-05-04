package officemoveremarks

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type officeMoveRemarksFetcher struct {
}

func NewOfficeMoveRemarksFetcher() services.OfficeMoveRemarksFetcher {
	return &officeMoveRemarksFetcher{}
}

func (o officeMoveRemarksFetcher) ListOfficeMoveRemarks(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.OfficeMoveRemarks, error) {

	officeMoveRemarks := models.OfficeMoveRemarks{}
	err := appCtx.DB().Q().EagerPreload("OfficeUser").
		Where("move_id = ?", moveID).All(&officeMoveRemarks)

	if err != nil {
		return nil, err
	}

	if len(officeMoveRemarks) == 0 {
		return nil, models.ErrFetchNotFound
	}

	return &officeMoveRemarks, nil
}
