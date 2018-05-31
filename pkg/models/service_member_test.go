package models_test

import (
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
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
	affiliation := internalmessages.AffiliationARMY
	rank := internalmessages.ServiceMemberRankE5
	firstName := "bob"
	lastName := "sally"
	telephone := "510 555-5555"
	email := "bobsally@gmail.com"
	fakeAddress, _ := testdatagen.MakeAddress(suite.db)
	fakeBackupAddress, _ := testdatagen.MakeAddress(suite.db)
	fakeID := uuid.Must(uuid.NewV4())

	servicemember := ServiceMember{
		UserID:                 user1.ID,
		Edipi:                  &edipi,
		Affiliation:            &affiliation,
		Rank:                   &rank,
		FirstName:              &firstName,
		LastName:               &lastName,
		Telephone:              &telephone,
		PersonalEmail:          &email,
		ResidentialAddress:     &fakeAddress,
		BackupMailingAddress:   &fakeBackupAddress,
		SocialSecurityNumberID: &fakeID,
		DutyStationID:          &fakeID,
	}

	// Then: IsProfileComplete should return false
	if servicemember.IsProfileComplete() != false {
		t.Error("Expected profile to be incomplete.")
	}
	// When: all required fields are set
	emailPreferred := true
	servicemember.EmailIsPreferred = &emailPreferred
	fakeBackupContact, _ := testdatagen.MakeBackupContact(suite.db)

	var backupContacts BackupContacts
	backupContacts = append(backupContacts, fakeBackupContact)
	servicemember.BackupContacts = backupContacts
	// Then: IsProfileComplete should return true
	if servicemember.IsProfileComplete() != true {
		t.Error("Expected profile to be complete.")
	}
}

func (suite *ModelSuite) TestFetchServiceMember() {
	user1, _ := testdatagen.MakeUser(suite.db)
	user2, _ := testdatagen.MakeUser(suite.db)

	firstName := "Oliver"
	resAddress, _ := testdatagen.MakeAddress(suite.db)
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.mustSave(&sm)

	// User is authorized to fetch service member
	session := &auth.Session{
		ApplicationName: auth.MyApp,
		UserID:          user1.ID,
		ServiceMemberID: sm.ID,
	}
	goodSm, err := FetchServiceMember(suite.db, session, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddress.ID, goodSm.ResidentialAddress.ID)
	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = FetchServiceMember(suite.db, session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}

	// User is forbidden from fetching order
	session.UserID = user2.ID
	session.ServiceMemberID = uuid.Nil
	_, err = FetchServiceMember(suite.db, session, sm.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}

}
