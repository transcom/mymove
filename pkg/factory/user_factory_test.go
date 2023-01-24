package factory

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildUser() {
	defaultEmail := "first.last@login.gov.test"
	customEmail := "leospaceman123@example.com"
	suite.Run("Successful creation of default user", func() {
		// Under test:      BuildUser
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values

		user := BuildUser(suite.DB(), nil, nil)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.False(user.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      BuildUser
		// Set up:          Create a User with a customized email and no trait
		// Expected outcome:User should be created with email and inactive status
		user := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
			},
		}, nil)
		suite.Equal(customEmail, user.LoginGovEmail)
		suite.False(user.Active)

	})

	suite.Run("Successful creation of user with trait", func() {
		// Under test:      BuildUser
		// Set up:          Create a User with a trait
		// Expected outcome:User should be created with default email and active status

		user := BuildUser(suite.DB(), nil,
			[]Trait{
				GetTraitActiveUser,
			})
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

	suite.Run("Successful creation of user with both", func() {
		// Under test:      BuildUser
		// Set up:          Create a User with a customized email and active trait
		// Expected outcome:User should be created with email and active status

		user := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
			}}, []Trait{
			GetTraitActiveUser,
		})
		suite.Equal(customEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

	suite.Run("Successful creation of stubbed user", func() {
		// Under test:      BuildUser
		// Set up:          Create a customized user, but don't pass in a db
		// Expected outcome:User should be created with email and active status
		//                  No user should be created in database
		precount, err := suite.DB().Count(&models.User{})
		suite.NoError(err)

		user := BuildUser(nil, []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
			}}, []Trait{
			GetTraitActiveUser,
		})

		suite.Equal(customEmail, user.LoginGovEmail)
		suite.True(user.Active)
		// Count how many users are in the DB, no new users should have been created.
		count, err := suite.DB().Count(&models.User{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}

func (suite *FactorySuite) TestBuildDefaultUser() {
	defaultEmail := "first.last@login.gov.test"
	suite.Run("Successful creation of default user", func() {
		// Under test:      BuildDefaultUser
		// Mocked:          None
		// Set up:          Use helper function BuildDefaultUser
		// Expected outcome:User should be created with GetTraitActiveUser

		user := BuildDefaultUser(suite.DB())
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

	suite.Run("Successful creation of stubbed default user", func() {
		// Under test:      BuildDefaultUser
		// Mocked:          None
		// Set up:          Use helper function BuildDefaultUser, but no db
		// Expected outcome:User should be created with GetTraitActiveUser

		user := BuildDefaultUser(nil)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

	suite.Run("Successful creation of User, Roles, and UsersRoles using BuildUserAndUsersRoles", func() {
		// Under test:      BuildUserAndUsersRoles
		// Mocked:          None
		// Set up:          Use helper function BuildUserAndUsersRoles
		// Expected outcome:User with correct roles, Role, and UsersRoles should be created

		precountRole, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)

		precountUsersRoles, err := suite.DB().Count(&models.UsersRoles{})
		suite.NoError(err)

		tioRole := roles.Role{
			RoleType: roles.RoleTypeTIO,
		}

		user := BuildUserAndUsersRoles(suite.DB(), []Customization{
			{
				Model: models.User{
					Roles: []roles.Role{tioRole},
				},
			},
		}, nil)

		// Check that the user has the office user role
		_, hasRole := user.Roles.GetRole(roles.RoleTypeTIO)
		suite.True(hasRole)

		// Count how many roles are in the DB, new role should have been created.
		count, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)
		suite.Equal(precountRole+1, count)

		// Count how many UsersRoles are in the DB, new UsersRoles should have been created.
		count, err = suite.DB().Count(&roles.Role{})
		suite.NoError(err)
		suite.Equal(precountUsersRoles+1, count)

	})
}
