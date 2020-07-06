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
func MakeAgency() *models.ServiceMemberAffiliation {
	agencies := [5]models.ServiceMemberAffiliation{
		models.AffiliationARMY,
		models.AffiliationNAVY,
		models.AffiliationMARINES,
		models.AffiliationAIRFORCE,
		models.AffiliationCOASTGUARD,
	}

	return &agencies[rand.Intn(len(agencies))]
}

// MakeCustomer creates a single Customer
func MakeCustomer(db *pop.Connection, assertions Assertions) models.ServiceMember {
	user := assertions.User
	aCustomer := assertions.Customer
	firstName := aCustomer.FirstName
	lastName := aCustomer.LastName
	agency := aCustomer.Affiliation
	currentAddressID := aCustomer.ResidentialAddressID
	currentAddress := aCustomer.ResidentialAddress
	email := aCustomer.PersonalEmail
	phoneNumber := aCustomer.Telephone

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
		newAddress := MakeAddress(db, Assertions{})
		currentAddress = &newAddress
	}
	customer := models.ServiceMember{
		Affiliation:          agency,
		ResidentialAddress:   currentAddress,
		ResidentialAddressID: &currentAddress.ID,
		Edipi:                swag.String(randomEdipi()),
		PersonalEmail:        email,
		FirstName:            firstName,
		LastName:             lastName,
		Telephone:            phoneNumber,
		User:                 user,
		UserID:               user.ID,
	}

	// Overwrite values with those from assertions
	mergeModels(&customer, aCustomer)
	mustCreate(db, &customer)

	return customer
}

// MakeDefaultCustomer makes a Customer with default values
func MakeDefaultCustomer(db *pop.Connection) models.ServiceMember {
	return MakeCustomer(db, Assertions{})
}
