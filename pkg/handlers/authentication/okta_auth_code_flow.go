package authentication

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	verifier "github.com/okta/okta-jwt-verifier-golang"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
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

func exchangeCode(code string, r *http.Request, appCtx appcontext.AppContext, provider okta.Provider, client HTTPClient) (Exchange, error) {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(provider.GetClientID() + ":" + provider.GetSecret()))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Set("code", code)
	q.Add("redirect_uri", provider.GetCallbackURL())

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

// Future functionality
/*
func revokeCurrentOktaToken(r *http.Request, provider okta.Provider, revokeURL, accessToken string) error {
	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(provider.GetClientID() + ":" + provider.GetSecret()))

	data := url.Values{}
	data.Set("token", accessToken)
	data.Set("token_type_hint", "access_token")

	req, err := http.NewRequest("POST", revokeURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creating the request: %v", err)
	}

	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		return fmt.Errorf("failed to revoke token, status: %v, body: %s", resp.Status, string(bodyBytes))
	}

	return nil
}
*/
