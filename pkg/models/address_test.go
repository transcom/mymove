package models_test

import (
	"github.com/go-openapi/swag"

	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicAddressInstantiation() {
	newAddress := Address{
		StreetAddress1: "street 1",
		StreetAddress2: swag.String("street 2"),
		StreetAddress3: swag.String("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
	}

	verrs, err := suite.DB().ValidateAndCreate(&newAddress)

	suite.Nil(err, "Error writing to the db.")
	suite.False(verrs.HasAny(), "Error validating model")
}

func (suite *ModelSuite) TestEmptyAddressInstantiation() {
	newAddress := Address{}

	expErrors := map[string][]string{
		"street_address1": {"StreetAddress1 can not be blank."},
		"city":            {"City can not be blank."},
		"state":           {"State can not be blank."},
		"postal_code":     {"PostalCode can not be blank."},
	}
	suite.verifyValidationErrors(&newAddress, expErrors)
}
