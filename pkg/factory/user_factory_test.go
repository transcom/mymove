package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MakerSuite) TestUserMaker() {
	defaultEmail := "first.last@login.gov.test"
	customEmail := "leospaceman123@example.com"
	suite.Run("Successful creation of default user", func() {
		// Under test:      UserMaker
		// Mocked:          None
		// Set up:          Create a User with no customizations or traits
		// Expected outcome:User should be created with default values

		user := BuildUser(suite.DB(), nil, nil)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.False(user.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      UserMaker
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
		// Under test:      UserMaker
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
		// Under test:      UserMaker
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
		// Under test:      UserMaker
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
