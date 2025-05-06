package address

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AddressSuite) TestCountrySearch() {
	suite.Run("Successfully search for US", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, models.StringPointer("us"))

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) == 1)
	})

	suite.Run("Successfully search for all", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, nil)

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) == 274)
	})

	suite.Run("Successfully search for starts with", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		countrySearcher := NewCountrySearcher()
		countries, err := countrySearcher.SearchCountries(appCtx, models.StringPointer("Uni"))

		suite.Nil(err)
		suite.NotNil(countries)
		suite.True(len(countries) > 0 && len(countries) < 274)
	})
}
