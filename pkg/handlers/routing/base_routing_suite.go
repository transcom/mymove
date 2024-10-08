package routing

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/csrf"
	"github.com/spf13/afero"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/handlers/authentication/okta"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	storageTest "github.com/transcom/mymove/pkg/storage/test"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type BaseRoutingSuite struct {
	handlers.BaseHandlerTestSuite
	port          int
	indexContent  string
	serverName    string
	routingConfig *Config
}

func NewBaseRoutingSuite() BaseRoutingSuite {
	return BaseRoutingSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(
			notifications.NewStubNotificationSender("milmovelocal"),
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
		port:         80,
		indexContent: "<html></html>",
		serverName:   "test-server",
	}
}

// override HandlerConfig to use the version saved in routing config
// so the same session manager(s) are used
func (suite *BaseRoutingSuite) HandlerConfig() handlers.HandlerConfig {
	return suite.RoutingConfig().HandlerConfig
}

// EqualDefaultIndex compares the response and ensures it has been
// served by the default index handler
func (suite *BaseRoutingSuite) EqualDefaultIndex(rr *httptest.ResponseRecorder) {
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(suite.indexContent, rr.Body.String())
}

func (suite *BaseRoutingSuite) EqualServerName(actualServerName string) {
	suite.Equal(suite.serverName, actualServerName)
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
	handlerConfig.SetNotificationSender(suite.TestNotificationSender())

	// Need this for any requests that will either retrieve or save files or their info.
	fakeS3 := storageTest.NewFakeS3Storage(true)
	handlerConfig.SetFileStorer(fakeS3)

	fakeOktaProvider := okta.NewOktaProvider(suite.Logger())
	authContext := authentication.NewAuthContext(suite.Logger(), *fakeOktaProvider, "http", suite.port)

	fakeFs := afero.NewMemMapFs()
	fakeBase := "fakebase"
	f, err := fakeFs.Create(path.Join(fakeBase, "index.html"))
	suite.FatalNoError(err)
	_, err = f.Write([]byte(suite.indexContent))
	suite.FatalNoError(err)
	suite.FatalNoError(f.Close())

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
		ServePPTAS:          true,
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
	return suite.SetupCustomSiteHandlerWithTelemetry(routingConfig, &telemetry.Config{})
}

func (suite *BaseRoutingSuite) SetupCustomSiteHandlerWithTelemetry(routingConfig *Config, telemetryConfig *telemetry.Config) http.Handler {
	siteHandler, err := InitRouting(suite.serverName, suite.AppContextForTest(), nil, routingConfig, telemetryConfig)
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

	suite.FatalNotNil(user.OktaID)
	suite.NotNil(user.OktaID)
	userIdentity, err := models.FetchUserIdentity(suite.DB(), user.OktaID)
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
	tokenHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
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

func (suite *BaseRoutingSuite) NewRequest(method string, hostname string, relativePath string, body io.Reader) *http.Request {
	req := httptest.NewRequest(method,
		fmt.Sprintf("http://%s%s", hostname, relativePath),
		body)
	// ensure the request has the suite logger
	return req.WithContext(logging.NewContext(req.Context(), suite.Logger()))
}

func (suite *BaseRoutingSuite) NewAdminRequest(method string, relativePath string, body io.Reader) *http.Request {
	return suite.NewRequest(method,
		suite.HandlerConfig().AppNames().AdminServername,
		relativePath,
		body)
}

func (suite *BaseRoutingSuite) NewAuthenticatedAdminRequest(method string, relativePath string, body io.Reader, adminUser models.AdminUser) *http.Request {
	req := suite.NewAdminRequest(method, relativePath, body)
	suite.SetupAdminRequestSession(req, adminUser)
	return req
}

func (suite *BaseRoutingSuite) NewMilRequest(method string, relativePath string, body io.Reader) *http.Request {
	return suite.NewRequest(method,
		suite.HandlerConfig().AppNames().MilServername,
		relativePath,
		body)
}

func (suite *BaseRoutingSuite) NewAuthenticatedMilRequest(method string, relativePath string, body io.Reader, serviceMember models.ServiceMember) *http.Request {
	req := suite.NewMilRequest(method, relativePath, body)
	suite.SetupMilRequestSession(req, serviceMember)
	return req
}

func (suite *BaseRoutingSuite) NewOfficeRequest(method string, relativePath string, body io.Reader) *http.Request {
	return suite.NewRequest(method,
		suite.HandlerConfig().AppNames().OfficeServername,
		relativePath,
		body)
}

func (suite *BaseRoutingSuite) NewAuthenticatedOfficeRequest(method string, relativePath string, body io.Reader, officeUser models.OfficeUser) *http.Request {
	req := suite.NewOfficeRequest(method, relativePath, body)
	suite.SetupOfficeRequestSession(req, officeUser)
	return req
}

func (suite *BaseRoutingSuite) NewPrimeRequest(method string, relativePath string, body io.Reader) *http.Request {
	return suite.NewRequest(method,
		suite.HandlerConfig().AppNames().PrimeServername,
		relativePath,
		body)
}

// the authentication.DevLocalPrimeMiddleware checks for the existence
// of a particular hash, so ensure that hash exists in the db and is
// associated with a user
func (suite *BaseRoutingSuite) NewAuthenticatedPrimeRequest(method string, relativePath string, body io.Reader, clientCert models.ClientCert) *http.Request {
	req := suite.NewPrimeRequest(method, relativePath, body)
	req.Header.Add("X-Devlocal-Cert-Hash", clientCert.Sha256Digest)
	return req
}

func (suite *BaseRoutingSuite) NewPPTASRequest(method string, relativePath string, body io.Reader) *http.Request {
	return suite.NewRequest(method,
		suite.HandlerConfig().AppNames().PPTASServerName,
		relativePath,
		body)
}
func (suite *BaseRoutingSuite) NewAuthenticatedPPTASRequest(method string, relativePath string, body io.Reader, clientCert models.ClientCert) *http.Request {
	req := suite.NewOfficeRequest(method, relativePath, body)
	req.Header.Add("X-Devlocal-Cert-Hash", clientCert.Sha256Digest)
	return req
}

// The ClientCertMiddleware looks at the TLS certificate on the
// request to make sure it matches something in the database. Fake the TLS
// info on the request
func (suite *BaseRoutingSuite) NewTLSAuthenticatedPrimeRequest(method string, relativePath string, body io.Reader) *http.Request {

	req := suite.NewPrimeRequest(method, relativePath, body)
	// runtime.Caller gets the path to the current file
	_, filename, _, ok := runtime.Caller(0)
	suite.FatalTrue(ok)
	dirname := filepath.Dir(filename)

	// Now load
	tlsDir, err := filepath.Abs(filepath.Join(dirname, "../../../config/tls"))
	suite.FatalNoError(err)
	certFile := filepath.Join(tlsDir, "devlocal-mtls.cer")
	keyFile := filepath.Join(tlsDir, "devlocal-mtls.key")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	suite.FatalNoError(err)
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	suite.FatalNoError(err)
	req.TLS = &tls.ConnectionState{
		PeerCertificates: []*x509.Certificate{x509Cert},
	}

	// make sure the this matches the devlocal hash in the db
	hash := sha256.Sum256(x509Cert.Raw)
	hashString := hex.EncodeToString(hash[:])
	devlocalCert := factory.FetchOrBuildDevlocalClientCert(suite.DB())
	suite.Equal(devlocalCert.Sha256Digest, hashString)
	return req
}

func (suite *BaseRoutingSuite) CreateFileWithContent(fpath string, fcontent string) {
	routingConfig := suite.RoutingConfig()
	dir := filepath.Dir(fpath)
	suite.NoError(routingConfig.FileSystem.MkdirAll(dir, 0600))
	f, err := routingConfig.FileSystem.Create(fpath)
	suite.NoError(err)
	_, err = f.WriteString(fcontent)
	suite.NoError(err)
	suite.NoError(f.Close())
}
