package address

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type addressCreator struct {
	checks []addressValidator
}

func NewAddressCreator() services.AddressCreator {
	return &addressCreator{
		checks: []addressValidator{
			checkID(),
		},
	}
}

func (f *addressCreator) CreateAddress(appCtx appcontext.AppContext, address *models.Address) (*models.Address, error) {
	err := validateAddress(appCtx, address, nil, f.checks...)
	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndCreate(address)
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "error creating an address")
		} else if err != nil {
			return apperror.NewQueryError("Address", err, "")
		}
		return nil
	})
	if txnErr != nil {
		return nil, txnErr
	}

	return address, nil
}
