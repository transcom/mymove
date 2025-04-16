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

	setupBasicUser := func(userRoleType roles.RoleType) models.User {
		userRole := roles.Role{
			RoleType: userRoleType,
		}

		user := factory.BuildUserAndUsersRoles(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{userRole},
				},
			},
		}, nil)

		return user
	}

	suite.Run("success - a simple user is deleted", func() {
		initialUserCount, _ := suite.DB().Count(&models.User{})
		initialServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		initialOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		initialAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		initialUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		initialUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})

		testUser := setupBasicUser(roles.RoleTypeTOO)

		// Verify the test data exists
		var user models.User
		err := suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected the user after setup")

		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles for the user after setup")

		var userPrivileges []models.UsersPrivileges
		// TODO: change or remove for service members
		//err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		//suite.NoError(err)
		//suite.Empty(userPrivileges, "Expected no privileges to remain for the user")

		setupUserCount, _ := suite.DB().Count(&models.User{})
		setupServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		setupOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		setupAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		setupUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		setupUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount+1, setupUserCount)
		suite.Equal(initialServiceMemberCount, setupServiceMemberCount)
		suite.Equal(initialOfficeUserCount, setupOfficeUserCount)
		suite.Equal(initialAdminUserCount, setupAdminUserCount)
		suite.Equal(initialUserRolesCount+1, setupUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, setupUserPrivilegesCount)

		// Delete the user
		err = deleter.DeleteUser(suite.AppContextForTest(), testUser.ID)
		suite.NoError(err)

		// Test that the user was deleted
		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")

		finalUserCount, _ := suite.DB().Count(&models.User{})
		finalServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		finalOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		finalAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		finalUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		finalUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount, finalUserCount)
		suite.Equal(initialServiceMemberCount, finalServiceMemberCount)
		suite.Equal(initialOfficeUserCount, finalOfficeUserCount)
		suite.Equal(initialAdminUserCount, finalAdminUserCount)
		suite.Equal(initialUserRolesCount, finalUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, finalUserPrivilegesCount)
	})

	suite.Run("Success - delete an Office User", func() {
		initialUserCount, _ := suite.DB().Count(&models.User{})
		initialServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		initialOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		initialAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		initialUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		initialUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})

		testUser := factory.BuildDefaultUser(suite.DB())
		status := models.OfficeUserStatusREQUESTED
		_ = factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Active: true,
					UserID: &testUser.ID,
					Email:  testUser.OktaEmail,
					Status: &status,
				},
			},
			{
				Model:    testUser,
				LinkOnly: true,
			},
		}, []roles.RoleType{roles.RoleTypeTOO})

		// Verify the test data exists
		var user models.User
		err := suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected the user after setup")

		var officeUser models.OfficeUser
		err = suite.DB().Where("user_id = ?", testUser.ID).First(&officeUser)
		suite.NoError(err)
		suite.NotEmpty(officeUser, "Expected the office user after setup")

		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles for the user after setup")

		// TODO: create privilages in the test data
		var userPrivileges []models.UsersPrivileges
		//err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		//suite.NoError(err)
		//suite.NotEmpty(userPrivileges, "Expected privileges for the user after setup")

		setupUserCount, _ := suite.DB().Count(&models.User{})
		setupServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		setupOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		setupAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		setupUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		setupUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount+1, setupUserCount)
		suite.Equal(initialServiceMemberCount, setupServiceMemberCount)
		suite.Equal(initialOfficeUserCount+1, setupOfficeUserCount)
		suite.Equal(initialAdminUserCount, setupAdminUserCount)
		suite.Equal(initialUserRolesCount+1, setupUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, setupUserPrivilegesCount)

		// delete the user
		err = deleter.DeleteUser(suite.AppContextForTest(), testUser.ID)
		suite.NoError(err)

		// Test that the user was deleted
		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		err = suite.DB().Where("user_id = ?", testUser.ID).First(&officeUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")

		finalUserCount, _ := suite.DB().Count(&models.User{})
		finalServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		finalOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		finalAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		finalUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		finalUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount, finalUserCount)
		suite.Equal(initialServiceMemberCount, finalServiceMemberCount)
		suite.Equal(initialOfficeUserCount, finalOfficeUserCount)
		suite.Equal(initialAdminUserCount, finalAdminUserCount)
		suite.Equal(initialUserRolesCount, finalUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, finalUserPrivilegesCount)
	})

	suite.Run("Success - delete an Admin User", func() {
		initialUserCount, _ := suite.DB().Count(&models.User{})
		initialServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		initialOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		initialAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		initialUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		initialUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})

		testUser := setupBasicUser(roles.RoleTypeHQ)
		_ = factory.BuildAdminUser(suite.DB(), []factory.Customization{
			{
				Model: models.AdminUser{
					Active: true,
					UserID: &testUser.ID,
					Email:  testUser.OktaEmail,
				},
			},
			{
				Model:    testUser,
				LinkOnly: true,
			},
		}, nil)

		// Verify the test data exists
		var user models.User
		err := suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected the user after setup")

		var adminUser models.AdminUser
		err = suite.DB().Where("user_id = ?", testUser.ID).First(&adminUser)
		suite.NoError(err)
		suite.NotEmpty(adminUser, "Expected the admin user after setup")

		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles for the user after setup")

		// TODO: create privilages in the test data
		var userPrivileges []models.UsersPrivileges
		//err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		//suite.NoError(err)
		//suite.NotEmpty(userPrivileges, "Expected privileges for the user after setup")

		setupUserCount, _ := suite.DB().Count(&models.User{})
		setupServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		setupOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		setupAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		setupUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		setupUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount+1, setupUserCount)
		suite.Equal(initialServiceMemberCount, setupServiceMemberCount)
		suite.Equal(initialOfficeUserCount, setupOfficeUserCount)
		suite.Equal(initialAdminUserCount+1, setupAdminUserCount)
		suite.Equal(initialUserRolesCount+1, setupUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, setupUserPrivilegesCount)

		// delete the user
		err = deleter.DeleteUser(suite.AppContextForTest(), testUser.ID)
		suite.NoError(err)

		// Test that the user was deleted
		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		err = suite.DB().Where("user_id = ?", testUser.ID).First(&adminUser)
		suite.Error(err)
		suite.Equal(sql.ErrNoRows, err, "sql: no rows in result set")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.Empty(userRoles, "Expected no roles to remain for the user")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		suite.NoError(err)
		suite.Empty(userPrivileges, "Expected no privileges to remain for the user")

		finalUserCount, _ := suite.DB().Count(&models.User{})
		finalServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		finalOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		finalAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		finalUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		finalUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount, finalUserCount)
		suite.Equal(initialServiceMemberCount, finalServiceMemberCount)
		suite.Equal(initialOfficeUserCount, finalOfficeUserCount)
		suite.Equal(initialAdminUserCount, finalAdminUserCount)
		suite.Equal(initialUserRolesCount, finalUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, finalUserPrivilegesCount)
	})

	suite.Run("error - a customer user has a move and cannot be deleted", func() {
		initialUserCount, _ := suite.DB().Count(&models.User{})
		initialServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		initialOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		initialAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		initialUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		initialUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})

		testUser := setupBasicUser(roles.RoleTypeCustomer)
		_ = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    testUser,
				LinkOnly: true,
			},
		}, nil)

		// Verify the test data exists
		var user models.User
		err := suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected the user after setup")

		var serviceMember models.ServiceMember
		err = suite.DB().Where("user_id = ?", testUser.ID).First(&serviceMember)
		suite.NoError(err)
		suite.NotEmpty(serviceMember, "Expected ServiceMember after setup")

		var userRoles []models.UsersRoles
		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles for the user after setup")

		//var userPrivileges []models.UsersPrivileges
		// TODO: change or remove for service members
		//err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		//suite.NoError(err)
		//suite.Empty(userPrivileges, "Expected no privileges to remain for the user")

		setupUserCount, _ := suite.DB().Count(&models.User{})
		setupServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		setupOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		setupAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		setupUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		setupUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(initialUserCount+1, setupUserCount)
		suite.Equal(initialServiceMemberCount+1, setupServiceMemberCount)
		suite.Equal(initialOfficeUserCount, setupOfficeUserCount)
		suite.Equal(initialAdminUserCount, setupAdminUserCount)
		suite.Equal(initialUserRolesCount+1, setupUserRolesCount)
		suite.Equal(initialUserPrivilegesCount, setupUserPrivilegesCount)

		// Delete the user
		err = deleter.DeleteUser(suite.AppContextForTest(), testUser.ID)
		suite.Error(err)

		// Test that the user was not deleted
		err = suite.DB().Where("id = ?", testUser.ID).First(&user)
		suite.NoError(err)
		suite.NotEmpty(user, "Expected the user remains after failed delete")

		err = suite.DB().Where("user_id = ?", testUser.ID).First(&serviceMember)
		suite.NoError(err)
		suite.NotEmpty(serviceMember, "Expected ServiceMember remains after failed delete")

		err = suite.DB().Where("user_id = ?", testUser.ID).All(&userRoles)
		suite.NoError(err)
		suite.NotEmpty(userRoles, "Expected roles for the user remain after failed delete")

		//err = suite.DB().Where("user_id = ?", testUser.ID).All(&userPrivileges)
		//suite.NoError(err)
		//suite.Empty(userPrivileges, "Expected no privileges to remain for the user")

		finalUserCount, _ := suite.DB().Count(&models.User{})
		finalServiceMemberCount, _ := suite.DB().Count(&models.ServiceMember{})
		finalOfficeUserCount, _ := suite.DB().Count(&models.OfficeUser{})
		finalAdminUserCount, _ := suite.DB().Count(&models.AdminUser{})
		finalUserRolesCount, _ := suite.DB().Count(&models.UsersRoles{})
		finalUserPrivilegesCount, _ := suite.DB().Count(&models.UsersPrivileges{})
		suite.Equal(setupUserCount, finalUserCount)
		suite.Equal(setupServiceMemberCount, finalServiceMemberCount)
		suite.Equal(setupOfficeUserCount, finalOfficeUserCount)
		suite.Equal(setupAdminUserCount, finalAdminUserCount)
		suite.Equal(setupUserRolesCount, finalUserRolesCount)
		suite.Equal(setupUserPrivilegesCount, finalUserPrivilegesCount)
	})

	suite.Run("error - a user is not found", func() {
		userID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(userID, "while looking for User")

		err := deleter.DeleteUser(suite.AppContextForTest(), userID)
		suite.Error(err)
		suite.Equal(err, expectedError)
	})

}
