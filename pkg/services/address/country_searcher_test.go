package address

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AddressSuite) TestCountrySearch() {
	suite.Run("Success - 2 characters search text 'US' - match on both country code and name", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, models.StringPointer("us"))

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) == 2)

		isFound := func(countries models.Countries, match string) bool {
			for _, country := range countries {
				if country.Country == match {
					return true
				}
			}
			return false
		}

		// matches on country code
		suite.True(isFound(countries, "US"))

		// matches on starts with country name "US VIRGIN ISLANDS"
		suite.True(isFound(countries, "VI"))
	})

	suite.Run("Success - nil searchQuery - search for all", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, nil)

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) == 274)
	})

	suite.Run("Success - empty string searchQuery - search for all", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, models.StringPointer("   "))

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) == 274)
	})

	suite.Run("Successfully search for starts with 'unit'", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, models.StringPointer("unit"))

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) == 3)

		isFound := func(countries models.Countries, match string) bool {
			for _, country := range countries {
				if country.CountryName == match {
					return true
				}
			}
			return false
		}

		suite.True(isFound(countries, "UNITED ARAB EMIRATES"))
		suite.True(isFound(countries, "UNITED KINGDOM"))
		suite.True(isFound(countries, "UNITED STATES"))
	})
}
