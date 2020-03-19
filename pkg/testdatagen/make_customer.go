package testdatagen

import (
	"fmt"
	"math/rand"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeFirstName chooses a random first name of a star wars character
func MakeFirstName() *string {
	firstNameArr := [10]string{"Darth", "Obi-Wan", "Luke", "Princess", "Han", "Baby", "Chew", "Jango", "R2", "Lando"}
	return &firstNameArr[rand.Intn(len(firstNameArr))]
}

// MakeLastName chooses a random last name of a star wars character
func MakeLastName() *string {
	lastNameArr := [10]string{"Vader", "Kenobi", "Skywalker", "Leia", "Solo", "Yoda", "Bacca", "Fett", "D2", "Calrissian"}
	return &lastNameArr[rand.Intn(len(lastNameArr))]
}

// MakeAgency chooses a random agency
func MakeAgency() *string {
	agencies := [5]string{
		"ARMY",
		"NAVY",
		"MARINES",
		"AIR_FORCE",
		"COAST_GUARD",
	}

	return &agencies[rand.Intn(len(agencies))]
}

// MakeCustomer creates a single Customer
func MakeCustomer(db *pop.Connection, assertions Assertions) models.Customer {
	user := assertions.User
	aCustomer := assertions.Customer
	firstName := aCustomer.FirstName
	lastName := aCustomer.LastName
	agency := aCustomer.Agency
	currentAddressID := aCustomer.CurrentAddressID
	currentAddress := aCustomer.CurrentAddress
	destinationAddressID := aCustomer.DestinationAddressID
	destinationAddress := aCustomer.DestinationAddress
	email := aCustomer.Email
	phoneNumber := aCustomer.PhoneNumber

	if firstName == nil {
		firstName = MakeFirstName()
	}

	if lastName == nil {
		lastName = MakeLastName()
	}

	if agency == nil {
		agency = MakeAgency()
	}
	if email == nil || *email == "" {
		e := fmt.Sprintf("%s%s@mail.com", *firstName, *lastName)
		email = &e
	}
	if phoneNumber == nil || *phoneNumber == "" {
		p := "212-123-456"
		phoneNumber = &p
	}

	if isZeroUUID(user.ID) {
		user = MakeUser(db, assertions)
	}
	if currentAddressID == nil || isZeroUUID(*currentAddressID) {
		currentAddress = MakeAddress(db, Assertions{})
	}
	if destinationAddressID == nil || isZeroUUID(*destinationAddressID) {
		destinationAddress = MakeAddress2(db, Assertions{})
	}
	customer := models.Customer{
		Agency:               agency,
		CurrentAddress:       currentAddress,
		CurrentAddressID:     &currentAddress.ID,
		DODID:                swag.String(randomEdipi()),
		DestinationAddress:   destinationAddress,
		DestinationAddressID: &destinationAddress.ID,
		Email:                email,
		FirstName:            firstName,
		LastName:             lastName,
		PhoneNumber:          phoneNumber,
		User:                 user,
		UserID:               user.ID,
	}

	// Overwrite values with those from assertions
	mergeModels(&customer, aCustomer)
	mustCreate(db, &customer)

	return customer
}

// MakeDefaultCustomer makes a Customer with default values
func MakeDefaultCustomer(db *pop.Connection) models.Customer {
	return MakeCustomer(db, Assertions{})
}
