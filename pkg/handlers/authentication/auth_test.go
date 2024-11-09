package authentication

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jarcoal/httpmock"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/handlers/ghcapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/random"
	"github.com/transcom/mymove/pkg/testingsuite"
)

// UserSessionCookieName is the key suffix at which we're storing our token cookie
const UserSessionCookieName = "session_token"

// This is a dumy private key that has no use and is not reflective of any real keys utilized. This key was generated
// specifically for the purpose of testing.
// #nosec G101 not real key- only used for testing
const DummyRSAPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDQ62hDHRRAduSuUQDxixn61bbRLj9iBBmRG03rW3PNnkSzrcof
9ytnKY2LX2DAPaSr/1Em7fvqiovzVg43ElfFHJBrCskJqWLphifv6qoGX1pwsPA/
Rb+MBqftMU1Zq7UC9Eis/Sje2QGx7k02JoQy+R/EP/kQq1B0p4/qCtR73QIDAQAB
AoGBALVzP+LKZsR2frdHc2JWRgIti9KyMCqZFPuKk2pOy41SYKkNz/djXTcESAM8
m3NcFqGr5nfBSoKyQkrd+wqpy7+8X15MpClVErfUeowoOpaFQBr0E5Yf8WuzWXV2
Daex1aeA+69OAPmYEiVJD4qY6m8vxHZZT0ISNEIW4ObhyQmhAkEA9SvhOADfQLp8
7vZXTWW/fhapi8NiKW8cWT4wDQhwnW2glGxyVJBWwj+VtcJ8j5mEfm6vInh7QAYl
2dV/sMaNNwJBANolofvurHjd56WcdHENctAJTxiWtTqA9RIrtIIzJW7cqR4ujQKL
ndD5v2nG+b2JdlcOBzNs0LVF+ItwYMTYKYsCQQC1LqhR6tMR0r9hGUuLNxY86CKD
1vBEDoi0qvB3sTUIImv5Q+t58vEqvDK3D/Nda+YuST3EC6WJuwFd6hljWlghAkBL
s9mVywrxWtijoTrLbMZWKZTYTJyRs+TYLHCU6ljoMw1BWxg2NOtMdQ8XDyTlwIlf
xo97Khz3e1O4WARM61LnAkAzTxo/AOHVKawAR45eq4rjz0rxyCgtcTGa1qaEt9Ap
WjqcmKEkxqxz6lX/Pj2GbyikMkDThcp1bd1DRSUDOxHP
-----END RSA PRIVATE KEY-----
`
const DummyRSAPublicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDQ62hDHRRAduSuUQDxixn61bbR
Lj9iBBmRG03rW3PNnkSzrcof9ytnKY2LX2DAPaSr/1Em7fvqiovzVg43ElfFHJBr
CskJqWLphifv6qoGX1pwsPA/Rb+MBqftMU1Zq7UC9Eis/Sje2QGx7k02JoQy+R/E
P/kQq1B0p4/qCtR73QIDAQAB
-----END PUBLIC KEY-----
`

const DummyRSAModulus = "0OtoQx0UQHbkrlEA8YsZ-tW20S4_YgQZkRtN61tzzZ5Es63KH_crZymNi19gwD2kq_9RJu376oqL81YONxJXxRyQawrJCali6YYn7-qqBl9acLDwP0W_jAan7TFNWau1AvRIrP0o3tkBse5NNiaEMvkfxD_5EKtQdKeP6grUe90"
const jwtKeyID = "keyID"
const officeProviderName = "officeProvider"

// SessionCookieName returns the session cookie name
func SessionCookieName(session *auth.Session) string {
	return fmt.Sprintf("%s_%s", string(session.ApplicationName), UserSessionCookieName)
}

func generateNonce() string {
	nonceBytes := make([]byte, 64)
	//RA Summary: gosec - G404 - Insecure random number source (rand)
	//RA: gosec detected use of the insecure package math/rand rather than the more secure cryptographically secure pseudo-random number generator crypto/rand.
	//RA: This particular usage is mitigated by sourcing the seed from crypto/rand in order to create the new random number using math/rand.
	//RA Developer Status: Mitigated
	//RA Validator: jneuner@mitre.org
	//RA Validator Status: Mitigated
	//RA Modified Severity: CAT III
	// #nosec G404
	randomInt := rand.New(random.NewCryptoSeededSource())
	for i := 0; i < 64; i++ {
		nonceBytes[i] = byte(randomInt.Int63() % 256)
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
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
	return NewAuthContext(suite.Logger(), *fakeOktaProvider(suite.Logger()),
		"http", suite.callbackPort)
}

func (suite *AuthSuite) urlForHost(host string) *url.URL {
	var u url.URL
	u.Scheme = "http"
	u.Host = fmt.Sprintf("%s:%d", host, suite.callbackPort)
	u.Path = "/"
	return &u
}

func TestAuthSuite(t *testing.T) {
	hs := &AuthSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func fakeOktaProvider(logger *zap.Logger) *okta.Provider {
	return okta.NewOktaProvider(logger)
}

func (suite *AuthSuite) SetupSessionContext(ctx context.Context, session *auth.Session, sessionManager auth.SessionManager) context.Context {
	ctx, err := sessionManager.Load(ctx, session.IDToken)
	suite.NoError(err)
	_, _, err = sessionManager.Commit(ctx)
	suite.NoError(err)
	sessionManager.Put(ctx, "session", session)
	return ctx
}

func (suite *AuthSuite) SetupSessionRequest(r *http.Request, session *auth.Session, sessionManager auth.SessionManager) *http.Request {
	ctx := suite.SetupSessionContext(r.Context(), session, sessionManager)
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
	// this sets up a user with a valid ID and Access token
	// calls /auth/logout which should clear those tokens
	// checks to make sure those values are not present in session

	OktaID := "2400c3c5-019d-4031-9c27-8a553e022297"

	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
	}
	suite.MustSave(&user)

	fakeToken := "some_token"
	fakeAccessToken := "some_access_token"

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	baseReq := httptest.NewRequest("POST", fmt.Sprintf("http://%s/auth/logout", appnames.OfficeServername), nil)
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		AccessToken:     fakeAccessToken,
		Hostname:        appnames.OfficeServername,
	}
	sessionManagers := handlerConfig.SessionManagers()
	officeSession := sessionManagers.Office
	authContext := suite.AuthContext()

	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("officeProvider", "OrgURL", "CallbackURL", fakeToken, "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)
	handler := officeSession.LoadAndSave(NewLogoutHandler(authContext, handlerConfig))

	rr := httptest.NewRecorder()
	req := suite.SetupSessionRequest(baseReq, &session, sessionManagers.Office)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	noIDTokenSession := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         "",
		AccessToken:     "",
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

func (suite *AuthSuite) TestLogoutOktaRedirectHandler() {
	// this sets up a user with a valid ID and Access token
	// calls /auth/logoutOktaRedirect which should clear those tokens
	// checks to make sure those values are not present in session

	OktaID := "2400c3c5-019d-4031-9c27-8a553e022297"

	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
	}
	suite.MustSave(&user)

	fakeToken := "some_token"
	fakeAccessToken := "some_access_token"

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	baseReq := httptest.NewRequest("POST", fmt.Sprintf("http://%s/auth/logoutOktaRedirect", appnames.MilServername), nil)
	session := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		AccessToken:     fakeAccessToken,
		Hostname:        appnames.MilServername,
	}
	sessionManagers := handlerConfig.SessionManagers()
	milSession := sessionManagers.Mil
	authContext := suite.AuthContext()

	oktaProvider := okta.NewOktaProvider(suite.Logger())
	err := oktaProvider.RegisterOktaProvider("milProvider", "OrgURL", "CallbackURL", fakeToken, "secret", []string{"openid", "profile", "email"})
	suite.NoError(err)
	handler := milSession.LoadAndSave(NewLogoutOktaRedirectHandler(authContext, handlerConfig))

	rr := httptest.NewRecorder()
	req := suite.SetupSessionRequest(baseReq, &session, sessionManagers.Office)
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")

	// Read and parse the body to extract the URL
	body := rr.Body.String()
	parsedURL, err := url.Parse(body)
	suite.NoError(err)

	rawQuery := parsedURL.RawQuery
	values, err := url.ParseQuery(rawQuery)
	suite.NoError(err)

	// parsing the redirect url which should end in auth/okta
	// this redirects the user back to the sign in page for Okta and not the MilMove home page
	postLogoutRedirectURI := values.Get("post_logout_redirect_uri")
	suite.NotEmpty(postLogoutRedirectURI, "post_logout_redirect_uri is empty")
	suite.True(strings.HasSuffix(postLogoutRedirectURI, "/auth/okta"), "post_logout_redirect_uri does not end with /auth/okta")
}

func (suite *AuthSuite) TestRequireAuthMiddleware() {
	// Given: a logged in user
	OktaID := ("2400c3c5-019d-4031-9c27-8a553e022297")
	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
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
	handler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		handlerSession = auth.SessionFromRequestContext(r)
	})
	sessionManager := scs.New()
	middleware := sessionManager.LoadAndSave(UserAuthMiddleware(suite.Logger())(handler))

	middleware.ServeHTTP(rr, req)

	// We should be not be redirected since we're logged in
	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, user.ID, "the authenticated user is different from expected")
}

func (suite *AuthSuite) TestCustomerAPIAuthMiddleware() {
	setUpRequest := func(endpoint string, serviceMember *models.ServiceMember, officeUser *models.OfficeUser) *http.Request {
		req := httptest.NewRequest("GET", endpoint, nil)

		session := auth.Session{
			IDToken:         "fake Token",
			ApplicationName: auth.MilApp,
		}

		if serviceMember != nil {
			session.UserID = serviceMember.User.ID
			session.ServiceMemberID = serviceMember.ID
		} else if officeUser != nil {
			session.UserID = officeUser.User.ID
			session.OfficeUserID = officeUser.ID
		} else {
			suite.Fail("No user provided")
		}

		ctx := auth.SetSessionInRequestContext(req, &session)

		return req.WithContext(ctx)
	}

	setUpHandlerAndMiddleware := func() http.Handler {
		handlerConfig := suite.HandlerConfig()

		api := internalapi.NewInternalAPI(handlerConfig)

		handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

		customerAPIAuthMiddleware := CustomerAPIAuthMiddleware(suite.AppContextForTest(), api)

		root := chi.NewRouter()
		root.Mount("/internal", api.Serve(customerAPIAuthMiddleware))

		return customerAPIAuthMiddleware(handler)
	}

	suite.Run("failure when route doesn't match an existing route", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		rr := httptest.NewRecorder()

		req := setUpRequest("/internal/does-not-exist", &serviceMember, nil)

		handlerWithMiddleware := setUpHandlerAndMiddleware()

		handlerWithMiddleware.ServeHTTP(rr, req)

		suite.Equal(http.StatusBadRequest, rr.Code, "handler returned wrong status code")
	})

	suite.Run("success when route is on allow list and user is a service member", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		rr := httptest.NewRecorder()

		// using an arbitrary ID here
		req := setUpRequest("/internal/moves/990fb790-df36-448d-aee0-682a23e60429", &serviceMember, nil)

		handlerWithMiddleware := setUpHandlerAndMiddleware()

		handlerWithMiddleware.ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	suite.Run("success when route is not on allow list and user is a service member", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), nil, nil)

		rr := httptest.NewRecorder()

		// using an arbitrary ID here
		req := setUpRequest("/internal/service_members/326de0c9-19e3-42a9-ba74-e11855ae27cd", &serviceMember, nil)

		handlerWithMiddleware := setUpHandlerAndMiddleware()

		handlerWithMiddleware.ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	suite.Run("success when route is on the allow list and user is an office user", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		rr := httptest.NewRecorder()

		// using an arbitrary ID here
		req := setUpRequest("/internal/moves/990fb790-df36-448d-aee0-682a23e60429", nil, &officeUser)

		handlerWithMiddleware := setUpHandlerAndMiddleware()

		handlerWithMiddleware.ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	})

	suite.Run("failure when route is not on the allow list and user is an office user", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)

		rr := httptest.NewRecorder()

		// using an arbitrary ID here
		req := setUpRequest("/internal/service_members/326de0c9-19e3-42a9-ba74-e11855ae27cd", nil, &officeUser)

		handlerWithMiddleware := setUpHandlerAndMiddleware()

		handlerWithMiddleware.ServeHTTP(rr, req)

		suite.Equal(http.StatusForbidden, rr.Code, "handler returned wrong status code")
	})
}

// Test permissions middleware with a user who will be ALLOWED POST access on the endpoint: ghc/v1/shipments/:shipmentID/approve
// role must have update.shipment permissions

func (suite *AuthSuite) TestRequirePermissionsMiddlewareAuthorized() {
	// TOO users have the proper permissions for our test - update.shipment
	tooOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})

	identity, err := models.FetchUserIdentity(suite.DB(), tooOfficeUser.User.OktaID)

	suite.NoError(err)

	rr := httptest.NewRecorder()
	// using an arbitrary ID here for the shipment
	req := httptest.NewRequest("POST", "/ghc/v1/shipments/123456/approve", nil)

	// And: the context contains the auth values
	handlerSession := auth.Session{
		UserID:          tooOfficeUser.User.ID,
		IDToken:         "fake Token",
		ApplicationName: "mil",
	}

	handlerSession.Roles = append(handlerSession.Roles, identity.Roles...)

	ctx := auth.SetSessionInRequestContext(req, &handlerSession)
	req = req.WithContext(ctx)

	handlerConfig := suite.HandlerConfig()
	api := ghcapi.NewGhcAPIHandler(handlerConfig)

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	middleware := PermissionsMiddleware(suite.AppContextForTest(), api)

	root := chi.NewRouter()
	root.Mount("/ghc/v1", api.Serve(middleware))

	middleware(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, tooOfficeUser.User.ID, "the authenticated user is different from expected")
}

// Test permissions middleware with a user who will be DENIED POST access on the endpoint: ghc/v1/shipments/:shipmentID/approve
// role must NOT have update.shipment permissions
func (suite *AuthSuite) TestRequirePermissionsMiddlewareUnauthorized() {
	// QAE users will be denied access as they lack the proper permissions for our test - update.shipment
	qaeOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeQae})

	identity, err := models.FetchUserIdentity(suite.DB(), qaeOfficeUser.User.OktaID)

	suite.NoError(err)

	rr := httptest.NewRecorder()
	// using an arbitrary ID here for the shipment
	req := httptest.NewRequest("POST", "/ghc/v1/shipments/123456/approve", nil)

	// And: the context contains the auth values
	handlerSession := auth.Session{
		UserID:          qaeOfficeUser.User.ID,
		IDToken:         "fake Token",
		ApplicationName: "mil",
	}

	handlerSession.Roles = append(handlerSession.Roles, identity.Roles...)

	ctx := auth.SetSessionInRequestContext(req, &handlerSession)
	req = req.WithContext(ctx)

	handlerConfig := suite.HandlerConfig()
	api := ghcapi.NewGhcAPIHandler(handlerConfig)

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	middleware := PermissionsMiddleware(suite.AppContextForTest(), api)

	root := chi.NewRouter()
	root.Mount("/ghc/v1", api.Serve(middleware))

	middleware(handler).ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code, "handler returned wrong status code")
	suite.Equal(handlerSession.UserID, qaeOfficeUser.User.ID, "the authenticated user is different from expected")
}

func (suite *AuthSuite) TestIsLoggedInWhenNoUserLoggedIn() {
	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	rr := httptest.NewRecorder()
	sessionManager := scs.New()
	handler := sessionManager.LoadAndSave(IsLoggedInMiddleware(suite.Logger(), false))

	handler.ServeHTTP(rr, req)

	// expects to return 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

	// expects to return that no one is logged in
	expected := "{\"isLoggedIn\":false}\n"
	suite.Equal(expected, rr.Body.String(), "handler returned wrong body")
}

func (suite *AuthSuite) TestUnderMaintenanceFlag() {
	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	rr := httptest.NewRecorder()
	sessionManager := scs.New()
	handler := sessionManager.LoadAndSave(IsLoggedInMiddleware(suite.Logger(), true))

	handler.ServeHTTP(rr, req)

	// expects to return 200 OK
	suite.Equal(http.StatusOK, rr.Code, "handler returned the wrong status code")

	// expects to return that no one is logged in
	expected := "{\"isLoggedIn\":false,\"underMaintenance\":true}\n"
	suite.Equal(expected, rr.Body.String(), "handler returned wrong body")
}

func (suite *AuthSuite) TestIsLoggedInWhenUserLoggedIn() {
	OktaID := "2400c3c5-019d-4031-9c27-8a553e022297"
	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
	}
	suite.MustSave(&user)

	req := httptest.NewRequest("GET", "/is_logged_in", nil)

	sessionManager := scs.New()
	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token"}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := sessionManager.LoadAndSave(IsLoggedInMiddleware(suite.Logger(), false))

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

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
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
	OktaID := "2400c3c5-019d-4031-9c27-8a553e022297"
	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
	}
	suite.MustSave(&user)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/v1/office-users", nil)

	// And: the context contains the auth values
	session := auth.Session{UserID: user.ID, IDToken: "fake Token", AdminUserID: uuid.Must(uuid.NewV4())}
	ctx := auth.SetSessionInRequestContext(req, &session)
	req = req.WithContext(ctx)

	var handlerSession *auth.Session
	handler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
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

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
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

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           "deactivated@example.com",
	}
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), &userIdentity, sessionManager)

	suite.Equal(authorizationResultUnauthorized, result, "authorizer did not recognize deactivated user")
}

func (suite *AuthSuite) TestAuthKnownSingleRoleOffice() {
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO})

	userIdentity, err := models.FetchUserIdentity(suite.DB(), officeUser.User.OktaID)
	suite.Assert().NoError(err)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		UserID:          *officeUser.UserID,
		Email:           officeUser.Email,
	}
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), userIdentity, sessionManager)

	suite.Equal(authorizationResultAuthorized, result)
	// Office app, so should only have office ID information
	suite.Equal(officeUser.ID, session.OfficeUserID)
	// Make sure session contains roles and permissions
	suite.NotEmpty(session.Roles)
	userRole, hasRole := officeUser.User.Roles.GetRole(roles.RoleTypeTIO)
	suite.True(hasRole)
	sessionRole, hasRole := session.Roles.GetRole(roles.RoleTypeTIO)
	suite.True(hasRole)
	suite.Equal(userRole.ID, sessionRole.ID)
	suite.NotEmpty(session.Permissions)
	suite.ElementsMatch(TIO.Permissions, session.Permissions)
}

func (suite *AuthSuite) TestAuthorizeDeactivateOfficeUser() {
	officeActive := false
	userIdentity := models.UserIdentity{
		OfficeActive: &officeActive,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           "deactivated@example.com",
	}
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), &userIdentity,
		sessionManager)

	suite.Equal(authorizationResultUnauthorized, result, "authorizer did not recognize deactivated office user")
}

func (suite *AuthSuite) TestRedirectOktaErrorMsg() {
	officeUserID := uuid.Must(uuid.NewV4())
	OktaID := ("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
	}
	suite.MustSave(&user)

	userIdentity := models.UserIdentity{
		ID:           user.ID,
		Active:       true,
		OfficeUserID: &officeUserID,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/okta/callback", appnames.OfficeServername), nil)

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
	}
	// okta.mil state cookie
	cookieName := StateCookieName(&session)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   "some mis-matched hash value",
		Path:    "/",
		Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
	}
	req.AddCookie(&cookie)

	authContext := suite.AuthContext()

	sessionManager := handlerConfig.SessionManagers().Office
	req = suite.SetupSessionRequest(req, &session, sessionManager)
	result := AuthorizeKnownUser(req.Context(), suite.AppContextWithSessionForTest(&session),
		&userIdentity, sessionManager)

	suite.Equal(authorizationResultAuthorized, result)
	// Set up mock callback handler for testing
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
		&MockHTTPClient{
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{"access_token": "mockToken", "id_token": "broken_id_token"}`))),
			},
			Err: nil,
		},
	}

	rr2 := httptest.NewRecorder()
	sessionManager.LoadAndSave(h).ServeHTTP(rr2, req)

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

	u := suite.urlForHost(appnames.OfficeServername)
	q := u.Query()
	q.Add("error", "SIGNIN_ERROR")
	u.RawQuery = q.Encode()
	suite.Equal(u.String(), rr2.Result().Header.Get("Location"))
}

// Test to make sure the full auth flow works, although we are using mock Okta endpoints
func (suite *AuthSuite) TestRedirectFromOktaForValidUser() {
	// build a real office user
	tioOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO})

	// Build provider
	provider, err := factory.BuildOktaProvider(officeProviderName)
	suite.NoError(err)

	// Mock the necessary Okta endpoints
	mockAndActivateOktaEndpoints(tioOfficeUser, provider)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
	}

	// okta.mil state cookie
	stateValue := "someStateValue"
	cookieName := StateCookieName(&session)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   shaAsString(stateValue),
		Path:    "/",
		Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/okta/callback?state=%s",
		appnames.OfficeServername, stateValue), nil)
	req.AddCookie(&cookie)

	authContext := suite.AuthContext()

	sessionManager := handlerConfig.SessionManagers().Office
	req = suite.SetupSessionRequest(req, &session, sessionManager)

	defer goth.ClearProviders()
	goth.UseProviders(provider)
	mockIDToken, err := generateJWTToken(provider.GetClientID(), provider.GetIssuerURL(), stateValue)
	suite.NoError(err)
	responseBody := fmt.Sprintf(`{"access_token": "mockToken", "id_token": "%s"}`, mockIDToken)
	// Create the callbackhandler with mock http client for testing
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
		&MockHTTPClient{
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
			},
			Err: nil,
		},
	}

	rr := httptest.NewRecorder()
	sessionManager.LoadAndSave(h).ServeHTTP(rr, req)

	suite.Equal(http.StatusTemporaryRedirect, rr.Code)

	suite.Equal(suite.urlForHost(appnames.OfficeServername).String(),
		rr.Result().Header.Get("Location"))
}

func (suite *AuthSuite) TestCallbackThatRequiresOktaParamsRedirect() {
	// build a real office user
	tioOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO})

	// Build provider
	provider, err := factory.BuildOktaProvider(officeProviderName)
	suite.NoError(err)

	// Mock the necessary Okta endpoints
	mockAndActivateOktaEndpoints(tioOfficeUser, provider)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Hostname:        appnames.OfficeServername,
	}

	// okta.mil state cookie
	stateValue := "someStateValue"
	cookieName := StateCookieName(&session)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   shaAsString(stateValue),
		Path:    "/",
		Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
	}
	errDescription := url.QueryEscape("The resource owner or authorization server denied the request.")
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/okta/callback?state=%s&error_description=%s",
		appnames.OfficeServername, stateValue, errDescription), nil)

	req.AddCookie(&cookie)

	authContext := suite.AuthContext()

	sessionManager := handlerConfig.SessionManagers().Office
	req = suite.SetupSessionRequest(req, &session, sessionManager)

	defer goth.ClearProviders()
	goth.UseProviders(provider)
	suite.NoError(err)
	// Create the callbackhandler with mock http client for testing
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
		&MockHTTPClient{
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			},
			Err: nil,
		},
	}

	rr := httptest.NewRecorder()
	sessionManager.LoadAndSave(h).ServeHTTP(rr, req)

	suite.Equal(http.StatusTemporaryRedirect, rr.Code)

	// this should clear the user's okta sessions and redirect them back to MM
	suite.Equal(suite.urlForHost(appnames.OfficeServername).String()+"sign-in"+"?okta_error=true",
		rr.Result().Header.Get("Location"))
}

func (suite *AuthSuite) TestCallbackThatLogsUserOutOfOkta() {
	// build a real office user
	tioOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO})

	// Build provider
	provider, err := factory.BuildOktaProvider(officeProviderName)
	suite.NoError(err)

	// Mock the necessary Okta endpoints
	mockAndActivateOktaEndpoints(tioOfficeUser, provider)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         "fake_token",
		Hostname:        appnames.OfficeServername,
	}

	// okta.mil state cookie
	stateValue := "someStateValue"
	errDescription := url.QueryEscape("The resource owner or authorization server denied the request.")
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/okta/callback?state=%s&error_description=%s",
		appnames.OfficeServername, stateValue, errDescription), nil)

	authContext := suite.AuthContext()

	sessionManager := handlerConfig.SessionManagers().Office
	req = suite.SetupSessionRequest(req, &session, sessionManager)

	defer goth.ClearProviders()
	goth.UseProviders(provider)
	suite.NoError(err)
	// Create the callbackhandler with mock http client for testing
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
		&MockHTTPClient{
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(""))),
			},
			Err: nil,
		},
	}

	rr := httptest.NewRecorder()
	sessionManager.LoadAndSave(h).ServeHTTP(rr, req)

	suite.Equal(http.StatusTemporaryRedirect, rr.Code)

	// since the ID token is in the session, we will use that to log the user out instead of it needing to clear the session
	actualURL, _ := url.Parse(rr.Result().Header.Get("Location"))
	redirectURI := url.QueryEscape(fmt.Sprintf("http://%s:1234/sign-in?okta_logged_out=true", appnames.OfficeServername))
	oktaLogoutURL := fmt.Sprintf("https://dummy.okta.com/oauth2/default/v1/logout?id_token_hint=%s&post_logout_redirect_uri=%s", session.IDToken, redirectURI)

	suite.Equal(actualURL.String(), oktaLogoutURL)
}

func generateJWTToken(aud, iss, nonce string) (string, error) {

	claims := jwt.MapClaims{
		"aud":   aud,
		"iss":   iss,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
		"iat":   time.Now().Unix(),
		"nonce": nonce,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = jwtKeyID

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(DummyRSAPrivateKey))
	if err != nil {
		return "", err
	}

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// Generate and activate Okta endpoints that will be using during the auth handlers.
func mockAndActivateOktaEndpoints(tioOfficeUser models.OfficeUser, provider *okta.Provider) {
	// Mock the OIDC .well-known openid-configuration endpoint
	jwksURL := provider.GetJWKSURL()
	openIDConfigURL := provider.GetOpenIDConfigURL()
	userInfoURL := provider.GetUserInfoURL()

	httpmock.RegisterResponder("GET", openIDConfigURL,
		httpmock.NewStringResponder(200, fmt.Sprintf(`{
            "jwks_uri": "%s"
        }`, jwksURL)))

	// Mock the JWKS endpoint to receive keys for JWT verification
	httpmock.RegisterResponder("GET", jwksURL,
		httpmock.NewStringResponder(200, fmt.Sprintf(`{
        "keys": [
            {
                "alg": "RS256",
                "kty": "RSA",
                "use": "sig",
                "n": "%s",
                "e": "AQAB",
                "kid": "%s"
            }
        ]
    }`, DummyRSAModulus, jwtKeyID)))

	// Mock the userinfo endpoint
	// Sub is the Okta user ID, it is not a UUID.
	tioOfficeOktaUserID := tioOfficeUser.User.OktaID

	httpmock.RegisterResponder("GET", userInfoURL,
		httpmock.NewStringResponder(200, fmt.Sprintf(`{
		"sub": "%s",
		"name": "name",
		"email": "name@okta.com"
	}`, tioOfficeOktaUserID)))

	httpmock.Activate()
}

// Test to make sure the full auth flow works, although we are using mock Okta endpoints
func (suite *AuthSuite) TestRedirectFromOktaForInvalidUser() {
	tioOfficeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTIO})
	suite.False(tioOfficeUser.Active)

	// Build provider
	provider, err := factory.BuildOktaProvider(officeProviderName)
	suite.NoError(err)

	// Mock the necessary Okta endpoints
	mockAndActivateOktaEndpoints(tioOfficeUser, provider)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
	}

	// okta.mil state cookie
	stateValue := "someStateValue"
	// The code value is what is used to retrieve the exchange token frin the code passed through the URL. As long as the code value exists,
	// then the exchangeCode function will not fail. HOWEVER, we still need to provide a mock exchange token in the body so that it
	// can be verified and a user can be pulled from it.
	codeValue := "someCodeValue"
	cookieName := StateCookieName(&session)
	cookie := http.Cookie{
		Name:    cookieName,
		Value:   shaAsString(stateValue),
		Path:    "/",
		Expires: auth.GetExpiryTimeFromMinutes(auth.SessionExpiryInMinutes),
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/okta/callback?state=%s&code=%s",
		appnames.OfficeServername, stateValue, codeValue), nil)
	req.AddCookie(&cookie)

	authContext := suite.AuthContext()

	sessionManager := handlerConfig.SessionManagers().Office
	req = suite.SetupSessionRequest(req, &session, sessionManager)

	defer goth.ClearProviders()
	goth.UseProviders(provider)
	mockIDToken, err := generateJWTToken(provider.GetClientID(), provider.GetIssuerURL(), stateValue)
	suite.NoError(err)
	responseBody := fmt.Sprintf(`{"access_token": "mockToken", "id_token": "%s"}`, mockIDToken)
	// Create the callbackhandler with mock http client for testing
	h := CallbackHandler{
		authContext,
		handlerConfig,
		setUpMockNotificationSender(),
		&MockHTTPClient{
			Response: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader([]byte(responseBody))),
			},
			Err: nil,
		},
	}

	rr := httptest.NewRecorder()
	sessionManager.LoadAndSave(h).ServeHTTP(rr, req)

	suite.Equal(http.StatusTemporaryRedirect, rr.Code)

	u := suite.urlForHost(appnames.OfficeServername)
	u.Path = "/invalid-permissions"
	suite.Equal(u.String(), rr.Result().Header.Get("Location"))
}

func (suite *AuthSuite) TestAuthKnownSingleRoleAdmin() {
	adminUserID := uuid.Must(uuid.NewV4())
	officeUserID := uuid.Must(uuid.NewV4())
	var adminUserRole models.AdminRole = "SYSTEM_ADMIN"
	OktaID := ("2400c3c5-019d-4031-9c27-8a553e022297")

	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
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

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
	}

	sessionManager := handlerConfig.SessionManagers().Admin
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), &userIdentity,
		sessionManager)
	suite.Equal(authorizationResultAuthorized, result)

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

	fakeToken := "some_token"
	session := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        appnames.MilServername,
	}
	sessionManager := handlerConfig.SessionManagers().Mil
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), &userIdentity,
		sessionManager)
	suite.Equal(authorizationResultAuthorized, result)

	foundUser, _ := models.GetUser(suite.DB(), user.ID)

	suite.NotEqual("", foundUser.CurrentMilSessionID)

	sessionStore := sessionManager.Store()
	_, existsBefore, _ := sessionStore.Find(foundUser.CurrentMilSessionID)
	suite.Equal(existsBefore, true)

	concurrentSession := auth.Session{
		ApplicationName: auth.MilApp,
		IDToken:         fakeToken,
		Hostname:        appnames.MilServername,
	}
	concurrentCtx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result = AuthorizeKnownUser(concurrentCtx, suite.AppContextWithSessionForTest(&concurrentSession),
		&userIdentity, sessionManager)
	suite.Equal(authorizationResultAuthorized, result)

	_, existsAfterConcurrentSession, _ := sessionStore.Find(foundUser.CurrentMilSessionID)
	suite.Equal(existsAfterConcurrentSession, false)
}

// TESTCASE SCENARIO
// What is being tested: authorizeUnknownUser function
// Mocked: oktaProvider, auth.Session, goth.User, scs.SessionManager
// Behaviour: The function gets passed in the following arguments:
// - an instance of goth.User: a struct with the okta ID and email
// - the callback handler
// - the session (instance of auth.Session)
// It should create the user using the okta ID and email, then create a
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
	sessionManager := handlerConfig.SessionManagers().Mil
	mockSender := setUpMockNotificationSender() // We should get an email for this activity

	// Prepare the goth.User to simulate the UUID and email that okta would
	// provide
	fakeUUID, _ := uuid.NewV4()
	user := models.OktaUser{
		Sub:   fakeUUID.String(),
		Email: "new_service_member@example.com",
	}
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), user,
		sessionManager, mockSender)
	suite.Equal(authorizationResultAuthorized, result)
	mockSender.(*mocks.NotificationSender).AssertNumberOfCalls(suite.T(), "SendNotification", 1)

	// Look up the user and service member in the test DB
	foundUser, _ := models.GetUserFromEmail(suite.DB(), user.Email)
	serviceMemberID := session.ServiceMemberID
	serviceMember, _ := models.FetchServiceMemberForUser(suite.DB(), &session, serviceMemberID)
	// Look up the session token in the session store (this test uses the memory store)
	sessionStore := sessionManager.Store()
	_, existsBefore, _ := sessionStore.Find(foundUser.CurrentMilSessionID)

	// Verify service member exists and its ID is populated in the session
	suite.NotEmpty(session.ServiceMemberID)

	// Verify session contains UserID that points to the newly-created user
	suite.Equal(foundUser.ID, session.UserID)

	// Verify user's OktaEmail and OktaID match the values passed in
	suite.Equal(user.Email, foundUser.OktaEmail)
	suite.Equal(user.Sub, foundUser.OktaID)

	// Verify that the user's CurrentMilSessionID is not empty. The value is
	// generated randomly, so we can't test for a specific string. Any string
	// except an empty string is acceptable.
	suite.NotEqual("", foundUser.CurrentMilSessionID)

	// Verify the session token also exists in the session store
	suite.Equal(true, existsBefore)

	// Verify the service member that was created is associated with the user
	// that was created
	suite.Equal(foundUser.ID, serviceMember.UserID)
}

func (suite *AuthSuite) TestAuthorizeDeactivateAdmin() {
	adminUserActive := false
	userIdentity := models.UserIdentity{
		AdminUserActive: &adminUserActive,
	}

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           "deactivated@example.com",
	}
	sessionManager := handlerConfig.SessionManagers().Admin
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), &userIdentity,
		sessionManager)

	suite.Equal(authorizationResultUnauthorized, result, "authorizer did not recognize deactivated admin user")
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
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}

	fakeUUID2, _ := uuid.NewV4()
	user := models.OktaUser{
		Sub:   fakeUUID2.String(),
		Email: officeUser.Email,
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), user,
		sessionManager, mockSender)
	suite.Equal(authorizationResultUnauthorized, result, "Office user is active")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeNotFound() {

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           "missing@email.com",
	}

	id, _ := uuid.NewV4()
	user := models.OktaUser{
		Sub:   id.String(),
		Email: "sample@email.com",
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), user,
		sessionManager, mockSender)
	suite.Equal(authorizationResultUnauthorized, result, "Office user not found")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeLogsIn() {
	user := factory.BuildDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	// no roles at all
	officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Active: true,
				UserID: &user.ID,
				Email:  user.OktaEmail,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	fakeToken := "some_token"

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}

	gothUser := models.OktaUser{
		Sub:   user.ID.String(),
		Email: officeUser.Email,
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), gothUser,
		sessionManager, mockSender)
	suite.Equal(authorizationResultAuthorized, result, "Office user should have been authorized")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), officeUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(officeUser.ID, session.OfficeUserID)
	suite.Equal(uuid.Nil, session.AdminUserID)
	suite.NotEqual("", foundUser.CurrentOfficeSessionID)
	// this user was created without roles or permissions
	suite.Empty(session.Roles)
	suite.Empty(session.Permissions)
}

func (suite *AuthSuite) TestAuthorizeUnknownUserOfficeLogsInWithPermissions() {
	user := factory.BuildDefaultUser(suite.DB())
	// user is in office_users but has never logged into the app
	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Active: true,
				UserID: &user.ID,
				Email:  user.OktaEmail,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, []roles.RoleType{roles.RoleTypeQae})

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()

	fakeToken := "some_token"

	session := auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.OfficeServername,
		Email:           officeUser.Email,
	}
	gothUser := models.OktaUser{
		Sub:   user.ID.String(),
		Email: officeUser.Email,
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Office
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), gothUser,
		sessionManager, mockSender)
	suite.Equal(authorizationResultAuthorized, result, "Office user should have been authorized")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), officeUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(officeUser.ID, session.OfficeUserID)
	suite.Equal(uuid.Nil, session.AdminUserID)
	suite.NotEqual("", foundUser.CurrentOfficeSessionID)
	// Make sure session contains roles and permissions
	suite.NotEmpty(session.Roles)
	userRole, hasRole := officeUser.User.Roles.GetRole(roles.RoleTypeQae)
	suite.True(hasRole)
	sessionRole, hasRole := session.Roles.GetRole(roles.RoleTypeQae)
	suite.True(hasRole)
	suite.Equal(userRole.ID, sessionRole.ID)
	suite.NotEmpty(session.Permissions)
	suite.ElementsMatch(QAE.Permissions, session.Permissions)
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
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		Hostname:        appnames.AdminServername,
		Email:           adminUser.Email,
	}

	fakeUUID2, _ := uuid.NewV4()
	user := models.OktaUser{
		Sub:   fakeUUID2.String(),
		Email: adminUser.Email,
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Admin
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), user,
		sessionManager, mockSender)
	suite.Equal(authorizationResultUnauthorized, result, "Admin user is active")
}

func (suite *AuthSuite) TestAuthorizeUnknownUserAdminNotFound() {
	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	// user not admin_users and has never logged into the app
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           "missing@email.com",
	}

	id, _ := uuid.NewV4()
	user := models.OktaUser{
		Sub:   id.String(),
		Email: "sample@email.com",
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Admin
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), user,
		sessionManager, mockSender)
	suite.Equal(authorizationResultUnauthorized, result, "Admin user not found")
}

func (suite *AuthSuite) TestAuthorizeKnownUserAdminNotFound() {
	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	// user exists in the DB, but not as an admin user
	fakeToken := "some_token"
	OktaID := "000"
	userID := uuid.Must(uuid.NewV4())
	serviceMemberID := uuid.Must(uuid.NewV4())

	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "email@example.com",
		Active:    true,
		ID:        userID,
	}
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          user.ID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           user.OktaEmail,
	}

	userIdentity := models.UserIdentity{
		ID:              user.ID,
		Active:          true,
		ServiceMemberID: &serviceMemberID,
	}

	sessionManager := handlerConfig.SessionManagers().Admin
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := AuthorizeKnownUser(ctx, suite.AppContextWithSessionForTest(&session), &userIdentity,
		sessionManager)
	suite.Equal(authorizationResultUnauthorized, result, "Admin user not found")
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
	fakeToken := "some_token"
	fakeUUID, _ := uuid.FromString("39b28c92-0506-4bef-8b57-e39519f42dc2")
	session := auth.Session{
		ApplicationName: auth.AdminApp,
		UserID:          fakeUUID,
		IDToken:         fakeToken,
		Hostname:        appnames.AdminServername,
		Email:           adminUser.Email,
	}

	gothUser := models.OktaUser{
		Sub:   user.ID.String(),
		Email: adminUser.Email,
	}

	mockSender := setUpMockNotificationSender()
	sessionManager := handlerConfig.SessionManagers().Admin
	ctx := suite.SetupSessionContext(context.Background(), &session, sessionManager)

	// Call the function under test
	result := authorizeUnknownUser(ctx, suite.AppContextWithSessionForTest(&session), gothUser,
		sessionManager, mockSender)
	suite.Equal(authorizationResultAuthorized, result, "Admin user should have been authorized")

	foundUser, _ := models.GetUserFromEmail(suite.DB(), adminUser.Email)

	// Office app, so should only have office ID information
	suite.Equal(adminUser.ID, session.AdminUserID)
	suite.Equal(uuid.Nil, session.OfficeUserID)
	suite.NotEqual("", foundUser.CurrentAdminSessionID)
}

func (suite *AuthSuite) TestoktaAuthenticatedRedirect() {
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
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/okta", appnames.OfficeServername), nil)
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
	clientCert := factory.FetchOrBuildDevlocalClientCert(suite.DB())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/prime/v1", appnames.PrimeServername), nil)

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
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

func (suite *AuthSuite) TestAuthorizePPTAS() {
	clientCert := factory.FetchOrBuildDevlocalClientCert(suite.DB())

	handlerConfig := suite.HandlerConfig()
	appnames := handlerConfig.AppNames()
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/pptas/v1", appnames.PrimeServername), nil)

	handler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})
	middleware := PPTASAuthorizationMiddleware(suite.Logger())(handler)
	rr := httptest.NewRecorder()

	ctx := SetClientCertInRequestContext(req, &clientCert)
	req = req.WithContext(ctx)
	middleware.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	// no cert in request
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", fmt.Sprintf("http://%s/pptas/v1", appnames.PrimeServername), nil)
	middleware.ServeHTTP(rr, req)

	suite.Equal(http.StatusUnauthorized, rr.Code)
}
