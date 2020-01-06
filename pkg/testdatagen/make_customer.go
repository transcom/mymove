package testdatagen

import (
	"math/rand"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeFirstName chooses a random first name of a star wars character
func MakeFirstName() string {
	firstNameArr := [10]string{"Darth", "Obi-Wan", "Luke", "Princess", "Han", "Baby", "Chew", "Jango", "R2", "Lando"}
	return firstNameArr[rand.Intn(len(firstNameArr))]
}

// MakeLastName chooses a random last name of a star wars character
func MakeLastName() string {
	lastNameArr := [10]string{"Vader", "Kenobi", "Skywalker", "Leia", "Solo", "Yoda", "Bacca", "Fett", "D2", "Calrissian"}
	return lastNameArr[rand.Intn(len(lastNameArr))]
}

// MakeAgency chooses a random agency
func MakeAgency() string {
	agencies := [5]string{"ARMY",
		"NAVY",
		"MARINES",
		"AIR_FORCE",
		"COAST_GUARD"}

	return agencies[rand.Intn(len(agencies))]
}

// MakeCustomer creates a single Customer
func MakeCustomer(db *pop.Connection, assertions Assertions) models.Customer {
	user := assertions.User
	firstName := assertions.Customer.FirstName
	lastName := assertions.Customer.LastName
	agency := assertions.Customer.Agency
	if firstName == "" {
		firstName = MakeFirstName()
	}

	if lastName == "" {
		lastName = MakeLastName()
	}

	if agency == "" {
		agency = MakeAgency()
	}

	if isZeroUUID(user.ID) {
		user = MakeUser(db, assertions)
	}
	customer := models.Customer{
		Agency:    agency,
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
