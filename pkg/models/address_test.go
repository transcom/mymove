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
		County:         "County",
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
		"county":          {"County can not be blank."},
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
		County:         "county",
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

func (suite *ModelSuite) TestAddressFormat() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)
	newAddress := &m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		County:         "County",
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
}
