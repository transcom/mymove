package routing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"time"

	"github.com/gorilla/csrf"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type BaseRoutingSuite struct {
	handlers.BaseHandlerTestSuite
	indexContent  string
	routingConfig *Config
}

func NewBaseRoutingSuite() BaseRoutingSuite {
	return BaseRoutingSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(
			notifications.NewStubNotificationSender("milmovelocal"),
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
		indexContent: "<html></html>",
	}
}

// override HandlerConfig to use the version saved in routing config
// so the same session manager(s) are used
func (suite *BaseRoutingSuite) HandlerConfig() handlers.HandlerConfig {
	return suite.RoutingConfig().HandlerConfig
}

func (suite *BaseRoutingSuite) RoutingConfig() *Config {
	// need to ensure we have only one routing config so the same
	// session managers are re-used
	if suite.routingConfig != nil {
		return suite.routingConfig
	}
	// ensure the routing config is reset when the test context is finished
	suite.T().Cleanup(func() { suite.routingConfig = nil })
	// Test that we can initialize routing and serve the index file
	handlerConfig := suite.BaseHandlerTestSuite.HandlerConfig()
	handlerConfig.SetAppNames(handlers.ApplicationTestServername())

	// Need this for any requests that will either retrieve or save files or their info.
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)

	fakeLoginGovProvider := authentication.NewLoginGovProvider("fakeHostname", "secret_key", suite.Logger())

	authContext := authentication.NewAuthContext(suite.Logger(), fakeLoginGovProvider, "http", 80)

	fakeFs := afero.NewMemMapFs()
	fakeBase := "fakebase"
	f, err := fakeFs.Create(path.Join(fakeBase, "index.html"))
	suite.FatalNoError(err)
	_, err = f.Write([]byte(suite.indexContent))
	suite.FatalNoError(err)

	// make a fake csrf auth key that would be completely insecure for
	// real world usage
	fakeCsrfAuthKey := make([]byte, 32)

	suite.routingConfig = &Config{
		FileSystem:    fakeFs,
		HandlerConfig: handlerConfig,
		AuthContext:   authContext,
		BuildRoot:     fakeBase,

		// include all these as true to increase test coverage
		ServeSwaggerUI:      true,
		ServePrime:          true,
		ServeSupport:        true,
		ServeDebugPProf:     true,
		ServeAPIInternal:    true,
		ServeAdmin:          true,
		ServePrimeSimulator: true,
		ServeGHC:            true,
		ServeDevlocalAuth:   true,
		ServeOrders:         true,

		CSRFMiddleware: InitCSRFMiddlware(fakeCsrfAuthKey, false, "/", auth.GorillaCSRFToken),

		// note that enabling devlocal auth also enables accessing the
		// Prime API without requiring mTLS. See
		// authentication.DevLocalClientCertMiddleware
	}

	return suite.routingConfig
}

func (suite *BaseRoutingSuite) SetupSiteHandler() http.Handler {
	return suite.SetupCustomSiteHandler(suite.RoutingConfig())
}

func (suite *BaseRoutingSuite) SetupCustomSiteHandler(routingConfig *Config) http.Handler {
	siteHandler, err := InitRouting(suite.AppContextForTest(), nil, routingConfig, &telemetry.Config{})
	suite.FatalNoError(err)
	return siteHandler
}

func (suite *BaseRoutingSuite) setupRequestSession(req *http.Request, user models.User, hostname string) {
	app, err := auth.ApplicationName(hostname, suite.HandlerConfig().AppNames())
	suite.FatalNoError(err)
	sessionManager := suite.HandlerConfig().SessionManagers().SessionManagerForApplication(app)

	fakeAuthContext, err := sessionManager.Load(context.Background(), "")
	suite.NoError(err)

	fakeSession := auth.Session{
		ApplicationName: app,
		Hostname:        hostname,
	}

	suite.FatalNotNil(user.LoginGovUUID)
	suite.FatalFalse(user.LoginGovUUID.IsNil())
	userIdentity, err := models.FetchUserIdentity(suite.DB(), user.LoginGovUUID.String())
	suite.FatalNoError(err)

	// use AuthorizeKnownUser which also sets up various things in the
	// Session, including Permissions
	authentication.AuthorizeKnownUser(fakeAuthContext, suite.AppContextWithSessionForTest(&fakeSession),
		userIdentity, sessionManager)

	sessionManager.Put(fakeAuthContext, "session", &fakeSession)
	// ignore expiry for this test
	// need to call commit ourselves to get the session token to put
	// into the cookie
	token, _, err := sessionManager.Commit(fakeAuthContext)
	suite.NoError(err)
	sessionCookie := sessionManager.SessionCookie()
	cookie := &http.Cookie{
		Name:     sessionCookie.Name,
		Value:    token,
		Path:     sessionCookie.Path,
		Domain:   sessionCookie.Domain,
		Secure:   sessionCookie.Secure,
		HttpOnly: sessionCookie.HttpOnly,
		SameSite: sessionCookie.SameSite,
	}
	cookie.Expires = time.Unix(1, 0)
	cookie.MaxAge = -1
	req.AddCookie(cookie)

	// set up CSRF cookie and headers
	maskedToken := ""
	tokenHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maskedToken = csrf.Token(r)
	})
	fakeReq := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	suite.routingConfig.CSRFMiddleware(tokenHandler).ServeHTTP(rr, fakeReq)
	for _, cookie := range rr.Result().Cookies() {
		req.AddCookie(cookie)
	}
	req.Header.Set("X-CSRF-Token", maskedToken)
}

func (suite *BaseRoutingSuite) SetupAdminRequestSession(req *http.Request, adminUser models.AdminUser) {
	suite.setupRequestSession(req, adminUser.User, suite.HandlerConfig().AppNames().AdminServername)
}

func (suite *BaseRoutingSuite) SetupMilRequestSession(req *http.Request, serviceMember models.ServiceMember) {
	suite.setupRequestSession(req, serviceMember.User, suite.HandlerConfig().AppNames().MilServername)
}

func (suite *BaseRoutingSuite) SetupOfficeRequestSession(req *http.Request, officeUser models.OfficeUser) {
	suite.setupRequestSession(req, officeUser.User, suite.HandlerConfig().AppNames().OfficeServername)
}

func (suite *BaseRoutingSuite) NewAdminRequest(method string, relativePath string, body io.Reader) *http.Request {
	return httptest.NewRequest(method,
		fmt.Sprintf("http://%s%s", suite.HandlerConfig().AppNames().AdminServername, relativePath),
		body)
}

func (suite *BaseRoutingSuite) NewAuthenticatedAdminRequest(method string, relativePath string, body io.Reader, adminUser models.AdminUser) *http.Request {
	req := suite.NewAdminRequest(method, relativePath, body)
	suite.SetupAdminRequestSession(req, adminUser)
	return req
}

func (suite *BaseRoutingSuite) NewMilRequest(method string, relativePath string, body io.Reader) *http.Request {
	return httptest.NewRequest(method,
		fmt.Sprintf("http://%s%s", suite.HandlerConfig().AppNames().MilServername, relativePath),
		body)
}

func (suite *BaseRoutingSuite) NewAuthenticatedMilRequest(method string, relativePath string, body io.Reader, serviceMember models.ServiceMember) *http.Request {
	req := suite.NewMilRequest(method, relativePath, body)
	suite.SetupMilRequestSession(req, serviceMember)
	return req
}

func (suite *BaseRoutingSuite) NewOfficeRequest(method string, relativePath string, body io.Reader) *http.Request {
	return httptest.NewRequest(method,
		fmt.Sprintf("http://%s%s", suite.HandlerConfig().AppNames().OfficeServername, relativePath),
		body)
}

func (suite *BaseRoutingSuite) NewAuthenticatedOfficeRequest(method string, relativePath string, body io.Reader, officeUser models.OfficeUser) *http.Request {
	req := suite.NewOfficeRequest(method, relativePath, body)
	suite.SetupOfficeRequestSession(req, officeUser)
	return req
}

func (suite *BaseRoutingSuite) NewPrimeRequest(method string, relativePath string, body io.Reader) *http.Request {
	return httptest.NewRequest(method,
		fmt.Sprintf("http://%s%s", suite.HandlerConfig().AppNames().PrimeServername, relativePath),
		body)
}

// the authentication.DevLocalPrimeMiddleware checks for the existence
// of a particular hash, so ensure that hash exists in the db and is
// associated with a user
func (suite *BaseRoutingSuite) NewAuthenticatedPrimeRequest(method string, relativePath string, body io.Reader, clientCert models.ClientCert) *http.Request {
	req := suite.NewMilRequest(method, relativePath, body)
	req.Header.Add("X-Devlocal-Cert-Hash", clientCert.Sha256Digest)
	return req
}
