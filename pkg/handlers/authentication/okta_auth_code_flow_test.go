package authentication

import (
	"net/url"
	"testing"

	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
)

func TestLogoutOktaUserURL(t *testing.T) {
	provider := &okta.Provider{}
	idToken := "mockIDToken"
	redirectURL := "https://example.com/"

	// Call the function being tested
	logoutURL, err := logoutOktaUserURL(provider, idToken, redirectURL)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Parse the returned URL to verify the query parameters
	parsedURL, err := url.Parse(logoutURL)
	if err != nil {
		t.Errorf("Failed to parse logout URL: %v", err)
	}

	// Check id_token_hint parameter
	idTokenHint := parsedURL.Query().Get("id_token_hint")
	if idTokenHint != idToken {
		t.Errorf("Expected id_token_hint parameter to be '%s', got '%s'", idToken, idTokenHint)
	}

	// Check post_logout_redirect_uri parameter
	postLogoutRedirectURI := parsedURL.Query().Get("post_logout_redirect_uri")
	expectedRedirectURL := redirectURL + "sign-in" + "?okta_logged_out=true"
	if postLogoutRedirectURI != expectedRedirectURL {
		t.Errorf("Expected post_logout_redirect_uri parameter to be '%s', got '%s'", expectedRedirectURL, postLogoutRedirectURI)
	}
}
