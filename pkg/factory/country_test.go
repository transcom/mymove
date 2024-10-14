package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildUSCountry() {
	suite.Run("Successful creation of US country", func() {
		defaultCountry := models.Country{
			Country:     "US",
			CountryName: "UNITED STATES",
		}

		country := FetchOrBuildCountry(suite.DB(), nil, nil)
		suite.Equal(defaultCountry.Country, country.Country)
		suite.Equal(defaultCountry.CountryName, country.CountryName)
	})
}
