package officeuser

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
					TOOAssignedID: &officeUser.ID,
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
		suite.ErrorContains(err, "violates foreign key constraint \"moves_too_assigned_id_fkey\" on table \"moves\"")

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

func (suite *OfficeUserServiceSuite) TestDeleteOktaAccount() {
	queryBuilder := query.NewQueryBuilder()
	deleter := NewOfficeUserDeleter(queryBuilder)

	// Setup Okta stuff
	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("adminProvider", "OrgURL", "CallbackURL", "fakeToken", "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	buildOfficeUser := func(user models.User) models.OfficeUser {
		status := models.OfficeUserStatusAPPROVED
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
		return officeUser
	}

	setupOktaUser := func() (models.User, models.OfficeUser) {
		user := factory.BuildDefaultUser(suite.DB())
		officeUser := buildOfficeUser(user)
		return user, officeUser
	}

	setupNonOktaUser := func() (models.User, models.OfficeUser) {
		user := factory.BuildNonOktaUser(suite.DB(), nil, nil)
		officeUser := buildOfficeUser(user)
		return user, officeUser
	}

	suite.Run("Success - No attempt to delete Okta account for user without an OktaId", func() {

		user, officeUser := setupNonOktaUser()
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

		err = deleter.DeleteOfficeUser(appCtx, officeUser.ID)
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

		user, officeUser := setupOktaUser()
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

		err = deleter.DeleteOfficeUser(appCtx, officeUser.ID)
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
		suite.Equal(2, deleteCallCount, "DELETE Okta user endpoint should be called twice for an active user")
	})

	suite.Run("Success - Okta account deleted for DEPROVISIONSED Okta user", func() {

		user, officeUser := setupOktaUser()
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

		err = deleter.DeleteOfficeUser(appCtx, officeUser.ID)
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

		user, officeUser := setupOktaUser()
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

		err = deleter.DeleteOfficeUser(appCtx, officeUser.ID)
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
