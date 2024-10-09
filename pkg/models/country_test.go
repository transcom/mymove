package models_test

import (
	"errors"

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
		suite.verifyValidationErrors(&newCountry, expErrors)
	})

	suite.Run("test empty Country", func() {
		emptyCountry := models.Country{}
		expErrors := map[string][]string{
			"country":      {"Country can not be blank."},
			"country_name": {"CountryName can not be blank."},
		}
		suite.verifyValidationErrors(&emptyCountry, expErrors)
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
