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
	excludedStates := [...]string{"AK", "HI"}

	suite.Run("Successfully search for location by zip", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, postalCode, excludedStates[:])

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
		suite.Contains((*address)[0].StateName, state)
		suite.Contains((*address)[0].UsprZipID, postalCode)
		suite.Contains((*address)[0].UsprcCountyNm, county)
	})

	suite.Run("Successfully search for location by city name", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, city, excludedStates[:])

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
	})

	suite.Run("Successfully search for location by city, state", func() {
		search := city + ", " + state
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:])

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
		suite.Contains((*address)[0].StateName, state)
	})

	suite.Run("Successfully search for location by city, state postalCode", func() {
		search := city + ", " + state + " " + postalCode
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:])

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
		suite.Contains((*address)[0].StateName, state)
		suite.Contains((*address)[0].UsprZipID, postalCode)
	})

	suite.Run("Search for excluded state returns nothing", func() {
		akSearch := "ANCHORAGE, AK 99503"
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, akSearch, excludedStates[:])
		suite.Nil(err)
		suite.Nil((*address))

		hiSearch := "HONOLULU, HI 96835"
		address, err = addressLookup.GetLocationsByZipCityState(appCtx, hiSearch, excludedStates[:])
		suite.Nil(err)
		suite.Nil((*address))
	})
}
