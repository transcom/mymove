package internalapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	oktaop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/okta_profile"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

func stringPtr(s string) *string {
	return &s
}

// TODO figure out how to write this test correctly - gettin 403 forbidden response
func (suite *HandlerSuite) TestGetOktaProfileHandler() {
	t := suite.T()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/okta-profile" {
			t.Errorf("Expected GET request to '/okta-profile', but got %s request to %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		responseJSON := `{
			"profile": {
				"login": "testuser@okta.mil",
				"email": "testuser@okta.mil",
				"firstName": "Test",
				"lastName": "User",
				"cac_edipi": "1231231231",
			}
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(responseJSON))
		if err != nil {
			t.Errorf("Error writing response: %v", err)
		}
	}))
	defer server.Close()
	loggedInUser := factory.BuildServiceMember(suite.DB(), nil, nil)

	// Given: A logged-in user
	user := internalmessages.UpdateOktaUserProfileData{
		Profile: &internalmessages.OktaUserProfileData{
			Login:     "testuser@okta.com",
			Email:     "testuser@okta.com",
			FirstName: "John",
			LastName:  "Doe",
			CacEdipi:  stringPtr("1234567890"),
		},
	}
	fmt.Print(user)

	// Create a mock HTTP request to your API
	req := httptest.NewRequest("GET", "/okta_profile", nil)
	req = suite.AuthenticateRequest(req, loggedInUser)

	params := oktaop.ShowOktaInfoParams{
		HTTPRequest: req,
	}

	handler := GetOktaProfileHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	suite.Assertions.IsType(nil, response)
}

// TODO figure out how to write this test correctly - gettin 403 forbidden response
func (suite *HandlerSuite) TestUpdateOktaProfileHandler() {
	t := suite.T()
	// Create a test server to emulate Okta's API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the request method and path
		if r.Method != http.MethodPost || r.URL.Path != "/okta-profile" {
			t.Errorf("Expected POST request to '/okta-profile', but got %s request to %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Emulate a successful response from Okta
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"profile": {
				"login": "testuser@example.com",
				"email": "testuser@example.com",
				"firstName": "John",
				"lastName": "Doe",
				"cac_edipi": "1234567890",
				"sub": "notARealNumber,
			}
		}`))
		if err != nil {
			t.Errorf("Error writing response: %v", err)
		}
	}))
	defer server.Close()

	defaultUser := factory.BuildDefaultUser(suite.DB())

	// Create a mock HTTP request using the test server's URL
	reqPayload := internalmessages.UpdateOktaUserProfileData{
		Profile: &internalmessages.OktaUserProfileData{
			Login:     "testuser",
			Email:     "testuser@example.com",
			FirstName: "John",
			LastName:  "Doe",
			CacEdipi:  stringPtr("1234567890"),
		},
	}

	body, _ := json.Marshal(reqPayload)

	// Create a mock HTTP request to your API
	req := httptest.NewRequest("POST", server.URL+"/okta-profile", bytes.NewReader(body))
	req = suite.AuthenticateUserRequest(req, defaultUser)

	params := oktaop.UpdateOktaInfoParams{
		HTTPRequest:               req,
		UpdateOktaUserProfileData: &reqPayload,
	}

	handler := UpdateOktaProfileHandler{suite.NewHandlerConfig()}
	response := handler.Handle(params)

	// TODO figure out how to write this test correctly
	suite.Assertions.IsType(nil, response)
}
