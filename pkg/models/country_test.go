package models_test

import (
	"errors"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCountryValidation() {
	suite.Run("test valid Country", func() {
		newCountry := models.Country{
			Country:     "US",
			CountryName: "UNITED STATES",
		}

		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&newCountry, expErrors, nil)
	})

	suite.Run("test empty Country", func() {
		emptyCountry := models.Country{}
		expErrors := map[string][]string{
			"country":      {"Country can not be blank."},
			"country_name": {"CountryName can not be blank."},
		}
		suite.verifyValidationErrors(&emptyCountry, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchCountryByCode() {

	suite.Run("successful fetch of country by code", func() {
		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)

		goodCountry, err := models.FetchCountryByCode(suite.DB(), country.Country)

		suite.NoError(err)
		suite.Equal(goodCountry.Country, country.Country)
		suite.Equal(goodCountry.CountryName, country.CountryName)
	})

	suite.Run("error when country doesn't exist", func() {
		_, err := models.FetchCountryByCode(suite.DB(), "AB")

		suite.Error(err)
		suite.True(errors.Is(err, models.ErrFetchNotFound), "Expected FETCH_NOT_FOUND error")
	})
}

func (suite *ModelSuite) TestFetchCountryByID() {

	suite.Run("successful fetch of country by id", func() {
		country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)

		goodCountry, err := models.FetchCountryByID(suite.DB(), country.ID)

		suite.NoError(err)
		suite.Equal(goodCountry.ID, country.ID)
		suite.Equal(goodCountry.Country, country.Country)
		suite.Equal(goodCountry.CountryName, country.CountryName)
	})

	suite.Run("error when country doesn't exist", func() {
		nonExistingID, err := uuid.NewV4()
		suite.NoError(err)

		_, err = models.FetchCountryByID(suite.DB(), nonExistingID)

		suite.Error(err)
		suite.Equal(models.ErrFetchNotFound, err)
	})
}

func (suite *ModelSuite) TestCountryIsEmpty() {

	emptyCountry := models.Country{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Holidays: models.CountryHolidays{
			models.CountryHoliday{},
		},
		Weekends: models.CountryWeekend{},
	}
	suite.True(emptyCountry.IsEmpty(), "country should be considered empty when it's missing an ID, Country, and CountryName")

	validCountry := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	suite.False(validCountry.IsEmpty(), "country should not be considered empty when it has required fields")

	countryEmptyId := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	countryEmptyId.ID = uuid.Nil
	suite.False(countryEmptyId.IsEmpty(), "country should not be considered empty when only it's ID is empty")

	countryEmptyCode := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	countryEmptyCode.Country = ""
	suite.False(countryEmptyCode.IsEmpty(), "country should not be considered empty when only it's Country field is empty")

	countryEmptyName := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	countryEmptyName.CountryName = ""
	suite.False(countryEmptyName.IsEmpty(), "country should not be considered empty when only it's CountryName is empty")
}
