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

	if address.PostalCode != "" {
		county, err := models.FindCountyByZipCode(appCtx.DB(), address.PostalCode)
		if err != nil {
			return nil, err
		}
		transformedAddress.County = county
	}

	if address.CountryId != nil {
		country, err := models.FetchCountryByID(appCtx.DB(), *address.CountryId)
		if err != nil {
			return nil, err
		}
		address.Country = &country
	} else {
		country, err := models.FetchCountryByCode(appCtx.DB(), "US")
		if err != nil {
			return nil, err
		}

		address.Country = &country
		address.CountryId = &country.ID
	}

	// Evaluate address and populate addresses isOconus value
	isOconus, err := models.IsAddressOconus(appCtx.DB(), transformedAddress)
	if err != nil {
		return nil, err
	}
	transformedAddress.IsOconus = &isOconus

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

	if transformedAddress.Country != nil && transformedAddress.Country.Country == "" {
		transformedAddress.Country = nil
	}

	return transformedAddress
}
