package authentication

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	verifier "github.com/okta/okta-jwt-verifier-golang"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
)

// ! See flow here:
// ! https://developer.okta.com/docs/guides/implement-grant-type/authcode/main/

func getProfileData(r *http.Request, appCtx appcontext.AppContext, hostname string) (map[string]string, error) {
	m := make(map[string]string)

	if appCtx.Session().AccessToken == "" {
		return m, nil
	}

	reqUrl := hostname + "/oauth2/default/v1/userinfo"

	req, _ := http.NewRequest("GET", reqUrl, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Bearer "+appCtx.Session().AccessToken)
	h.Add("Accept", "application/json")

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	err := json.Unmarshal(body, &m)
	if err != nil {
		appCtx.Logger().Error("get profile data", zap.Error(err))
		return nil, err
	}

	return m, nil
}

// ! Refactor after chamber is modified
func verifyToken(t string, nonce string, session *auth.Session, orgURL string) (*verifier.Jwt, error) {

	// Gather Okta information
	clientID := os.Getenv("OKTA_CUSTOMER_CLIENT_ID")
	if session.IsOfficeApp() {
		clientID = os.Getenv("OKTA_OFFICE_CLIENT_ID")
	} else if session.IsAdminApp() {
		clientID = os.Getenv("OKTA_ADMIN_CLIENT_ID")
	}

	tv := map[string]string{}
	tv["nonce"] = nonce
	tv["aud"] = clientID

	issuer := orgURL + "/oauth2/default"
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

// ! Refactor once chamber is holding new secrets
func exchangeCode(code string, r *http.Request, appCtx appcontext.AppContext, hash string) (Exchange, error) {
	session := auth.SessionFromRequestContext(r)

	appType := "CUSTOMER"
	if session.IsOfficeApp() {
		appType = "OFFICE"
	} else if session.IsAdminApp() {
		appType = "ADMIN"
	}

	authHeader := base64.StdEncoding.EncodeToString(
		[]byte(os.Getenv("OKTA_"+appType+"_CLIENT_ID") + ":" + os.Getenv("OKTA_"+appType+"_SECRET_KEY")))

	q := r.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Set("code", code)
	// TODO: Replace os.Getenv
	q.Add("redirect_uri", os.Getenv("OKTA_"+appType+"_CALLBACK_URL"))

	// TODO: Replace os.Getenv
	url := os.Getenv("OKTA_"+appType+"_HOSTNAME") + "/oauth2/default/v1/token?" + q.Encode()

	req, _ := http.NewRequest("POST", url, bytes.NewReader([]byte("")))
	h := req.Header
	h.Add("Authorization", "Basic "+authHeader)
	h.Add("Accept", "application/json")
	h.Add("Content-Type", "application/x-www-form-urlencoded")
	h.Add("Connection", "close")
	h.Add("Content-Length", "0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		appCtx.Logger().Error("Code exchange", zap.Error(err))
	}
	fmt.Println("t")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		appCtx.Logger().Error("Code exchange", zap.Error(err))
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
	IdToken          string `json:"id_token,omitempty"`
}
