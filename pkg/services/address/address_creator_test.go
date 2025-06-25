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
	oConusState := "AK"

	suite.Run("Successfully creates a CONUS address", func() {
		country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NoError(err)
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     postalCode,
			CountryId:      &country.ID,
			Country:        &country,
		})

		suite.Nil(err)
		suite.NotNil(address)
		suite.NotNil(address.ID)
		suite.Equal(streetAddress1, address.StreetAddress1)
		suite.Equal(city, address.City)
		suite.Equal(state, address.State)
		suite.Equal(postalCode, address.PostalCode)
		suite.False(*address.IsOconus)
		suite.Nil(address.StreetAddress2)
		suite.NotNil(address.Country)
		suite.NotNil(address.UsPostRegionCity)
		suite.NotNil(address.UsPostRegionCityID)
	})

	suite.Run("Successfully creates an OCONUS address with AK state", func() {
		country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NoError(err)
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: streetAddress1,
			City:           city,
			State:          oConusState,
			PostalCode:     postalCode,
			IsOconus:       models.BoolPointer(true),
			CountryId:      &country.ID,
			Country:        &country,
		})

		suite.Nil(err)
		suite.NotNil(address)
		suite.NotNil(address.ID)
		suite.Equal(streetAddress1, address.StreetAddress1)
		suite.Equal(city, address.City)
		suite.Equal(oConusState, address.State)
		suite.Equal(postalCode, address.PostalCode)
		suite.True(*address.IsOconus)
		suite.Nil(address.StreetAddress2)
		suite.NotNil(address.Country)
		suite.NotNil(address.UsPostRegionCity)
		suite.NotNil(address.UsPostRegionCityID)
	})

	suite.Run("Successfully creates an address with empty strings for optional fields", func() {
		country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NoError(err)

		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: streetAddress1,
			StreetAddress2: models.StringPointer(""),
			StreetAddress3: models.StringPointer(""),
			City:           city,
			State:          state,
			PostalCode:     postalCode,
			CountryId:      &country.ID,
			Country:        &country,
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
		suite.NotNil(address.Country)
		suite.NotNil(address.UsPostRegionCity)
		suite.NotNil(address.UsPostRegionCityID)
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

		usprc, err := models.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "35007", "ALABASTER")
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{UsPostRegionCityID: &usprc.ID})

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

	suite.Run("Fails when zip code is invalid", func() {
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: streetAddress1,
			City:           city,
			State:          state,
			PostalCode:     "11111",
		})

		suite.Nil(address)
		suite.NotNil(err)
		suite.Equal("No county found for provided zip code 11111.", err.Error())
	})

	suite.Run("Successfully creates a CONUS address", func() {
		country, err := models.FetchCountryByCode(suite.DB(), "US")
		suite.NoError(err)
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: "7645 Ballinshire N",
			City:           "Indianapolis",
			State:          "IN",
			PostalCode:     "46254",
			Country:        &country,
			CountryId:      &country.ID,
		})

		suite.False(*address.IsOconus)
		suite.NotNil(address.ID)
		suite.Nil(err)
		suite.NotNil(address.Country)
	})

	suite.Run("Successfully creates a CONUS address", func() {
		country := &models.Country{}
		country.Country = "US"
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: "7645 Ballinshire N",
			City:           "Indianapolis",
			State:          "IN",
			PostalCode:     "46254",
			Country:        country,
		})

		suite.False(*address.IsOconus)
		suite.NotNil(address.ID)
		suite.Nil(err)
		suite.NotNil(address.Country)
	})
	suite.Run("Fails when us_post_region_city is not found", func() {
		country := &models.Country{}
		country.Country = "US"
		addressCreator := NewAddressCreator()
		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1: "7645 Ballinshire N",
			City:           "Charlotte",
			State:          "IN",
			PostalCode:     "46254",
			Country:        country,
		})

		suite.NotNil(err)
		suite.Nil(address)
		suite.Equal("No UsPostRegionCity found for provided zip code 46254 and city Charlotte.", err.Error())
	})

	suite.Run("returns error when address has an invalid USPRC assignment", func() {
		country := &models.Country{}
		country.Country = "US"
		addressCreator := NewAddressCreator()

		usprc, err := models.FindByZipCodeAndCity(suite.DB(), "29229", "Columbia")
		suite.NoError(err)

		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1:     "7645 Ballinshire N",
			City:               "Indianapolis",
			State:              "IN",
			PostalCode:         "46254",
			Country:            country,
			UsPostRegionCityID: &usprc.ID,
			UsPostRegionCity:   usprc,
		})

		suite.Nil(address)
		suite.Error(err, "error creating an address")
		errors := err.(apperror.InvalidInputError)
		suite.Contains(errors.ValidationErrors.Keys(), "us_post_region_city_id")
	})

	suite.Run("returns error when USPRC validation fails", func() {
		country := &models.Country{}
		country.Country = "US"
		addressCreator := NewAddressCreator()

		usprc, err := models.FindByZipCodeAndCity(suite.DB(), "29229", "Columbia")
		suite.NoError(err)

		address, err := addressCreator.CreateAddress(suite.AppContextForTest(), &models.Address{
			StreetAddress1:     "7645 Ballinshire N",
			City:               "Indianapolis",
			State:              "IN",
			PostalCode:         "29229",
			Country:            country,
			UsPostRegionCityID: &usprc.ID,
			UsPostRegionCity:   usprc,
		})

		suite.Nil(address)
		suite.Error(err, "No UsPostRegionCity found for provided zip code 29229 and city Indianapolis.")
	})
}
