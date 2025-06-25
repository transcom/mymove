package models_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/markbates/goth"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
)

func setupOktaMilProvider(suite *ModelSuite) *okta.Provider {
	provider, err := factory.BuildOktaProvider(okta.MilProviderName)
	suite.NoError(err)
	return provider
}

func setupOktaAdminProvider(suite *ModelSuite) *okta.Provider {
	provider, err := factory.BuildOktaProvider(okta.AdminProviderName)
	suite.NoError(err)
	return provider
}

func (suite *ModelSuite) TestSearchForExistingOktaUsers() {
	provider := setupOktaMilProvider(suite)
	httpmock.Activate()
	mockOktaGETEndpointExistingUserNoError(provider)
	oktaEmail := "test@example.com"
	oktaEdipi := "1234567890"

	users, err := models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", oktaEmail, &oktaEdipi, nil)

	suite.NoError(err)
	suite.Equal(1, len(users))
	suite.Equal("fakeOktaID", users[0].ID)
}

func (suite *ModelSuite) TestSearchForExistingOktaUsersValidation() {
	provider := setupOktaMilProvider(suite)

	// invalid email format
	_, err := models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", "invalid-email", nil, nil)
	suite.Error(err)
	suite.Contains(err.Error(), "invalid email format")

	// empty email
	_, err = models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", "", nil, nil)
	suite.Error(err)
	suite.Contains(err.Error(), "email is required")

	// invalid EDIPI format (not 10 digits)
	invalidEdipi := "12345"
	_, err = models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", "test@example.com", &invalidEdipi, nil)
	suite.Error(err)
	suite.Contains(err.Error(), "invalid EDIPI format")
}

func (suite *ModelSuite) TestCreateOktaUser() {
	provider := setupOktaMilProvider(suite)
	payload := models.OktaUserPayload{
		Profile: models.OktaProfile{
			FirstName:   "New",
			LastName:    "User",
			Email:       "newuser@example.com",
			Login:       "newuser@example.com",
			MobilePhone: "555-555-5555",
			CacEdipi:    "9876543210",
		},
		GroupIds: []string{"group-123"},
	}

	httpmock.Activate()
	mockOktaPOSTEndpointsNoError(provider)

	createdUser, err := models.CreateOktaUser(suite.AppContextForTest(), provider, "fakeKey", payload)

	suite.NoError(err)
	suite.Equal("newFakeOktaID", createdUser.ID)
}

func (suite *ModelSuite) TestGetOktaUserGroups_Success() {
	const milProviderName = "milProvider"
	provider, err := factory.BuildOktaProvider(milProviderName)
	suite.NoError(err)
	userID := "fakeUserID"

	// create a JSON response that returns two groups
	groupsJSON := `[
	{"id": "group1", "profile": { "name": "Test Group 1", "description": "Description 1" }},
	{"id": "group2", "profile": { "name": "Test Group 2", "description": "Description 2" }}
	]`

	groupsEndpoint := provider.GetUserGroupsURL(userID)
	httpmock.RegisterResponder("GET", groupsEndpoint,
		httpmock.NewStringResponder(200, groupsJSON))
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groups, err := models.GetOktaUserGroups(suite.AppContextForTest(), provider, "fakeKey", userID)
	suite.NoError(err)
	suite.Equal(2, len(groups))
	suite.Equal("group1", groups[0].ID)
	suite.Equal("Test Group 1", groups[0].Profile.Name)
}

func (suite *ModelSuite) TestAddOktaUserToGroup_Success() {
	const milProviderName = "milProvider"
	provider, err := factory.BuildOktaProvider(milProviderName)
	suite.NoError(err)
	groupID := "group123"
	userID := "user456"

	// okta returns a 204 with an empty body when successful
	url := provider.AddUserToGroupURL(groupID, userID)
	httpmock.RegisterResponder("PUT", url, httpmock.NewStringResponder(204, ""))
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	err = models.AddOktaUserToGroup(suite.AppContextForTest(), provider, "fakeKey", groupID, userID)
	suite.NoError(err)
}

func (suite *ModelSuite) TestAddOktaUserToGroup_Failure() {
	const milProviderName = "milProvider"
	provider, err := factory.BuildOktaProvider(milProviderName)
	suite.NoError(err)
	groupID := "group123"
	userID := "user456"

	// simulate an error response from Okta.
	errorResponse := `{
		"errorCode": "E0000001",
		"errorSummary": "Invalid group",
		"errorLink": "http://example.com",
		"errorId": "abc123",
		"errorCauses": []
	}`

	url := provider.AddUserToGroupURL(groupID, userID)
	httpmock.RegisterResponder("PUT", url, httpmock.NewStringResponder(200, errorResponse))
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// call the function and verify that it returns an error
	err = models.AddOktaUserToGroup(suite.AppContextForTest(), provider, "fakeKey", groupID, userID)
	suite.Error(err)
	suite.Contains(err.Error(), "Invalid group")
}

func (suite *ModelSuite) TestDeleteOktaUser() {
	const oktaID = "fakeOktaID"
	provider := setupOktaMilProvider(suite)

	httpmock.Activate()
	mockOktaGetUserEndpointNoError(provider, oktaID, models.OktaStatusActive)
	mockOktaDeleteEndpointNoError(provider, oktaID)

	err := models.DeleteOktaUser(suite.AppContextForTest(), provider, "fakeOktaID", "fakeKey")
	suite.NoError(err)
}

func (suite *ModelSuite) TestDeleteOktaUserHandled() {
	provider := setupOktaAdminProvider(suite)
	goth.UseProviders(provider)
	expectedOktaUsersURL := provider.GetUsersURL()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	session := &auth.Session{
		ApplicationName: auth.AdminApp,
		Hostname:        "adminlocal",
	}

	suite.Run("Success - No attempt to delete Okta account for user without an OktaId", func() {
		user := factory.BuildNonOktaUser(suite.DB(), nil, nil)
		suite.Empty(user.OktaID)

		mockOktaGetUserEndpointNoError(provider, user.OktaID, models.OktaStatusActive)
		mockOktaDeleteEndpointNoError(provider, user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		models.DeleteOktaUserHandled(appCtx, user.OktaID)

		// verify calls to okta
		callInfo := httpmock.GetCallCountInfo()
		getEndpoint := expectedOktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]
		deleteEndpoint := expectedOktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(0, getCallCount, "GET Okta user endpoint should NOT be called for an user with an empty oktaID")
		suite.Equal(0, deleteCallCount, "DELETE Okta user endpoint should NOT be called for an user with an empty oktaID")
	})

	suite.Run("Success - Okta account deleted for ACTIVE Okta user", func() {

		user := factory.BuildUser(suite.DB(), nil, nil)
		suite.NotNil(user.OktaID)

		mockOktaGetUserEndpointNoError(provider, user.OktaID, models.OktaStatusActive)
		mockOktaDeleteEndpointNoError(provider, user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		models.DeleteOktaUserHandled(appCtx, user.OktaID)

		// Get the call count info
		callInfo := httpmock.GetCallCountInfo()

		// Check if the GET endpoint was called
		getEndpoint := expectedOktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]

		// Check if the DELETE endpoint was called
		deleteEndpoint := expectedOktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(1, getCallCount, "GET Okta user endpoint should be called once")
		suite.Equal(2, deleteCallCount, "DELETE Okta user endpoint should be called twice for an active user")
	})

	suite.Run("Success - Okta account deleted for DEPROVISIONED Okta user", func() {

		user := factory.BuildUser(suite.DB(), nil, nil)
		suite.NotNil(user.OktaID)

		mockOktaGetUserEndpointNoError(provider, user.OktaID, models.OktaStatusDeprovisioned)
		mockOktaDeleteEndpointNoError(provider, user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)
		appCtx := appcontext.NewAppContext(suite.DB(), suite.AppContextForTest().Logger(), session, request)

		models.DeleteOktaUserHandled(appCtx, user.OktaID)

		// Get the call count info
		callInfo := httpmock.GetCallCountInfo()

		// Check if the GET endpoint was called
		getEndpoint := expectedOktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]

		// Check if the DELETE endpoint was called
		deleteEndpoint := expectedOktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(1, getCallCount, "GET Okta user endpoint should be called once")
		suite.Equal(1, deleteCallCount, "DELETE Okta user endpoint should be called once for a deprovisioned user")
	})

	suite.Run("Success - Okta account not deleted - Okta user not found", func() {

		user := factory.BuildUser(suite.DB(), nil, nil)
		suite.NotNil(user.OktaID)

		mockOktaGetUserEndpointError(provider, user.OktaID)
		mockOktaDeleteEndpointNoError(provider, user.OktaID)

		request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", user.ID.String()), nil)

		ctx := auth.SetSessionInRequestContext(request, session)
		request = request.WithContext(ctx)

		// Create an observed logger
		observedZapCore, observedLogs := observer.New(zap.InfoLevel)
		testLogger := suite.Logger()
		observedLogger := testLogger.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewTee(core, observedZapCore)
		}))
		appCtx := appcontext.NewAppContext(suite.DB(), observedLogger, session, request)

		models.DeleteOktaUserHandled(appCtx, user.OktaID)

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
		getEndpoint := expectedOktaUsersURL + user.OktaID
		getCallCount := callInfo[http.MethodGet+" "+getEndpoint]
		deleteEndpoint := expectedOktaUsersURL + user.OktaID
		deleteCallCount := callInfo[http.MethodDelete+" "+deleteEndpoint]

		suite.Equal(1, getCallCount, "GET Okta user endpoint should be called once")
		suite.Equal(0, deleteCallCount, "DELETE Okta user endpoint should NOT be called for user not found")
	})
}

func mockOktaGETEndpointExistingUserNoError(provider *okta.Provider) {
	getUsersEndpoint := provider.GetUsersURL()
	oktaID := "fakeOktaID"

	response := fmt.Sprintf(`[
		{
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
		}
	]`, oktaID, models.OktaStatusProvisioned)

	httpmock.RegisterResponder(http.MethodGet, getUsersEndpoint,
		httpmock.NewStringResponder(200, response))
}

func mockOktaPOSTEndpointsNoError(provider *okta.Provider) {
	activate := "true"
	createUserEndpoint := provider.GetCreateUserURL(activate)
	oktaID := "newFakeOktaID"

	httpmock.RegisterResponder(http.MethodPost, createUserEndpoint,
		httpmock.NewStringResponder(200, fmt.Sprintf(`{
		"id": "%s",
		"profile": {
			"firstName": "First",
			"lastName": "Last",
			"email": "email@email.com",
			"login": "email@email.com"
		}
	}`, oktaID)))
}

func mockOktaGetUserEndpointNoError(provider *okta.Provider, oktaID string, status models.OktaStatus) {
	getUserEndpoint := provider.GetUserURL(oktaID)

	httpmock.RegisterResponder(http.MethodGet, getUserEndpoint,
		httpmock.NewStringResponder(200, fmt.Sprintf(`{
		"id": "%s",
		"status": "%s",
		"profile": {
			"firstName": "First",
			"lastName": "Last",
			"email": "email@email.com",
			"login": "email@email.com"
		}
	}`, oktaID, status)))
}

func mockOktaDeleteEndpointNoError(provider *okta.Provider, oktaID string) {
	deleteUserEndpoint := provider.GetUserURL(oktaID)

	httpmock.RegisterResponder(http.MethodDelete, deleteUserEndpoint,
		httpmock.NewStringResponder(204, ""))
}

func mockOktaGetUserEndpointError(provider *okta.Provider, oktaID string) {
	getUsersEndpoint := provider.GetUserURL(oktaID)
	response := `[
			{
				"errorSummary": "didn't find the okta user"
			}
		]`

	httpmock.RegisterResponder(http.MethodGet, getUsersEndpoint,
		httpmock.NewStringResponder(404, response))
}
