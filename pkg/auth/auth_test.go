package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth/context"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

type AuthSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *AuthSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *AuthSuite) mustSave(model interface{}) {
	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		suite.T().Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func TestAuthSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &AuthSuite{db: db, logger: logger}
	suite.Run(t, hs)
}

func fakeLoginGovProvider(logger *zap.Logger) LoginGovProvider {
	return NewLoginGovProvider("fakeHostname", "secret_key", "client_id", logger)
}

func getHandlerParamsWithToken(ss string, expiry time.Time) (*httptest.ResponseRecorder, *http.Request) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)

	// Set a secure cookie on the request
	cookie := http.Cookie{
		Name:    UserSessionCookieName,
		Value:   ss,
		Path:    "/",
		Expires: expiry,
	}
	req.AddCookie(&cookie)

	return rr, req
}

func createRandomRSAPEM() (s string, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		err = errors.Wrap(err, "failed to generate key")
		return
	}

	asn1 := x509.MarshalPKCS1PrivateKey(priv)
	privBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: asn1,
	})
	s = string(privBytes[:])

	return
}

func (suite *AuthSuite) TestGenerateNonce() {
	t := suite.T()
	nonce := generateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}

func (suite *AuthSuite) TestAuthorizationLogoutHandler() {
	t := suite.T()
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	testHostname := "hostname"
	responsePattern := regexp.MustCompile(`href="(.+)"`)
	req, err := http.NewRequest("GET", "/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	authContext := NewAuthContext(fmt.Sprintf("http://%s", testHostname), suite.logger, fakeLoginGovProvider(suite.logger))
	handler := AuthorizationLogoutHandler(authContext)

	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, fakeUUID, fakeToken)

	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusTemporaryRedirect)
	}

	redirectURL, err := url.Parse(responsePattern.FindStringSubmatch(rr.Body.String())[1])
	if err != nil {
		t.Fatal(err)
	}
	params := redirectURL.Query()

	postRedirectURI, err := url.Parse(params["post_logout_redirect_uri"][0])
	if err != nil {
		t.Fatal(err)
	}

	if testHostname != postRedirectURI.Host {
		t.Errorf("handler returned wrong redirect URI hostname: got %v wanted %v", postRedirectURI.Host, testHostname)
	}

	if token := params["id_token_hint"][0]; token != fakeToken {
		t.Errorf("handler returned wrong id_token: got %v wanted %v", token, fakeToken)
	}
}

func (suite *AuthSuite) TestTokenParsingMiddlewareWithBadToken() {
	t := suite.T()
	fakeToken := "some_token"
	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Error("error creating RSA key", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := TokenParsingMiddleware(suite.logger, pem, false)(handler)

	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	rr, req := getHandlerParamsWithToken(fakeToken, expiry)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be no token passed through
	if incomingToken, ok := context.GetIDToken(req.Context()); ok {
		t.Errorf("expected id_token to be nil, got %v", incomingToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestTokenParsingMiddlewareWithValidToken() {
	t := suite.T()
	email := "some_email@domain.com"
	idToken := "fake_id_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	// Brand new token, shouldn't be renewed
	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	ss, err := signTokenStringWithUserInfo(fakeUUID, email, idToken, expiry, pem)
	if err != nil {
		t.Fatal(err)
	}
	rr, req := getHandlerParamsWithToken(ss, expiry)

	var handledRequest *http.Request
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handledRequest = r
	})
	middleware := TokenParsingMiddleware(suite.logger, pem, false)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be an ID token in the request context
	if incomingToken, ok := context.GetIDToken(handledRequest.Context()); !ok || incomingToken != idToken {
		t.Errorf("handler returned wrong id_token: got %v, wanted %v", incomingToken, idToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestTokenParsingMiddlewareWithRenewalToken() {
	t := suite.T()
	email := "some_email@domain.com"
	idToken := "fake_id_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	// Token will expire in 1 minute, should be renewed
	expiry := getExpiryTimeFromMinutes(1)
	ss, err := signTokenStringWithUserInfo(fakeUUID, email, idToken, expiry, pem)
	if err != nil {
		t.Fatal(err)
	}
	rr, req := getHandlerParamsWithToken(ss, expiry)

	var handledRequest *http.Request
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handledRequest = r
	})
	middleware := TokenParsingMiddleware(suite.logger, pem, false)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be an ID token in the request context
	if incomingToken, ok := context.GetIDToken(handledRequest.Context()); !ok || incomingToken != idToken {
		t.Errorf("handler returned wrong id_token: got %v, wanted %v", incomingToken, idToken)
	}

	// And the cookie should be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 1 {
		t.Errorf("expected 1 cookie to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestTokenParsingMiddlewareWithExpiredToken() {
	t := suite.T()
	email := "some_email@domain.com"
	idToken := "fake_id_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")

	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Fatal(err)
	}

	expiry := getExpiryTimeFromMinutes(-1)
	ss, err := signTokenStringWithUserInfo(fakeUUID, email, idToken, expiry, pem)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := TokenParsingMiddleware(suite.logger, pem, false)(handler)

	rr, req := getHandlerParamsWithToken(ss, expiry)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be no token passed through
	if incomingToken, ok := context.GetIDToken(req.Context()); ok {
		t.Errorf("expected id_token to be nil, got %v", incomingToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestRequireAuthMiddleware() {
	t := suite.T()

	// Given: a logged in user
	userUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  userUUID,
		LoginGovEmail: "email@example.com",
		Type:          internalmessages.UserTypeUNKNOWN,
	}
	suite.mustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	// And: the context contains the auth values
	ctx := req.Context()
	ctx = context.PopulateAuthContext(ctx, user.ID, "fake token")
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := RequireAuthMiddleware(handler)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're logged in
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}
}

func (suite *AuthSuite) TestRequireAuthMiddlewareUnauthorized() {
	t := suite.T()

	// Given: No logged in users
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := RequireAuthMiddleware(handler)

	middleware.ServeHTTP(rr, req)

	// We should receive an unauthorized response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusUnauthorized)
	}
}
