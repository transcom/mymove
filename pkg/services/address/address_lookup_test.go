package address

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
)

func (suite *AddressSuite) TestAddressLookup() {
	city := "DEERFIELD"
	state := "NH"
	postalCode := "03037"
	county := "ROCKINGHAM"

	suite.Run("Successfully search for location by zip", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, postalCode)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
		suite.Contains((*address)[0].StateName, state)
		suite.Contains((*address)[0].UsprZipID, postalCode)
		suite.Contains((*address)[0].UsprcCountyNm, county)
	})

	suite.Run("Successfully search for location by city name", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, city)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
	})

	suite.Run("Successfully search for location by city, state", func() {
		search := city + ", " + state
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{})
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
		suite.Contains((*address)[0].StateName, state)
	})
}
