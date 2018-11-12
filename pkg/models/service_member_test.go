package models_test

import (
	"github.com/gofrs/uuid"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestBasicServiceMemberInstantiation() {
	servicemember := &ServiceMember{}

	expErrors := map[string][]string{
		"user_id": {"UserID can not be blank."},
	}

	suite.verifyValidationErrors(servicemember, expErrors)
}

func (suite *ModelSuite) TestIsProfileCompleteWithIncompleteSM() {
	t := suite.T()
	// Given: a user and a service member
	user1 := User{
		LoginGovUUID:  uuid.Must(uuid.NewV4()),
		LoginGovEmail: "whoever@example.com",
	}
	verrs, err := suite.db.ValidateAndCreate(&user1)
	if verrs.HasAny() || err != nil {
		t.Error(verrs, err)
	}

	// And: a service member is incompletely initialized with almost all required values
	edipi := "12345567890"
	affiliation := AffiliationARMY
	rank := ServiceMemberRankE5
	firstName := "bob"
	lastName := "sally"
	telephone := "510 555-5555"
	email := "bobsally@gmail.com"
	fakeAddress := testdatagen.MakeDefaultAddress(suite.db)
	fakeBackupAddress := testdatagen.MakeDefaultAddress(suite.db)
	station := testdatagen.MakeDefaultDutyStation(suite.db)

	serviceMember := ServiceMember{
		UserID:                 user1.ID,
		Edipi:                  &edipi,
		Affiliation:            &affiliation,
		Rank:                   &rank,
		FirstName:              &firstName,
		LastName:               &lastName,
		Telephone:              &telephone,
		PersonalEmail:          &email,
		ResidentialAddressID:   &fakeAddress.ID,
		BackupMailingAddressID: &fakeBackupAddress.ID,
		DutyStationID:          &station.ID,
	}

	// Then: IsProfileComplete should return false
	if serviceMember.IsProfileComplete() != false {
		t.Error("Expected profile to be incomplete.")
	}
	// When: all required fields are set
	emailPreferred := true
	serviceMember.EmailIsPreferred = &emailPreferred

	newSsn := SocialSecurityNumber{}
	newSsn.SetEncryptedHash("555-55-5555")
	suite.mustSave(&newSsn)
	serviceMember.SocialSecurityNumber = &newSsn
	serviceMember.SocialSecurityNumberID = &newSsn.ID

	suite.mustSave(&serviceMember)

	contactAssertions := testdatagen.Assertions{
		BackupContact: BackupContact{
			ServiceMember:   serviceMember,
			ServiceMemberID: serviceMember.ID,
		},
	}
	testdatagen.MakeBackupContact(suite.db, contactAssertions)

	if err = suite.db.Load(&serviceMember); err != nil {
		t.Errorf("Could not load BackupContacts for serviceMember: %v", err)
	}

	// Then: IsProfileComplete should return true
	if serviceMember.IsProfileComplete() != true {
		t.Error("Expected profile to be complete.")
	}
}
