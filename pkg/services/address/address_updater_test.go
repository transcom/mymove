package address

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *AddressSuite) TestAddressUpdater() {
	createOriginalAddress := func() *models.Address {
		originalAddress := testdatagen.MakeAddress(suite.AppContextForTest().DB(), testdatagen.Assertions{})
		return &originalAddress
	}

	streetAddress1 := "288 SW Sunset Way"
	city := "Elizabethtown"
	state := "KY"
	postalCode := "42701"

	suite.Run("Successfully updates an address", func() {
		originalAddress := createOriginalAddress()

		addressUpdater := NewAddressUpdater()
		desiredAddress := &models.Address{
			ID:             originalAddress.ID,
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.NotNil(updatedAddress)
		suite.Nil(err)
		suite.Equal(originalAddress.ID, updatedAddress.ID)
		suite.Equal(desiredAddress.StreetAddress1, updatedAddress.StreetAddress1)
		suite.Equal(desiredAddress.StreetAddress2, updatedAddress.StreetAddress2)
		suite.Nil(updatedAddress.StreetAddress2)
		suite.Equal(desiredAddress.City, updatedAddress.City)
		suite.Equal(desiredAddress.State, updatedAddress.State)
		suite.Equal(desiredAddress.PostalCode, updatedAddress.PostalCode)
		suite.Equal(desiredAddress.Country, updatedAddress.Country)
		suite.Nil(updatedAddress.Country)
	})

	suite.Run("Fails to updates an address because of missing fields (eg. of failure in ValidateAndUpdate)", func() {
		originalAddress := createOriginalAddress()

		addressUpdater := NewAddressUpdater()
		desiredAddress := &models.Address{
			ID:             originalAddress.ID,
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(time.Now()))

		suite.Nil(updatedAddress)
		suite.NotNil(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Fails to updates because of stale etag", func() {
		originalAddress := createOriginalAddress()

		addressUpdater := NewAddressUpdater()
		desiredAddress := &models.Address{
			ID: originalAddress.ID,
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.Nil(updatedAddress)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("invalid input while updating an address", err.Error())
		errors := err.(apperror.InvalidInputError)
		suite.Len(errors.ValidationErrors.Errors, 4)
		suite.Contains(errors.ValidationErrors.Keys(), "street_address1")
		suite.Contains(errors.ValidationErrors.Keys(), "city")
		suite.Contains(errors.ValidationErrors.Keys(), "state")
		suite.Contains(errors.ValidationErrors.Keys(), "postal_code")
	})

	suite.Run("Fails to update an address because of invalid ID", func() {
		originalAddress := createOriginalAddress()

		addressUpdater := NewAddressUpdater()
		desiredAddress := &models.Address{
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.Nil(updatedAddress)
		suite.NotNil(err)
		suite.IsType(&apperror.BadDataError{}, err)
		expectedError := fmt.Sprintf("Data received from requester is bad: %s: invalid ID used for address", apperror.BadDataCode)
		suite.Equal(expectedError, err.Error())
	})
}
