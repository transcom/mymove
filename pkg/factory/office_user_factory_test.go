package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *FactorySuite) TestBuildOfficeUser() {
	defaultEmail := "first.last@login.gov.test"
	customEmail := "leospaceman123@example.com"
	suite.Run("Successful creation of default office user", func() {
		// Under test:      BuildOfficeUser
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values

		officeUser := BuildOfficeUser(suite.DB(), nil, nil)
		suite.Equal(defaultEmail, officeUser.User.LoginGovEmail)
		suite.False(officeUser.User.Active)
	})

	suite.Run("Successful creation of officeUser with matched email", func() {
		// Under test:      BuildOfficeUser
		// Mocked:          None
		// Set up:          Create a User but pass in a trait that sets
		//                  both the officeuser and user email to a random
		//                  value, as officeuser has uniqueness constraints
		// Expected outcome:OfficeUser should have the same random email as User

		officeUser := BuildOfficeUser(suite.DB(), nil, []Trait{
			GetTraitOfficeUserEmail,
		})
		suite.Equal(officeUser.Email, officeUser.User.LoginGovEmail)
		suite.False(officeUser.User.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      BuildOfficeUser
		// Set up:          Create an officeUser and pass in specified emails
		// Expected outcome:officeUser and User should be created with specified emails
		customEmail2 := "leospaceman456@example.com"
		officeUser := BuildOfficeUser(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
			},
			{
				Model: models.OfficeUser{
					Email: customEmail2,
				},
			},
		}, nil)
		suite.Equal(customEmail, officeUser.User.LoginGovEmail)
		suite.Equal(customEmail2, officeUser.Email)
		suite.False(officeUser.User.Active)
	})
}
func (suite *FactorySuite) TestBuildOfficeUserExtra() {

	suite.Run("Successful creation of TIO Office User", func() {

		tioRole, _ := testdatagen.LookupOrMakeRole(suite.DB(), roles.RoleTypeTIO, "Transportation Invoicing Officer")

		officeUser := BuildOfficeUser(suite.DB(), []Customization{
			{
				Model: models.User{
					Roles: []roles.Role{tioRole},
				},
				Type: &User,
			},
		}, []Trait{
			GetTraitOfficeUserEmail,
		})
		// Check that the email trait worked
		suite.Equal(officeUser.Email, officeUser.User.LoginGovEmail)
		suite.False(officeUser.User.Active)
		// Check that the user has the office user role
		_, hasRole := officeUser.User.Roles.GetRole(roles.RoleTypeTIO)
		suite.True(hasRole)
	})

	suite.Run("Successful creation of OfficeUser with linked User", func() {
		// Under test:       BuildOfficeUser
		// Set up:           Create an officeUser and pass in a precreated user
		// Expected outcome: officeUser should link in the precreated user

		user := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					CurrentAdminSessionID: "breathe",
				},
			},
		}, nil)
		officeUser := BuildOfficeUser(suite.DB(), []Customization{
			{
				Model:    user,
				LinkOnly: true,
			},
		}, []Trait{
			GetTraitOfficeUserEmail,
		})
		// Check that the linked user was used
		suite.Equal(user.ID, *officeUser.UserID)
		suite.Equal(user.ID, officeUser.User.ID)
		suite.Equal("breathe", officeUser.User.CurrentAdminSessionID)
		suite.False(officeUser.Active)

	})
	suite.Run("Successful creation of OfficeUser with forced id User", func() {
		// Under test:       BuildOfficeUser
		// Set up:           Create an officeUser and pass in a precreated user
		// Expected outcome: officeUser and User should be created with specified emails
		defaultLoginGovEmail := "first.last@login.gov.test"
		uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
		officeUser := BuildOfficeUser(suite.DB(), []Customization{
			{
				Model: models.User{
					ID: uuid,
				},
			},
		}, []Trait{
			GetTraitOfficeUserEmail,
		})
		// Check that the forced ID was used
		suite.Equal(uuid, *officeUser.UserID)
		suite.Equal(uuid, officeUser.User.ID)

		// Check that id can be found in DB
		foundUser := models.User{}
		err := suite.DB().Find(&foundUser, uuid)
		suite.NoError(err)

		// Check that email was applied to user
		suite.NotContains(defaultLoginGovEmail, officeUser.User.LoginGovEmail)
		suite.Equal(officeUser.Email, officeUser.User.LoginGovEmail)
	})

	suite.Run("Successful creation of stubbed OfficeUser with forced id User", func() {
		// Under test:       BuildOfficeUser
		// Set up:           Create an officeUser and pass in a precreated user
		// Expected outcome: officeUser and User should be created with specified emails
	})
}
