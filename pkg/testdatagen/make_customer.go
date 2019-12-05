package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeCustomer creates a single Customer
func MakeCustomer(db *pop.Connection, assertions Assertions) models.Customer {

	customer := models.Customer{
		FirstName: "Carmen",
		LastName:  "Cruz",
	}

	// Overwrite values with those from assertions
	mergeModels(&customer, assertions.Customer)

	mustCreate(db, &customer)

	return customer
}

// MakeDefaultCustomer makes a Customer with default values
func MakeDefaultCustomer(db *pop.Connection) models.Customer {
	return MakeCustomer(db, Assertions{})
}
