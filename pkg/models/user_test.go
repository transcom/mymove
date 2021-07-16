package models_test

import (
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	. "github.com/transcom/mymove/pkg/models"
	userroles "github.com/transcom/mymove/pkg/services/users_roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestUserValidation() {
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc1")
	userEmail := "sally@government.gov"

	newUser := User{
		LoginGovUUID:  &fakeUUID,
		LoginGovEmail: userEmail,
	}

	verrs, err := newUser.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
	suite.Equal(userEmail, newUser.LoginGovEmail)
	suite.Equal(fakeUUID, *newUser.LoginGovUUID)
}

func (suite *ModelSuite) TestUserCreationWithoutValues() {
	newUser := &User{}

	expErrors := map[string][]string{
		"login_gov_email": {"LoginGovEmail can not be blank."},
	}

	suite.verifyValidationErrors(newUser, expErrors)
}

func (suite *ModelSuite) TestUserCreationDuplicateUUID() {
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	userEmail := "sally@government.gov"

	newUser := User{
		LoginGovUUID:  &fakeUUID,
		LoginGovEmail: userEmail,
	}

	sameUser := User{
		LoginGovUUID:  &fakeUUID,
		LoginGovEmail: userEmail,
	}

	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
	//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
	//RA: in which this would be considered a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	suite.DB().Create(&newUser)
	err := suite.DB().Create(&sameUser)

	suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "constraint_name"), "Db should have errored on unique constraint for UUID")
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

	bob := testdatagen.MakeDefaultServiceMember(suite.DB())
	identity, err = FetchUserIdentity(suite.DB(), bob.User.LoginGovUUID.String())
	suite.Nil(err, "loading bob's identity")
	suite.NotNil(identity)
	suite.Equal(bob.UserID, identity.ID)
	suite.Equal(bob.User.LoginGovEmail, identity.Email)
	suite.Equal(bob.ID, *identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)

	carolUser := testdatagen.MakeDefaultUser(suite.DB())
	carol := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: OfficeUser{
			UserID: &carolUser.ID,
			User:   carolUser,
		},
	})
	identity, err = FetchUserIdentity(suite.DB(), carol.User.LoginGovUUID.String())
	suite.Nil(err, "loading carol's identity")
	suite.NotNil(identity)
	suite.Equal(*carol.UserID, identity.ID)
	suite.Equal(carol.User.LoginGovEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Equal(carol.ID, *identity.OfficeUserID)

	adminUser := testdatagen.MakeDefaultUser(suite.DB())
	systemAdmin := testdatagen.MakeAdminUser(suite.DB(), testdatagen.Assertions{
		AdminUser: AdminUser{
			User:   adminUser,
			UserID: &adminUser.ID,
		},
	})
	identity, err = FetchUserIdentity(suite.DB(), systemAdmin.User.LoginGovUUID.String())
	suite.Nil(err, "loading systemAdmin's identity")
	suite.NotNil(identity)
	suite.Equal(*systemAdmin.UserID, identity.ID)
	suite.Equal(systemAdmin.User.LoginGovEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)

	rs := []roles.Role{{
		ID:       uuid.FromStringOrNil("ed2d2cd7-d427-412a-98bb-a9b391d98d32"),
		RoleType: roles.RoleTypeCustomer,
	}, {
		ID:       uuid.FromStringOrNil("9dc423b6-33b8-493a-a59b-6a823660cb07"),
		RoleType: roles.RoleTypeTOO,
	},
	}
	suite.NoError(suite.DB().Create(&rs))
	customerRole := rs[0]
	patUUID := uuid.Must(uuid.NewV4())
	pat := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
		User: User{
			LoginGovUUID: &patUUID,
			Active:       true,
			Roles:        []roles.Role{customerRole},
		},
	})

	identity, err = FetchUserIdentity(suite.DB(), pat.LoginGovUUID.String())
	suite.Nil(err, "loading pat's identity")
	suite.NotNil(identity)
	suite.Equal(len(identity.Roles), 1)

	tooRole := rs[1]
	billyUUID := uuid.Must(uuid.NewV4())
	billy := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{
		User: User{
			LoginGovUUID: &billyUUID,
			Active:       true,
			Roles:        []roles.Role{tooRole},
		},
	})

	suite.DB().MigrationURL()
	identity, err = FetchUserIdentity(suite.DB(), billy.LoginGovUUID.String())
	suite.Nil(err, "loading billy's identity")
	suite.NotNil(identity)
	suite.Equal(len(identity.Roles), 1)
	suite.Equal(identity.Roles[0].RoleType, tooRole.RoleType)
}

func (suite *ModelSuite) TestFetchUserIdentityDeletedRoles() {
	// creates a custom comparison function for testing the role type
	compareRoleTypeLists := func(expectedList roles.Roles, actualList roles.Roles) func() (success bool) {
		return func() (success bool) {
			// compare length first
			if len(expectedList) != len(actualList) {
				return false
			}

			// then compare the role type
			// types are unique so we shouldn't run into duplicates
			for _, expectedRole := range expectedList {
				roleMatches := false
				for _, actualRole := range actualList {
					if expectedRole.RoleType == actualRole.RoleType {
						roleMatches = true
						break
					}
				}

				if !roleMatches {
					return false
				}
			}

			return true
		}
	}

	/*
		Test that user identity is properly fetched
	*/
	// this creates a user with TOO, TIO, and Services Counselor roles
	multiRoleUser := testdatagen.MakeOfficeUserWithMultipleRoles(suite.DB(), testdatagen.Assertions{})
	identity, err := FetchUserIdentity(suite.DB(), multiRoleUser.User.LoginGovUUID.String())
	suite.Nil(err, "failed to fetch user identity")
	suite.Equal(*multiRoleUser.UserID, identity.ID)
	suite.Condition(compareRoleTypeLists(multiRoleUser.User.Roles, identity.Roles))

	/*
		Test that user identity is properly fetched after deleting roles
	*/
	// then update user roles to soft delete
	userRoles := userroles.NewUsersRolesCreator(suite.DB())
	// we'll be soft deleting the services counselor role
	updateToRoles := []roles.RoleType{
		roles.RoleTypeTOO,
		roles.RoleTypeTIO,
	}
	_, err = userRoles.UpdateUserRoles(*multiRoleUser.UserID, updateToRoles)
	suite.NoError(err)

	// re-fetch user identity and check roles
	identity, err = FetchUserIdentity(suite.DB(), multiRoleUser.User.LoginGovUUID.String())
	suite.Nil(err, "failed to fetch user identity")
	suite.Equal(*multiRoleUser.UserID, identity.ID)

	expectedRoles := roles.Roles{
		roles.Role{RoleType: roles.RoleTypeTOO},
		roles.Role{RoleType: roles.RoleTypeTIO},
	}
	suite.Condition(compareRoleTypeLists(expectedRoles, identity.Roles))
}

func (suite *ModelSuite) TestFetchAppUserIdentities() {
	suite.Run("default user no profile", func() {
		testdatagen.MakeStubbedUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.Empty(identities)
	})

	suite.Run("service member", func() {

		// Regular service member
		testdatagen.MakeDefaultServiceMember(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.NotNil(identities[0].ServiceMemberID)
			suite.Nil(identities[0].OfficeUserID)
		}

		// Service member is super user
		testdatagen.MakeDefaultServiceMember(suite.DB())
		identities, err = FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(2, len(identities))

		if len(identities) == 2 {
			suite.NotNil(identities[1].ServiceMemberID)
			suite.Nil(identities[1].OfficeUserID)
		}
	})

	// In the following tests you won't see extra users returned. Eeach query is
	// limited by the app it expects to be run in.

	suite.Run("office user", func() {
		testdatagen.MakeDefaultOfficeUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.OfficeApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.Nil(identities[0].ServiceMemberID)
			suite.NotNil(identities[0].OfficeUserID)
		}
	})

	suite.Run("admin user", func() {
		testdatagen.MakeDefaultAdminUser(suite.DB())
		identities, err := FetchAppUserIdentities(suite.DB(), auth.AdminApp, 5)
		suite.Nil(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.Nil(identities[0].ServiceMemberID)
			suite.Nil(identities[0].OfficeUserID)
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
