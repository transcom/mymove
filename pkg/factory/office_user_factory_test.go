package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildOfficeUser() {
	suite.Run("Successful creation of default office user", func() {
		// Under test:      BuildOfficeUser
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values
		defaultUserEmail := "first.last@login.gov.test"
		defaultTransportationOffice := "JPPSO Testy McTest"

		defaultOffice := models.OfficeUser{
			FirstName: "Leo",
			LastName:  "Spaceman",
			Email:     "leo_spaceman_office@example.com",
			Telephone: "415-555-1212",
		}

		officeUser := BuildOfficeUser(suite.DB(), nil, nil)
		suite.Equal(defaultUserEmail, officeUser.User.LoginGovEmail)
		suite.False(officeUser.User.Active)
		suite.Equal(defaultOffice.FirstName, officeUser.FirstName)
		suite.Equal(defaultOffice.LastName, officeUser.LastName)
		suite.Equal(defaultOffice.Email, officeUser.Email)
		suite.Equal(defaultOffice.Telephone, officeUser.Telephone)
		suite.Equal(defaultTransportationOffice, officeUser.TransportationOffice.Name)
		suite.False(officeUser.Active)

	})

	suite.Run("Successful creation of officeUser with trait", func() {
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
		customOffice := models.OfficeUser{
			Email:     "mycustom@example.com",
			FirstName: "Jason",
			LastName:  "Ash",
		}
		transportationOffice := models.TransportationOffice{
			Name:  "Test transportaion office",
			Gbloc: "TEST",
		}
		customEmail := "leospaceman456@example.com"
		officeUser := BuildOfficeUser(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
			},
			{Model: customOffice},
			{Model: transportationOffice},
		}, nil)
		suite.Equal(customEmail, officeUser.User.LoginGovEmail)
		suite.Equal(customOffice.Email, officeUser.Email)
		suite.Equal(customOffice.FirstName, officeUser.FirstName)
		suite.Equal(customOffice.LastName, officeUser.LastName)
		suite.Equal(transportationOffice.Name, officeUser.TransportationOffice.Name)
		suite.Equal(transportationOffice.Gbloc, officeUser.TransportationOffice.Gbloc)
		suite.False(officeUser.User.Active)
	})
}
func (suite *FactorySuite) TestBuildOfficeUserExtra() {
	// Under test:      BuildOfficeUser
	// Mocked:          None
	// Set up:          Create a OfficeUser but pass in a role
	// Expected outcome:Created User should have the associated Role

	suite.Run("Successful creation of TIO Office User", func() {

		// Create the TIO Role
		tioRole := roles.Role{
			RoleType: roles.RoleTypeTIO,
			RoleName: "Transportation Invoicing Officer",
		}
		_, err := FetchOrBuildRole(suite.DB(), tioRole.RoleType, tioRole.RoleName)
		suite.NoError(err)

		// FUNCTION UNDER TEST
		officeUser := BuildOfficeUser(suite.DB(), nil, []Trait{
			GetTraitOfficeUserTIO,
		})

		// VALIDATE RESULT
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
		//                   No new user should be created

		// SETUP
		user := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					CurrentOfficeSessionID: "breathe",
				},
			},
		}, nil)
		// Count how many users we have
		precount, err := suite.DB().Count(&models.User{})
		suite.NoError(err)

		// FUNCTION UNDER TEST
		officeUser := BuildOfficeUser(suite.DB(), []Customization{
			{
				Model:    user,
				LinkOnly: true,
			},
		}, []Trait{
			GetTraitOfficeUserEmail,
		})

		// VALIDATION
		// Check that no new user was created
		count, err := suite.DB().Count(&models.User{})
		suite.NoError(err)
		suite.Equal(precount, count)

		// Check that the linked user was used
		suite.Equal(user.ID, *officeUser.UserID)
		suite.Equal(user.ID, officeUser.User.ID)
		suite.Equal("breathe", officeUser.User.CurrentOfficeSessionID)
		suite.False(officeUser.Active)

	})
	suite.Run("Successful creation of OfficeUser with forced id User", func() {
		// Under test:       BuildOfficeUser
		// Set up:           Create an officeUser and pass in an ID for User
		// Expected outcome: officeUser and User should be created
		//                   User should have specified ID

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
		// Set up:           Create a stubbed officeUser and pass in a precreated user
		// Expected outcome: officeUser and User should be returned as expected
		uuid := uuid.FromStringOrNil("6f97d298-1502-4d8c-9472-f8b5b2a63a10")
		officeUser := BuildOfficeUser(nil, []Customization{
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

		// Check that id cannot be found in DB
		foundUser := models.User{}
		err := suite.DB().Find(&foundUser, uuid)
		suite.Error(err)

		// Check that email was applied to user
		suite.Equal(officeUser.Email, officeUser.User.LoginGovEmail)
	})
}
