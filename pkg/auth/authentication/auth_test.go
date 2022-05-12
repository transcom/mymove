package authentication

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/notifications/mocks"

	"github.com/transcom/mymove/pkg/notifications"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gofrs/uuid"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const (
	// OfficeTestHost
	OfficeTestHost string = "office.example.com"
	// MilTestHost
	MilTestHost string = "mil.example.com"
	// OrdersTestHost
	OrdersTestHost string = "orders.example.com"
	// DpsTestHost
	DpsTestHost string = "dps.example.com"
	// SddcTestHost
	SddcTestHost string = "sddc.example.com"
	// AdminTestHost
	AdminTestHost string = "admin.example.com"
)

// UserSessionCookieName is the key suffix at which we're storing our token cookie
const UserSessionCookieName = "session_token"

// SessionCookieName returns the session cookie name
func SessionCookieName(session *auth.Session) string {
	return fmt.Sprintf("%s_%s", string(session.ApplicationName), UserSessionCookieName)
}

// ApplicationTestServername is a collection of the test servernames
func ApplicationTestServername() auth.ApplicationServername {
	appnames := auth.ApplicationServername{
		MilServername:    MilTestHost,
		OfficeServername: OfficeTestHost,
		OrdersServername: OrdersTestHost,
		DpsServername:    DpsTestHost,
		SddcServername:   SddcTestHost,
		AdminServername:  AdminTestHost,
	}
	return appnames
}

type AuthSuite struct {
	testingsuite.PopTestSuite
}

func (suite *AuthSuite) SetupTest() {
	gob.Register(auth.Session{})
}

func TestAuthSuite(t *testing.T) {
	hs := &AuthSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func fakeLoginGovProvider(logger *zap.Logger) LoginGovProvider {
	return NewLoginGovProvider("fakeHostname", "secret_key", logger)
}

func setupScsSession(ctx context.Context, session *auth.Session, sessionManager *scs.SessionManager) context.Context {
	values := make(map[string]interface{})
	values["session"] = session
	expiry := time.Now().Add(30 * time.Minute).UTC()
	b, _ := sessionManager.Codec.Encode(expiry, values)

	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
	//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
	//RA: in which this would be considered a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	sessionManager.Store.Commit("session_token", b, expiry)
	scsContext, _ := sessionManager.Load(ctx, "session_token")
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
	//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
	//RA: in which this would be considered a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: Mitigated
	//RA Modified Severity: N/A
	// nolint:errcheck
	sessionManager.Commit(scsContext)
	return scsContext
}

func setupSessionManagers() [3]*scs.SessionManager {
	var milSession, adminSession, officeSession *scs.SessionManager
	store := memstore.New()
	milSession = scs.New()
	milSession.Store = store
	milSession.Cookie.Name = "mil_session_token"

	adminSession = scs.New()
	adminSession.Store = store
	adminSession.Cookie.Name = "admin_session_token"

	officeSession = scs.New()
	officeSession.Store = store
	officeSession.Cookie.Name = "office_session_token"

	return [3]*scs.SessionManager{milSession, adminSession, officeSession}
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
	t := suite.T()
	loginGovUUID, _ := uuid.FromString("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "email@example.com",
		Active:        true,
	}
	suite.MustSave(&user)

	fakeToken := "some_token"
	callbackPort := 1234

	req := httptest.NewRequest("POST", fmt.Sprintf("http://%s/auth/logout", OfficeTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)
	sessionManagers := setupSessionManagers()
	officeSession := sessionManagers[2]
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	handler := officeSession.LoadAndSave(NewLogoutHandler(authContext, handlerConfig))

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
	suite.NoError(err)
	suite.Equal(OfficeTestHost, postRedirectURI.Hostname())
	suite.Equal(strconv.Itoa(callbackPort), postRedirectURI.Port())
	token := params["id_token_hint"][0]
	suite.Equal(fakeToken, token, "handler id_token")
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
	req := httptest.NewRequest("GET", "/admin/v1/office_users", nil)

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
	req := httptest.NewRequest("GET", "/admin/v1/office_users", nil)

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

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", OfficeTestHost), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           "deactivated@example.com",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "authorizer did not recognize deactivated user")
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

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", OfficeTestHost), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	officeSession := sessionManagers[2]
	scsContext := setupScsSession(ctx, &session, officeSession)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(scsContext), "")

	// Office app, so should only have office ID information
	suite.Equal(officeUserID, session.OfficeUserID)
}

func (suite *AuthSuite) TestAuthorizeDeactivateOfficeUser() {
	officeActive := false
	userIdentity := models.UserIdentity{
		OfficeActive: &officeActive,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", OfficeTestHost), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           "deactivated@example.com",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "authorizer did not recognize deactivated office user")
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

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", OfficeTestHost), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
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

	ctx := auth.SetSessionInRequestContext(req, &session)

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	officeSession := sessionManagers[2]
	scsContext := setupScsSession(ctx, &session, officeSession)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(scsContext), "")

	rr2 := httptest.NewRecorder()
	officeSession.LoadAndSave(h).ServeHTTP(rr2, req.WithContext(scsContext))

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

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", AdminTestHost), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		IDToken:         fakeToken,
		Hostname:        AdminTestHost,
	}

	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	adminSession := sessionManagers[1]
	scsContext := setupScsSession(ctx, &session, adminSession)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(scsContext), "")

	// admin app, so should only have admin ID information
	suite.Equal(userIdentity.ID, session.UserID)
	suite.Equal(adminUserID, session.AdminUserID)
	suite.Equal(uuid.Nil, session.OfficeUserID)
	suite.True(session.IsAdminUser())
	suite.True(session.IsSystemAdmin())
	suite.False(session.IsProgramAdmin())
}

func (suite *AuthSuite) TestAuthKnownServiceMember() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	userID := uuid.Must(uuid.NewV4())

	userIdentity := models.UserIdentity{
		ID:              user.ID,
		ServiceMemberID: &userID,
		Active:          true,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/authorize", MilTestHost), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        MilTestHost,
	}

	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	milSession := sessionManagers[0]
	scsContext := setupScsSession(ctx, &session, milSession)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(scsContext), "")

	foundUser, _ := models.GetUser(suite.DB(), user.ID)

	suite.NotEqual("", foundUser.CurrentMilSessionID)

	sessionStore := milSession.Store
	_, existsBefore, _ := sessionStore.Find(foundUser.CurrentMilSessionID)
	suite.Equal(existsBefore, true)

	concurrentSession := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        MilTestHost,
	}
	concurrentCtx := auth.SetSessionInRequestContext(req, &concurrentSession)
	concurrentScsContext := setupScsSession(concurrentCtx, &concurrentSession, milSession)
	authorizeKnownUser(suite.AppContextWithSessionForTest(&concurrentSession), &userIdentity, h, rr, req.WithContext(concurrentScsContext), "")

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

	// Prepare the session and session manager
	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        MilTestHost,
	}
	sessionManagers := setupSessionManagers()
	milSession := sessionManagers[0]

	// Prepare the request and set the session in the request context
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", MilTestHost), nil)
	ctx := auth.SetSessionInRequestContext(req, &session)
	scsContext := setupScsSession(ctx, &session, milSession)

	// Prepare the callback handler
	callbackPort := 1234
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	mockSender := setUpMockNotificationSender() // We should get an email for this activity
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req.WithContext(scsContext), h.landingURL(&session))
	mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)

	// Look up the user and service member in the test DB
	foundUser, _ := models.GetUserFromEmail(suite.DB(), user.Email)
	serviceMemberID := session.ServiceMemberID
	serviceMember, _ := models.FetchServiceMemberForUser(suite.DB(), &session, serviceMemberID)
	// Look up the session token in the session store (this test uses the memory store)
	sessionStore := milSession.Store
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
	suite.Equal(fmt.Sprintf("http://%s:1234/", MilTestHost), rr.Result().Header.Get("Location"))
}

func (suite *AuthSuite) TestAuthorizeDeactivateAdmin() {
	adminUserActive := false
	userIdentity := models.UserIdentity{
		AdminUserActive: &adminUserActive,
	}

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", AdminTestHost), nil)

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        AdminTestHost,
		Email:           "deactivated@example.com",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()
	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "authorizer did not recognize deactivated admin user")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeDeactivated() {
	// deactivated office user exists, but user has never logged it (and therefore first need to create a new user).
	officeUser := testdatagen.MakeOfficeUserWithNoUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			Active: false,
		},
	})

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", OfficeTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Hostname:        OfficeTestHost,
		Email:           officeUser.Email,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	fakeUUID2, _ := uuid.NewV4()
	user := goth.User{
		UserID: fakeUUID2.String(),
		Email:  officeUser.Email,
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "Office user is active")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeNotFound() {

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", OfficeTestHost), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           "missing@email.com",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	id, _ := uuid.NewV4()
	user := goth.User{
		UserID: id.String(),
		Email:  "sample@email.com",
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "Office user not found")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeLogsIn() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			Active: true,
			UserID: &user.ID,
		},
		User: user,
	})

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov/callback", OfficeTestHost), nil)
	fakeToken := "some_token"

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           officeUser.Email,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	gothUser := goth.User{
		UserID: user.ID.String(),
		Email:  officeUser.Email,
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	officeSession := sessionManagers[2]
	scsContext := setupScsSession(ctx, &session, officeSession)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), gothUser, h, rr, req.WithContext(scsContext), "")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), officeUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(officeUser.ID, session.OfficeUserID)
	suite.Equal(uuid.Nil, session.AdminUserID)
	suite.NotEqual("", foundUser.CurrentOfficeSessionID)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminDeactivated() {
	// user is in office_users but is inactive and has never logged into the app
	adminUser := testdatagen.MakeAdminUserWithNoUser(suite.DB(), testdatagen.Assertions{})

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", AdminTestHost), nil)
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		Hostname:        AdminTestHost,
		Email:           adminUser.Email,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	fakeUUID2, _ := uuid.NewV4()
	user := goth.User{
		UserID: fakeUUID2.String(),
		Email:  adminUser.Email,
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "Admin user is active")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminNotFound() {
	// user not admin_users and has never logged into the app
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", AdminTestHost), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        AdminTestHost,
		Email:           "missing@email.com",
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	id, _ := uuid.NewV4()
	user := goth.User{
		UserID: id.String(),
		Email:  "sample@email.com",
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), user, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "Admin user not found")
}

func (suite *AuthSuite) TestAuthorizeKnownUserAdminNotFound() {
	// user exists in the DB, but not as an admin user
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/login-gov", AdminTestHost), nil)
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
		Hostname:        AdminTestHost,
		Email:           user.LoginGovEmail,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	userIdentity := models.UserIdentity{
		ID:              user.ID,
		Active:          true,
		ServiceMemberID: &serviceMemberID,
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeKnownUser(suite.AppContextWithSessionForTest(&session), &userIdentity, h, rr, req.WithContext(ctx), "")

	suite.Equal(http.StatusForbidden, rr.Code, "Admin user not found")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminLogsIn() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	// user is in admin_users but has not logged into the app before
	adminUser := testdatagen.MakeAdminUser(suite.DB(), testdatagen.Assertions{
		AdminUser: models.AdminUser{
			Active: true,
			UserID: &user.ID,
		},
		User: user,
	})

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/auth/logout", AdminTestHost), nil)
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        AdminTestHost,
		Email:           adminUser.Email,
	}
	ctx := auth.SetSessionInRequestContext(req, &session)

	gothUser := goth.User{
		UserID: user.ID.String(),
		Email:  adminUser.Email,
	}

	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)

	adminSession := sessionManagers[1]
	scsContext := setupScsSession(ctx, &session, adminSession)

	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
	}
	rr := httptest.NewRecorder()

	authorizeUnknownUser(suite.AppContextWithSessionForTest(&session), gothUser, h, rr, req.WithContext(scsContext), "")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), adminUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(adminUser.ID, session.AdminUserID)
	suite.Equal(uuid.Nil, session.OfficeUserID)
	suite.NotEqual("", foundUser.CurrentAdminSessionID)
}

func (suite *AuthSuite) TestLoginGovAuthenticatedRedirect() {
	user := testdatagen.MakeDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	officeUser := testdatagen.MakeOfficeUser(suite.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			Active: true,
			UserID: &user.ID,
		},
		User: user,
	})

	fakeToken := "some_token"

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        OfficeTestHost,
		Email:           officeUser.Email,
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/login-gov", OfficeTestHost), nil)
	ctx := auth.SetSessionInRequestContext(req, &session)
	callbackPort := 1234
	sessionManagers := setupSessionManagers()
	authContext := NewAuthContext(suite.Logger(), fakeLoginGovProvider(suite.Logger()), "http", callbackPort, sessionManagers)
	handlerConfig := handlers.NewHandlerConfig(suite.DB(), suite.Logger())
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
}
