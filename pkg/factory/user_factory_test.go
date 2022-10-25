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

		user, err := BuildUser(suite.DB(), nil, nil)
		suite.NoError(err)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.False(user.Active)
	})

	suite.Run("Successful creation of user with customization", func() {
		// Under test:      UserMaker
		// Set up:          Create a User with a customized email and no trait
		// Expected outcome:User should be created with email and inactive status
		user, err := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
				Type: User,
			},
		}, nil)
		suite.NoError(err)
		suite.Equal(customEmail, user.LoginGovEmail)
		suite.False(user.Active)

	})

	suite.Run("Successful creation of user with trait", func() {
		// Under test:      UserMaker
		// Set up:          Create a User with a trait
		// Expected outcome:User should be created with default email and active status

		user, err := BuildUser(suite.DB(), nil,
			[]Trait{
				GetTraitActiveUser,
			})
		suite.NoError(err)
		suite.Equal(defaultEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

	suite.Run("Successful creation of user with both", func() {
		// Under test:      UserMaker
		// Set up:          Create a User with a customized email and active trait
		// Expected outcome:User should be created with email and active status

		user, err := BuildUser(suite.DB(), []Customization{
			{
				Model: models.User{
					LoginGovEmail: customEmail,
				},
				Type: User,
			}}, []Trait{
			GetTraitActiveUser,
		})
		suite.NoError(err)
		suite.Equal(customEmail, user.LoginGovEmail)
		suite.True(user.Active)
	})

}
