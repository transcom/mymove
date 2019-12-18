package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeCustomer creates a single Customer
func MakeCustomer(db *pop.Connection, assertions Assertions) models.Customer {
	user := assertions.User
	firstName := assertions.Customer.FirstName
	lastName := assertions.Customer.LastName

	if firstName == "" {
		firstName = "Bob"
	}

	if lastName == "" {
		lastName = "Vance"
	}

	if isZeroUUID(user.ID) {
		user = MakeUser(db, assertions)
	}
	customer := models.Customer{
		FirstName: firstName,
		LastName:  lastName,
		User:      user,
		UserID:    user.ID,
		DODID:     randomEdipi(),
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
