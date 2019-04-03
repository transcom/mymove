package authentication

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const (
	// TspTestHost
	TspTestHost string = "tsp.example.com"
	// OfficeTestHost
	OfficeTestHost string = "office.example.com"
	// MilTestHost
	MilTestHost string = "mil.example.com"
)

type AuthSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *AuthSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestAuthSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &AuthSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func fakeLoginGovProvider(logger Logger) LoginGovProvider {
	return NewLoginGovProvider("fakeHostname", "secret_key", logger)
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
	callbackPort := 1234

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/auth/logout", OfficeTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	authContext := NewAuthContext(suite.logger, fakeLoginGovProvider(suite.logger), "http", callbackPort)
	handler := LogoutHandler{authContext, "fake key", false, false}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %d wanted %d", status, http.StatusOK)
	}

	redirectURL, err := url.Parse(rr.Body.String())
	if err != nil {
		t.Fatal(err)
	}
	params := redirectURL.Query()

	postRedirectURI, err := url.Parse(params["post_logout_redirect_uri"][0])
	suite.Nil(err)
	suite.Equal(OfficeTestHost, postRedirectURI.Hostname())
	suite.Equal(strconv.Itoa(callbackPort), postRedirectURI.Port())
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
	suite.MustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token"}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	var handlerSession *auth.Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerSession = auth.SessionFromRequestContext(r)
	})
	middleware := UserAuthMiddleware(suite.logger)(handler)

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
	middleware := UserAuthMiddleware(suite.logger)(handler)

	middleware.ServeHTTP(rr, req)

	// We should receive an unauthorized response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusUnauthorized)
	}
}

func (suite *AuthSuite) TestAuthorizeDisableUser() {
	userIdentity := models.UserIdentity{
		Disabled: true,
	}

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	errStatus, err := authorizeSession(&session, &userIdentity)

	if suite.Error(err) {
		suite.Equal(http.StatusForbidden, errStatus, "authorizer did not recognize disabled user")
	}

}

func (suite *AuthSuite) TestAuthorizeNonOfficeUserForbidden() {
	userID := uuid.Must(uuid.NewV4())
	userIdentity := models.UserIdentity{
		Disabled:     true,
		ID:           userID,
		OfficeUserID: nil,
	}

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          userID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	errStatus, err := authorizeSession(&session, &userIdentity)

	if suite.Error(err) {
		suite.Equal(http.StatusForbidden, errStatus, "authorizer did not recognize disabled user")
		suite.Equal(uuid.Nil, session.OfficeUserID)
	}
}

func (suite *AuthSuite) TestAuthorizeNonTSPUserForbidden() {
	userID := uuid.Must(uuid.NewV4())
	userIdentity := models.UserIdentity{
		Disabled:  true,
		ID:        userID,
		TspUserID: nil,
	}

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          userID,
		IDToken:         fakeToken,
		Hostname:        TspTestHost,
	}
	errStatus, err := authorizeSession(&session, &userIdentity)

	if suite.Error(err) {
		suite.Equal(http.StatusForbidden, errStatus, "authorizer did not recognize disabled user")
		suite.Equal(uuid.Nil, session.TspUserID)
	}
}

func (suite *AuthSuite) TestAuthorizeAllowedUser() {
	userID := uuid.Must(uuid.NewV4())
	officeUserID := uuid.Must(uuid.NewV4())
	smID := uuid.Must(uuid.NewV4())
	dpsID := uuid.Must(uuid.NewV4())
	fName := "fname"
	mName := "mname"
	lName := "lname"
	userIdentity := models.UserIdentity{
		Disabled:               false,
		ID:                     userID,
		OfficeUserID:           &officeUserID,
		ServiceMemberFirstName: &fName,
		ServiceMemberMiddle:    &mName,
		ServiceMemberLastName:  &lName,
		ServiceMemberID:        &smID,
		DpsUserID:              &dpsID,
	}

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          userID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	_, err := authorizeSession(&session, &userIdentity)

	if suite.NoError(err) {
		suite.Equal(userID, session.UserID)
		suite.Equal(officeUserID, session.OfficeUserID)
		suite.Equal(smID, session.ServiceMemberID)
		suite.Equal(dpsID, session.DpsUserID)
		suite.Equal(fName, session.FirstName)
		suite.Equal(mName, session.Middle)
		suite.Equal(lName, session.LastName)
	}
}
