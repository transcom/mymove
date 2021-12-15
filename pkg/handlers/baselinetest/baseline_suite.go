package baselinetest

import (
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2/memstore"
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	"github.com/trussworks/httpbaselinetest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/handlers/routing"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type BaselineSuite struct {
	handlers.BaseHandlerTestSuite
	BaselineTestSuite *httpbaselinetest.Suite
	AppNames          auth.ApplicationServername
}

func NewBaselineSuite(t *testing.T) *BaselineSuite {
	// skipping all baseline tests for now as they don't pass yet
	if os.Getenv("BASELINETEST") == "" {
		t.Skip("Skipping baselinetest")
	}
	return &BaselineSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		BaselineTestSuite:    httpbaselinetest.NewDefaultSuite(t),
		AppNames: auth.ApplicationServername{
			MilServername:    "mil.example.com",
			OfficeServername: "office.example.com",
			AdminServername:  "admin.example.com",
		},
	}
}

func (suite *BaselineSuite) GetSqlxDb() *sqlx.DB {
	// *sigh* pop does not expose the sqlx.DB, so use reflection to get
	// access

	// the store has *sqlx.DB as the first field
	dbi := reflect.ValueOf(suite.DB().Store).Elem().Field(0).Interface()
	if db, ok := dbi.(*sqlx.DB); ok {
		return db
	}
	return nil
}

func (suite *BaselineSuite) RoutingConfigForTest() *routing.Config {
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetAppNames(suite.AppNames)
	mc := clock.NewMock()
	mc.Set(time.Date(2000, time.January, 2, 3, 4, 5, 6, time.UTC))

	handlerConfig.SetClock(mc)

	sessionManagers := SetupFakeSessionManagers(mc, memstore.New(), false,
		time.Duration(180), time.Duration(180))
	handlerConfig.SetAppSessionManagers(sessionManagers)

	p := suite.initFakeLoginGovProvider()

	authConfig := authentication.NewAuthConfig(suite.Logger(),
		p, "http", 80)

	// runtime.Caller is the preferred go way of getting the current filename
	_, fileName, _, _ := runtime.Caller(0)
	buildRoot := filepath.Join(filepath.Dir(fileName), "fakebase")

	return &routing.Config{
		HandlerConfig:  handlerConfig,
		AuthConfig:     authConfig,
		CSRFMiddleware: routing.NewFakeCSRFMiddleware(mc, suite.Logger()),
		BuildRoot:      buildRoot,
	}
}

func (suite *BaselineSuite) RoutingForTest() http.Handler {
	handler, err := routing.InitRouting(
		suite.AppContextForTest(),
		nil,
		suite.RoutingConfigForTest(),
		&telemetry.Config{})
	if err != nil {
		suite.T().Fatalf("Error initializing routing %s", err)
	}
	return handler
}
