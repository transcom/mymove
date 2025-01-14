package models_test

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	m "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicAddressInstantiation() {
	newAddress := &m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		County:         m.StringPointer("County"),
	}

	verrs, err := newAddress.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestEmptyAddressInstantiation() {
	newAddress := m.Address{}

	expErrors := map[string][]string{
		"street_address1": {"StreetAddress1 can not be blank."},
		"city":            {"City can not be blank."},
		"state":           {"State can not be blank."},
		"postal_code":     {"PostalCode can not be blank."},
	}
	suite.verifyValidationErrors(&newAddress, expErrors)
}

func (suite *ModelSuite) TestAddressCountryCode() {
	noCountry := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		County:         m.StringPointer("county"),
	}

	var expected *string
	countryCode, err := noCountry.CountryCode()
	suite.NoError(err)
	suite.Equal(expected, countryCode)

	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	usCountry := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		Country:        &country,
	}
	countryCode, err = usCountry.CountryCode()
	suite.NoError(err)
	suite.Equal("US", *countryCode)
}

func (suite *ModelSuite) TestIsAddressOconusNoCountry() {
	address := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "SC",
		PostalCode:     "29229",
		County:         m.StringPointer("county"),
	}

	result, err := m.IsAddressOconus(suite.DB(), address)
	suite.NoError(err)

	suite.Equal(false, result)
}

// Test IsOconus logic for an address with no country and a state of AK
func (suite *ModelSuite) TestIsAddressOconusForAKState() {
	address := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "Anchorage",
		State:          "AK",
		PostalCode:     "99502",
		County:         m.StringPointer("county"),
	}

	result, err := m.IsAddressOconus(suite.DB(), address)
	suite.NoError(err)

	suite.Equal(true, result)
}

func (suite *ModelSuite) TestAddressFormat() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	newAddress := &m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		County:         m.StringPointer("County"),
		Country:        &country,
		CountryId:      &country.ID,
	}

	verrs, err := newAddress.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")

	formattedAddress := newAddress.Format()

	suite.Equal("street 1\nstreet 2\nstreet 3\ncity, state 90210", formattedAddress)

	formattedAddress = newAddress.LineFormat()

	suite.Equal("street 1, street 2, street 3, city, state, 90210, UNITED STATES", formattedAddress)

	formattedAddress = newAddress.LineDisplayFormat()

	suite.Equal("street 1 street 2 street 3, city, state 90210", formattedAddress)
}

func (suite *ModelSuite) TestPartialAddressFormat() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	newAddress := &m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: nil,
		StreetAddress3: nil,
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		County:         m.StringPointer("County"),
		Country:        &country,
		CountryId:      &country.ID,
	}

	verrs, err := newAddress.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")

	formattedAddress := newAddress.Format()

	suite.Equal("street 1\ncity, state 90210", formattedAddress)

	formattedAddress = newAddress.LineFormat()

	suite.Equal("street 1, city, state, 90210, UNITED STATES", formattedAddress)

	formattedAddress = newAddress.LineDisplayFormat()

	suite.Equal("street 1, city, state 90210", formattedAddress)
}

func (suite *ModelSuite) TestAddressIsEmpty() {
	suite.Run("empty whitespace address", func() {
		testAddress := m.Address{
			StreetAddress1: " ",
			State:          " ",
			PostalCode:     " ",
		}
		suite.True(m.IsAddressEmpty(&testAddress))
	})
	suite.Run("empty n/a address", func() {
		testAddress := m.Address{
			StreetAddress1: "n/a",
			State:          "n/a",
			PostalCode:     "n/a",
		}
		suite.True(m.IsAddressEmpty(&testAddress))
	})
	suite.Run("nonempty address", func() {
		testAddress := m.Address{
			StreetAddress1: "street 1",
			State:          "state",
			PostalCode:     "90210",
		}
		suite.False(m.IsAddressEmpty(&testAddress))
	})
}

func (suite *ModelSuite) Test_FetchDutyLocationGblocForAK() {
	setupDataForOconusDutyLocation := func(postalCode string) (m.ReRateArea, m.OconusRateArea, m.UsPostRegionCity, m.DutyLocation) {
		usprc, err := m.FindByZipCode(suite.AppContextForTest().DB(), postalCode)
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: m.Address{
					IsOconus:           m.BoolPointer(true),
					UsPostRegionCityID: &usprc.ID,
				},
			},
		}, nil)
		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model:    address,
				LinkOnly: true,
				Type:     &factory.Addresses.DutyLocationAddress,
			},
		}, nil)

		contract := m.ReContract{
			Code: "Test_create_oconus_order_code",
			Name: "Test_create_oconus_order",
		}
		verrs, err := suite.AppContextForTest().DB().ValidateAndSave(&contract)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(verrs.Error())
		}

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := m.ReRateArea{
			ID:         uuid.Must(uuid.NewV4()),
			ContractID: contract.ID,
			IsOconus:   true,
			Code:       rateAreaCode,
			Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
			Contract:   contract,
		}
		verrs, err = suite.DB().ValidateAndCreate(&rateArea)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		us_country, err := m.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		oconusRateArea := m.OconusRateArea{
			ID:                 uuid.Must(uuid.NewV4()),
			RateAreaId:         rateArea.ID,
			CountryId:          us_country.ID,
			UsPostRegionCityId: usprc.ID,
			Active:             true,
		}
		verrs, err = suite.DB().ValidateAndCreate(&oconusRateArea)
		if verrs.HasAny() {
			suite.Fail(verrs.Error())
		}
		if err != nil {
			suite.Fail(err.Error())
		}

		return rateArea, oconusRateArea, *usprc, originDutyLocation
	}

	suite.Run("fetches duty location GBLOC for AK address, Zone II AirForce", func() {
		_, oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99707")

		airForce := m.AffiliationAIRFORCE
		defaultDepartmentIndicator := "AIR_AND_SPACE_FORCE"
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &airForce,
				},
			},
		}, nil)

		jppsoRegion := m.JppsoRegions{
			Name: "JPPSO Elmendorf-Richardson",
			Code: "MBFL",
		}
		suite.MustSave(&jppsoRegion)

		gblocAors := m.GblocAors{
			JppsoRegionID:       jppsoRegion.ID,
			OconusRateAreaID:    oconusRateArea.ID,
			DepartmentIndicator: &defaultDepartmentIndicator,
		}
		suite.MustSave(&gblocAors)

		gbloc, err := m.FetchOconusAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(gbloc.Gbloc, "MBFL")
	})

	suite.Run("fetches duty location GBLOC for AK address, Zone II Army", func() {
		_, oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99707")

		army := m.AffiliationARMY
		defaultDepartmentIndicator := "ARMY"
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		jppsoRegion := m.JppsoRegions{
			Name: "JPPSO-Northwest",
			Code: "JEAT",
		}
		suite.MustSave(&jppsoRegion)

		gblocAors := m.GblocAors{
			JppsoRegionID:       jppsoRegion.ID,
			OconusRateAreaID:    oconusRateArea.ID,
			DepartmentIndicator: &defaultDepartmentIndicator,
		}
		suite.MustSave(&gblocAors)

		gbloc, err := m.FetchOconusAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(gbloc.Gbloc, "JEAT")
	})

	suite.Run("fetches duty location GBLOC for AK Cordova address, Zone IV", func() {
		_, oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99574")

		army := m.AffiliationARMY
		defaultDepartmentIndicator := "ARMY"
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		jppsoRegion := m.JppsoRegions{
			Name: "USCG Base Kodiak",
			Code: "MAPS",
		}
		suite.MustSave(&jppsoRegion)

		gblocAors := m.GblocAors{
			JppsoRegionID:       jppsoRegion.ID,
			OconusRateAreaID:    oconusRateArea.ID,
			DepartmentIndicator: &defaultDepartmentIndicator,
		}
		suite.MustSave(&gblocAors)

		gbloc, err := m.FetchOconusAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(gbloc.Gbloc, "MAPS")
	})

	suite.Run("fetches duty location GBLOC for AK NOT Cordova address, Zone IV", func() {
		_, oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99803")

		army := m.AffiliationARMY
		defaultDepartmentIndicator := "ARMY"
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		jppsoRegion := m.JppsoRegions{
			Name: "USCG Base Ketchikan",
			Code: "MAPK",
		}
		suite.MustSave(&jppsoRegion)

		gblocAors := m.GblocAors{
			JppsoRegionID:       jppsoRegion.ID,
			OconusRateAreaID:    oconusRateArea.ID,
			DepartmentIndicator: &defaultDepartmentIndicator,
		}
		suite.MustSave(&gblocAors)

		gbloc, err := m.FetchOconusAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(gbloc.Gbloc, "MAPK")
	})
}
