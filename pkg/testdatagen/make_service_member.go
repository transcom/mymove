package testdatagen

import (
	"math/rand"
	"strconv"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// randomEdipi creates a random Edipi for a service member
func randomEdipi() string {
	low := 1000000000
	high := 9999999999
	// #nosec G404 TODO needs review
	return strconv.Itoa(low + rand.Intn(high-low))
}

// MakeServiceMember creates a single ServiceMember with associated data.
func MakeServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	aServiceMember := assertions.ServiceMember
	user := aServiceMember.User
	agency := aServiceMember.Affiliation
	email := "leo_spaceman_sm@example.com"
	currentAddressID := aServiceMember.ResidentialAddressID
	currentAddress := aServiceMember.ResidentialAddress

	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.ServiceMember.UserID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		user = MakeDefaultUser(db)
	}
	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	if agency == nil {
		army := models.AffiliationARMY
		agency = &army
	}

	if currentAddressID == nil || isZeroUUID(*currentAddressID) {
		newAddress := MakeDefaultAddress(db)
		currentAddressID = &newAddress.ID
		currentAddress = &newAddress
	}

	serviceMember := models.ServiceMember{
		UserID:               user.ID,
		User:                 user,
		Edipi:                swag.String(randomEdipi()),
		Affiliation:          agency,
		FirstName:            swag.String("Leo"),
		LastName:             swag.String("Spacemen"),
		Telephone:            swag.String("212-123-4567"),
		PersonalEmail:        &email,
		ResidentialAddressID: currentAddressID,
		ResidentialAddress:   currentAddress,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceMember, assertions.ServiceMember)

	mustCreate(db, &serviceMember, assertions.Stub)

	return serviceMember
}

// MakeDefaultServiceMember returns a service member with default options
func MakeDefaultServiceMember(db *pop.Connection) models.ServiceMember {
	return MakeServiceMember(db, Assertions{})
}

// MakeExtendedServiceMember creates a single ServiceMember and associated User, Addresses,
// and Backup Contact.
func MakeExtendedServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	affiliation := assertions.ServiceMember.Affiliation
	if affiliation == nil {
		army := models.AffiliationARMY
		affiliation = &army
	}
	residentialAddress := MakeDefaultAddress(db)
	backupMailingAddress := MakeDefaultAddress(db)
	e1 := models.ServiceMemberRankE1
	station := FetchOrMakeDefaultCurrentDutyStation(db)

	// Combine extended SM defaults with assertions
	smDefaults := models.ServiceMember{
		Edipi:                  swag.String(randomEdipi()),
		Rank:                   &e1,
		Affiliation:            affiliation,
		ResidentialAddressID:   &residentialAddress.ID,
		BackupMailingAddressID: &backupMailingAddress.ID,
		DutyStationID:          &station.ID,
		DutyStation:            station,
		EmailIsPreferred:       swag.Bool(true),
		Telephone:              swag.String("555-555-5555"),
	}

	mergeModels(&smDefaults, assertions.ServiceMember)

	assertions.ServiceMember = smDefaults

	serviceMember := MakeServiceMember(db, assertions)

	contactAssertions := Assertions{
		BackupContact: models.BackupContact{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
	}
	backupContact := MakeBackupContact(db, contactAssertions)
	serviceMember.BackupContacts = append(serviceMember.BackupContacts, backupContact)
	mustSave(db, &serviceMember)

	return serviceMember
}
