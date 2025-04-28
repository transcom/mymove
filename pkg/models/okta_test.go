package models_test

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
)

func setupOktaProvider(suite *ModelSuite) *okta.Provider {
	const milProviderName = "milProvider"
	provider, err := factory.BuildOktaProvider(milProviderName)
	suite.NoError(err)
	return provider
}

func (suite *ModelSuite) TestSearchForExistingOktaUsers() {
	provider := setupOktaProvider(suite)
	httpmock.Activate()
	mockAndActivateOktaGETEndpointExistingUserNoError(provider)
	oktaEmail := "test@example.com"
	oktaEdipi := "1234567890"

	users, err := models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", oktaEmail, &oktaEdipi, nil)

	suite.NoError(err)
	suite.Equal(1, len(users))
	suite.Equal("fakeOktaID", users[0].ID)
}

func (suite *ModelSuite) TestSearchForExistingOktaUsersValidation() {
	provider := setupOktaProvider(suite)

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
	provider := setupOktaProvider(suite)
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
	mockAndActivateOktaPOSTEndpointsNoError(provider)

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
	provider := setupOktaProvider(suite)

	httpmock.Activate()
	mockAndActivateOktaGetUserEndpointNoError(provider)
	mockAndActivateOktaDeleteEndpointNoError(provider)

	err := models.DeleteOktaUser(suite.AppContextForTest(), provider, "fakeOktaID", "fakeKey")
	suite.NoError(err)
}

func mockAndActivateOktaGETEndpointExistingUserNoError(provider *okta.Provider) {
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

func mockAndActivateOktaPOSTEndpointsNoError(provider *okta.Provider) {
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

func mockAndActivateOktaGetUserEndpointNoError(provider *okta.Provider) {
	oktaID := "fakeOktaID"
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
	}`, oktaID, models.OktaStatusActive)))
}

func mockAndActivateOktaDeleteEndpointNoError(provider *okta.Provider) {
	oktaID := "fakeOktaID"
	deleteUserEndpoint := provider.GetUserURL(oktaID)

	httpmock.RegisterResponder(http.MethodDelete, deleteUserEndpoint,
		httpmock.NewStringResponder(204, ""))
}
