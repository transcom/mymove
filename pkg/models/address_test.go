package models_test

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	m "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
	suite.verifyValidationErrors(&newAddress, expErrors, nil)
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

func (suite *ModelSuite) Test_FetchDutyLocationGblocForAK() {
	setupDataForOconusDutyLocation := func(postalCode string) (m.OconusRateArea, m.UsPostRegionCity, m.DutyLocation) {
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

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: m.ReRateArea{
				ContractID: contract.ID,
				IsOconus:   true,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			},
		})
		suite.NotNil(rateArea)
		suite.Nil(err)

		us_country, err := m.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		oconusRateArea, err := m.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
		suite.NotNil(oconusRateArea)
		suite.Nil(err)

		return *oconusRateArea, *usprc, originDutyLocation
	}

	suite.Run("fetches duty location GBLOC for AK address, Zone II AirForce", func() {
		oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99707")

		airForce := m.AffiliationAIRFORCE
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &airForce,
				},
			},
		}, nil)

		jppsoRegion, err := m.FetchJppsoRegionByCode(suite.DB(), "MBFL")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := m.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, m.DepartmentIndicatorAIRANDSPACEFORCE.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		gbloc, err := m.FetchAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(string(*gbloc), "MBFL")
	})

	suite.Run("fetches duty location GBLOC for AK address, Zone II Army", func() {
		oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99707")

		army := m.AffiliationARMY
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		jppsoRegion, err := m.FetchJppsoRegionByCode(suite.DB(), "JEAT")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := m.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, m.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		gbloc, err := m.FetchAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(string(*gbloc), "JEAT")
	})

	suite.Run("fetches duty location GBLOC for AK Cordova address, Zone IV", func() {
		usprc, err := m.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "99574", "CORDOVA")
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

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})

		rateAreaCode := uuid.Must(uuid.NewV4()).String()[0:5]
		rateArea := testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: m.ReRateArea{
				ContractID: contract.ID,
				IsOconus:   true,
				Name:       fmt.Sprintf("Alaska-%s", rateAreaCode),
				Contract:   contract,
			},
		})
		suite.NotNil(rateArea)

		us_country, err := m.FetchCountryByCode(suite.DB(), "US")
		suite.NotNil(us_country)
		suite.Nil(err)

		oconusRateArea, err := m.FetchOconusRateAreaByCityId(suite.DB(), usprc.ID.String())
		suite.NotNil(oconusRateArea)
		suite.Nil(err)
		army := m.AffiliationARMY
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		jppsoRegion, err := m.FetchJppsoRegionByCode(suite.DB(), "MAPS")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := m.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, m.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		gbloc, err := m.FetchAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(string(*gbloc), "MAPS")
	})

	suite.Run("fetches duty location GBLOC for AK NOT Cordova address, Zone IV", func() {
		oconusRateArea, _, originDutyLocation := setupDataForOconusDutyLocation("99803")

		army := m.AffiliationARMY
		serviceMember := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: m.ServiceMember{
					Affiliation: &army,
				},
			},
		}, nil)

		jppsoRegion, err := m.FetchJppsoRegionByCode(suite.DB(), "MAPK")
		suite.NotNil(jppsoRegion)
		suite.Nil(err)

		gblocAors, err := m.FetchGblocAorsByJppsoCodeRateAreaDept(suite.DB(), jppsoRegion.ID, oconusRateArea.ID, m.DepartmentIndicatorARMY.String())
		suite.NotNil(gblocAors)
		suite.Nil(err)

		gbloc, err := m.FetchAddressGbloc(suite.DB(), originDutyLocation.Address, serviceMember)
		suite.NoError(err)
		suite.NotNil(gbloc)
		suite.Equal(string(*gbloc), "MAPK")
	})
}

func (suite *ModelSuite) TestIsAddressAlaska() {
	var address *m.Address
	bool1, err := address.IsAddressAlaska()
	suite.Error(err)
	suite.Equal("address is nil", err.Error())
	suite.Equal(false, bool1)

	address = &m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		PostalCode:     "90210",
		County:         m.StringPointer("County"),
	}

	bool2, err := address.IsAddressAlaska()
	suite.NoError(err)
	suite.Equal(m.BoolPointer(false), &bool2)

	address.State = "MT"
	bool3, err := address.IsAddressAlaska()
	suite.NoError(err)
	suite.Equal(m.BoolPointer(false), &bool3)

	address.State = "AK"
	bool4, err := address.IsAddressAlaska()
	suite.NoError(err)
	suite.Equal(m.BoolPointer(true), &bool4)
}
