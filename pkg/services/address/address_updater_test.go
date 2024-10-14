package address

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AddressSuite) TestAddressUpdater() {
	createOriginalAddress := func() *models.Address {
		originalAddress := factory.BuildAddress(suite.AppContextForTest().DB(), nil, nil)
		return &originalAddress
	}

	streetAddress1 := "288 SW Sunset Way"
	city := "Elizabethtown"
	state := "KY"
	postalCode := "42701"
	county := "HARDIN"

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
		suite.Equal(desiredAddress.City, updatedAddress.City)
		suite.Equal(desiredAddress.State, updatedAddress.State)
		suite.Equal(desiredAddress.PostalCode, updatedAddress.PostalCode)
		suite.NotNil(updatedAddress.StreetAddress2)
		suite.Equal(originalAddress.StreetAddress2, updatedAddress.StreetAddress2)
		suite.NotNil(updatedAddress.StreetAddress3)
		suite.Equal(originalAddress.StreetAddress3, updatedAddress.StreetAddress3)
		suite.NotNil(updatedAddress.Country)
		suite.Equal(county, desiredAddress.County)
	})

	suite.Run("Fails to updates because of stale etag", func() {
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

	suite.Run("Fails to updates an address because of invalid input (eg. of failure in ValidateAndUpdate)", func() {
		originalAddress := createOriginalAddress()

		addressUpdater := NewAddressUpdater()
		desiredAddress := &models.Address{
			ID:             originalAddress.ID,
			StreetAddress1: " ",
			City:           " ",
			State:          " ",
			PostalCode:     postalCode, // Provide postal code here because it is not explicitly input error
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.Nil(updatedAddress)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("invalid input while updating an address", err.Error())
		errors := err.(apperror.InvalidInputError)
		suite.Len(errors.ValidationErrors.Errors, 3)
		suite.Contains(errors.ValidationErrors.Keys(), "street_address1")
		suite.Contains(errors.ValidationErrors.Keys(), "city")
		suite.Contains(errors.ValidationErrors.Keys(), "state")
	})

	suite.Run("Fails to updates an address because of invalid county", func() {
		originalAddress := createOriginalAddress()

		addressUpdater := NewAddressUpdater()
		// Street address, city, and state are not related to how a postal code gets its county at this time
		desiredAddress := &models.Address{
			ID:             originalAddress.ID,
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     " ",
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.Nil(updatedAddress)
		suite.NotNil(err)
		suite.Equal("No county found for provided zip code  .", err.Error())
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

	suite.Run("Able to update when providing US country value in updated address", func() {
		originalAddress := createOriginalAddress()
		addressUpdater := NewAddressUpdater()

		desiredAddress := &models.Address{
			ID:             originalAddress.ID,
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
			Country:        &models.Country{Country: "US"},
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.NoError(err)
		suite.NotNil(updatedAddress)
		suite.Equal(updatedAddress.Country.Country, "US")
	})

	suite.Run("Receives an error when trying to update to an international address", func() {
		originalAddress := createOriginalAddress()
		addressUpdater := NewAddressUpdater()

		desiredAddress := &models.Address{
			ID:             originalAddress.ID,
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
			Country:        &models.Country{Country: "GB"},
		}
		updatedAddress, err := addressUpdater.UpdateAddress(suite.AppContextForTest(), desiredAddress, etag.GenerateEtag(originalAddress.UpdatedAt))

		suite.Error(err)
		suite.Nil(updatedAddress)
		suite.Equal("- the country GB is not supported at this time - only US is allowed", err.Error())
	})
}
