package baselinetest

import (
	"context"
	"net/http"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2/memstore"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/afero"
	"github.com/trussworks/httpbaselinetest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/authentication"
	"github.com/transcom/mymove/pkg/handlers/routing"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/telemetry"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type BaselineSuite struct {
	handlers.BaseHandlerTestSuite
	BaselineTestSuite *httpbaselinetest.Suite
	AppNames          auth.ApplicationServername
	RoutingConfig     *routing.Config
}

func NewBaselineSuite(t *testing.T) *BaselineSuite {

	// this is kinda dangerous as it overrides the global
	// behavior, but that's all that gofrs/uuid provides
	// so ensure the original is restored at the end of the test
	origUUIDGenerator := uuid.DefaultGenerator
	uuidCleanup := func() {
		uuid.DefaultGenerator = origUUIDGenerator
	}
	uuid.DefaultGenerator = NewFakeGenerator()
	t.Cleanup(uuidCleanup)

	baseHandlerSuite := handlers.NewBaseHandlerTestSuite(notifications.NewStubNotificationSender("milmovelocal"), testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction())

	pop.SetNowFunc(func() time.Time {
		return baseHandlerSuite.HandlerConfig().Clock().Now()
	})
	nowCleanup := func() {
		pop.SetNowFunc(time.Now)
	}
	t.Cleanup(nowCleanup)

	return &BaselineSuite{
		BaseHandlerTestSuite: baseHandlerSuite,
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

const indexContent = "<html></html>"

func (suite *BaselineSuite) RoutingConfigForTest() *routing.Config {
	if suite.RoutingConfig != nil {
		return suite.RoutingConfig
	}
	handlerConfig := suite.HandlerConfig()
	handlerConfig.SetAppNames(suite.AppNames)

	memStore := memstore.New()
	sessionManagers := SetupFakeSessionManagers(handlerConfig.Clock(), memStore, false,
		time.Duration(180), time.Duration(180))
	handlerConfig.SetAppSessionManagers(sessionManagers)

	p := suite.initFakeLoginGovProvider()

	authConfig := authentication.NewAuthConfig(suite.Logger(),
		p, "http", 80)

	fakeFs := afero.NewMemMapFs()
	fakeBase := "fakebase"
	f, err := fakeFs.Create(path.Join(fakeBase, "index.html"))
	suite.NoError(err)
	_, err = f.Write([]byte(indexContent))
	suite.NoError(err)

	suite.RoutingConfig = &routing.Config{
		FileSystem:           fakeFs,
		HandlerConfig:        handlerConfig,
		AuthConfig:           authConfig,
		CSRFMiddleware:       routing.NewFakeCSRFMiddleware(handlerConfig.Clock(), suite.Logger()),
		MaskedCSRFMiddleware: routing.NewFakeMaskedCSRFMiddleware(suite.Logger()),
		BuildRoot:            fakeBase,

		ServeAPIInternal: true,
		ServePrime:       true,
		ServeGHC:         true,
		ServeAdmin:       true,
	}

	return suite.RoutingConfig
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

func (suite *BaselineSuite) CookiesForUser(user models.User) []http.Cookie {
	milSession := suite.RoutingConfigForTest().HandlerConfig.AppSessionManagers().MilSessionManager()
	fakeMilSession, ok := milSession.(*FakeSessionManager)
	suite.FatalFalse(!ok, "Cannot convert to FakeSessionManager")

	sessionToken, err := fakeMilSession.generateToken()
	if err != nil {
		suite.FatalNoError(err, "Cannot generate session token")
	}
	session := auth.Session{
		ApplicationName: auth.MilApp,
		UserID:          user.ID,
		IDToken:         "fake_openid_token",
	}
	// milSession.Put(ctx, "session", session)
	// _, _, err = milSession.Commit(ctx)
	// if err != nil {
	// 	suite.FatalNoError(err, "Error committing session")
	// }

	memStore, ok := fakeMilSession.ScsStore.(*memstore.MemStore)
	suite.FatalFalse(!ok, "Cannot convert scs store to memstore")

	deadline := fakeMilSession.clock.Now().Add(fakeMilSession.Lifetime).UTC()
	values := make(map[string]interface{})
	values["session"] = session
	b, err := fakeMilSession.Codec.Encode(deadline, values)
	suite.FatalNoError(err, "Cannot encode deadline+values")

	// set expiry far in the future because scs.Store uses time.Now
	err = memStore.Commit(sessionToken, b, time.Now().AddDate(10, 0, 0))
	suite.FatalNoError(err, "Cannot commit to memstore")

	ctx := context.Background()
	sessionCookie := fakeMilSession.GetSessionCookie(ctx, sessionToken, time.Time{})

	csrfCookies := routing.GetFakeCSRFCookies()

	return append(csrfCookies, *sessionCookie)
}
