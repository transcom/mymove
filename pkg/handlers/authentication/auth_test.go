package authentication

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ghcapi"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// UserSessionCookieName is the key suffix at which we're storing our token cookie
const UserSessionCookieName = "session_token"

// SessionCookieName returns the session cookie name
func SessionCookieName(session *auth.Session) string {
	return fmt.Sprintf("%s_%s", string(session.ApplicationName), UserSessionCookieName)
}

type AuthSuite struct {
	handlers.BaseHandlerTestSuite
	callbackPort int
}

func (suite *AuthSuite) SetupTest() {
	gob.Register(auth.Session{})
	suite.callbackPort = 1234
}

// AuthContext returns a testing auth context
func (suite *AuthSuite) AuthContext() Context {
	return NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()),
		"http", suite.callbackPort)
}

func TestAuthSuite(t *testing.T) {
	hs := &AuthSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func fakeLoginGovProvider(logger *zap.Logger) LoginGovProvider {
	return NewLoginGovProvider("fakeHostname", "secret_key", logger)
}

func (suite *AuthSuite) SetupSessionRequest(r *http.Request, session *auth.Session, sessionManager auth.SessionManager) *http.Request {
	ctx, err := sessionManager.Load(r.Context(), session.IDToken)
	suite.NoError(err)
	_, _, err = sessionManager.Commit(ctx)
	suite.NoError(err)
	sessionManager.Put(ctx, "session", session)
	ctx = auth.SetSessionInRequestContext(r.WithContext(ctx), session)
	return r.WithContext(ctx)
}

func setUpMockNotificationSender() notifications.NotificationSender {
	// We need a NotificationSender for sending user activity emails to system admins.
	// If an unknown Service Member tries to log in, we'll create a new account for them, and that's requires an email.
	// This function allows us to set up a fresh mock for each test so we can check the number of calls it has.
	mockSender := mocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.UserAccountModified"),
	).Return(nil)

	return &mockSender
}

func (suite *AuthSuite) TestGenerateNonce() {
	t := suite.T()
	nonce := generateNonce()

	if (nonce == "") || (len(nonce) < 1) {
		t.Error("No nonce was returned.")
	}
}

func (suite *AuthSuite) TestAuthorizationLogoutHandler() {
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	fakeToken := "some_token"

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	baseReq := httptest.NewRequest("POST", fmt.Sprintf("http://%s/auth/logout", appnames.OfficeServername), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
	}
	sessionManagers := handlerConfig.SessionManagers()
	officeSession := sessionManagers.Office
	authContext := suite.AuthContext()
	fakeProvider := openidConnect.Provider{
		ClientKey: "some_token",
	}
	fakeProvider.SetName("officeProvider")
	goth.UseProviders(&fakeProvider)

	handler := officeSession.LoadAndSave(NewLogoutHandler(authContext, handlerConfig))

	rr := httptest.NewRecorder()
	req := suite.SetupSessionRequest(baseReq, &session, sessionManagers.Office)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	redirectURL, err := url.Parse(rr.Body.String())
	suite.FatalNoError(err)
	params := redirectURL.Query()

	postRedirectURI, err := url.Parse(params["post_logout_redirect_uri"][0])
	suite.NoError(err)
	suite.Equal(appnames.OfficeServername, postRedirectURI.Hostname())
	suite.Equal(strconv.Itoa(suite.callbackPort), postRedirectURI.Port())
	token := params["client_id"][0]
	suite.Equal(fakeToken, token, "handler id_token")

	noIDTokenSession := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         "",
		Hostname:        appnames.OfficeServername,
	}

	req = suite.SetupSessionRequest(baseReq, &noIDTokenSession, sessionManagers.Office)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	redirectURI, err := url.Parse(rr.Body.String())
	suite.NoError(err)
	suite.Equal(appnames.OfficeServername, redirectURI.Hostname())
	suite.Equal(strconv.Itoa(suite.callbackPort), redirectURI.Port())
}

func (suite *AuthSuite) TestRequireAuthMiddleware() {
	// Given: a logged in user
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token", ApplicationName: "mil"}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)
	cookieName := SessionCookieName(&session)
	cookie := http.Cookie{
		Name:  cookieName,
		Value: "some randomly generated string",
		Path:  "/",
	}
	req.AddCookie(&cookie)

	var handlerSession *auth.Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerSession = auth.SessionFromRequestContext(r)
	})
	sessionManager := scs.New()
	middleware := sessionManager.LoadAndSave(UserAuthMiddleware(suite.Logger())(handler))

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're logged in
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, user.ID, "the authenticated user is different from expected")
}

// Test permissions middleware with a user who will be ALLOWED POST access on the endpoint: ghc/v1/shipments/:shipmentID/approve
// role must have update.shipment permissions
func (suite *AuthSuite) TestRequirePermissionsMiddlewareAuthorized() {
	// TIO users have the proper permissions for our test - update.shipment
	tioOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})

	identity, err := models.FetchUserIdentity(suite.DB(), tioOfficeUser.User.LoginGovUUID.String())

	suite.NoError(err)

	rr := httptest.NewRecorder()
	// using an arbitrary ID here for the shipment
	req := httptest.NewRequest("POST", "/ghc/v1/shipments/123456/approve", nil)

	// And: the context contains the auth values
	handlerSession := auth.Session{
		UserID:          tioOfficeUser.User.ID,
		IDToken:         "fake Token",
		ApplicationName: "mil",
	}

	handlerSession.Roles = append(handlerSession.Roles, identity.Roles...)

	ctx := auth.SetSessionInRequestContext(req, &handlerSession)
	req = req.WithContext(ctx)

	handlerConfig := suite.HandlerConfig()
	api := ghcapi.NewGhcAPIHandler(handlerConfig)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	middleware := PermissionsMiddleware(suite.AppContextForTest(), api)

	root := mux.NewRouter()
	ghcMux := root.PathPrefix("/ghc/v1/").Subrouter()
	ghcMux.PathPrefix("/").Handler(api.Serve(middleware))

	middleware(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, tioOfficeUser.User.ID, "the authenticated user is different from expected")
}

// Test permissions middleware with a user who will be DENIED POST access on the endpoint: ghc/v1/shipments/:shipmentID/approve
// role must NOT have update.shipment permissions
func (suite *AuthSuite) TestRequirePermissionsMiddlewareUnauthorized() {
	// QAECSR users will be denied access as they lack the proper permissions for our test - update.shipment
	qaeCsrOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQaeCsr})

	identity, err := models.FetchUserIdentity(suite.DB(), qaeCsrOfficeUser.User.LoginGovUUID.String())

	suite.NoError(err)

	rr := httptest.NewRecorder()
	// using an arbitrary ID here for the shipment
	req := httptest.NewRequest("POST", "/ghc/v1/shipments/123456/approve", nil)

	// And: the context contains the auth values
	handlerSession := auth.Session{
		UserID:          qaeCsrOfficeUser.User.ID,
		IDToken:         "fake Token",
		ApplicationName: "mil",
	}

	handlerSession.Roles = append(handlerSession.Roles, identity.Roles...)

	ctx := auth.SetSessionInRequestContext(req, &handlerSession)
	req = req.WithContext(ctx)

	handlerConfig := suite.HandlerConfig()
	api := ghcapi.NewGhcAPIHandler(handlerConfig)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	middleware := PermissionsMiddleware(suite.AppContextForTest(), api)

	root := mux.NewRouter()
	ghcMux := root.PathPrefix("/ghc/v1/").Subrouter()
	ghcMux.PathPrefix("/").Handler(api.Serve(middleware))

	middleware(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, qaeCsrOfficeUser.User.ID, "the authenticated user is different from expected")
}

func (suite *AuthSuite) TestIsLoggedInWhenNoUserLoggedIn() {
	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	rr := httptest.NewRecorder()
	sessionManager := scs.New()
	handler := sessionManager.LoadAndSave(IsLoggedInMiddleware(suite.Logger()))

	handler.ServeHTTP(rr, req)

	// expects to return 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

	// expects to return that no one is logged in
	expected := "{\"isLoggedIn\":false}\n"
	suite.Equal(expected, rr.Body.String(), "handler returned wrong body")
}

func (suite *AuthSuite) TestIsLoggedInWhenUserLoggedIn() {
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	sessionManager := scs.New()
	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token"}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := sessionManager.LoadAndSave(IsLoggedInMiddleware(suite.Logger()))

	handler.ServeHTTP(rr, req)

	// expects to return 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

	// expects to return that no one is logged in
	expected := "{\"isLoggedIn\":true}\n"
	suite.Equal(expected, rr.Body.String(), "handler returned wrong body")
}

func (suite *AuthSuite) TestRequireAuthMiddlewareUnauthorized() {
	t := suite.T()

	// Given: No logged in users
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/moves", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	sessionManager := scs.New()
	middleware := sessionManager.LoadAndSave(UserAuthMiddleware(suite.Logger())(handler))

	middleware.ServeHTTP(rr, req)

	// We should receive an unauthorized response
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusUnauthorized)
	}
}

func (suite *AuthSuite) TestRequireAdminAuthMiddleware() {
	// Given: a logged in user
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/v1/office-users", nil)

	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token", AdminUserID: uuid.Must(uuid.NewV4())}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	var handlerSession *auth.Session
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerSession = auth.SessionFromRequestContext(r)
	})

	middleware := AdminAuthMiddleware(suite.Logger())(handler)

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're logged in
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, user.ID, "the authenticated user is different from expected")
}

func (suite *AuthSuite) TestRequireAdminAuthMiddlewareUnauthorized() {
	t := suite.T()

	// Given: No logged in users
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/v1/office-users", nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := AdminAuthMiddleware(suite.Logger())(handler)

	middleware.ServeHTTP(rr, req)

	// We should receive an unauthorized response
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v wanted %v", status, http.StatusUnauthorized)
	}
}

func (suite *AuthSuite) TestAuthorizeDeactivateUser() {
	userIdentity := models.UserIdentity{
		Active: false,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", appnames.OfficeServername), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           "deactivated@example.com",
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Office)
	authContext := suite.AuthContext()

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "authorizer did not recognize deactivated user")
}

func (suite *AuthSuite) TestAuthKnownSingleRoleOffice() {
	officeUserID := uuid.Must(uuid.NewV4())
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	userIdentity := models.UserIdentity{
		ID:           user.ID,
		Active:       true,
		OfficeUserID: &officeUserID,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", appnames.OfficeServername), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
	}
	authContext := suite.AuthContext()

	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Office)

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	// Office app, so should only have office ID information
	suite.Equal(officeUserID, session.OfficeUserID)
}

func (suite *AuthSuite) TestAuthorizeDeactivateOfficeUser() {
	officeActive := false
	userIdentity := models.UserIdentity{
		OfficeActive: &officeActive,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", appnames.OfficeServername), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           "deactivated@example.com",
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Office)
	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}

	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "authorizer did not recognize deactivated office user")
}

func (suite *AuthSuite) TestRedirectLoginGovErrorMsg() {
	officeUserID := uuid.Must(uuid.NewV4())
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	userIdentity := models.UserIdentity{
		ID:           user.ID,
		Active:       true,
		OfficeUserID: &officeUserID,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", appnames.OfficeServername), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
	}
	// login.gov state cookie
	cookieName := StateCookieName(&session)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   "some mis-matched hash value",
		Path:    "/",
		Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
	}
	req.AddCookie(&cookie)

	authContext := suite.AuthContext()

	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Office)

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	rr2 := httptest.NewRecorder()
	sessionManagers.Office.LoadAndSave(h).ServeHTTP(rr2, req)

	// Office app, so should only have office ID information
	suite.Equal(officeUserID, session.OfficeUserID)

	suite.Equal(2, len(rr2.Result().Cookies()))
	// check for blank value for cookie login gov state value and the session cookie value
	for _, cookie := range rr2.Result().Cookies() {
		if cookie.Name == cookieName || cookie.Name == "office_session_token" {
			suite.Equal("", cookie.Value)
			suite.Equal("/", cookie.Path)
		}
	}

	suite.Equal("http://office.example.com:1234/?error=SIGNIN_ERROR", rr2.Result().Header.Get("Location"))
}

func (suite *AuthSuite) TestAuthKnownSingleRoleAdmin() {
	adminUserID := uuid.Must(uuid.NewV4())
	officeUserID := uuid.Must(uuid.NewV4())
	var adminUserRole models.AdminRole = "SYSTEM_ADMIN"
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	userIdentity := models.UserIdentity{
		ID:            user.ID,
		Active:        true,
		OfficeUserID:  &officeUserID,
		AdminUserID:   &adminUserID,
		AdminUserRole: &adminUserRole,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", appnames.AdminServername), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
	}

	authContext := suite.AuthContext()

	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Admin)

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	// admin app, so should only have admin ID information
	suite.Equal(userIdentity.ID, session.UserID)
	suite.Equal(adminUserID, session.AdminUserID)
	suite.Equal(uuid.Nil, session.OfficeUserID)
	suite.True(session.IsAdminUser())
	suite.True(session.IsSystemAdmin())
	suite.False(session.IsProgramAdmin())
}

func (suite *AuthSuite) TestAuthKnownServiceMember() {
	user := factory.BuildDefaultUser(suite.DB())
	userID := uuid.Must(uuid.NewV4())

	userIdentity := models.UserIdentity{
		ID:              user.ID,
		ServiceMemberID: &userID,
		Active:          true,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	baseReq := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", appnames.MilServername), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        appnames.MilServername,
	}

	authContext := suite.AuthContext()

	sessionManagers := handlerConfig.SessionManagers()
	req := suite.SetupSessionRequest(baseReq, &session, sessionManagers.Mil)

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	foundUser, _ := models.GetUser(suite.DB(), user.ID)

	suite.NotEqual("", foundUser.CurrentMilSessionID)

	sessionStore := sessionManagers.Mil.Store()
	_, existsBefore, _ := sessionStore.Find(foundUser.CurrentMilSessionID)
	suite.Equal(existsBefore, true)

	concurrentSession := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        appnames.MilServername,
	}
	concurentReq := suite.SetupSessionRequest(baseReq, &concurrentSession, sessionManagers.Mil)
	authorizeKnownUser(suite.AppContextWithSessionForTest(&concurrentSession), &userIdentity, h, rr, concurentReq, "")

	_, existsAfterConcurrentSession, _ := sessionStore.Find(foundUser.CurrentMilSessionID)
	suite.Equal(existsAfterConcurrentSession, false)
}

// TESTCASE SCENARIO
// What is being tested: authorizeUnknownUser function
// Mocked: LoginGovProvider, auth.Session, goth.User, scs.SessionManager
// Behaviour: The function gets passed in the following arguments:
// - an instance of goth.User: a struct with the login.gov UUID and email
// - the callback handler
// - the session (instance of auth.Session)
// - the http ResponseWriter
// - the http Request with a context that includes the session
// - the landing URL string (where to redirect the user after successful auth)
// It should create the user using the login.gov UUID and email, then create a
// service member associated with the user, and populate the session with the ID
// of the service member in the `ServiceMemberID` key.
func (suite *AuthSuite) TestAuthUnknownServiceMember() {
	// Set up: Prepare the session, goth.User, callback handler, http response
	//         and request, landing URL, and pass them into authorizeUnknownUser

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	// Prepare the session and session manager
	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        appnames.MilServername,
	}
	// Prepare the request and set the session in the request context
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", appnames.MilServername), nil)
	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Mil)

	// Prepare the callback handler
	authContext := suite.AuthContext()

	mockSender := setUpMockNotificationSender() // We should get an email for this activity
	h := CallbackHandler{
		authContext,
		handlerConfig,
		mockSender,
	}

	// Prepare the request and response writer
	rr := httptest.NewRecorder()

	// Prepare the goth.User to simulate the UUID and email that login.gov would
	// provide
	fakeUUID, _ := uuid.NewV4()
	user := goth.User{
		UserID: fakeUUID.String(),
		Email:  "new_service_member@example.com",
	}

	// Call the function under test
	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr,
		req, h.landingURL(&session))
	mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)

	// Look up the user and service member in the test DB
	foundUser, _ := models.GetUserFromEmail(suite.DB(), user.Email)
	serviceMemberID := session.ServiceMemberID
	serviceMember, _ := models.FetchServiceMemberForUser(suite.DB(), &session, serviceMemberID)
	// Look up the session token in the session store (this test uses the memory store)
	sessionStore := sessionManagers.Mil.Store()
	_, existsBefore, _ := sessionStore.Find(foundUser.CurrentMilSessionID)

	// Verify service member exists and its ID is populated in the session
	suite.NotEmpty(session.ServiceMemberID)

	// Verify session contains UserID that points to the newly-created user
	suite.Equal(foundUser.ID, session.UserID)

	// Verify user's LoginGovEmail and LoginGovUUID match the values passed in
	suite.Equal(user.Email, foundUser.LoginGovEmail)
	suite.Equal(user.UserID, foundUser.LoginGovUUID.String())

	// Verify that the user's CurrentMilSessionID is not empty. The value is
	// generated randomly, so we can't test for a specific string. Any string
	// except an empty string is acceptable.
	suite.NotEqual("", foundUser.CurrentMilSessionID)

	// Verify the session token also exists in the session store
	suite.Equal(true, existsBefore)

	// Verify the service member that was created is associated with the user
	// that was created
	suite.Equal(foundUser.ID, serviceMember.UserID)

	// Verify handler redirects to landing URL
	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "handler did not redirect")
	suite.Equal(fmt.Sprintf("http://%s:1234/", appnames.MilServername), rr.Result().Header.Get("Location"))
}

func (suite *AuthSuite) TestAuthorizeDeactivateAdmin() {
	adminUserActive := false
	userIdentity := models.UserIdentity{
		AdminUserActive: &adminUserActive,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", appnames.AdminServername), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           "deactivated@example.com",
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Admin)
	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "authorizer did not recognize deactivated admin user")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeDeactivated() {
	// deactivated office user exists, but user has never logged it (and therefore first need to create a new user).

	// Create office user with no user
	office := factory.BuildTransportationOffice(suite.DB(), nil, nil)
	officeUser := models.OfficeUser{
		TransportationOffice:   office,
		TransportationOfficeID: office.ID,
		FirstName:              "Leo",
		LastName:               "Spaceman",
		Email:                  "leospaceman12345@example.com",
		Telephone:              "415-555-1212",
		Active:                 false,
	}
	verrs, err := suite.DB().ValidateAndCreate(&officeUser)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", appnames.OfficeServername), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Office)

	fakeUUID2, _ := uuid.NewV4()
	user := goth.User{
		UserID: fakeUUID2.String(),
		Email:  officeUser.Email,
	}

	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "Office user is active")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeNotFound() {

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", appnames.OfficeServername), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           "missing@email.com",
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Office)

	id, _ := uuid.NewV4()
	user := goth.User{
		UserID: id.String(),
		Email:  "sample@email.com",
	}

	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "Office user not found")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeLogsIn() {
	user := factory.BuildDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Active: true,
				UserID: &user.ID,
				Email:  user.LoginGovEmail,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", appnames.OfficeServername), nil)
	fakeToken := "some_token"

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}

	gothUser := goth.User{
		UserID: user.ID.String(),
		Email:  officeUser.Email,
	}

	authContext := suite.AuthContext()

	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Office)

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), gothUser, h,
		rr, req, "")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), officeUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(officeUser.ID, session.OfficeUserID)
	suite.Equal(uuid.Nil, session.AdminUserID)
	suite.NotEqual("", foundUser.CurrentOfficeSessionID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeLogsInWithPermissions() {
	user := factory.BuildDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Active: true,
				UserID: &user.ID,
				Email:  user.LoginGovEmail,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, []roles.RoleType{roles.RoleTypeQaeCsr})

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", appnames.OfficeServername), nil)
	fakeToken := "some_token"

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}
	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Office)
	gothUser := goth.User{
		UserID: user.ID.String(),
		Email:  officeUser.Email,
	}

	authContext := suite.AuthContext()

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), gothUser, h, rr, req, "")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), officeUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(officeUser.ID, session.OfficeUserID)
	suite.Equal(uuid.Nil, session.AdminUserID)
	suite.NotEqual("", foundUser.CurrentOfficeSessionID)
	// Make sure session contains roles and permissions
	suite.NotEmpty(session.Roles)
	userRole, hasRole := officeUser.User.Roles.GetRole(roles.RoleTypeQaeCsr)
	suite.True(hasRole)
	sessionRole, hasRole := session.Roles.GetRole(roles.RoleTypeQaeCsr)
	suite.True(hasRole)
	suite.Equal(userRole.ID, sessionRole.ID)
	suite.NotEmpty(session.Permissions)
	suite.ElementsMatch(QAECSR.Permissions, session.Permissions)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminDeactivated() {
	// Create an admin user that is inactive and has never logged into the app
	adminUser := models.AdminUser{
		FirstName: "Leo",
		LastName:  "Spaceman",
		Email:     "leo_spaceman_admin@example.com",
		Role:      "SYSTEM_ADMIN",
	}
	verrs, err := suite.DB().ValidateAndCreate(&adminUser)
	suite.NoError(err)
	suite.False(verrs.HasAny())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", appnames.AdminServername), nil)
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		Hostname:        appnames.AdminServername,
		Email:           adminUser.Email,
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Admin)

	fakeUUID2, _ := uuid.NewV4()
	user := goth.User{
		UserID: fakeUUID2.String(),
		Email:  adminUser.Email,
	}

	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "Admin user is active")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminNotFound() {
	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	// user not admin_users and has never logged into the app
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", appnames.AdminServername), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           "missing@email.com",
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Admin)

	id, _ := uuid.NewV4()
	user := goth.User{
		UserID: id.String(),
		Email:  "sample@email.com",
	}

	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "Admin user not found")
}

func (suite *AuthSuite) TestAuthorizeKnownUserAdminNotFound() {
	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	// user exists in the DB, but not as an admin user
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/login-gov", appnames.AdminServername), nil)
	fakeToken := "some_token"
	loginGovUUID := uuid.Must(uuid.NewV4())
	userID := uuid.Must(uuid.NewV4())
	serviceMemberID := uuid.Must(uuid.NewV4())

	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
		ID:            userID,
	}
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           user.LoginGovEmail,
	}
	req = suite.SetupSessionRequest(req, &session, handlerConfig.SessionManagers().Admin)

	userIdentity := models.UserIdentity{
		ID:              user.ID,
		Active:          true,
		ServiceMemberID: &serviceMemberID,
	}

	authContext := suite.AuthContext()
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req, "")

	suite.Equal(http.StatusTemporaryRedirect, rr.Code, "Admin user not found")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminLogsIn() {
	// user is in admin_users but has not logged into the app before
	adminUser := factory.BuildAdminUser(suite.DB(), []factory.Customization{
		{
			Model: models.AdminUser{
				Active: true,
			},
		},
	}, []factory.Trait{
		factory.GetTraitActiveUser,
		factory.GetTraitAdminUserEmail,
	})
	user := adminUser.User

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", appnames.AdminServername), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           adminUser.Email,
	}

	gothUser := goth.User{
		UserID: user.ID.String(),
		Email:  adminUser.Email,
	}

	authContext := suite.AuthContext()

	sessionManagers := handlerConfig.SessionManagers()
	req = suite.SetupSessionRequest(req, &session, sessionManagers.Admin)

	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), gothUser, h, rr, req, "")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), adminUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(adminUser.ID, session.AdminUserID)
	suite.Equal(uuid.Nil, session.OfficeUserID)
	suite.NotEqual("", foundUser.CurrentAdminSessionID)
}

func (suite *AuthSuite) TestLoginGovAuthenticatedRedirect() {
	user := factory.BuildDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Active: true,
				UserID: &user.ID,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	fakeToken := "some_token"

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov", appnames.OfficeServername), nil)
	ctx := auth.SetSessionInRequestContext(req, &session)
	authContext := suite.AuthContext()
	h := RedirectHandler{
		authContext,
		handlerConfig,
		false,
	}
	rr := httptest.NewRecorder()
	req = req.WithContext(ctx)
	h.ServeHTTP(rr, req)

	suite.Equal(http.StatusTemporaryRedirect, rr.Code,
		"handler returned wrong status code")
	redirectURL, err := rr.Result().Location()
	suite.NoError(err)
	suite.Equal(appnames.OfficeServername, redirectURL.Hostname())
	suite.Equal(strconv.Itoa(suite.callbackPort), redirectURL.Port())
	suite.Equal("/", redirectURL.EscapedPath())
}

func (suite *AuthSuite) TestAuthorizePrime() {
	user := factory.BuildDefaultUser(suite.DB())
	clientCert := testdatagen.MakeDevClientCert(suite.DB(), testdatagen.Assertions{
		ClientCert: models.ClientCert{
			UserID: user.ID,
		},
	})

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/prime/v1", appnames.PrimeServername), nil)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	middleware := PrimeAuthorizationMiddleware(suite.Logger())(handler)
	rr := httptest.NewRecorder()

	ctx := SetClientCertInRequestContext(req, &clientCert)
	req = req.WithContext(ctx)
	middleware.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	// no cert in request
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("http://%s/prime/v1", appnames.PrimeServername), nil)
	middleware.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}
