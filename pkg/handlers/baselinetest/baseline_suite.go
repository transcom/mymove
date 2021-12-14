package baselinetest

import (
	"net/http"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2/memstore"
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
	appNames          auth.ApplicationServername
}

func NewBaselineSuite(t *testing.T) *BaselineSuite {
	// skipping all baseline tests for now as they don't pass yet
	if os.Getenv("BASELINETEST") == "" {
		t.Skip("Skipping baselinetest")
	}
	return &BaselineSuite{
		BaseHandlerTestSuite: handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		BaselineTestSuite:    httpbaselinetest.NewDefaultSuite(t),
		appNames: auth.ApplicationServername{
			MilServername:    "mil.example.com",
			OfficeServername: "office.example.com",
			AdminServername:  "admin.example.com",
		},
	}
}

func (suite *BaselineSuite) getSqlxDb() *sqlx.DB {
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
	handlerConfig.SetAppNames(suite.appNames)

	sessionManagers := auth.SetupSessionManagers(true, memstore.New(), false,
		time.Duration(180), time.Duration(180))
	handlerConfig.SetSessionManagers(sessionManagers)

	p := suite.initFakeLoginGovProvider()

	authContext := authentication.NewAuthContext(suite.Logger(),
		p, "http", 80,
		handlerConfig.GetSessionManagers())

	return &routing.Config{
		HandlerConfig:  handlerConfig,
		AuthContext:    authContext,
		CSRFMiddleware: routing.NewFakeCSRFMiddleware(suite.Logger()),
		BuildRoot:      "fakebase",
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
