package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2/memstore"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type RoutingSuite struct {
	handlers.BaseHandlerTestSuite
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

func (suite *RoutingSuite) setupRouting() *Config {
	// Test that we can initialize routing and serve the index file
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetAppNames(appNames)

	sessionManagers := auth.SetupSessionManagers(true, memstore.New(), false,
		time.Duration(180), time.Duration(180))
	handlerConfig.SetSessionManagers(sessionManagers)

	fakeLoginGovProvider := authentication.NewLoginGovProvider("fakeHostname", "secret_key", suite.Logger())

	authContext := authentication.NewAuthContext(suite.Logger(), fakeLoginGovProvider, "http", 80, sessionManagers)

	fakeFs := afero.NewMemMapFs()
	fakeBase := "fakebase"
	f, err := fakeFs.Create(path.Join(fakeBase, "index.html"))
	suite.NoError(err)
	_, err = f.Write([]byte(indexContent))
	suite.NoError(err)

	rConfig := &Config{
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
	}

	return rConfig

}

func (suite *RoutingSuite) TestBasicRoutingInit() {

	rConfig := suite.setupRouting()
	h, err := InitRouting(suite.AppContextForTest(), nil, rConfig, &telemetry.Config{})
	suite.NoError(err)

	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/", appNames.MilServername), nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal(indexContent, rr.Body.String())
}

func (suite *RoutingSuite) TestServeGHC() {

	rConfig := suite.setupRouting()
	h, err := InitRouting(suite.AppContextForTest(), nil, rConfig, &telemetry.Config{})
	suite.NoError(err)

	user := testdatagen.MakeUser(suite.DB(), testdatagen.Assertions{})

	// make the request without auth
	// getting auth working here in the test is a good bit more work.
	// Will have that in a future PR
	req := httptest.NewRequest("GET", fmt.Sprintf("http://%s/ghc/v1/customer/%s", appNames.MilServername, user.ID.String()), nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	suite.Equal(http.StatusUnauthorized, rr.Code)

	// make the request with GHC routing turned off
	rConfig = suite.setupRouting()
	rConfig.ServeGHC = false
	h, err = InitRouting(suite.AppContextForTest(), nil, rConfig, &telemetry.Config{})
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
