package address

import (
	"strings"

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

func (suite *AddressSuite) TestOconusAddressLookup() {
	country := "GB"
	city := "SANDRIDGE"

	suite.Run("Successfully search for location by principal division", func() {
		principalDivision := "HERTFORDSHIRE"
		principalDivisionForSearch := ", HERTFORDSHIRE"
		principalCity := "BAYFORD"

		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVIntlLocation()
		address, err := addressLookup.GetOconusLocations(appCtx, country, principalDivisionForSearch, false)

		suite.Nil(err)
		suite.NotNil(address)
		returnCity := (*address)[0].CityName
		suite.Contains(strings.ToUpper(*returnCity), strings.ToUpper(principalCity))
		returnedPrincipalDivision := (*address)[0].CountryPrnDivName
		suite.Contains(strings.ToUpper(*returnedPrincipalDivision), strings.ToUpper(principalDivision))
	})

	suite.Run("Successfully search for location by city name", func() {
		principalDivision := "HERTFORDSHIRE"

		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVIntlLocation()
		address, err := addressLookup.GetOconusLocations(appCtx, country, city, false)

		suite.Nil(err)
		suite.NotNil(address)
		returnCity := (*address)[0].CityName
		suite.Contains(strings.ToUpper(*returnCity), strings.ToUpper(city))
		returnedPrincipalDivision := (*address)[0].CountryPrnDivName
		suite.Contains(strings.ToUpper(*returnedPrincipalDivision), strings.ToUpper(principalDivision))
	})
}
