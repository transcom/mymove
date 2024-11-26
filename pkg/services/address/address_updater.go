package address

import (
	"fmt"

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

	// Fetch new county if postal code has been modified from its original
	if originalAddress.PostalCode != address.PostalCode || address.County == nil {
		county, err := models.FindCountyByZipCode(appCtx.DB(), address.PostalCode)
		if err != nil {
			return nil, err
		}
		address.County = county
	}

	mergedAddress := mergeAddress(*address, *originalAddress)

	err := validateAddress(appCtx, &mergedAddress, originalAddress, f.checks...)
	if err != nil {
		return nil, err
	}

	// until international moves are supported, we will default the country for created addresses to "US"
	if mergedAddress.Country != nil && mergedAddress.Country.Country != "US" {
		return nil, fmt.Errorf("- the country %s is not supported at this time - only US is allowed", mergedAddress.Country.Country)
	}
	// first we will check to see if the country values have changed at all
	// until international moves are supported, we will default the country for created addresses to "US"
	if mergedAddress.Country != nil && mergedAddress.Country.Country != "" && mergedAddress.Country != originalAddress.Country {
		country, err := models.FetchCountryByCode(appCtx.DB(), address.Country.Country)
		if err != nil {
			return nil, err
		}
		mergedAddress.Country = &country
		mergedAddress.CountryId = &country.ID
	} else if mergedAddress.Country == nil {
		country, err := models.FetchCountryByCode(appCtx.DB(), "US")
		if err != nil {
			return nil, err
		}
		mergedAddress.Country = &country
		mergedAddress.CountryId = &country.ID
	}

	// Evaluate address and populate addresses isOconus value
	isOconus, err := models.IsAddressOconus(appCtx.DB(), mergedAddress)
	if err != nil {
		return nil, err
	}
	mergedAddress.IsOconus = &isOconus

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().ValidateAndUpdate(&mergedAddress)
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

	return &mergedAddress, nil
}

func mergeAddress(address, originalAddress models.Address) models.Address {
	mergedAddress := originalAddress
	if address.StreetAddress1 != "" {
		mergedAddress.StreetAddress1 = address.StreetAddress1
	}
	if address.City != "" {
		mergedAddress.City = address.City
	}
	if address.State != "" {
		mergedAddress.State = address.State
	}
	if address.PostalCode != "" {
		mergedAddress.PostalCode = address.PostalCode
	}
	if address.County != nil && *address.County != "" {
		mergedAddress.County = address.County
	}

	mergedAddress.StreetAddress2 = services.SetOptionalStringField(address.StreetAddress2, mergedAddress.StreetAddress2)
	mergedAddress.StreetAddress3 = services.SetOptionalStringField(address.StreetAddress3, mergedAddress.StreetAddress3)
	if address.Country != nil {
		mergedAddress.Country.Country = *services.SetOptionalStringField(&address.Country.Country, &mergedAddress.Country.Country)
	}
	return mergedAddress
}
