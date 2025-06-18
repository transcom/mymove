package address

import (
	"slices"
	"strings"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AddressSuite) TestAddressLookup() {
	city := "DEERFIELD"
	state := "NH"
	postalCode := "03037"
	county := "ROCKINGHAM"
	excludedStates := [...]string{"AK", "HI"}

	POBoxCity := "FREDERIKSTED"
	POBoxState := "VI"
	POBoxPostalCode := "00841"
	POBoxCounty := "SAINT CROIX"

	suite.Run("Successfully search for location by zip", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, postalCode, excludedStates[:], true)

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
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, city, excludedStates[:], true)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
	})

	suite.Run("Successfully search for location by city, state", func() {
		search := city + ", " + state
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:], true)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, city)
		suite.Contains((*address)[0].StateName, state)
	})

	suite.Run("Successfully search for location by city, state postalCode", func() {
		search := city + ", " + state + " " + postalCode
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:], true)

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
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, akSearch, excludedStates[:], true)
		suite.Nil(err)
		suite.Nil((*address))

		hiSearch := "HONOLULU, HI 96835"
		address, err = addressLookup.GetLocationsByZipCityState(appCtx, hiSearch, excludedStates[:], true)
		suite.Nil(err)
		suite.Nil((*address))
	})

	suite.Run("Successfully search for PO Box location by zip", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, POBoxPostalCode, excludedStates[:], true)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Contains((*address)[0].CityName, POBoxCity)
		suite.Contains((*address)[0].StateName, POBoxState)
		suite.Contains((*address)[0].UsprZipID, POBoxPostalCode)
		suite.Contains((*address)[0].UsprcCountyNm, POBoxCounty)
	})

	suite.Run("Successfully return nothing in search for PO Box location by zip when PO Boxes are excluded", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, POBoxPostalCode, excludedStates[:], false)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Equal(len(*address), 0)
	})

	suite.Run("Successfully search by city excludes PO Boxes when PO Boxes are excluded", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		addresses, err := addressLookup.GetLocationsByZipCityState(appCtx, POBoxCity, excludedStates[:], false)

		suite.Nil(err)
		suite.NotNil(addresses)

		includesPOBoxLocation := slices.ContainsFunc((*addresses), func(l models.VLocation) bool {
			return l.UsprZipID == POBoxPostalCode
		})

		suite.Equal(includesPOBoxLocation, false)
	})

	suite.Run("Successfully search by city and state includes PO Boxes when PO Boxes are included", func() {
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		search := POBoxCity + ", " + POBoxState
		addresses, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:], true)

		suite.Nil(err)
		suite.NotNil(addresses)

		includesPOBoxLocation := slices.ContainsFunc((*addresses), func(l models.VLocation) bool {
			return l.UsprZipID == POBoxPostalCode
		})

		suite.Equal(includesPOBoxLocation, true)
	})

	suite.Run("Successfully return nothing in search for PO Box location by exact match when PO Boxes are excluded", func() {
		search := POBoxCity + ", " + POBoxState + " " + POBoxPostalCode
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:], false, true)

		suite.Nil(err)
		suite.NotNil(address)
		suite.Equal(len(*address), 0)
	})

	suite.Run("Successfully return match in search for PO Box location by exact match when PO Boxes are included", func() {
		search := POBoxCity + ", " + POBoxState + " " + POBoxPostalCode
		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVLocation()
		address, err := addressLookup.GetLocationsByZipCityState(appCtx, search, excludedStates[:], true, true)

		suite.Nil(err)
		suite.NotNil(address)
		suite.GreaterOrEqual(len(*address), 1)
	})
}

func (suite *AddressSuite) TestOconusAddressLookup() {
	country := "GB"
	city := "LONDON"

	suite.Run("Successfully search for location by principal division", func() {
		principalDivision := "CARDIFF"
		principalDivisionForSearch := "LONDON, CARDIFF"

		appCtx := appcontext.NewAppContext(suite.AppContextForTest().DB(), suite.AppContextForTest().Logger(), &auth.Session{}, nil)
		addressLookup := NewVIntlLocation()
		address, err := addressLookup.GetOconusLocations(appCtx, country, principalDivisionForSearch, false)

		suite.Nil(err)
		suite.NotNil(address)
		returnCity := (*address)[0].CityName
		suite.Contains(strings.ToUpper(*returnCity), strings.ToUpper(city))
		returnedPrincipalDivision := (*address)[0].CountryPrnDivName
		suite.Contains(strings.ToUpper(*returnedPrincipalDivision), strings.ToUpper(principalDivision))
	})

	suite.Run("Successfully search for location by city name", func() {
		principalDivision := "ABERDEEN CITY"

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
