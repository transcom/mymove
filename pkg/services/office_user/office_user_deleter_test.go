package officeuser

import (
	"database/sql"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *OfficeUserServiceSuite) TestDeleteOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	deleter := NewOfficeUserDeleter(queryBuilder)
	setupTestData := func() (models.User, models.OfficeUser) {
		user := factory.BuildDefaultUser(suite.DB())
		status := models.OfficeUserStatusREQUESTED
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					UserID: &user.ID,
					Email:  user.OktaEmail,
					Status: &status,
				},
			},
			{
				Model:    user,
				LinkOnly: true,
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		return user, officeUser
	}

	suite.Run("success - a requested office user is deleted", func() {
		testUser, testOfficeUser := setupTestData()

		err := deleter.DeleteOfficeUser(suite.AppContextForTest(), testOfficeUser.ID)
		suite.NoError(err)

		var user models.User
		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		var officeUser models.OfficeUser
		err = suite.DB().Where("user_id = ?", testUser.ID).First(&officeUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		// .All does not return a sql no rows error, so we will verify that the struct is empty
		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		var userPrivileges []models.UsersPrivileges
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")
	})
}
