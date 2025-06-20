package models_test

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	m "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicAddressInstantiation() {

	usprc, err := models.FindByZipCodeAndCity(suite.DB(), "90210", "BEVERLY HILLS")
	suite.NoError(err)

	newAddress := &m.Address{
		StreetAddress1:     "street 1",
		StreetAddress2:     m.StringPointer("street 2"),
		StreetAddress3:     m.StringPointer("street 3"),
		City:               "BEVERLY HILLS",
		State:              "CA",
		PostalCode:         "90210",
		County:             m.StringPointer("County"),
		UsPostRegionCityID: &usprc.ID,
	}

	verrs, err := newAddress.Validate(suite.DB())

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestAddressInstantiationWithIncorrectUsPostRegionCityID() {

	usPostRegionCityID := uuid.Must(uuid.NewV4())
	newAddress := &m.Address{
		StreetAddress1:     "street 1",
		City:               "BEVERLY HILLS",
		State:              "CA",
		PostalCode:         "90210",
		County:             m.StringPointer("County"),
		UsPostRegionCityID: &usPostRegionCityID,
	}

	expErrors := map[string][]string{
		"us_post_region_city_id": {"UsPostRegionCityID is invalid."},
	}

	suite.verifyValidationErrors(newAddress, expErrors, suite.AppContextForTest())
}

func (suite *ModelSuite) TestEmptyAddressInstantiation() {

	var usprc models.UsPostRegionCity
	newAddress := m.Address{
		UsPostRegionCityID: &usprc.ID,
	}

	expErrors := map[string][]string{
		"street_address1":        {"StreetAddress1 can not be blank."},
		"city":                   {"City can not be blank."},
		"state":                  {"State can not be blank."},
		"postal_code":            {"PostalCode can not be blank."},
		"us_post_region_city_id": {"UsPostRegionCityID can not be blank."},
	}
	suite.verifyValidationErrors(&newAddress, expErrors, nil)
}

func (suite *ModelSuite) TestAddressCountryCode() {
	usprc, err := models.FindByZipCodeAndCity(suite.DB(), "90210", "BEVERLY HILLS")
	suite.NoError(err)

	noCountry := m.Address{
		StreetAddress1:     "street 1",
		StreetAddress2:     m.StringPointer("street 2"),
		StreetAddress3:     m.StringPointer("street 3"),
		City:               "BEVERLY HILLS",
		State:              "CA",
		PostalCode:         "90210",
		County:             m.StringPointer("county"),
		UsPostRegionCityID: &usprc.ID,
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
		City:           "BEVERLY HILLS",
		State:          "CA",
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
	usprc, err := m.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "90210", "BEVERLY HILLS")
	suite.NotNil(usprc)
	suite.FatalNoError(err)
	newAddress := &m.Address{
		StreetAddress1:     "street 1",
		StreetAddress2:     m.StringPointer("street 2"),
		StreetAddress3:     m.StringPointer("street 3"),
		City:               "BEVERLY HILLS",
		State:              "CA",
		PostalCode:         "90210",
		County:             m.StringPointer("County"),
		Country:            &country,
		CountryId:          &country.ID,
		UsPostRegionCityID: &usprc.ID,
		UsPostRegionCity:   usprc,
	}

	verrs, err := newAddress.Validate(suite.DB())

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")

	formattedAddress := newAddress.Format()

	suite.Equal("street 1\nstreet 2\nstreet 3\nBEVERLY HILLS, CA 90210", formattedAddress)

	formattedAddress = newAddress.LineFormat()

	suite.Equal("street 1, street 2, street 3, BEVERLY HILLS, CA, 90210, UNITED STATES", formattedAddress)

	formattedAddress = newAddress.LineDisplayFormat()

	suite.Equal("street 1 street 2 street 3, BEVERLY HILLS, CA 90210", formattedAddress)
}

func (suite *ModelSuite) TestPartialAddressFormat() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	usprc, err := m.FindByZipCodeAndCity(suite.AppContextForTest().DB(), "90210", "BEVERLY HILLS")
	suite.NotNil(usprc)
	suite.FatalNoError(err)
	newAddress := &m.Address{
		StreetAddress1:     "street 1",
		StreetAddress2:     nil,
		StreetAddress3:     nil,
		City:               "BEVERLY HILLS",
		State:              "CA",
		PostalCode:         "90210",
		County:             m.StringPointer("County"),
		Country:            &country,
		CountryId:          &country.ID,
		UsPostRegionCityID: &usprc.ID,
		UsPostRegionCity:   usprc,
	}

	verrs, err := newAddress.Validate(suite.DB())

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")

	formattedAddress := newAddress.Format()

	suite.Equal("street 1\nBEVERLY HILLS, CA 90210", formattedAddress)

	formattedAddress = newAddress.LineFormat()

	suite.Equal("street 1, BEVERLY HILLS, CA, 90210, UNITED STATES", formattedAddress)

	formattedAddress = newAddress.LineDisplayFormat()

	suite.Equal("street 1, BEVERLY HILLS, CA 90210", formattedAddress)
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
					City:               usprc.USPostRegionCityNm,
					PostalCode:         usprc.UsprZipID,
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
					PostalCode:         usprc.UsprZipID,
					City:               usprc.USPostRegionCityNm,
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
		City:           "BEVERLY HILLS",
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

func (suite *ModelSuite) TestValidateUSPRCAssignment() {

	suite.Run("returns false for invalid assignment", func() {
		incorrectUSPRCID := uuid.Must(uuid.NewV4())

		newAddress := &m.Address{
			StreetAddress1:     "street 1",
			StreetAddress2:     m.StringPointer("street 2"),
			StreetAddress3:     m.StringPointer("street 3"),
			City:               "BEVERLY HILLS",
			State:              "CA",
			PostalCode:         "90210",
			County:             m.StringPointer("County"),
			UsPostRegionCityID: &incorrectUSPRCID,
		}

		valid, err := m.ValidateUsPostRegionCityID(suite.DB(), *newAddress)
		suite.NoError(err)
		suite.Equal(false, valid)
	})

	suite.Run("returns true for valid assignment", func() {

		newAddress := &m.Address{
			StreetAddress1: "street 1",
			StreetAddress2: m.StringPointer("street 2"),
			StreetAddress3: m.StringPointer("street 3"),
			City:           "BEVERLY HILLS",
			State:          "CA",
			PostalCode:     "90210",
			County:         m.StringPointer("County"),
		}

		expectedUSPRC, err := m.FindByZipCodeAndCity(suite.DB(), newAddress.PostalCode, newAddress.City)
		suite.NoError(err)

		newAddress.UsPostRegionCityID = &expectedUSPRC.ID

		valid, err := m.ValidateUsPostRegionCityID(suite.DB(), *newAddress)
		suite.NoError(err)
		suite.Equal(true, valid)
	})

	suite.Run("returns error when fails to lookup USPRC", func() {

		uuid := uuid.Must(uuid.NewV4())
		newAddress := &m.Address{
			StreetAddress1:     "street 1",
			StreetAddress2:     m.StringPointer("street 2"),
			StreetAddress3:     m.StringPointer("street 3"),
			City:               "BEVERLY HILLS",
			State:              "CA",
			PostalCode:         "29229",
			County:             m.StringPointer("County"),
			UsPostRegionCityID: &uuid,
		}

		valid, err := m.ValidateUsPostRegionCityID(suite.DB(), *newAddress)
		suite.Error(err, "No UsPostRegionCity found for provided zip code 29229 and city BEVERLY HILLS.")
		suite.Equal(false, valid)
	})
}

func (suite *ModelSuite) TestValidPostalCode() {

	suite.Run("returns true or false if a postal code is valid or not", func() {

		testCases := []struct {
			name        string
			input       string
			expected    bool
			expectedErr bool
		}{
			{"5 digit postal code", "90201", true, false},
			{"5 digit postal code not in the UsPostRegionCity table", "76334", false, false},
			{"Not a 5 digit postal code", "33", false, false},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				isValid, err := m.ValidPostalCode(suite.DB(), tc.input)
				if tc.expectedErr {
					suite.Error(err)
				} else {
					suite.NoError(err)
				}

				suite.Equal(tc.expected, isValid)
			})
		}
	})
}
