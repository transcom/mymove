package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestUserCreation() {
	t := suite.T()

	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"

	newUser := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}

	if verrs, err := suite.DB().ValidateAndCreate(&newUser); err != nil || verrs.HasAny() {
		t.Fatal("Didn't create user in db.")
	}

	if newUser.ID == uuid.Nil {
		t.Error("Didn't get an id back for user.")
	}

	if (newUser.LoginGovEmail != userEmail) &&
		(newUser.LoginGovUUID != fakeUUID) {
		t.Error("Required values didn't get set.")
	}
}

func (suite *ModelSuite) TestUserCreationWithoutValues() {
	newUser := &User{}

	expErrors := map[string][]string{
		"login_gov_email": {"LoginGovEmail can not be blank."},
		"login_gov_uuid":  {"LoginGovUUID can not be blank."},
	}

	suite.verifyValidationErrors(newUser, expErrors)
}

func (suite *ModelSuite) TestUserCreationDuplicateUUID() {
	t := suite.T()

	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	userEmail := "sally@government.gov"

	newUser := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}

	sameUser := User{
		LoginGovUUID:  fakeUUID,
		LoginGovEmail: userEmail,
	}

	suite.DB().Create(&newUser)
	err := suite.DB().Create(&sameUser)

	if err.Error() != `pq: duplicate key value violates unique constraint "constraint_name"` {
		t.Fatal("Db should have errored on unique constraint for UUID")
	}
}

func (suite *ModelSuite) TestCreateUser() {
	const testEmail = "Sally@GoVernment.gov"
	const expectedEmail = "sally@government.gov"
	const goodUUID = "39b28c92-0506-4bef-8b57-e39519f42dc2"
	const badUUID = "39xnfc92-0506-4bef-8b57-e39519f42dc2"

	sally, err := CreateUser(suite.DB(), goodUUID, testEmail)
	suite.Nil(err, "No error for good create")
	suite.Equal(expectedEmail, sally.LoginGovEmail, "should convert email to lower case")
	suite.NotEqual(sally.ID, uuid.Nil)

	fail, err := CreateUser(suite.DB(), expectedEmail, badUUID)
	suite.NotNil(err, "should get and error from bad uuid")
	suite.Nil(fail, "no user with bad uuid")
}

func (suite *ModelSuite) TestFetchUserIdentity() {
	const goodUUID = "39b28c92-0506-4bef-8b57-e39519f42dc2"

	// First check that it all works with no record
	identity, err := FetchUserIdentity(suite.DB(), goodUUID)
	suite.Equal(ErrFetchNotFound, err, "Expected not to find missing Identity")
	suite.Nil(identity)

	alice := testdatagen.MakeDefaultUser(suite.DB())
	identity, err = FetchUserIdentity(suite.DB(), alice.LoginGovUUID.String())
	suite.Nil(err, "loading alice's identity")
	suite.NotNil(identity)
	suite.Equal(alice.ID, identity.ID)
	suite.Equal(alice.LoginGovEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)
	suite.Nil(identity.TspUserID)

	bob := testdatagen.MakeDefaultServiceMember(suite.DB())
	identity, err = FetchUserIdentity(suite.DB(), bob.User.LoginGovUUID.String())
	suite.Nil(err, "loading bob's identity")
	suite.NotNil(identity)
	suite.Equal(bob.UserID, identity.ID)
	suite.Equal(bob.User.LoginGovEmail, identity.Email)
	suite.Equal(bob.ID, *identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)
	suite.Nil(identity.TspUserID)

	carol := testdatagen.MakeDefaultOfficeUser(suite.DB())
	identity, err = FetchUserIdentity(suite.DB(), carol.User.LoginGovUUID.String())
	suite.Nil(err, "loading carol's identity")
	suite.NotNil(identity)
	suite.Equal(*carol.UserID, identity.ID)
	suite.Equal(carol.User.LoginGovEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Equal(carol.ID, *identity.OfficeUserID)
	suite.Nil(identity.TspUserID)

	danielle := testdatagen.MakeDefaultTspUser(suite.DB())
	identity, err = FetchUserIdentity(suite.DB(), danielle.User.LoginGovUUID.String())
	suite.Nil(err, "loading danielle's identity")
	suite.NotNil(identity)
	suite.Equal(*danielle.UserID, identity.ID)
	suite.Equal(danielle.User.LoginGovEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)
	suite.Equal(danielle.ID, *identity.TspUserID)

	systemAdmin := testdatagen.MakeDefaultAdminUser(suite.DB())
	identity, err = FetchUserIdentity(suite.DB(), systemAdmin.User.LoginGovUUID.String())
	suite.Nil(err, "loading systemAdmin's identity")
	suite.NotNil(identity)
	suite.Equal(*systemAdmin.UserID, identity.ID)
	suite.Equal(systemAdmin.User.LoginGovEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)
	suite.Nil(identity.TspUserID)
}

func (suite *ModelSuite) TestFetchAppUserIdentities() {

	suite.T().Run("default user no profile", func(t *testing.T) {
		testdatagen.MakeDefaultUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.Empty(identities)
	})

	suite.T().Run("service member", func(t *testing.T) {

		// Regular service member
		testdatagen.MakeDefaultServiceMember(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.NotNil(identities[0].ServiceMemberID)
			suite.Nil(identities[0].OfficeUserID)
			suite.Nil(identities[0].TspUserID)
		}

		// Service member is super user
		testdatagen.MakeServiceMember(suite.DB(), testdatagen.Assertions{
			User: User{},
		})
		identities, err = FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(2, len(identities))

		if len(identities) == 2 {
			suite.NotNil(identities[1].ServiceMemberID)
			suite.Nil(identities[1].OfficeUserID)
			suite.Nil(identities[1].TspUserID)
		}
	})

	// In the following tests you won't see extra users returned. Eeach query is
	// limited by the app it expects to be run in.

	suite.T().Run("office user", func(t *testing.T) {
		testdatagen.MakeDefaultOfficeUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.OfficeApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.Nil(identities[0].ServiceMemberID)
			suite.NotNil(identities[0].OfficeUserID)
			suite.Nil(identities[0].TspUserID)
		}
	})

	suite.T().Run("tsp user", func(t *testing.T) {
		testdatagen.MakeDefaultTspUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.TspApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.Nil(identities[0].ServiceMemberID)
			suite.Nil(identities[0].OfficeUserID)
			suite.NotNil(identities[0].TspUserID)
		}
	})

	suite.T().Run("admin user", func(t *testing.T) {
		testdatagen.MakeDefaultAdminUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.AdminApp, 5)
		suite.Nil(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.Nil(identities[0].ServiceMemberID)
			suite.Nil(identities[0].OfficeUserID)
			suite.Nil(identities[0].TspUserID)
			suite.NotNil(identities[0].AdminUserID)
		}
	})
}

func (suite *ModelSuite) TestGetUser() {

	alice := testdatagen.MakeDefaultUser(suite.DB())

	user1, err := GetUserFromEmail(suite.DB(), alice.LoginGovEmail)
	suite.Nil(err, "loading alice's user")
	suite.NotNil(user1)
	if err == nil && user1 != nil {
		suite.Equal(alice.ID, user1.ID)
		suite.Equal(alice.LoginGovEmail, user1.LoginGovEmail)
	}

	user2, err := GetUser(suite.DB(), alice.ID)
	suite.Nil(err, "loading alice's user")
	suite.NotNil(user2)
	if err == nil && user2 != nil {
		suite.Equal(alice.ID, user2.ID)
	}
}
