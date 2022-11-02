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
	transformedAddress := transformNilValuesForOptionalFields(*address)
	err := validateAddress(appCtx, &transformedAddress, nil, f.checks...)
	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndCreate(&transformedAddress)
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

	return &transformedAddress, nil
}

func transformNilValuesForOptionalFields(address models.Address) models.Address {
	transformedAddress := address

	if transformedAddress.StreetAddress2 != nil && *transformedAddress.StreetAddress2 == "" {
		transformedAddress.StreetAddress2 = nil
	}

	if transformedAddress.StreetAddress3 != nil && *transformedAddress.StreetAddress3 == "" {
		transformedAddress.StreetAddress3 = nil
	}

	if transformedAddress.Country != nil && *transformedAddress.Country == "" {
		transformedAddress.Country = nil
	}

	return transformedAddress
}
