package customersupportremarks

import (
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type customerSupportRemarksFetcher struct {
}

func NewCustomerSupportRemarks() services.CustomerSupportRemarksFetcher {
	return &customerSupportRemarksFetcher{}
}

func (o customerSupportRemarksFetcher) ListCustomerSupportRemarks(appCtx appcontext.AppContext, moveCode string) (*models.CustomerSupportRemarks, error) {
	var move models.Move
	err := appCtx.DB().Q().Where("locator = ?", moveCode).First(&move)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, models.ErrFetchNotFound
		}
		return nil, err
	}

	customerSupportRemarks := models.CustomerSupportRemarks{}
	err = appCtx.DB().Q().Scope(utilities.ExcludeDeletedScope()).EagerPreload("OfficeUser").
		Where("move_id = ?", move.ID).Order("created_at desc").All(&customerSupportRemarks)

	if err != nil {
		return nil, err
	}

	return &customerSupportRemarks, nil
}
