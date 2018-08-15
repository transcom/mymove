package testdatagen

import (
	"github.com/gobuffalo/pop"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

// MakeServiceMember creates a single ServiceMember with associated data.
func MakeServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	user := assertions.ServiceMember.User
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.ServiceMember.UserID) {
		user = MakeUser(db, assertions)
	}

	serviceMember := models.ServiceMember{
		UserID:        user.ID,
		User:          user,
		FirstName:     models.StringPointer("Leo"),
		LastName:      models.StringPointer("Spacemen"),
		PersonalEmail: models.StringPointer("leo@example.com"),
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
	Army := internalmessages.AffiliationARMY
	E1 := internalmessages.ServiceMemberRankE1

	station := MakeDefaultDutyStation(db)
	emailPreferred := true
	// Combine extended SM defaults with assertions
	smDefaults := models.ServiceMember{
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

	serviceMemberAssertions := Assertions{
		ServiceMember: smDefaults,
	}

	serviceMember := MakeServiceMember(db, serviceMemberAssertions)

	contactAssertions := Assertions{
		BackupContact: models.BackupContact{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
	}
	MakeBackupContact(db, contactAssertions)

	return serviceMember
}

// MakeServiceMemberData created 5 ServiceMembers (and in turn a User for each)
func MakeServiceMemberData(db *pop.Connection) {
	for i := 0; i < 5; i++ {
		MakeDefaultServiceMember(db)
	}
}
