package officeuser

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *OfficeUserServiceSuite) TestDeleteOfficeUser() {
	queryBuilder := query.NewQueryBuilder()
	deleter := NewOfficeUserDeleter(queryBuilder)
	setupTestUser := func(status models.OfficeUserStatus) (models.User, models.OfficeUser) {
		user := factory.BuildDefaultUser(suite.DB())
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

	setupTestUserWithAssignedMove := func() (models.User, models.OfficeUser) {
		user, officeUser := setupTestUser(models.OfficeUserStatusAPPROVED)
		_ = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					TOOTaskOrderAssignedID: &officeUser.ID,
				},
			},
			{
				Model:    officeUser,
				LinkOnly: true,
			},
		}, nil)
		return user, officeUser
	}

	suite.Run("success - a requested office user is deleted", func() {
		testUser, testOfficeUser := setupTestUser(models.OfficeUserStatusREQUESTED)

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

	suite.Run("success - an active office user is deleted", func() {
		testUser, testOfficeUser := setupTestUser(models.OfficeUserStatusAPPROVED)

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

	suite.Run("error - an office user assigned to a move", func() {
		testUser, testOfficeUser := setupTestUserWithAssignedMove()

		err := deleter.DeleteOfficeUser(suite.AppContextForTest(), testOfficeUser.ID)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.ErrorContains(err, "violates foreign key constraint \"moves_too_task_order_assigned_id_fkey\" on table \"moves\"")

		var user models.User
		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected user to remain after failed delete")

		var officeUser models.OfficeUser
		err = suite.DB().Where("user_id = ?", testUser.ID).First(&officeUser)
		suite.NoError(err)
		suite.NotEmpty(officeUser, "Expected office user to remain after failed delete")

		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles to remain after failed delete")
	})

	suite.Run("error - an office user is not found", func() {
		officeUserID := uuid.Must(uuid.NewV4())

		err := deleter.DeleteOfficeUser(suite.AppContextForTest(), officeUserID)
		suite.Error(err)
	})
}
