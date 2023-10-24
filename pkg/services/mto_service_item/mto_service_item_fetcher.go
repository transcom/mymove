package mtoserviceitem

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type mtoServiceItemFetcher struct {
}

// NewMTOServiceItemFetcher creates a new MTOServiceItemFetcher struct
func NewMTOServiceItemFetcher() services.MTOServiceItemFetcher {
	return &mtoServiceItemFetcher{}
}

// searches the database and returns the service item based on the passed in uuid
func (p mtoServiceItemFetcher) GetServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error) {
	var serviceItem models.MTOServiceItem
	findServiceItemQuery := appCtx.DB().Q()

	err := findServiceItemQuery.Find(&serviceItem, serviceItemID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(serviceItemID, "while looking for service item")
		default:
			return nil, apperror.NewQueryError("ServiceItem", err, "")
		}
	}

	return &serviceItem, nil
}
