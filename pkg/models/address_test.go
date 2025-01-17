package models_test

import (
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
