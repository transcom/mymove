package address

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AddressSuite) TestAddressCreator() {
	streetAddress1 := "288 SW Sunset Way"
	city := "Elizabethtown"
	state := "KY"
	postalCode := "42701"

	suite.Run("Successfully creates an address", func() {
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		})

		suite.Nil(err)
		suite.NotNil(address)
		suite.NotNil(address.ID)
		suite.Equal(streetAddress1, address.StreetAddress1)
		suite.Equal(city, address.City)
		suite.Equal(state, address.State)
		suite.Equal(postalCode, address.PostalCode)
		suite.Nil(address.StreetAddress2)
		suite.Nil(address.Country)
	})

	suite.Run("Successfully creates an address with empty strings for optional fields", func() {
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: models.StringPointer(""),
			StreetAddress3: models.StringPointer(""),
			City:           city,
			State:          state,
			PostalCode:     postalCode,
			Country:        models.StringPointer(""),
		})

		suite.Nil(err)
		suite.NotNil(address)
		suite.NotNil(address.ID)
		suite.Equal(streetAddress1, address.StreetAddress1)
		suite.Equal(city, address.City)
		suite.Equal(state, address.State)
		suite.Equal(postalCode, address.PostalCode)
		suite.Nil(address.StreetAddress2)
		suite.Nil(address.StreetAddress3)
		suite.Nil(address.Country)
	})

	suite.Run("Fails to add an address because an ID is passed (fails to pass rules check)", func() {
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			ID:             uuid.FromStringOrNil("06c82380-4fc3-469f-803d-76763e6f87dd"),
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
		})

		suite.Nil(address)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the address.", err.Error())
		errors := err.(apperror.InvalidInputError)
		suite.Len(errors.ValidationErrors.Errors, 1)
		suite.Contains(errors.ValidationErrors.Keys(), "ID")
	})

	suite.Run("Fails because of missing field", func() {
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{})

		suite.Nil(address)
		suite.NotNil(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("error creating an address", err.Error())
		errors := err.(apperror.InvalidInputError)
		suite.Len(errors.ValidationErrors.Errors, 4)
		suite.Contains(errors.ValidationErrors.Keys(), "street_address1")
		suite.Contains(errors.ValidationErrors.Keys(), "city")
		suite.Contains(errors.ValidationErrors.Keys(), "state")
		suite.Contains(errors.ValidationErrors.Keys(), "postal_code")
	})
}
