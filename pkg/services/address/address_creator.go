package address

import (
	"fmt"

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

	// until international moves are supported, we will default the country for created addresses to "US"
	if address.Country != nil && address.Country.Country != "US" {
		return nil, fmt.Errorf("- the country %s is not supported at this time - only US is allowed", address.Country.Country)
	}

	if address.Country != nil && address.Country.Country != "" {
		country, err := models.FetchCountryByCode(appCtx.DB(), address.Country.Country)
		if err != nil {
			return nil, err
		}
		transformedAddress.Country = &country
		transformedAddress.CountryId = &country.ID
	} else {
		country, err := models.FetchCountryByCode(appCtx.DB(), "US")
		if err != nil {
			return nil, err
		}
		transformedAddress.Country = &country
		transformedAddress.CountryId = &country.ID
		transformedAddress.Country = &country
	}

	// use the data we have first, if it's not nil
	if transformedAddress.Country != nil {
		country := transformedAddress.Country
		if country.Country != "US" || country.Country == "US" && transformedAddress.State == "AK" || country.Country == "US" && transformedAddress.State == "HI" {
			boolTrueVal := true
			transformedAddress.IsOconus = &boolTrueVal
		} else {
			boolFalseVal := false
			transformedAddress.IsOconus = &boolFalseVal
		}
	} else if transformedAddress.CountryId != nil {
		country, err := models.FetchCountryByID(appCtx.DB(), *transformedAddress.CountryId)
		if err != nil {
			return nil, err
		}
		if country.Country != "US" || country.Country == "US" && transformedAddress.State == "AK" || country.Country == "US" && transformedAddress.State == "HI" {
			boolTrueVal := true
			transformedAddress.IsOconus = &boolTrueVal
		} else {
			boolFalseVal := false
			transformedAddress.IsOconus = &boolFalseVal
		}
	} else {
		boolFalseVal := false
		transformedAddress.IsOconus = &boolFalseVal
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

	if transformedAddress.Country != nil && transformedAddress.Country.Country == "" {
		transformedAddress.Country = nil
	}

	return transformedAddress
}
