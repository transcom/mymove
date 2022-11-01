package address

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type addressUpdater struct {
	checks []addressValidator
}

func NewAddressUpdater() services.AddressUpdater {
	return &addressUpdater{
		checks: []addressValidator{
			checkID(),
		},
	}
}

func (f *addressUpdater) UpdateAddress(appCtx appcontext.AppContext, address *models.Address, eTag string) (*models.Address, error) {
	originalAddress := models.FetchAddressByID(appCtx.DB(), &address.ID)
	if originalAddress == nil {
		return nil, apperror.NewBadDataError("invalid ID used for address")
	}

	// verify ETag
	if etag.GenerateEtag(originalAddress.UpdatedAt) != eTag {
		return nil, apperror.NewPreconditionFailedError(originalAddress.ID, nil)
	}

	// tk handle merge?

	err := validateAddress(appCtx, address, originalAddress, f.checks...)
	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().ValidateAndUpdate(address)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(address.ID, err, verrs, "invalid input while updating an address")
		} else if err != nil {
			return apperror.NewQueryError("Address update", err, "")
		}
		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return address, nil
}
