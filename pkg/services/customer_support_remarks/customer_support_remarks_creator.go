package customersupportremarks

import (
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type customerSupportRemarksCreator struct {
}

func NewCustomerSupportRemarksCreator() services.CustomerSupportRemarksCreator {
	return &customerSupportRemarksCreator{}
}

func (o customerSupportRemarksCreator) CreateCustomerSupportRemark(appCtx appcontext.AppContext, customerSupportRemark *models.CustomerSupportRemark, moveCode string) (*models.CustomerSupportRemark, error) {

	// Need to get the MoveID from the MoveCode
	var move models.Move
	err := appCtx.DB().Q().Where("locator = ?", moveCode).First(&move)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, models.ErrFetchNotFound
		}
		return nil, err
	}
	customerSupportRemark.MoveID = move.ID

	verrs, err := appCtx.DB().ValidateAndCreate(customerSupportRemark)
	if verrs.Count() != 0 || err != nil {
		return nil, err
	}

	return customerSupportRemark, nil
}
