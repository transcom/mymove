package authentication

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	verifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/models"
)

// ! See flow here:
// ! https://developer.okta.com/docs/guides/implement-grant-type/authcode/main/

func getProfileData(appCtx appcontext.AppContext, provider okta.Provider) (models.OktaUser, error) {
	user := models.OktaUser{}

	if appCtx.Session().AccessToken == "" {
		return user, nil
	}

	reqURL := provider.GetUserInfoURL()

	req, _ := http.NewRequest("GET", reqURL, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+appCtx.Session().AccessToken)
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("could not execute request", zap.Error(err))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("could not read response body", zap.Error(err))
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &user)
	if err != nil {
		appCtx.Logger().Error("could not unmarshal body", zap.Error(err))
		return user, err
	}

	return user, nil
}

func verifyToken(t string, nonce string, provider okta.Provider) (*verifier.Jwt, error) {
	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = provider.GetClientID()

	issuer := provider.GetIssuerURL()
	jv := verifier.JwtVerifier{
		Issuer:           issuer,
		ClaimsToValidate: tv,
	}

	result, err := jv.New().VerifyIdToken(t)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	if result != nil {
		return result, nil
	}

	return nil, fmt.Errorf("token could not be verified: %s", "")
}

type OktaProfile struct {
	ID string `json:"id,omitempty"`
}

func clearOktaUserSessions(appCtx appcontext.AppContext, r *http.Request, provider okta.Provider, client HTTPClient) (string, error) {
	// if the user makes it here that means they have an existing okta session and MM is no longer storing it
	// Okta should still have the cookie saved as a cookie, which will find and clear all user sessions
	var oktaSessionToken string
	for _, c := range r.Cookies() {
		if c.Name == "office_okta_state" {
			oktaSessionToken = c.Value
			break
		}
	}
	appCtx.Logger().Info(oktaSessionToken)

	// setting viper so we can access the api key in the env vars
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	apiKey := v.GetString(cli.OktaAPIKeyFlag)

	// getting the api call url from provider.go
	getUserURL := provider.GetUserURLWithToken()

	// we need to know the user's okta ID so we can clear the sessions
	// https://developer.okta.com/docs/reference/api/users/#get-user
	req, _ := http.NewRequest("GET", getUserURL, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "SSWS "+apiKey)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/json")
	h.Add("Cookie", oktaSessionToken)

	// make the request and snag the okta ID
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("Error making http response", zap.Error(err))
		return "error", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("Error reading response body", zap.Error(err))
		return "error", err
	}
	defer resp.Body.Close()
	var profile OktaProfile
	err = json.Unmarshal(body, &profile)
	if err != nil {
		appCtx.Logger().Error("JSON unmarshal error", zap.Error(err))
		return "error", err
	}

	if profile.ID == "" {
		if err != nil {
			appCtx.Logger().Error("Okta profile ID not found or returned from Okta", zap.Error(err))
			return "error", err
		}
	}

	// now we will clear all the user's sessions with the DELETE endpoint using the Okta ID
	// https://developer.okta.com/docs/reference/api/users/#user-sessions
	clearSessionURL := provider.ClearUserSessionsURL(profile.ID)
	req2, _ := http.NewRequest("DELETE", clearSessionURL, bytes.NewReader([]byte("")))
	headers := req2.Header
	headers.Add("Accept", "application/json")
	headers.Add("Content-Type", "application/json")
	headers.Add("Authorization", "SSWS "+apiKey)

	resp2, err := client.Do(req2)
	if err != nil {
		appCtx.Logger().Error("Error clearing okta user session", zap.Error(err))
		return "error", err
	}
	defer resp2.Body.Close()

	if resp2.StatusCode == http.StatusNoContent {
		appCtx.Logger().Info("Response has no content (204 No Content)")
		return "success", nil
	}

	return "error", err
}

func exchangeCode(code string, r *http.Request, appCtx appcontext.AppContext, provider okta.Provider, client HTTPClient) (Exchange, error) {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(provider.GetClientID() + ":" + provider.GetSecret()))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Set("code", code)
	q.Add("redirect_uri", provider.GetCallbackURL())
	q.Add("scope", "openid email profile")

	url := provider.GetTokenURL() + "?" + q.Encode()

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	if err != nil {
		appCtx.Logger().Error("Post request generate", zap.Error(err))
		return Exchange{}, err
	}
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("Exchange client request", zap.Error(err))
		return Exchange{}, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("Exchange response body", zap.Error(err))
		return Exchange{}, err
	}
	defer resp.Body.Close()
	var exchange Exchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		appCtx.Logger().Error("get profile data", zap.Error(err))
		return Exchange{}, err
	}

	return exchange, nil
}

type Exchange struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	AccessToken      string `json:"access_token,omitempty"`
	TokenType        string `json:"token_type,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Scope            string `json:"scope,omitempty"`
	IDToken          string `json:"id_token,omitempty"`
}

// logging a user out of okta requires calling the /logout API endpoint
// it is a GET request and clears the browser session
// the URL will need to be built using the ID token and a redirect URI
func logoutOktaUserURL(provider *okta.Provider, idToken string, redirectURL string) (string, error) {
	// baseURL will end in /logout
	baseURL := provider.GetLogoutURL()

	// Parse URL
	logoutURL, err := url.Parse(baseURL)
	if err != nil {
		return "Failed to parse logout URL", err
	}

	// add params required by Okta to successfully sign a user out
	params := logoutURL.Query()
	params.Set("id_token_hint", idToken)
	params.Set("post_logout_redirect_uri", redirectURL+"sign-in"+"?okta_logged_out=true")

	logoutURL.RawQuery = params.Encode()

	oktaLogoutURL := logoutURL.String()

	return oktaLogoutURL, err
}
