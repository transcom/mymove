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

	"github.com/gorilla/context"
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

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
	testHostname := "hostname"
	responsePattern := regexp.MustCompile(`href="(.+)"`)
	req, err := http.NewRequest("GET", "/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AuthorizationLogoutHandler(fmt.Sprintf("http://%s", testHostname)))

	context.Set(req, "id_token", fakeToken)

	handler.ServeHTTP(rr, req)

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

func (suite *AuthSuite) TestEnforceUserAuthMiddlewareWithBadToken() {
	t := suite.T()
	fakeToken := "some_token"
	pem, err := createRandomRSAPEM()
	if err != nil {
		t.Error("error creating RSA key", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(pem)(handler)

	expiry := getExpiryTimeFromMinutes(sessionExpiryInMinutes)
	rr, req := getHandlerParamsWithToken(fakeToken, expiry)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be no token passed through
	if incomingToken, ok := context.Get(req, "id_token").(string); ok {
		t.Errorf("expected id_token to be nil, got %v", incomingToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestUserAuthMiddlewareWithValidToken() {
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

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(pem)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be an ID token in the request context
	if incomingToken, ok := context.Get(req, "id_token").(string); !ok || incomingToken != idToken {
		t.Errorf("handler returned wrong id_token: got %v, wanted %v", incomingToken, idToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestUserAuthMiddlewareWithRenewalToken() {
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

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := UserAuthMiddleware(pem)(handler)

	middleware.ServeHTTP(rr, req)

	// We should get a 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be an ID token in the request context
	if incomingToken, ok := context.Get(req, "id_token").(string); !ok || incomingToken != idToken {
		t.Errorf("handler returned wrong id_token: got %v, wanted %v", incomingToken, idToken)
	}

	// And the cookie should be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 1 {
		t.Errorf("expected 1 cookie to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestPassiveUserAuthMiddlewareWithExpiredToken() {
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
	middleware := UserAuthMiddleware(pem)(handler)

	rr, req := getHandlerParamsWithToken(ss, expiry)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're not enforcing auth
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusOK)
	}

	// And there should be no token passed through
	if incomingToken, ok := context.Get(req, "id_token").(string); ok {
		t.Errorf("expected id_token to be nil, got %v", incomingToken)
	}

	// And the cookie should not be renewed
	if setCookies := rr.HeaderMap["Set-Cookie"]; len(setCookies) != 0 {
		t.Errorf("expected no cookies to be set, got %v", len(setCookies))
	}
}

func (suite *AuthSuite) TestGetOrCreateUser() {
	t := suite.T()

	// When: login gov UUID is passed to create user func
	userData := map[string]interface{}{}
	userData["sub"] = "39b28c92-0506-4bef-8b57-e39519f42dc2"
	userData["email"] = "sally@government.gov"
	loginGovUUID, _ := uuid.FromString(userData["sub"].(string))

	// And: user does not yet exist in the db
	newUser, err := getOrCreateUser(suite.db, userData)
	if err != nil {
		t.Error("error querying or creating user.")
	}

	// Then: expect fields to be set on returned user
	if newUser.LoginGovEmail != userData["email"] {
		t.Error("expected email to be set")
	}
	if newUser.LoginGovUUID != loginGovUUID {
		t.Error("expected uuid to be set")
	}

	// When: The same UUID is passed in func
	sameUser, err := getOrCreateUser(suite.db, userData)
	if err != nil {
		t.Error("error querying or creating user.")
	}

	// Then: expect the existing user to be returned
	if sameUser.LoginGovEmail != newUser.LoginGovEmail {
		t.Error("expected existing user to have been returned")
	}

	// And: no new user to have been created
	query := suite.db.Where("login_gov_uuid = ?", loginGovUUID)
	var users []models.User
	queryErr := query.All(&users)
	if queryErr != nil {
		t.Error("DB Query Error", zap.Error(err))
	}
	if len(users) > 1 {
		t.Error("1 user should have been returned")
	}
}
