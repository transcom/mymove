package models_test

import (
	"github.com/go-openapi/swag"
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicAddressInstantiation() {
	t := suite.T()

	newAddress := Address{
		StreetAddress1: "street 1",
		StreetAddress2: swag.String("street 2"),
		StreetAddress3: swag.String("street 3"),
		City:           "city",
		State:          "state",
		PostalCode:     "90210",
	}

	verrs, err := suite.db.ValidateAndCreate(&newAddress)

	if err != nil {
		t.Fatal("Error writing to the db.", err)
	}

	if verrs.HasAny() {
		t.Fatal("Error validating model")
	}
}
