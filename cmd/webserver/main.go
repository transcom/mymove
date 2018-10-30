package main

import (
	"bytes"
	"crypto/tls"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/transcom/mymove/pkg/config"
	"github.com/transcom/mymove/pkg/di"
	"github.com/transcom/mymove/pkg/handlers/dpsapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"go.uber.org/dig"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/authentication"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/publicapi"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
	documentServices "github.com/transcom/mymove/pkg/services/document"
	userServices "github.com/transcom/mymove/pkg/services/user"
	"github.com/transcom/mymove/pkg/storage"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"
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

func newDatabase(cfg *config.DatabaseConfig) (*pop.Connection, error) {
	if err := pop.AddLookupPaths(cfg.ConfigDir); err != nil {
		return nil, err
	}
	return pop.Connect(cfg.Environment)
}

func initializeHoneycomb(cfg *config.HoneycombConfig, logger *zap.Logger) {
	// For now, this should be part of the config parsing and validation
	if cfg.Enabled != nil && cfg.APIKey != nil && cfg.DataSet != nil && *cfg.Enabled && len(*cfg.APIKey) > 0 && len(*cfg.DataSet) > 0 {
		cfg.UseHoneycomb = true
	}

	if cfg.UseHoneycomb {
		logger.Debug("Honeycomb Integration enabled", zap.String("honeycomb-dataset", *cfg.DataSet))
		beeline.Init(beeline.Config{
			WriteKey: *cfg.APIKey,
			Dataset:  *cfg.DataSet,
			Debug:    *cfg.Debug,
		})
	} else {
		logger.Debug("Honeycomb Integration disabled")
	}

}

// Assert that our secret keys can be parsed into actual private keys
// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
func validateKeys(cCfg *server.SessionCookieConfig, lgConfig *authentication.LoginGovConfig, l *zap.Logger) {
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
	services.FetchServiceMember
	services.FetchUpload
	services.FetchDocument
	notifications.NotificationSender
	route.Planner
	storage.FileStorer
	Cookie *server.SessionCookieConfig
}

// FOR NOW - Once handlers are implemented like handlers.internalapi.ShowLoggedInUserHandler and have explicit
// dependencies we shouldn't need the big single HandlersContext
func populateHandlerContext(ctxt handlers.HandlerContext, p PopulateHandlerContextParams) {

	ctxt.SetCookieSecret(p.Cookie.Secret)
	if p.Cookie.NoTimeout {
		ctxt.SetNoSessionTimeout()
	}

	ctxt.SetNotificationSender(p.NotificationSender)
	ctxt.SetPlanner(p.Planner)
	ctxt.SetFileStorer(p.FileStorer)
	ctxt.SetFetchServiceMember(p.FetchServiceMember)
	ctxt.SetFetchDocument(p.FetchDocument)
	ctxt.SetFetchUpload(p.FetchUpload)
}

// fileHandler serves up a single file
func fileHandler(entryPoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entryPoint)
	}
}

// IndexHandlerFunc is a type marker for DI.
type IndexHandlerFunc http.HandlerFunc

// indexHandler injects New Relic client code and credentials into index.html
// and returns a handler that will serve the resulting content
func indexHandler(localEnv *server.LocalEnvConfig, cfg *config.NewRelicConfig, l *zap.Logger) IndexHandlerFunc {
	data := map[string]string{
		"NewRelicApplicationID": cfg.AppID,
		"NewRelicLicenseKey":    cfg.Key,
	}
	newRelicTemplate, err := template.ParseFiles(path.Join(localEnv.SiteDir, "new_relic.html"))
	if err != nil {
		l.Fatal("could not load new_relic.html template: run make client_build", zap.Error(err))
	}
	newRelicHTML := bytes.NewBuffer([]byte{})
	if err := newRelicTemplate.Execute(newRelicHTML, data); err != nil {
		l.Fatal("could not render new_relic.html template", zap.Error(err))
	}

	indexPath := path.Join(localEnv.SiteDir, "index.html")
	// #nosec - indexPath does not come from user input
	indexHTML, err := ioutil.ReadFile(indexPath)
	if err != nil {
		l.Fatal("could not read index.html template: run make client_build", zap.Error(err))
	}
	mergedHTML := bytes.Replace(indexHTML, []byte(`<script type="new-relic-placeholder"></script>`), newRelicHTML.Bytes(), 1)

	stat, err := os.Stat(indexPath)
	if err != nil {
		l.Fatal("could not stat index.html template", zap.Error(err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "index.html", stat.ModTime(), bytes.NewReader(mergedHTML))
	}
}

// BuildSiteParams contains the dependencies needed to configure the main site http.Handler
type BuildSiteParams struct {
	dig.In
	server.LogRequestMiddleware
	server.SessionCookieMiddleware
	server.AppDetectorMiddleware
	authentication.UserAuthMiddleware
	IndexHandlerFunc
	Env                   *server.LocalEnvConfig
	Swagger               *config.SwaggerConfig
	Hosts                 *server.HostsConfig
	S3Config              *storage.S3StorerConfig
	OrdersAPIHandler      ordersapi.Handler
	PublicAPIHandler      publicapi.Handler
	InternalAPIHandler    internalapi.Handler
	DPSAPIHandler         dpsapi.Handler
	AuthContext           *authentication.Context
	AuthCallbackHandler   *authentication.CallbackHandler
	AuthLogoutHandler     *authentication.LogoutHandler
	AuthUserListHandler   *authentication.UserListHandler
	AuthAssignUserHandler *authentication.AssignUserHandler
	AuthCreateUserHandler *authentication.CreateUserHandler
	HoneycombConfig       *config.HoneycombConfig
}

// SiteHandler is the DI marker for the main site http.Handler
type SiteHandler http.Handler

// buildSite creates the top level http.Handler for the site
func buildSite(p BuildSiteParams, l *zap.Logger) (SiteHandler, error) {

	// Base routes
	site := goji.NewMux()

	// Add site-wide middleware: they are evaluated in the reverse order in which
	// they are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(httpsComplianceMiddleware)
	site.Use(securityHeadersMiddleware)
	site.Use(limitBodySizeMiddleware)

	// Stub health check
	site.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, r *http.Request) {})

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(p.Env.SiteDir))

	// Allow public content through without any auth or app checks
	site.Handle(pat.Get("/static/*"), clientHandler)
	site.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	site.Handle(pat.Get("/downloads/*"), clientHandler)
	site.Handle(pat.Get("/favicon.ico"), clientHandler)

	// /orders/* has specific authentication controls
	ordersMux := goji.SubMux()
	ordersDetectionMiddleware := server.HostnameDetectorMiddleware(l, p.Hosts.OrdersName)
	ordersMux.Use(ordersDetectionMiddleware)
	ordersMux.Use(noCacheMiddleware)
	ordersMux.Handle(pat.Get("/swagger.yaml"), fileHandler(p.Swagger.Orders))
	ordersMux.Handle(pat.Get("/docs"), fileHandler(path.Join(p.Env.SiteDir, "swagger-ui", "orders.html")))
	ordersMux.Handle(pat.New("/*"), p.OrdersAPIHandler)
	site.Handle(pat.Get("/orders/v0/*"), ordersMux)

	dpsMux := goji.SubMux()
	dpsDetectionMiddleware := server.HostnameDetectorMiddleware(l, p.Hosts.DPSName)
	dpsMux.Use(dpsDetectionMiddleware)
	dpsMux.Use(noCacheMiddleware)
	dpsMux.Handle(pat.Get("/swagger.yaml"), fileHandler(p.Swagger.Orders))
	dpsMux.Handle(pat.Get("/docs"), fileHandler(path.Join(p.Env.SiteDir, "swagger-ui", "dps.html")))
	dpsMux.Handle(pat.New("/*"), p.DPSAPIHandler)
	site.Handle(pat.New("/dps/v0/*"), dpsMux)

	root := goji.NewMux()
	root.Use(p.SessionCookieMiddleware)
	root.Use(p.AppDetectorMiddleware) // Comes after the sessionCookieMiddleware as it sets session state
	root.Use(p.LogRequestMiddleware)
	site.Handle(pat.New("/*"), root)

	// /api/* - Public API
	publicMux := goji.SubMux()
	root.Handle(pat.New("/api/v1/*"), publicMux)
	publicMux.Handle(pat.Get("/swagger.yaml"), fileHandler(p.Swagger.API))
	publicMux.Handle(pat.Get("/docs"), fileHandler(path.Join(p.Env.SiteDir, "swagger-ui", "api.html")))

	publicAPIMux := goji.SubMux()
	publicMux.Handle(pat.New("/*"), publicAPIMux)
	publicAPIMux.Use(noCacheMiddleware)
	publicAPIMux.Use(p.UserAuthMiddleware)
	publicAPIMux.Handle(pat.New("/*"), p.PublicAPIHandler)

	// /internal/* - Internal API
	internalMux := goji.SubMux()
	root.Handle(pat.New("/internal/*"), internalMux)
	internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(p.Swagger.Internal))
	internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(p.Env.SiteDir, "swagger-ui", "internal.html")))

	internalAPIMux := goji.SubMux()
	internalMux.Handle(pat.New("/*"), internalAPIMux)
	internalAPIMux.Use(noCacheMiddleware)
	internalAPIMux.Use(p.UserAuthMiddleware)
	internalAPIMux.Handle(pat.New("/*"), p.InternalAPIHandler)

	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), &authentication.RedirectHandler{Context: *p.AuthContext})
	authMux.Handle(pat.Get("/login-gov/callback"), p.AuthCallbackHandler)
	authMux.Handle(pat.Get("/logout"), p.AuthLogoutHandler)

	if p.Env.Environment == "development" || p.Env.Environment == "test" {
		zap.L().Info("Enabling devlocal auth")
		localAuthMux := goji.SubMux()
		root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
		localAuthMux.Handle(pat.Get("/login"), p.AuthUserListHandler)
		localAuthMux.Handle(pat.Post("/login"), p.AuthAssignUserHandler)
		localAuthMux.Handle(pat.Post("/new"), p.AuthCreateUserHandler)
	}

	if p.S3Config == nil { // Using local filesystem
		// Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler("tmp")
		root.Handle(pat.Get("/storage/*"), fs)
	}

	// Serve index.html to all requests that haven't matches a previous route,
	root.HandleFunc(pat.Get("/*"), p.IndexHandlerFunc)

	// return site, wrapping in honeycomb if needed
	if p.HoneycombConfig == nil {
		return site, nil
	}
	return hnynethttp.WrapHandler(site), nil
}

func serveSite(cfg *config.ListenerConfig, hostsConfig *server.HostsConfig, siteHandler SiteHandler, l *zap.Logger) error {
	errChan := make(chan error)

	moveMilCerts := []server.TLSCert{
		{
			//Append move.mil cert with CA certificate chain
			CertPEMBlock: bytes.Join([][]byte{
				[]byte(cfg.DoDTLSCert),
				[]byte(cfg.DoDCACert)},
				[]byte("\n"),
			),
			KeyPEMBlock: []byte(cfg.DoDTLSKey),
		},
	}
	pkcs7Package, err := ioutil.ReadFile(cfg.DoDCACertPackage)
	if err != nil {
		l.Fatal("Failed to read DoD CA certificate package", zap.Error(err))
	}
	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		l.Fatal("Failed to parse DoD CA certificate package", zap.Error(err))
	}

	go func() {
		noTLSServer := server.Server{
			ListenAddress: hostsConfig.ListenInterface,
			HTTPHandler:   siteHandler,
			Logger:        l,
			Port:          cfg.NoTLSPort,
		}
		errChan <- noTLSServer.ListenAndServe()
	}()
	go func() {
		tlsServer := server.Server{
			ClientAuthType: tls.NoClientCert,
			ListenAddress:  hostsConfig.ListenInterface,
			HTTPHandler:    siteHandler,
			Logger:         l,
			Port:           cfg.TLSPort,
			TLSCerts:       moveMilCerts,
		}
		errChan <- tlsServer.ListenAndServeTLS()
	}()
	go func() {
		mutualTLSServer := server.Server{
			// Ensure that any DoD-signed client certificate can be validated,
			// using the package of DoD root and intermediate CAs provided by DISA
			CaCertPool:     dodCACertPool,
			ClientAuthType: tls.RequireAndVerifyClientCert,
			ListenAddress:  hostsConfig.ListenInterface,
			HTTPHandler:    siteHandler,
			Logger:         l,
			Port:           cfg.MutualTLSPort,
			TLSCerts:       moveMilCerts,
		}
		errChan <- mutualTLSServer.ListenAndServeTLS()
	}()
	return <-errChan
}

func dependencies() *di.Container {

	c := di.NewContainer(config.ParseConfig)
	c.MustInvoke(zap.ReplaceGlobals)
	// Now hook up the DI providers
	c.MustProvide(newDatabase)
	server.AddProviders(c)
	models.AddProviders(c)
	userServices.AddProviders(c)
	documentServices.AddProviders(c)
	handlers.AddProviders(c)
	internalapi.AddProviders(c)
	publicapi.AddProviders(c)
	ordersapi.AddProviders(c)
	dpsapi.AddProviders(c)
	authentication.AddProviders(c)
	route.AddProviders(c)
	notifications.AddProviders(c)
	storage.AddProviders(c)
	c.MustProvide(indexHandler)
	return c
}

func main() {

	// Set up the DI context and logging
	diContext := dependencies()

	// Initialize honeycomb
	diContext.MustInvoke(initializeHoneycomb)

	//  Validate that the keys used for RSA encryption are well formed
	diContext.MustInvoke(validateKeys)

	// FOR NOW configure handler context. This should not be necessary once each handler
	// declares its own dependencies.
	diContext.MustInvoke(populateHandlerContext)

	// Construct the main site handler
	diContext.MustProvide(buildSite)

	// And run the services
	diContext.MustInvoke(serveSite)
}
