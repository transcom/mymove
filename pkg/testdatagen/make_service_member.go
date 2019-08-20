package testdatagen

import (
	"math/rand"
	"strconv"

	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// randomEdipi creates a random Edipi for a service member
func randomEdipi() string {
	low := 1000000000
	high := 9999999999
	return strconv.Itoa(low + rand.Intn(high-low))
}

// MakeServiceMember creates a single ServiceMember with associated data.
func MakeServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	user := assertions.ServiceMember.User
	email := "leo_spaceman_sm@example.com"

	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.ServiceMember.UserID) {
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		user = MakeUser(db, assertions)
	}
	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	serviceMember := models.ServiceMember{
		UserID:        user.ID,
		User:          user,
		Edipi:         models.StringPointer(randomEdipi()),
		FirstName:     models.StringPointer("Leo"),
		LastName:      models.StringPointer("Spacemen"),
		PersonalEmail: models.StringPointer(email),
	}

	// Overwrite values with those from assertions
	mergeModels(&serviceMember, assertions.ServiceMember)

	mustCreate(db, &serviceMember)

	return serviceMember
}

// MakeDefaultServiceMember returns a service member with default options
func MakeDefaultServiceMember(db *pop.Connection) models.ServiceMember {
	return MakeServiceMember(db, Assertions{})
}

// MakeExtendedServiceMember creates a single ServiceMember and associated User, Addresses,
// and Backup Contact.
func MakeExtendedServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	residentialAddress := MakeDefaultAddress(db)
	backupMailingAddress := MakeDefaultAddress(db)
	Army := models.AffiliationARMY
	E1 := models.ServiceMemberRankE1

	station := FetchOrMakeDefaultCurrentDutyStation(db)
	emailPreferred := true
	// Combine extended SM defaults with assertions
	smDefaults := models.ServiceMember{
		Edipi:                  models.StringPointer(randomEdipi()),
		Rank:                   &E1,
		Affiliation:            &Army,
		ResidentialAddressID:   &residentialAddress.ID,
		BackupMailingAddressID: &backupMailingAddress.ID,
		DutyStationID:          &station.ID,
		DutyStation:            station,
		EmailIsPreferred:       &emailPreferred,
		Telephone:              models.StringPointer("555-555-5555"),
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

// MakeServiceMemberData created 5 ServiceMembers (and in turn a User for each)
func MakeServiceMemberData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeDefaultServiceMember(db)
	}
}
