package models_test

import (
	"github.com/gofrs/uuid"

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

	usaCountry := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		CountryId:      m.UUIDPointer(uuid.Must(uuid.NewV4())),
	}
	countryCode, err = usaCountry.CountryCode()
	suite.NoError(err)
	suite.Equal("USA", *countryCode)

	usCountry := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		CountryId:      m.UUIDPointer(uuid.Must(uuid.NewV4())),
		County:         "county",
	}
	countryCode, err = usCountry.CountryCode()
	suite.NoError(err)
	suite.Equal("USA", *countryCode)

	notUsaCountry := m.Address{
		StreetAddress1: "street 1",
		StreetAddress2: m.StringPointer("street 2"),
		StreetAddress3: m.StringPointer("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
		County:         "county",
		CountryId:      m.UUIDPointer(uuid.Must(uuid.NewV4())),
	}

	countryCode, err = notUsaCountry.CountryCode()
	suite.Nil(countryCode)
	suite.Error(err)
	suite.Equal("NotImplementedCountryCode: Country 'Ireland'", err.Error())

}
