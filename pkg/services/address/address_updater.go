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

	mergedAddress := mergeAddress(*address, *originalAddress)

	err := validateAddress(appCtx, &mergedAddress, originalAddress, f.checks...)
	if err != nil {
		return nil, err
	}

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

	mergedAddress.StreetAddress2 = services.SetOptionalStringField(address.StreetAddress2, mergedAddress.StreetAddress2)
	mergedAddress.StreetAddress3 = services.SetOptionalStringField(address.StreetAddress3, mergedAddress.StreetAddress3)
	mergedAddress.Country = services.SetOptionalStringField(address.Country, mergedAddress.Country)

	return mergedAddress
}
