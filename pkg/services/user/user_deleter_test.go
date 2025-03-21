package user

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *UserServiceSuite) TestDeleteUser() {
	queryBuilder := query.NewQueryBuilder()
	deleter := NewUserDeleter(queryBuilder)
	setupTestData := func() models.User {
		customerRole := roles.Role{
			RoleType: roles.RoleTypeCustomer,
		}
		user := factory.BuildUserAndUsersRoles(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{customerRole},
				},
			},
		}, nil)
		return user
	}

	suite.Run("success - a customer user is deleted", func() {
		testUser := setupTestData()
		var user models.User
		err := suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected the user after setup")

		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles for the user after setup")

		var userPrivileges []models.UsersPrivileges

		err = deleter.DeleteUser(suite.AppContextForTest(), testUser.ID)
		suite.NoError(err)

		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")
	})

	suite.Run("error - a user is not found", func() {
		userID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(userID, "while looking for User")

		err := deleter.DeleteUser(suite.AppContextForTest(), userID)
		suite.Error(err)
		suite.Equal(err, expectedError)
	})
}
