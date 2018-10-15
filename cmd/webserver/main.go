package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"go.uber.org/dig"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	dep "github.com/transcom/mymove/pkg/dependencies"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/publicapi"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
	"github.com/transcom/mymove/pkg/storage"
)

// max request body size is 20 mb
const maxBodySize int64 = 200 * 1000 * 1000

func limitBodySizeMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("limitBodySizeMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func noCacheMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("noCacheMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func httpsComplianceMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("httpsComplianceMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		// set the HSTS header using values recommended by OWASP
		// https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet#Examples
		w.Header().Set("strict-transport-security", "max-age=31536000; includeSubdomains; preload")
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func securityHeadersMiddleware(inner http.Handler) http.Handler {
	zap.L().Debug("securityHeadersMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		// Sets headers to prevent rendering our page in an iframe, prevents clickjacking
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Frame-Options
		w.Header().Set("X-Frame-Options", "deny")
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy/frame-ancestors
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func newDatabase(cfg *DatabaseConfig) (*pop.Connection, error) {
	if err := pop.AddLookupPaths(cfg.configDir); err != nil {
		return nil, err
	}
	return pop.Connect(cfg.environment)
}

func dependencies() *dep.Container {

	c := dep.NewContainer(parseConfig)
	c.MustInvoke(zap.ReplaceGlobals)
	// Now hook up the DI providers
	c.MustProvide(newDatabase)
	logging.AddProviders(c)
	models.AddProviders(c)
	services.AddProviders(c)
	handlers.AddProviders(c)
	auth.AddProviders(c)
	route.AddProviders(c)
	notifications.AddProviders(c)
	storage.AddProviders(c)
	return c
}

func initializeHoneycomb(config *HoneycombConfig, logger *zap.Logger) {
	// For now, this should be part of the config parsing and validation
	if config.enabled != nil && config.apiKey != nil && config.dataSet != nil && *config.enabled && len(*config.apiKey) > 0 && len(*config.dataSet) > 0 {
		config.useHoneycomb = true
	}

	if config.useHoneycomb {
		logger.Debug("Honeycomb Integration enabled", zap.String("honeycomb-dataset", *config.dataSet))
		beeline.Init(beeline.Config{
			WriteKey: *config.apiKey,
			Dataset:  *config.dataSet,
			Debug:    *config.debug,
		})
	} else {
		logger.Debug("Honeycomb Integration disabled")
	}

}

// Assert that our secret keys can be parsed into actual private keys
// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
func validateKeys(cCfg *auth.SessionCookieConfig, lgConfig *authentication.LoginGovConfig, l *zap.Logger) {
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(cCfg.Secret)); err != nil {
		l.Fatal("Client auth private key", zap.Error(err))
	}
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(lgConfig.Secret)); err != nil {
		l.Fatal("Login.gov private key", zap.Error(err))
	}
}

// PopulateHandlerContextParams is the list of dependencies needed to populate the HandlerContext
type PopulateHandlerContextParams struct {
	dig.In
	cookie *auth.SessionCookieConfig
	sender notifications.NotificationSender
	router route.Planner
	storer storage.FileStorer
}

// FOR NOW - Once handlers are implemented like handlers.internalapi.ShowLoggedInUserHandler and have explicit
// dependencies we shouldn't need the big single HandlersContext
func populateHandlerContext(ctxt handlers.HandlerContext, p PopulateHandlerContextParams) {

	ctxt.SetCookieSecret(p.cookie.Secret)
	if p.cookie.NoTimeout {
		ctxt.SetNoSessionTimeout()
	}

	ctxt.SetNotificationSender(p.sender)
	ctxt.SetPlanner(p.router)
	ctxt.SetFileStorer(p.storer)
}

type BuildSiteParams struct {
	dig.In
	env *server.LocalEnvConfig
	logRequest logging.LogRequestMiddleware
}

func buildSite(p BuildSiteParams, l *zap.Logger) (http.Handler, error) {

	// Base routes
	site := goji.NewMux()

	// Add sitewide middleware: they are evaluated in the reverse order in which
	// they are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(httpsComplianceMiddleware)
	site.Use(securityHeadersMiddleware)
	site.Use(limitBodySizeMiddleware)

	// Stub health check
	site.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, r *http.Request) {})

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(p.env.SiteDir))
	// Allow public content through without any auth or app checks
	site.Handle(pat.Get("/static/*"), clientHandler)
	site.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	site.Handle(pat.Get("/downloads/*"), clientHandler)
	site.Handle(pat.Get("/favicon.ico"), clientHandler)

	// /orders/* has specific authentication controls
	ordersMux := goji.SubMux()
	ordersDetectionMiddleware := auth.OrdersDetectorMiddleware(l, *ordersHostname)
	ordersMux.Use(ordersDetectionMiddleware)
	ordersMux.Use(noCacheMiddleware)
	ordersMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*ordersSwagger))
	ordersMux.Handle(pat.Get("/docs"), fileHandler(path.Join(p.env.SiteDir, "swagger-ui", "orders.html")))
	ordersMux.Handle(pat.New("/*"), ordersapi.NewOrdersAPIHandler(handlerContext))
	site.Handle(pat.Get("/orders/v0/*"), ordersMux)

	root := goji.NewMux()
	root.Use(sessionCookieMiddleware)
	root.Use(appDetectionMiddleware) // Comes after the sessionCookieMiddleware as it sets session state
	root.Use(p.logRequest)
	site.Handle(pat.New("/*"), root)

}




apiMux := goji.SubMux()
root.Handle(pat.New("/api/v1/*"), apiMux)
apiMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*apiSwagger))
apiMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "api.html")))

externalAPIMux := goji.SubMux()
apiMux.Handle(pat.New("/*"), externalAPIMux)
externalAPIMux.Use(noCacheMiddleware)
externalAPIMux.Use(userAuthMiddleware)
externalAPIMux.Handle(pat.New("/*"), publicapi.NewPublicAPIHandler(handlerContext))

internalMux := goji.SubMux()
root.Handle(pat.New("/internal/*"), internalMux)
internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*internalSwagger))
internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "internal.html")))

// Mux for internal API that enforces auth
internalAPIMux := goji.SubMux()
internalMux.Handle(pat.New("/*"), internalAPIMux)
internalAPIMux.Use(userAuthMiddleware)
internalAPIMux.Use(noCacheMiddleware)
internalAPIMux.Handle(pat.New("/*"), internalapi.NewInternalAPIHandler(handlerContext))

authContext := authentication.NewAuthContext(logger, loginGovProvider, *loginGovCallbackProtocol, *loginGovCallbackPort)
authMux := goji.SubMux()
root.Handle(pat.New("/auth/*"), authMux)
authMux.Handle(pat.Get("/login-gov"), authentication.RedirectHandler{Context: authContext})
authMux.Handle(pat.Get("/login-gov/callback"), authentication.NewCallbackHandler(authContext, dbConnection, *clientAuthSecretKey, *noSessionTimeout))
authMux.Handle(pat.Get("/logout"), authentication.NewLogoutHandler(authContext, *clientAuthSecretKey, *noSessionTimeout))

if *env == "development" || *env == "test" {
zap.L().Info("Enabling devlocal auth")
localAuthMux := goji.SubMux()
root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
localAuthMux.Handle(pat.Get("/login"), authentication.NewUserListHandler(authContext, dbConnection))
localAuthMux.Handle(pat.Post("/login"), authentication.NewAssignUserHandler(authContext, dbConnection, *clientAuthSecretKey, *noSessionTimeout))
localAuthMux.Handle(pat.Post("/new"), authentication.NewCreateUserHandler(authContext, dbConnection, *clientAuthSecretKey, *noSessionTimeout))
}

if *storageBackend == "filesystem" {
// Add a file handler to provide access to files uploaded in development
fs := storage.NewFilesystemHandler("tmp")
root.Handle(pat.Get("/storage/*"), fs)
}

// Serve index.html to all requests that haven't matches a previous route,
root.HandleFunc(pat.Get("/*"), indexHandler(*build, *newRelicApplicationID, *newRelicLicenseKey, logger))

var httpHandler http.Handler
if useHoneycomb {
httpHandler = hnynethttp.WrapHandler(site)
} else {
httpHandler = site
}

func main() {

	// Set up the DI context and logging
	diContext := dependencies()

	// Initialize honeycomb
	diContext.MustInvoke(initializeHoneycomb)

	//  Validate that the keys used for RSA encryption are well formed
	diContext.MustInvoke(validateKeys)

	// FOR NOW configure handler context
	diContext.MustInvoke(populateHandlerContext)


	errChan := make(chan error)
	moveMilCerts := []server.TLSCert{
		{
			//Append move.mil cert with CA certificate chain
			CertPEMBlock: bytes.Join([][]byte{
				[]byte(*moveMilDODTLSCert),
				[]byte(*moveMilDODCACert)},
				[]byte("\n"),
			),
			KeyPEMBlock: []byte(*moveMilDODTLSKey),
		},
	}
	go func() {
		noTLSServer := server.Server{
			ListenAddress: *listenInterface,
			HTTPHandler:   httpHandler,
			Logger:        logger,
			Port:          *noTLSPort,
		}
		errChan <- noTLSServer.ListenAndServe()
	}()
	go func() {
		tlsServer := server.Server{
			ClientAuthType: tls.NoClientCert,
			ListenAddress:  *listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         logger,
			Port:           *tlsPort,
			TLSCerts:       moveMilCerts,
		}
		errChan <- tlsServer.ListenAndServeTLS()
	}()
	go func() {
		mutualTLSServer := server.Server{
			// Only allow certificates validated by the specified
			// client certificate CA.
			ClientAuthType: tls.RequireAndVerifyClientCert,
			CACertPEMBlock: []byte(*moveMilDODCACert),
			ListenAddress:  *listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         logger,
			Port:           *mutualTLSPort,
			TLSCerts:       moveMilCerts,
		}
		errChan <- mutualTLSServer.ListenAndServeTLS()
	}()
	logger.Fatal("listener error", zap.Error(<-errChan))
}

// fileHandler serves up a single file
func fileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}

// indexHandler injects New Relic client code and credentials into index.html
// and returns a handler that will serve the resulting content
func indexHandler(buildDir, newRelicApplicationID, newRelicLicenseKey string, logger *zap.Logger) http.HandlerFunc {
	data := map[string]string{
		"NewRelicApplicationID": newRelicApplicationID,
		"NewRelicLicenseKey":    newRelicLicenseKey,
	}
	newRelicTemplate, err := template.ParseFiles(path.Join(buildDir, "new_relic.html"))
	if err != nil {
		logger.Fatal("could not load new_relic.html template: run make client_build", zap.Error(err))
	}
	newRelicHTML := bytes.NewBuffer([]byte{})
	if err := newRelicTemplate.Execute(newRelicHTML, data); err != nil {
		logger.Fatal("could not render new_relic.html template", zap.Error(err))
	}

	indexPath := path.Join(buildDir, "index.html")
	// #nosec - indexPath does not come from user input
	indexHTML, err := ioutil.ReadFile(indexPath)
	if err != nil {
		logger.Fatal("could not read index.html template: run make client_build", zap.Error(err))
	}
	mergedHTML := bytes.Replace(indexHTML, []byte(`<script type="new-relic-placeholder"></script>`), newRelicHTML.Bytes(), 1)

	stat, err := os.Stat(indexPath)
	if err != nil {
		logger.Fatal("could not stat index.html template", zap.Error(err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "index.html", stat.ModTime(), bytes.NewReader(mergedHTML))
	}
}
