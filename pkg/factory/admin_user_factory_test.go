package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildAdminUser() {
	defaultEmail := "first.last@okta.mil"
	suite.Run("Successful creation of default admin user", func() {
		// Under test:      BuildAdminUser
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values
		defaultAdmin := models.AdminUser{
			FirstName: "Leo",
			LastName:  "Spaceman",
			Email:     "leo_spaceman_admin@example.com",
			Role:      "SYSTEM_ADMIN",
		}

		adminUser := BuildAdminUser(suite.DB(), nil, nil)
		suite.Equal(defaultEmail, adminUser.User.OktaEmail)
		suite.False(adminUser.User.Active)
		suite.Equal(defaultAdmin.FirstName, adminUser.FirstName)
		suite.Equal(defaultAdmin.LastName, adminUser.LastName)
		suite.Equal(defaultAdmin.Email, adminUser.Email)
		suite.Equal(defaultAdmin.Role, adminUser.Role)
		suite.False(adminUser.Active)
	})

	suite.Run("Successful creation of adminUser with trait", func() {
		// Under test:      BuildAdminUser
		// Mocked:          None
		// Set up:          Create a User but pass in a trait that sets
		//                  both the adminuser and user email to a random
		//                  value, as adminuser has uniqueness constraints
		// Expected outcome:AdminUser should have the same random email as User

		adminUser := BuildAdminUser(suite.DB(), nil, []Trait{
			GetTraitAdminUserEmail,
		})
		suite.Equal(adminUser.Email, adminUser.User.OktaEmail)
		suite.False(adminUser.User.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      BuildAdminUser
		// Set up:          Create an adminUser and pass in specified emails
		// Expected outcome:adminUser and User should be created with specified emails
		customAdmin := models.AdminUser{
			Email:     "mycustom@example.com",
			FirstName: "Jason",
			LastName:  "Ash",
			Role:      "SYSTEM_ADMIN",
		}
		customEmail := "leospaceman456@example.com"
		adminUser := BuildAdminUser(suite.DB(), []Customization{
			{
				Model: models.User{
					OktaEmail: customEmail,
				},
			},
			{Model: customAdmin},
		}, nil)
		suite.Equal(customEmail, adminUser.User.OktaEmail)
		suite.Equal(customAdmin.Email, adminUser.Email)
		suite.Equal(customAdmin.FirstName, adminUser.FirstName)
		suite.Equal(customAdmin.LastName, adminUser.LastName)
		suite.Equal(customAdmin.Role, adminUser.Role)
		suite.False(adminUser.User.Active)
	})
}
func (suite *FactorySuite) TestBuildAdminUserExtra() {
	// Under test:      BuildAdminUser
	// Mocked:          None
	// Set up:          Create a AdminUser but pass in a role
	// Expected outcome:Created User should have the associated Role

	suite.Run("Successful creation of TIO Admin User", func() {

		// Create the TIO Role
		tioRole := roles.Role{
			ID:       uuid.Must(uuid.NewV4()),
			RoleType: roles.RoleTypeTIO,
			RoleName: "Task Invoicing Officer",
		}
		verrs, err := suite.DB().ValidateAndCreate(&tioRole)
		suite.NoError(err)
		suite.False(verrs.HasAny())

		// FUNCTION UNDER TEST
		adminUser := BuildAdminUser(suite.DB(), []Customization{
			{
				Model: models.User{
					Roles: []roles.Role{tioRole},
				},
				Type: &User,
			},
		}, []Trait{
			GetTraitAdminUserEmail,
		})

		// VALIDATE RESULT
		// Check that the email trait worked
		suite.Equal(adminUser.Email, adminUser.User.OktaEmail)
		suite.False(adminUser.User.Active)
		// Check that the user has the admin user role
		_, hasRole := adminUser.User.Roles.GetRole(roles.RoleTypeTIO)
		suite.True(hasRole)
	})

	suite.Run("Successful creation of AdminUser with linked User", func() {
		// Under test:       BuildAdminUser
		// Set up:           Create an adminUser and pass in a precreated user
		// Expected outcome: adminUser should link in the precreated user
		//                   No new user should be created

		// SETUP
		user := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					CurrentAdminSessionID: "breathe",
				},
			},
		}, nil)
		// Count how many users we have
		precount, err := suite.DB().Count(&models.User{})
		suite.NoError(err)

		// FUNCTION UNDER TEST
		adminUser := BuildAdminUser(suite.DB(), []Customization{
			{
				Model:    user,
				LinkOnly: true,
			},
		}, []Trait{
			GetTraitAdminUserEmail,
		})

		// VALIDATION
		// Check that no new user was created
		count, err := suite.DB().Count(&models.User{})
		suite.NoError(err)
		suite.Equal(precount, count)

		// Check that the linked user was used
		suite.Equal(user.ID, *adminUser.UserID)
		suite.Equal(user.ID, adminUser.User.ID)
		suite.Equal("breathe", adminUser.User.CurrentAdminSessionID)
		suite.False(adminUser.Active)

	})
	suite.Run("Successful creation of AdminUser with forced id User", func() {
		// Under test:       BuildAdminUser
		// Set up:           Create an adminUser and pass in an ID for User
		// Expected outcome: adminUser and User should be created
		//                   User should have specified ID

		defaultOktaEmail := "first.last@okta.mil"
		uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
		adminUser := BuildAdminUser(suite.DB(), []Customization{
			{
				Model: models.User{
					ID: uuid,
				},
			},
		}, []Trait{
			GetTraitAdminUserEmail,
		})
		// Check that the forced ID was used
		suite.Equal(uuid, *adminUser.UserID)
		suite.Equal(uuid, adminUser.User.ID)

		// Check that id can be found in DB
		foundUser := models.User{}
		err := suite.DB().Find(&foundUser, uuid)
		suite.NoError(err)

		// Check that email was applied to user
		suite.NotContains(defaultOktaEmail, adminUser.User.OktaEmail)
		suite.Equal(adminUser.Email, adminUser.User.OktaEmail)
	})

	suite.Run("Successful creation of stubbed AdminUser with forced id User", func() {
		// Under test:       BuildAdminUser
		// Set up:           Create a stubbed adminUser and pass in a precreated user
		// Expected outcome: adminUser and User should be returned as expected
		uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
		adminUser := BuildAdminUser(nil, []Customization{
			{
				Model: models.User{
					ID: uuid,
				},
			},
		}, []Trait{
			GetTraitAdminUserEmail,
		})
		// Check that the forced ID was used
		suite.Equal(uuid, *adminUser.UserID)
		suite.Equal(uuid, adminUser.User.ID)

		// Check that id cannot be found in DB
		foundUser := models.User{}
		err := suite.DB().Find(&foundUser, uuid)
		suite.Error(err)

		// Check that email was applied to user
		suite.Equal(adminUser.Email, adminUser.User.OktaEmail)
	})
}
