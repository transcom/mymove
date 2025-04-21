package user

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jarcoal/httpmock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"
)

const oktaUsersURL = "OrgURL/api/v1/users/"

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

func (suite *UserServiceSuite) TestDeleteOktaAccount() {
	queryBuilder := query.NewQueryBuilder()
	deleter := NewUserDeleter(queryBuilder)

	//Setup Okta stuff
	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	suite.Run("Success - No attempt to delete Okta account for user without an OktaId", func() {

		user := factory.BuildNonOktaUser(suite.DB(), nil, nil)
		suite.Empty(user.OktaID)

		mockOktaGetEndpointNoError(user.OktaID, models.OktaStatusActive)
		mockOktaDeleteEndpointNoError(user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		err = deleter.DeleteUser(appCtx, user.ID)
		suite.NoError(err)

		// verify calls to okta
		callInfo := httpmock.GetCallCountInfo()
		getEndpoint := oktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]
		deleteEndpoint := oktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(0, getCallCount, "GET Okta user endpoint should NOT be called for an user with an empty oktaID")
		suite.Equal(0, deleteCallCount, "DELETE Okta user endpoint should NOT be called for an user with an empty oktaID")
	})

	suite.Run("Success - Okta account deleted for ACTIVE Okta user", func() {

		user := factory.BuildUser(suite.DB(), nil, nil)
		suite.NotNil(user.OktaID)

		mockOktaGetEndpointNoError(user.OktaID, models.OktaStatusActive)
		mockOktaDeleteEndpointNoError(user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		err = deleter.DeleteUser(appCtx, user.ID)
		suite.NoError(err)

		// Get the call count info
		callInfo := httpmock.GetCallCountInfo()

		// Check if the GET endpoint was called
		getEndpoint := oktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]

		// Check if the DELETE endpoint was called
		deleteEndpoint := oktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(1, getCallCount, "GET Okta user endpoint should not be called")
		suite.Equal(2, deleteCallCount, "DELETE Okta user endpoint should be called twice for an active user")
	})

	suite.Run("Success - Okta account deleted for DEPROVISIONSED Okta user", func() {

		user := factory.BuildUser(suite.DB(), nil, nil)
		suite.NotNil(user.OktaID)

		mockOktaGetEndpointNoError(user.OktaID, models.OktaStatusDeprovisioned)
		mockOktaDeleteEndpointNoError(user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		err = deleter.DeleteUser(appCtx, user.ID)
		suite.NoError(err)

		// Get the call count info
		callInfo := httpmock.GetCallCountInfo()

		// Check if the GET endpoint was called
		getEndpoint := oktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]

		// Check if the DELETE endpoint was called
		deleteEndpoint := oktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(1, getCallCount, "GET Okta user endpoint should be called once")
		suite.Equal(1, deleteCallCount, "DELETE Okta user endpoint should be called once for a deprovisioned user")
	})

	suite.Run("Success - Okta account not deleted - Okta user not found", func() {

		oktaProvider := okta.NewOktaProvider(suite.Logger())
		err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
		suite.NoError(err)

		user := factory.BuildUser(suite.DB(), nil, nil)
		suite.NotNil(user.OktaID)

		mockOktaGetEndpointError(user.OktaID)
		mockOktaDeleteEndpointNoError(user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			Hostname:        "adminlocal",
		}

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)

		// Create an observed logger
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		testLogger := suite.Logger()
		observedLogger := testLogger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, observedZapCore)
		}))
		appCtx := appcontext.NewAppContext(suite.DB(), observedLogger, session, request)

		err = deleter.DeleteUser(appCtx, user.ID)
		suite.NoError(err)

		expectedMessage := "error deleting user from okta"
		foundLog := false
		for _, log := range observedLogs.All() {
			if log.Level == zap.ErrorLevel && strings.Contains(log.Message, expectedMessage) {
				foundLog = true
				break
			}
		}
		suite.Assert().True(foundLog, "Expected error log message not found")

		callInfo := httpmock.GetCallCountInfo()
		getEndpoint := oktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]
		deleteEndpoint := oktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(1, getCallCount, "GET Okta user endpoint should be called once")
		suite.Equal(0, deleteCallCount, "DELETE Okta user endpoint should NOT be called for user not found")
	})
}

func mockOktaGetEndpointNoError(oktaID string, status models.OktaStatus) {
	getUsersEndpoint := "OrgURL/api/v1/users/" + oktaID
	response := fmt.Sprintf(`{
			"id": "%s",
			"status": "%s",
			"created": "2025-02-07T20:39:47.000Z",
			"activated": "2025-02-07T20:39:47.000Z",
			"profile": {
				"firstName": "First",
				"lastName": "Last",
				"mobilePhone": "555-555-5555",
				"secondEmail": "",
				"login": "email@email.com",
				"email": "email@email.com",
				"cac_edipi": "1234567890"
			}
		}`, oktaID, status)

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
}

func mockOktaGetEndpointError(oktaID string) {
	getUsersEndpoint := "OrgURL/api/v1/users/" + oktaID
	response := `[
			{
				"errorSummary": "didn't find the okta user"
			}
		]`

	httpmock.RegisterResponder("GET", getUsersEndpoint,
		httpmock.NewStringResponder(404, response))
}

func mockOktaDeleteEndpointNoError(oktaID string) {
	// oktaID := "fakeOktaID"
	deleteUserEndpoint := oktaUsersURL + oktaID

	httpmock.RegisterResponder(http.MethodDelete, deleteUserEndpoint,
		httpmock.NewStringResponder(204, ""))
}
