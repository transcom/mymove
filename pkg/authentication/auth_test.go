package authentication

import (
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/server"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

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
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		log.Panic(err)
	}
	if verrs.Count() > 0 {
		t.Fatalf("errors encountered saving %v: %v", model, verrs)
	}
}

func TestAuthSuite(t *testing.T) {
	configLocation := "../../../config"
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

func fakeLoginGovProvider(logger *zap.Logger) *LoginGovProvider {
	return &LoginGovProvider{
		"fakeHostname",
		"secret_key",
		logger,
	}
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
	hostsConfig := server.HostsConfig{
		MyName:     "my.move.host",
		OfficeName: "office.move.host",
		TspName:    "tsp.move.host",
	}
	loginGovConfig := LoginGovConfig{
		Host:             "login.gov",
		CallbackProtocol: "https://",
		CallbackPort:     "1234",
	}
	responsePattern := regexp.MustCompile(`href="(.+)"`)

	req, err := http.NewRequest("GET", "/auth/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Host = hostsConfig.OfficeName
	session := server.Session{UserID: fakeUUID, IDToken: fakeToken}
	ctx := server.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	authContext := NewAuthContext(&loginGovConfig, fakeLoginGovProvider(suite.logger), suite.logger)
	handler := NewLogoutHandler(authContext, &server.SessionCookieConfig{Secret: "fake key", NoTimeout: false})
	wrappedHandler := server.NewAppDetectorMiddleware(&hostsConfig, suite.logger)(handler)
	rr := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusTemporaryRedirect {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusTemporaryRedirect)
	}

	redirectURL, err := url.Parse(responsePattern.FindStringSubmatch(rr.Body.String())[1])
	if err != nil {
		t.Fatal(err)
	}
	params := redirectURL.Query()

	postRedirectURI, err := url.Parse(params["post_logout_redirect_uri"][0])

	suite.Nil(err)
	suite.logger.Info(postRedirectURI.String())
	suite.Equal(hostsConfig.OfficeName, postRedirectURI.Hostname())
	suite.Equal(loginGovConfig.CallbackPort, postRedirectURI.Port())
	token := params["id_token_hint"][0]
	suite.Equal(fakeToken, token, "handler id_token")
}

func (suite *AuthSuite) TestRequireAuthMiddleware() {
	// Given: a logged in user
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  loginGovUUID,
		LoginGovEmail: "email@example.com",
	}
	suite.mustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	// And: the context contains the auth values
	session := server.Session{UserID: user.ID, IDToken: "fake Token"}
	ctx := server.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	var handlerSession *server.Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerSession = server.SessionFromRequestContext(r)
	})
	middleware := NewUserAuthMiddleware(suite.logger)(handler)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're logged in
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, user.ID, "the authenticated user is different from expected")
}

func (suite *AuthSuite) TestRequireAuthMiddlewareUnauthorized() {
	t := suite.T()

	// Given: No logged in users
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := NewUserAuthMiddleware(suite.logger)(handler)

	middleware.ServeHTTP(rr, req)

	// We should receive an unauthorized response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusUnauthorized)
	}
}
