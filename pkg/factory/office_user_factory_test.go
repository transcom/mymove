package factory

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MakerSuite) TestBuildOfficeUser() {
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
func (suite *MakerSuite) TestBuildOfficeUserRoles() {

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

}
