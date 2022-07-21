package testdatagen

import (
	"log"
	"strconv"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/random"
)

// RandomEdipi creates a random Edipi for a service member
func RandomEdipi() string {
	low := 1000000000
	high := 9999999999
	randInt, err := random.GetRandomIntAddend(low, high)
	if err != nil {
		log.Panicf("Failure to generate random Edipi %v", err)
	}
	return strconv.Itoa(low + int(randInt))
}

// MakeServiceMember creates a single ServiceMember
// If not provided, it will also create an associated
// - User
// - ResidentialAddress
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

	randomEdipi := RandomEdipi()
	rank := models.ServiceMemberRankE1

	serviceMember := models.ServiceMember{
		UserID:               user.ID,
		User:                 user,
		Edipi:                swag.String(randomEdipi),
		Affiliation:          agency,
		FirstName:            swag.String("Leo"),
		LastName:             swag.String("Spacemen"),
		Telephone:            swag.String("212-123-4567"),
		PersonalEmail:        &email,
		ResidentialAddressID: currentAddressID,
		ResidentialAddress:   currentAddress,
		Rank:                 &rank,
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceMember, assertions.ServiceMember)

	mustCreate(db, &serviceMember, assertions.Stub)

	return serviceMember
}

// MakeDefaultServiceMember returns a service member with default options
// It will also create an associated
//   - User
//   - ResidentialAddress
func MakeDefaultServiceMember(db *pop.Connection) models.ServiceMember {
	return MakeServiceMember(db, Assertions{})
}

// MakeExtendedServiceMember creates a single ServiceMember
// If not provided it will also create an associated
//   - User,
//   - ResidentialAddress
//   - BackupMailingAddress
//   - DutyLocation
//   - BackupContact
func MakeExtendedServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	affiliation := assertions.ServiceMember.Affiliation
	if affiliation == nil {
		army := models.AffiliationARMY
		affiliation = &army
	}
	residentialAddress := MakeDefaultAddress(db)
	backupMailingAddress := MakeAddress2(db, assertions)
	e1 := models.ServiceMemberRankE1
	dutyLocation := FetchOrMakeDefaultCurrentDutyLocation(db)

	// Combine extended SM defaults with assertions
	smDefaults := models.ServiceMember{
		Edipi:                  swag.String(RandomEdipi()),
		Rank:                   &e1,
		Affiliation:            affiliation,
		ResidentialAddressID:   &residentialAddress.ID,
		BackupMailingAddressID: &backupMailingAddress.ID,
		DutyLocationID:         &dutyLocation.ID,
		DutyLocation:           dutyLocation,
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
	if !assertions.Stub {
		MustSave(db, &serviceMember)
	}

	return serviceMember
}

// MakeStubbedServiceMember returns a stubbed service member that is not stored in the DB
func MakeStubbedServiceMember(db *pop.Connection) models.ServiceMember {
	user := MakeStubbedUser(db)

	return MakeServiceMember(db, Assertions{
		ServiceMember: models.ServiceMember{
			ID:     uuid.Must(uuid.NewV4()),
			User:   user,
			UserID: user.ID,
		},
		Stub: true,
	})
}
