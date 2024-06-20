package okta_test

import (
	"testing"

	"github.com/markbates/goth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type OktaSuite struct {
	handlers.BaseHandlerTestSuite
}

func TestAuthSuite(t *testing.T) {
	hs := &OktaSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

// TODO: Flesh out with actually accessing the data. See error that can pop up when okta provider is not wrapped correctly.
func (suite *OktaSuite) TestConvertGothProviderToOktaProvider() {
	oktaProvider := &okta.Provider{}
	gothProvider := goth.Provider(oktaProvider)

	provider, err := okta.ConvertGothProviderToOktaProvider(gothProvider)
	suite.NoError(err)
	suite.Equal(oktaProvider, provider)
}

func TestRegisterProviders(t *testing.T) {
	// Create a new instance of Viper for testing
	v := viper.New()

	// Set mock configuration values in Viper
	v.Set(cli.OktaTenantOrgURLFlag, "mock-okta-org-url")
	v.Set(cli.OktaCustomerCallbackURL, "mock-customer-callback-url")
	v.Set(cli.OktaCustomerClientIDFlag, "mock-customer-client-id")
	v.Set(cli.OktaCustomerSecretKeyFlag, "mock-customer-secret-key")
	v.Set(cli.OktaOfficeCallbackURL, "mock-office-callback-url")
	v.Set(cli.OktaOfficeClientIDFlag, "mock-office-client-id")
	v.Set(cli.OktaOfficeSecretKeyFlag, "mock-office-secret-key")
	v.Set(cli.OktaAdminCallbackURL, "mock-admin-callback-url")
	v.Set(cli.OktaAdminClientIDFlag, "mock-admin-client-id")
	v.Set(cli.OktaAdminSecretKeyFlag, "mock-admin-secret-key")

	// Create a Provider instance for testing
	provider := &okta.Provider{}

	// Call the RegisterProviders function with the mock Viper configuration
	err := provider.RegisterProviders(v)

	// Check if there was an error during registration
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Add additional assertions to verify that the providers were registered correctly,
	// for example, by checking if the providers have been added to op.providers.

}

func TestRegisterOktaProvider(t *testing.T) {
	// Create a new instance of your Provider
	op := &okta.Provider{}

	// Mock data for registering the Okta provider
	name := "mock-provider"
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-id"
	secret := "mock-secret"
	scope := []string{"openid", "email", "profile"}

	// Call the RegisterOktaProvider function
	err := op.RegisterOktaProvider(name, orgURL, callbackURL, clientID, secret, scope)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestGetTokenURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetTokenURL()

	expectedURL := orgURL + "/oauth2/default/v1/token"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetAuthURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetAuthURL()

	expectedURL := orgURL + "/oauth2/default/v1/authorize"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetUserInfoURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetUserInfoURL()

	expectedURL := orgURL + "/oauth2/default/v1/userinfo"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetIssuerURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetIssuerURL()

	expectedURL := orgURL + "/oauth2/default"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetLogoutURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetLogoutURL()

	expectedURL := orgURL + "/oauth2/default/v1/logout"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetRevokeURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetRevokeURL()

	expectedURL := orgURL + "/oauth2/default/v1/revoke"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetSessionsURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetSessionsURL()

	expectedURL := orgURL + "/oauth2/default/v1/sessions"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetJWKSURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetJWKSURL()

	expectedURL := orgURL + "/oauth2/default/.well-known/jwks.json"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetOpenIDConfigURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	// Call the GetTokenURL function
	url := provider.GetOpenIDConfigURL()

	expectedURL := orgURL + "/oauth2/default/.well-known/openid-configuration"
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGetUserURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	oktaUserID := "mock-user-id"

	// Call the GetTokenURL function
	url := provider.GetUserURL(oktaUserID)

	expectedURL := orgURL + "/api/v1/users/" + oktaUserID
	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestCreateAccountURL(t *testing.T) {
	// Create a new instance of your Provider with the desired orgURL
	orgURL := "https://mock-okta-org-url.com"
	callbackURL := "https://mock-callback-url.com"
	clientID := "mock-client-ID"
	secret := "mock-secret"
	provider := okta.NewProvider(orgURL, callbackURL, clientID, secret)

	activate := "true"

	// Call the GetCreateAccountURL function
	url := provider.GetCreateAccountURL(activate)

	expectedURL := orgURL + "/api/v1/users/?activate=" + activate

	if url != expectedURL {
		t.Errorf("Expected URL to be '%s', got: '%s'", expectedURL, url)
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce := okta.GenerateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}
