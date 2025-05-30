package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	m "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	r "github.com/transcom/mymove/pkg/models/roles"
	userroles "github.com/transcom/mymove/pkg/services/users_roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestUserValidation() {
	oktaID := "abcdefghijklmnopqrst"
	userEmail := "sally@government.gov"

	newUser := m.User{
		OktaID:    oktaID,
		OktaEmail: userEmail,
	}

	verrs, err := newUser.Validate(nil)

	suite.NoError(err)
	suite.False(verrs.HasAny(), "Error validating model")
	suite.Equal(userEmail, newUser.OktaEmail)
	suite.Equal(oktaID, newUser.OktaID)
}

func (suite *ModelSuite) TestUserCreationWithoutValues() {
	newUser := &m.User{}

	expErrors := map[string][]string{
		"okta_email": {"OktaEmail can not be blank."},
	}

	suite.verifyValidationErrors(newUser, expErrors, nil)
}

func (suite *ModelSuite) TestCreateUser() {
	const testEmail = "Sally@GoVernment.gov"
	const expectedEmail = "sally@government.gov"
	oktaID := factory.MakeRandomString(20)

	sally, err := m.CreateUser(suite.DB(), oktaID, testEmail)
	suite.Nil(err, "No error for good create")
	suite.Equal(expectedEmail, sally.OktaEmail, "should convert email to lower case")
	suite.NotEqual(sally.ID, uuid.Nil)
}

func (suite *ModelSuite) TestFetchUserIdentity() {
	oktaID := factory.MakeRandomString(20)
	// First check that it all works with no record
	identity, err := m.FetchUserIdentity(suite.DB(), oktaID)
	suite.Equal(m.ErrFetchNotFound, err, "Expected not to find missing Identity")
	suite.Nil(identity)

	alice := factory.BuildDefaultUser(suite.DB())
	identity, err = m.FetchUserIdentity(suite.DB(), alice.OktaID)
	suite.Nil(err, "loading alice's identity")
	suite.NotNil(identity)
	suite.Equal(alice.ID, identity.ID)
	suite.Equal(alice.OktaEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)

	bob := factory.BuildServiceMember(suite.DB(), nil, nil)
	identity, err = m.FetchUserIdentity(suite.DB(), bob.User.OktaID)
	suite.Nil(err, "loading bob's identity")
	suite.NotNil(identity)
	suite.Equal(bob.UserID, identity.ID)
	suite.Equal(bob.User.OktaEmail, identity.Email)
	suite.Equal(bob.ID, *identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)

	carolUser := factory.BuildDefaultUser(suite.DB())

	carol := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
		{
			Model: m.OfficeUser{
				UserID: &carolUser.ID,
			},
		},
		{
			Model:    carolUser,
			LinkOnly: true,
		},
	}, nil)
	identity, err = m.FetchUserIdentity(suite.DB(), carol.User.OktaID)
	suite.Nil(err, "loading carol's identity")
	suite.NotNil(identity)
	suite.Equal(*carol.UserID, identity.ID)
	suite.Equal(carol.User.OktaEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Equal(carol.ID, *identity.OfficeUserID)

	systemAdmin := factory.BuildDefaultAdminUser(suite.DB())
	identity, err = m.FetchUserIdentity(suite.DB(), systemAdmin.User.OktaID)
	suite.Nil(err, "loading systemAdmin's identity")
	suite.NotNil(identity)
	suite.Equal(*systemAdmin.UserID, identity.ID)
	suite.Equal(systemAdmin.User.OktaEmail, identity.Email)
	suite.Nil(identity.ServiceMemberID)
	suite.Nil(identity.OfficeUserID)
	customerRole := roles.Role{
		RoleType: roles.RoleTypeCustomer,
	}
	tooRole := roles.Role{
		RoleType: roles.RoleTypeTOO,
	}
	pat := factory.BuildUserAndUsersRoles(suite.DB(), []factory.Customization{
		{
			Model: models.User{
				Roles: []roles.Role{customerRole},
			},
		},
	}, nil)

	identity, err = m.FetchUserIdentity(suite.DB(), pat.OktaID)
	suite.Nil(err, "loading pat's identity")
	suite.NotNil(identity)
	suite.Equal(len(identity.Roles), 1)
	billy := factory.BuildUserAndUsersRoles(suite.DB(), []factory.Customization{
		{
			Model: models.User{
				Roles: []roles.Role{tooRole},
			},
		},
	}, nil)

	suite.DB().MigrationURL()
	identity, err = m.FetchUserIdentity(suite.DB(), billy.OktaID)
	suite.Nil(err, "loading billy's identity")
	suite.NotNil(identity)
	suite.Equal(len(identity.Roles), 1)
	suite.Equal(identity.Roles[0].RoleType, tooRole.RoleType)

	supervisorPrivilege := factory.FetchOrBuildPrivilegeByPrivilegeType(suite.DB(), r.PrivilegeTypeSupervisor)

	sueOktaID := factory.MakeRandomString(20)
	sue := factory.BuildUser(suite.DB(), []factory.Customization{
		{
			Model: m.User{
				OktaID:     sueOktaID,
				Active:     true,
				Privileges: []roles.Privilege{supervisorPrivilege},
			},
		},
	}, nil)

	identity, err = m.FetchUserIdentity(suite.DB(), sue.OktaID)
	suite.Nil(err, "loading sue's identity")
	suite.NotNil(identity)
	suite.Equal(len(identity.Privileges), 1)
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
	multiRoleUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
	identity, err := m.FetchUserIdentity(suite.DB(), multiRoleUser.User.OktaID)
	suite.Nil(err, "failed to fetch user identity")
	suite.Equal(*multiRoleUser.UserID, identity.ID)
	suite.Condition(compareRoleTypeLists(multiRoleUser.User.Roles, identity.Roles))

	/*
		Test that user identity is properly fetched after deleting roles
	*/
	// then update user roles to soft delete
	userRoles := userroles.NewUsersRolesCreator()
	// we'll be soft deleting the services counselor role
	updateToRoles := []roles.RoleType{
		roles.RoleTypeTOO,
		roles.RoleTypeTIO,
	}
	_, _, err = userRoles.UpdateUserRoles(suite.AppContextForTest(), *multiRoleUser.UserID, updateToRoles)
	suite.NoError(err)

	// re-fetch user identity and check roles
	identity, err = m.FetchUserIdentity(suite.DB(), multiRoleUser.User.OktaID)
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
		identities, err := m.FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)
		suite.NoError(err)
		suite.Empty(identities)
	})

	suite.Run("service member", func() {

		// Create a user email that won't be filtered out of the devlocal user query w/ a default value of
		// first.last@okta.mil
		user := factory.BuildUser(suite.DB(), []factory.Customization{
			{
				Model: m.User{
					OktaEmail: "test@example.com",
				},
			}}, nil)

		factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model:    user,
				LinkOnly: true,
			},
		}, nil)

		// This service member will be filtered out from the result because we haven't overridden the default email
		factory.BuildServiceMember(suite.DB(), nil, nil)

		identities, err := m.FetchAppUserIdentities(suite.DB(), auth.MilApp, 5)

		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))
		suite.NotNil(identities[0].ServiceMemberID)
		suite.Equal("test@example.com", identities[0].Email)
		suite.Nil(identities[0].OfficeUserID)
	})

	// In the following tests you won't see extra users returned. Eeach query is
	// limited by the app it expects to be run in.

	suite.Run("office user", func() {
		factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		identities, err := m.FetchAppUserIdentities(suite.DB(), auth.OfficeApp, 5)
		suite.NoError(err)
		suite.NotEmpty(identities)
		suite.Equal(1, len(identities))

		if len(identities) > 1 {
			suite.Nil(identities[0].ServiceMemberID)
			suite.NotNil(identities[0].OfficeUserID)
		}
	})

	suite.Run("admin user", func() {
		factory.BuildDefaultAdminUser(suite.DB())
		identities, err := m.FetchAppUserIdentities(suite.DB(), auth.AdminApp, 5)
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
	alice := factory.BuildDefaultUser(suite.DB())

	user1, err := m.GetUserFromEmail(suite.DB(), alice.OktaEmail)
	suite.Nil(err, "loading alice's user")
	suite.NotNil(user1)
	if err == nil && user1 != nil {
		suite.Equal(alice.ID, user1.ID)
		suite.Equal(alice.OktaEmail, user1.OktaEmail)
	}

	user2, err := m.GetUser(suite.DB(), alice.ID)
	suite.Nil(err, "loading alice's user")
	suite.NotNil(user2)
	if err == nil && user2 != nil {
		suite.Equal(alice.ID, user2.ID)
	}
}

func (suite *ModelSuite) TestGetUserFromOktaID() {
	user := factory.BuildDefaultUser(suite.DB())

	foundUser, err := m.GetUserFromOktaID(suite.DB(), user.OktaID)
	suite.Nil(err)
	suite.NotNil(foundUser)
	if err == nil && foundUser != nil {
		suite.Equal(user.ID, foundUser.ID)
		suite.Equal(user.OktaEmail, foundUser.OktaEmail)
	}
}
