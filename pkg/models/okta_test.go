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

	users, err := models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", oktaEmail, &oktaEdipi)

	suite.NoError(err)
	suite.Equal(1, len(users))
	suite.Equal("fakeOktaID", users[0].ID)
}

func (suite *ModelSuite) TestSearchForExistingOktaUsersValidation() {
	provider := setupOktaProvider(suite)

	// invalid email format
	_, err := models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", "invalid-email", nil)
	suite.Error(err)
	suite.Contains(err.Error(), "invalid email format")

	// empty email
	_, err = models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", "", nil)
	suite.Error(err)
	suite.Contains(err.Error(), "email is required")

	// invalid EDIPI format (not 10 digits)
	invalidEdipi := "12345"
	_, err = models.SearchForExistingOktaUsers(suite.AppContextForTest(), provider, "fakeKey", "test@example.com", &invalidEdipi)
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
