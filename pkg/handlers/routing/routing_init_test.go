package routing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type RoutingSuite struct {
	handlers.BaseHandlerTestSuite
	routingConfig *Config
}

func TestRoutingSuite(t *testing.T) {
	hs := &RoutingSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

var appNames = auth.ApplicationServername{
	MilServername:    "mil.example.com",
	OfficeServername: "office.example.com",
	AdminServername:  "admin.example.com",
	PrimeServername:  "prime.example.com",
}

const indexContent = "<html></html>"

func (suite *RoutingSuite) SetupTest() {
	// Test that we can initialize routing and serve the index file
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetAppNames(appNames)

	fakeLoginGovProvider := authentication.NewLoginGovProvider("fakeHostname", "secret_key", suite.Logger())

	authContext := authentication.NewAuthContext(suite.Logger(), fakeLoginGovProvider, "http", 80)

	fakeFs := afero.NewMemMapFs()
	fakeBase := "fakebase"
	f, err := fakeFs.Create(path.Join(fakeBase, "index.html"))
	suite.NoError(err)
	_, err = f.Write([]byte(indexContent))
	suite.NoError(err)

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
	}
}

func (suite *RoutingSuite) TestBasicRoutingInit() {

	h, err := InitRouting(suite.AppContextForTest(), nil, suite.routingConfig, &telemetry.Config{})
	suite.NoError(err)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/", appNames.MilServername), nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(indexContent, rr.Body.String())
}

func (suite *RoutingSuite) TestServeGHC() {

	h, err := InitRouting(suite.AppContextForTest(), nil, suite.routingConfig, &telemetry.Config{})
	suite.NoError(err)

	user := factory.BuildUser(suite.DB(), nil, nil)

	// make the request without auth
	// getting auth working here in the test is a good bit more work.
	// Will have that in a future PR
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/ghc/v1/customer/%s", appNames.MilServername, user.ID.String()), nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	suite.Equal(http.StatusUnauthorized, rr.Code)

	// make the request with GHC routing turned off
	suite.routingConfig.ServeGHC = false
	h, err = InitRouting(suite.AppContextForTest(), nil, suite.routingConfig, &telemetry.Config{})
	suite.NoError(err)
	req = httptest.NewRequest("GET", fmt.Sprintf("http://%s/ghc/v1/customer/%s", appNames.MilServername, user.ID.String()), nil)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	// if the API is not enabled, the routing will be served by the
	// SPA handler, sending back the index page, which will have the
	// javascript SPA routing
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(indexContent, rr.Body.String())
}

func (suite *RoutingSuite) setupOfficeRequestSession(req *http.Request, officeUser models.OfficeUser) {

	fakeAuthContext := context.Background()
	sessionManager := suite.routingConfig.HandlerConfig.SessionManagers().Office
	fakeAuthContext, err := sessionManager.Load(fakeAuthContext, "")
	suite.NoError(err)

	fakeSession := auth.Session{
		ApplicationName: auth.Application(auth.OfficeApp),
		Hostname:        suite.HandlerConfig().AppNames().OfficeServername,
	}

	userIdentity, err := models.FetchUserIdentity(suite.DB(), officeUser.User.LoginGovUUID.String())
	suite.Assert().NoError(err)

	// use AuthorizeKnownUser which also sets up various things in the
	// Session, including Permissions
	authentication.AuthorizeKnownUser(fakeAuthContext, suite.AppContextWithSessionForTest(&fakeSession),
		userIdentity, sessionManager)

	sessionManager.Put(fakeAuthContext, "session", fakeSession)
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

	req.Header.Add("cookie", cookie.String())
}

func (suite *RoutingSuite) TestOfficeLoggedInEndpoint() {
	h, err := InitRouting(suite.AppContextForTest(), nil, suite.routingConfig, &telemetry.Config{})
	suite.NoError(err)

	officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
		[]roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/internal/users/logged_in", appNames.OfficeServername), nil)
	suite.setupOfficeRequestSession(req, officeUser)

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)

	var userPayload internalmessages.LoggedInUserPayload
	suite.NoError(json.Unmarshal(rr.Body.Bytes(), &userPayload))
	suite.Equal(officeUser.UserID.String(), userPayload.ID.String())
	suite.NotNil(userPayload.OfficeUser)
	suite.Equal(officeUser.ID.String(), userPayload.OfficeUser.ID.String())
	suite.NotEmpty(userPayload.Permissions)
}
