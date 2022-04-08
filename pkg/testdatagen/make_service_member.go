package testdatagen

import (
	"strconv"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/jaswdr/faker"

	"github.com/transcom/mymove/pkg/models"
)

// randomEdipi creates a random Edipi for a service member
func randomEdipi() string {
	fake := faker.New()

	low := 1000000000
	high := 9999999999

	return strconv.Itoa(fake.IntBetween(low, high))
}

// randomServiceMemberAffiliation returns a random service member affiliation
func randomServiceMemberAffiliation() models.ServiceMemberAffiliation {
	fake := faker.New()

	affiliations := []string{
		models.AffiliationARMY.String(),
		models.AffiliationNAVY.String(),
		models.AffiliationMARINES.String(),
		models.AffiliationAIRFORCE.String(),
		models.AffiliationCOASTGUARD.String(),
	}

	affiliation := fake.RandomStringElement(affiliations)

	return models.ServiceMemberAffiliation(affiliation)
}

func randomServiceMemberRank() models.ServiceMemberRank {
	fake := faker.New()

	ranks := []string{
		models.ServiceMemberRankE1.String(),
		models.ServiceMemberRankE2.String(),
		models.ServiceMemberRankE3.String(),
		models.ServiceMemberRankE4.String(),
		models.ServiceMemberRankE5.String(),
		models.ServiceMemberRankE6.String(),
		models.ServiceMemberRankE7.String(),
		models.ServiceMemberRankE8.String(),
		models.ServiceMemberRankE9.String(),
		models.ServiceMemberRankE9SPECIALSENIORENLISTED.String(),
		models.ServiceMemberRankO1ACADEMYGRADUATE.String(),
		models.ServiceMemberRankO2.String(),
		models.ServiceMemberRankO3.String(),
		models.ServiceMemberRankO4.String(),
		models.ServiceMemberRankO5.String(),
		models.ServiceMemberRankO6.String(),
		models.ServiceMemberRankO7.String(),
		models.ServiceMemberRankO8.String(),
		models.ServiceMemberRankO9.String(),
		models.ServiceMemberRankO10.String(),
		models.ServiceMemberRankW1.String(),
		models.ServiceMemberRankW2.String(),
		models.ServiceMemberRankW3.String(),
		models.ServiceMemberRankW4.String(),
		models.ServiceMemberRankW5.String(),
		models.ServiceMemberRankAVIATIONCADET.String(),
		models.ServiceMemberRankCIVILIANEMPLOYEE.String(),
		models.ServiceMemberRankACADEMYCADET.String(),
		models.ServiceMemberRankMIDSHIPMAN.String(),
	}

	rank := fake.RandomStringElement(ranks)

	return models.ServiceMemberRank(rank)
}

// MakeServiceMember creates a single ServiceMember
// If not provided, it will also create an associated
// - User
// - ResidentialAddress
func MakeServiceMember(db *pop.Connection, assertions Assertions) models.ServiceMember {
	fake := faker.New()

	user := assertions.ServiceMember.User
	email := fake.Internet().Email()

	// ID is required because it must be populated for Eager saving to work.
	if assertions.ServiceMember.UserID.IsNil() { // TODO: This should be checking assertions.User.ID, but I won't change it right now because it would require changing any place this is called incorrectly.
		if assertions.User.LoginGovEmail == "" {
			assertions.User.LoginGovEmail = email
		}
		user = MakeDefaultUser(db)
	}
	if assertions.User.LoginGovEmail != "" {
		email = assertions.User.LoginGovEmail
	}

	serviceMember := models.ServiceMember{
		UserID:        user.ID,
		User:          user,
		Edipi:         models.StringPointer(randomEdipi()),
		Affiliation:   randomServiceMemberAffiliation().Pointer(),
		FirstName:     models.StringPointer(fake.Person().FirstName()),
		LastName:      models.StringPointer(fake.Person().LastName()),
		Telephone:     models.StringPointer(fake.Phone().Number()),
		PersonalEmail: &email,
		Rank:          randomServiceMemberRank().Pointer(),
	}

	if assertions.ServiceMember.ResidentialAddressID == nil || assertions.ServiceMember.ResidentialAddressID.IsNil() {
		newAddress := MakeDefaultAddress(db) // TODO: This isn't passing along assertions, so Stub won't work, but it's also doing other things in there so it'll take more to refactor.

		serviceMember.ResidentialAddressID = &newAddress.ID
		serviceMember.ResidentialAddress = &newAddress
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
	fake := faker.New()

	backupMailingAddress := MakeAddress2(db, assertions)
	dutyLocation := FetchOrMakeDefaultCurrentDutyLocation(db)

	smDefaults := models.ServiceMember{
		BackupMailingAddressID: &backupMailingAddress.ID,
		DutyLocationID:         &dutyLocation.ID,
		DutyLocation:           dutyLocation,
		EmailIsPreferred:       swag.Bool(true),
		Telephone:              models.StringPointer(fake.Phone().Number()),
	}

	mergeModels(&smDefaults, assertions.ServiceMember)

	assertions.ServiceMember = smDefaults

	serviceMember := MakeServiceMember(db, assertions)

	contactAssertions := Assertions{
		BackupContact: models.BackupContact{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
		Stub: assertions.Stub,
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
