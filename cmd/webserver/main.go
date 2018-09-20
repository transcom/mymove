package main

import (
	"bytes"
	"crypto/tls"
	"errors"
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
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
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

func main() {

	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	listenInterface := flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")

	myHostname := flag.String("http_my_server_name", "localhost", "Hostname according to environment.")
	officeHostname := flag.String("http_office_server_name", "officelocal", "Hostname according to environment.")
	tspHostname := flag.String("http_tsp_server_name", "tsplocal", "Hostname according to environment.")
	ordersHostname := flag.String("http_orders_server_name", "orderslocal", "Hostname according to environment.")

	httpPort := flag.String("port", "8080", "the HTTP `port` to listen on.")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
	clientAuthSecretKey := flag.String("client_auth_secret_key", "", "Client auth secret JWT key.")
	noSessionTimeout := flag.Bool("no_session_timeout", false, "whether user sessions should timeout.")

	httpServerLogLinePrefixDrop := flag.String("http_server_log_line_prefix_drop", "", "drop any log lines with this prefix")

	internalSwagger := flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	apiSwagger := flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	ordersSwagger := flag.String("orders-swagger", "swagger/orders.yaml", "The location of the Orders API swagger definition")

	httpsClientAuthPort := flag.String("https_client_auth_port", "9443", "The `port` for the HTTPS listener requiring client authentication.")
	httpsClientAuthCACert := flag.String("https_client_auth_ca_cert", "", "the CA certificate for the HTTPS listener requiring client authentication.")

	httpsPort := flag.String("https_port", "8443", "the `port` to listen on.")
	httpsCert := flag.String("https_cert", "", "TLS certificate.")
	httpsKey := flag.String("https_key", "", "TLS private key.")

	loginGovCallbackProtocol := flag.String("login_gov_callback_protocol", "https://", "Protocol for non local environments.")
	loginGovCallbackPort := flag.String("login_gov_callback_port", "443", "The port for callback urls.")
	loginGovSecretKey := flag.String("login_gov_secret_key", "", "Login.gov auth secret JWT key.")
	loginGovMyClientID := flag.String("login_gov_my_client_id", "", "Client ID registered with login gov.")
	loginGovOfficeClientID := flag.String("login_gov_office_client_id", "", "Client ID registered with login gov.")
	loginGovTSPClientID := flag.String("login_gov_tsp_client_id", "", "Client ID registered with login gov.")
	loginGovHostname := flag.String("login_gov_hostname", "", "Hostname for communicating with login gov.")

	/* For bing Maps use the following
	bingMapsEndpoint := flag.String("bing_maps_endpoint", "", "URL for the Bing Maps Truck endpoint to use")
	bingMapsKey := flag.String("bing_maps_key", "", "Authentication key to use for the Bing Maps endpoint")
	*/
	hereGeoEndpoint := flag.String("here_maps_geocode_endpoint", "", "URL for the HERE maps geocoder endpoint")
	hereRouteEndpoint := flag.String("here_maps_routing_endpoint", "", "URL for the HERE maps routing endpoint")
	hereAppID := flag.String("here_maps_app_id", "", "HERE maps App ID for this application")
	hereAppCode := flag.String("here_maps_app_code", "", "HERE maps App API code")

	storageBackend := flag.String("storage_backend", "filesystem", "Storage backend to use, either filesystem or s3.")
	emailBackend := flag.String("email_backend", "local", "Email backend to use, either SES or local")
	s3Bucket := flag.String("aws_s3_bucket_name", "", "S3 bucket used for file storage")
	s3Region := flag.String("aws_s3_region", "", "AWS region used for S3 file storage")
	s3KeyNamespace := flag.String("aws_s3_key_namespace", "", "Key prefix for all objects written to S3")
	awsSesRegion := flag.String("aws_ses_region", "", "AWS region used for SES")

	newRelicApplicationID := flag.String("new_relic_application_id", "", "App ID for New Relic Browser")
	newRelicLicenseKey := flag.String("new_relic_license_key", "", "License key for New Relic Browser")

	honeycombEnabled := flag.Bool("honeycomb_enabled", false, "Honeycomb enabled")
	honeycombAPIKey := flag.String("honeycomb_api_key", "", "API Key for Honeycomb")
	honeycombDataset := flag.String("honeycomb_dataset", "", "Dataset for Honeycomb")
	honeycombDebug := flag.Bool("honeycomb_debug", false, "Debug honeycomb using stdout.")

	flag.Parse()

	logger, err := logging.Config(*env, *debugLogging)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// Honeycomb
	useHoneycomb := false
	if honeycombEnabled != nil && honeycombAPIKey != nil && honeycombDataset != nil && *honeycombEnabled && len(*honeycombAPIKey) > 0 && len(*honeycombDataset) > 0 {
		useHoneycomb = true
	}
	if useHoneycomb {
		zap.L().Debug("Honeycomb Integration enabled")
		beeline.Init(beeline.Config{
			WriteKey: *honeycombAPIKey,
			Dataset:  *honeycombDataset,
			Debug:    *honeycombDebug,
		})
	} else {
		zap.L().Debug("Honeycomb Integration disabled")
	}

	// Assert that our secret keys can be parsed into actual private keys
	// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*loginGovSecretKey)); err != nil {
		logger.Fatal("Login.gov private key", zap.Error(err))
	}
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*clientAuthSecretKey)); err != nil {
		logger.Fatal("Client auth private key", zap.Error(err))
	}
	if *loginGovHostname == "" {
		log.Fatal("Must provide the Login.gov hostname parameter, exiting")
	}

	//DB connection
	err = pop.AddLookupPaths(*config)
	if err != nil {
		logger.Fatal("Adding Pop config path", zap.Error(err))
	}
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	// Register Login.gov authentication provider for My.(move.mil)
	loginGovProvider := authentication.NewLoginGovProvider(*loginGovHostname, *loginGovSecretKey, logger)
	err = loginGovProvider.RegisterProvider(*myHostname, *loginGovMyClientID, *officeHostname, *loginGovOfficeClientID, *tspHostname, *loginGovTSPClientID, *loginGovCallbackProtocol, *loginGovCallbackPort)
	if err != nil {
		logger.Fatal("Registering login provider", zap.Error(err))
	}

	// Session management and authentication middleware
	sessionCookieMiddleware := auth.SessionCookieMiddleware(logger, *clientAuthSecretKey, *noSessionTimeout)
	appDetectionMiddleware := auth.DetectorMiddleware(logger, *myHostname, *officeHostname, *tspHostname)
	userAuthMiddleware := authentication.UserAuthMiddleware(logger)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)
	handlerContext.SetCookieSecret(*clientAuthSecretKey)
	if *noSessionTimeout {
		handlerContext.SetNoSessionTimeout()
	}

	if *emailBackend == "ses" {
		// Setup Amazon SES (email) service
		// TODO: This might be able to be combined with the AWS Session that we're using for S3 down
		// below.
		sesSession, err := awssession.NewSession(&aws.Config{
			Region: aws.String(*awsSesRegion),
		})
		if err != nil {
			logger.Fatal("Failed to create a new AWS client config provider", zap.Error(err))
		}
		sesService := ses.New(sesSession)
		handlerContext.SetNotificationSender(notifications.NewNotificationSender(sesService, logger))
	} else {
		handlerContext.SetNotificationSender(notifications.NewStubNotificationSender(logger))
	}

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(*build))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := route.NewHEREPlanner(logger, hereGeoEndpoint, hereRouteEndpoint, hereAppID, hereAppCode)
	handlerContext.SetPlanner(routePlanner)

	var storer storage.FileStorer
	if *storageBackend == "s3" {
		zap.L().Info("Using s3 storage backend")
		if len(*s3Bucket) == 0 {
			log.Fatalln(errors.New("must provide aws_s3_bucket_name parameter, exiting"))
		}
		if *s3Region == "" {
			log.Fatalln(errors.New("Must provide aws_s3_region parameter, exiting"))
		}
		if *s3KeyNamespace == "" {
			log.Fatalln(errors.New("Must provide aws_s3_key_namespace parameter, exiting"))
		}
		aws := awssession.Must(awssession.NewSession(&aws.Config{
			Region: s3Region,
		}))

		storer = storage.NewS3(*s3Bucket, *s3KeyNamespace, logger, aws)
	} else {
		zap.L().Info("Using filesystem storage backend")
		fsParams := storage.DefaultFilesystemParams(logger)
		storer = storage.NewFilesystem(fsParams)
	}
	handlerContext.SetFileStorer(storer)

	// Base routes
	site := goji.NewMux()
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(httpsComplianceMiddleware)
	site.Use(securityHeadersMiddleware)
	site.Use(limitBodySizeMiddleware)

	// Stub health check
	site.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, r *http.Request) {})

	// Allow public content through without any auth or app checks
	site.Handle(pat.Get("/static/*"), clientHandler)
	site.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	site.Handle(pat.Get("/downloads/*"), clientHandler)
	site.Handle(pat.Get("/favicon.ico"), clientHandler)

	ordersMux := goji.SubMux()
	ordersDetectionMiddleware := auth.OrdersDetectorMiddleware(logger, *ordersHostname)
	ordersMux.Use(ordersDetectionMiddleware)
	ordersMux.Use(noCacheMiddleware)
	ordersMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*ordersSwagger))
	ordersMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "orders.html")))
	ordersMux.Handle(pat.New("/*"), ordersapi.NewOrdersAPIHandler(handlerContext))
	site.Handle(pat.Get("/orders/v0/*"), ordersMux)

	root := goji.NewMux()
	root.Use(sessionCookieMiddleware)
	root.Use(appDetectionMiddleware) // Comes after the sessionCookieMiddleware as it sets session state
	root.Use(logging.LogRequestMiddleware)
	site.Handle(pat.New("/*"), root)

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

	errChan := make(chan error)
	localhostCert := server.TLSCert{
		CertPEMBlock: []byte(*httpsCert),
		KeyPEMBlock:  []byte(*httpsKey),
	}
	go func() {
		httpServer := server.Server{
			ListenAddress: *listenInterface,
			HTTPHandler:   httpHandler,
			Logger:        logger,
			Port:          *httpPort,
			PrefixToDrop:  *httpServerLogLinePrefixDrop,
		}
		errChan <- httpServer.ListenAndServe()
	}()
	go func() {
		tlsServer := server.Server{
			ClientAuthType: tls.NoClientCert,
			ListenAddress:  *listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         logger,
			Port:           *httpsPort,
			TLSCerts:       []server.TLSCert{localhostCert},
			PrefixToDrop:   *httpServerLogLinePrefixDrop,
		}
		errChan <- tlsServer.ListenAndServeTLS()
	}()
	go func() {
		mutualTLSServer := server.Server{
			// Only allow certificates validated by the specified
			// client certificate CA.
			ClientAuthType: tls.RequireAndVerifyClientCert,
			CACertPEMBlock: []byte(*httpsClientAuthCACert),
			ListenAddress:  *listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         logger,
			Port:           *httpsClientAuthPort,
			TLSCerts:       []server.TLSCert{localhostCert},
			PrefixToDrop:   *httpServerLogLinePrefixDrop,
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
