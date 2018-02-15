package models

import (
	"github.com/go-openapi/swag"
	"testing"
)

func TestBasicAddressInstantiation(t *testing.T) {
	newAddress := Address{
		StreetAddress1: "street 1",
		StreetAddress2: swag.String("street 2"),
		City:           "city",
		State:          "state",
		Zip:            "90210",
	}

	verrs, err := dbConnection.ValidateAndCreate(&newAddress)

	if err != nil {
		t.Fatal("Error writing to the db.", err)
	}

	if verrs.HasAny() {
		t.Fatal("Error validating model")
	}
}
