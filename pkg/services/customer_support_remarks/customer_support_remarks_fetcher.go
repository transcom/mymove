package customersupportremarks

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type customerSupportRemarksFetcher struct {
}

func NewCustomerSupportRemarks() services.CustomerSupportRemarksFetcher {
	return &customerSupportRemarksFetcher{}
}

func (o customerSupportRemarksFetcher) ListCustomerSupportRemarks(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.CustomerSupportRemarks, error) {

	customerSupportRemarks := models.CustomerSupportRemarks{}
	err := appCtx.DB().Q().EagerPreload("OfficeUser").
		Where("move_id = ?", moveID).All(&customerSupportRemarks)

	if err != nil {
		return nil, err
	}

	if len(customerSupportRemarks) == 0 {
		return nil, models.ErrFetchNotFound
	}

	return &customerSupportRemarks, nil
}
