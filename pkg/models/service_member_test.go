package models_test

import (
	"github.com/gobuffalo/uuid"
	"github.com/transcom/mymove/pkg/server"

	"github.com/transcom/mymove/pkg/auth"
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

func (suite *ModelSuite) TestFetchServiceMemberForUser() {
	user1 := testdatagen.MakeDefaultUser(suite.db)
	user2 := testdatagen.MakeDefaultUser(suite.db)

	firstName := "Oliver"
	resAddress := testdatagen.MakeDefaultAddress(suite.db)
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.mustSave(&sm)

	// User is authorized to fetch service member
	session := &server.Session{
		ApplicationName: auth.MyApp,
		UserID:          user1.ID,
		ServiceMemberID: sm.ID,
	}
	goodSm, err := FetchServiceMemberForUser(suite.db, session, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddress.ID, goodSm.ResidentialAddress.ID)
	}

	// Wrong ServiceMember
	wrongID, _ := uuid.NewV4()
	_, err = FetchServiceMemberForUser(suite.db, session, wrongID)
	if suite.Error(err) {
		suite.Equal(ErrFetchNotFound, err)
	}

	// User is forbidden from fetching order
	session.UserID = user2.ID
	session.ServiceMemberID = uuid.Nil
	_, err = FetchServiceMemberForUser(suite.db, session, sm.ID)
	if suite.Error(err) {
		suite.Equal(ErrFetchForbidden, err)
	}
}

func (suite *ModelSuite) TestFetchServiceMemberNotForUser() {
	user1 := testdatagen.MakeDefaultUser(suite.db)

	firstName := "Nino"
	resAddress := testdatagen.MakeDefaultAddress(suite.db)
	sm := ServiceMember{
		User:                 user1,
		UserID:               user1.ID,
		FirstName:            &firstName,
		ResidentialAddressID: &resAddress.ID,
		ResidentialAddress:   &resAddress,
	}
	suite.mustSave(&sm)

	goodSm, err := FetchServiceMember(suite.db, sm.ID)
	if suite.NoError(err) {
		suite.Equal(sm.FirstName, goodSm.FirstName)
		suite.Equal(sm.ResidentialAddressID, goodSm.ResidentialAddressID)
	}
}
