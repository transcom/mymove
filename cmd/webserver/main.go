package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/dpsauth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/dpsapi"
	"github.com/transcom/mymove/pkg/handlers/internalapi"
	"github.com/transcom/mymove/pkg/handlers/ordersapi"
	"github.com/transcom/mymove/pkg/handlers/publicapi"
	"github.com/transcom/mymove/pkg/iws"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/notifications"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/server"
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

func initFlags(flag *pflag.FlagSet) {

	flag.String("build", "build", "the directory to serve static files from.")
	flag.String("config-dir", "config", "The location of server config files")
	flag.String("env", "development", "The environment to run in, which configures the database.")
	flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	flag.String("service-name", "app", "The service name identifies the application for instrumentation.")

	flag.String("http-my-server-name", "localhost", "Hostname according to environment.")
	flag.String("http-office-server-name", "officelocal", "Hostname according to environment.")
	flag.String("http-tsp-server-name", "tsplocal", "Hostname according to environment.")
	flag.String("http-orders-server-name", "orderslocal", "Hostname according to environment.")
	flag.String("http-dps-server-name", "dpslocal", "Hostname according to environment.")

	// SDDC + DPS Auth config
	flag.String("http-sddc-server-name", "sddclocal", "Hostname according to envrionment.")
	flag.String("http-sddc-protocol", "https", "Protocol for sddc")
	flag.String("http-sddc-port", "", "The port for sddc")
	flag.String("dps-auth-secret-key", "", "DPS auth JWT secret key")
	flag.String("dps-redirect-url", "", "DPS url to redirect to")
	flag.String("dps-cookie-name", "", "Name of the DPS cookie")
	flag.String("dps-cookie-domain", "sddclocal", "Domain of the DPS cookie")

	// Initialize Swagger
	flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	flag.String("orders-swagger", "swagger/orders.yaml", "The location of the Orders API swagger definition")
	flag.String("dps-swagger", "swagger/dps.yaml", "The location of the DPS API swagger definition")

	flag.Bool("debug-logging", false, "log messages at the debug level.")
	flag.String("client-auth-secret-key", "", "Client auth secret JWT key.")
	flag.Bool("no-session-timeout", false, "whether user sessions should timeout.")

	flag.String("dod-ca-package", "", "Path to PKCS#7 package containing certificates of all DoD root and intermediate CAs")
	flag.String("move-mil-dod-ca-cert", "", "The DoD CA certificate used to sign the move.mil TLS certificate.")
	flag.String("move-mil-dod-tls-cert", "", "The DoD-signed TLS certificate for various move.mil services.")
	flag.String("move-mil-dod-tls-key", "", "The private key for the DoD-signed TLS certificate for various move.mil services.")

	// Ports to listen to
	flag.Int("mutual-tls-port", 9443, "The `port` for the mutual TLS listener.")
	flag.Int("tls-port", 8443, "the `port` for the server side TLS listener.")
	flag.Int("no-tls-port", 8080, "the `port` for the listener not requiring any TLS.")

	// Login.Gov config
	flag.String("login-gov-callback-protocol", "https://", "Protocol for non local environments.")
	flag.Int("login-gov-callback-port", 443, "The port for callback urls.")
	flag.String("login-gov-secret-key", "", "Login.gov auth secret JWT key.")
	flag.String("login-gov-my-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-office-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-tsp-client-id", "", "Client ID registered with login gov.")
	flag.String("login-gov-hostname", "", "Hostname for communicating with login gov.")

	/* For bing Maps use the following
	bingMapsEndpoint := flag.String("bing_maps_endpoint", "", "URL for the Bing Maps Truck endpoint to use")
	bingMapsKey := flag.String("bing_maps_key", "", "Authentication key to use for the Bing Maps endpoint")
	*/

	// HERE Maps Config
	flag.String("here-maps-geocode-endpoint", "", "URL for the HERE maps geocoder endpoint")
	flag.String("here-maps-routing-endpoint", "", "URL for the HERE maps routing endpoint")
	flag.String("here-maps-app-id", "", "HERE maps App ID for this application")
	flag.String("here-maps-app-code", "", "HERE maps App API code")

	// EDI Invoice Config
	flag.Bool("send-prod-invoice", false, "Flag (bool) for EDI Invoices to signify if they should go to production GEX")

	flag.String("storage-backend", "filesystem", "Storage backend to use, either filesystem or s3.")
	flag.String("email-backend", "local", "Email backend to use, either SES or local")
	flag.String("aws-s3-bucket-name", "", "S3 bucket used for file storage")
	flag.String("aws-s3-region", "", "AWS region used for S3 file storage")
	flag.String("aws-s3-key-namespace", "", "Key prefix for all objects written to S3")
	flag.String("aws-ses-region", "", "AWS region used for SES")

	// New Relic Config
	flag.String("new-relic-application-id", "", "App ID for New Relic Browser")
	flag.String("new-relic-license-key", "", "License key for New Relic Browser")

	// Honeycomb Config
	flag.Bool("honeycomb-enabled", false, "Honeycomb enabled")
	flag.String("honeycomb-api-key", "", "API Key for Honeycomb")
	flag.String("honeycomb-dataset", "", "Dataset for Honeycomb")
	flag.Bool("honeycomb-debug", false, "Debug honeycomb using stdout.")

	// IWS
	flag.String("iws-rbs-host", "", "Hostname for the IWS RBS")

	// DB Config
	flag.String("db-name", "dev_db", "Database Name")
	flag.String("db-host", "localhost", "Database Hostname")
	flag.Int("db-port", 5432, "Database Port")
	flag.String("db-user", "postgres", "Database Username")
	flag.String("db-password", "", "Database Password")
}

func initDODCertificates(v *viper.Viper, logger *zap.Logger) ([]server.TLSCert, *x509.CertPool, error) {

	moveMilCerts := []server.TLSCert{
		server.TLSCert{
			//Append move.mil cert with CA certificate chain
			CertPEMBlock: bytes.Join([][]byte{
				[]byte(v.GetString("move-mil-dod-tls-cert")),
				[]byte(v.GetString("move-mil-dod-ca-cert"))},
				[]byte("\n"),
			),
			KeyPEMBlock: []byte(v.GetString("move-mil-dod-tls-key")),
		},
	}

	pkcs7Package, err := ioutil.ReadFile(v.GetString("dod-ca-package")) // #nosec
	if err != nil {
		return moveMilCerts, nil, errors.Wrap(err, "Failed to read DoD CA certificate package")
	}

	dodCACertPool, err := server.LoadCertPoolFromPkcs7Package(pkcs7Package)
	if err != nil {
		return moveMilCerts, dodCACertPool, errors.Wrap(err, "Failed to parse DoD CA certificate package")
	}

	return moveMilCerts, dodCACertPool, nil
}

func initRoutePlanner(v *viper.Viper, logger *zap.Logger) route.Planner {
	return route.NewHEREPlanner(
		logger,
		v.GetString("here-maps-geocode-endpoint"),
		v.GetString("here-maps-routing-endpoint"),
		v.GetString("here-maps-app-id"),
		v.GetString("here-maps-app-code"))
}

func initHoneycomb(v *viper.Viper, logger *zap.Logger) bool {

	honeycombAPIKey := v.GetString("honeycomb-api-key")
	honeycombDataset := v.GetString("honeycomb-dataset")
	honeycombServiceName := v.GetString("service-name")

	if v.GetBool("honeycomb-enabled") && len(honeycombAPIKey) > 0 && len(honeycombDataset) > 0 {
		logger.Debug("Honeycomb Integration enabled", zap.String("honeycomb-dataset", honeycombDataset))
		beeline.Init(beeline.Config{
			WriteKey:    honeycombAPIKey,
			Dataset:     honeycombDataset,
			Debug:       v.GetBool("honeycomb-debug"),
			ServiceName: honeycombServiceName,
		})
		return true
	}

	logger.Debug("Honeycomb Integration disabled")
	return false
}

func initRealTimeBrokerService(v *viper.Viper, logger *zap.Logger) (*iws.RealTimeBrokerService, error) {
	return iws.NewRealTimeBrokerService(
		v.GetString("iws-rbs-host"),
		v.GetString("dod-ca-package"),
		v.GetString("move-mil-dod-tls-cert"),
		v.GetString("move-mil-dod-tls-key"))
}

func initDatabase(v *viper.Viper, logger *zap.Logger) (*pop.Connection, error) {

	env := v.GetString("env")
	dbName := v.GetString("db-name")
	dbHost := v.GetString("db-host")
	dbPort := strconv.Itoa(v.GetInt("db-port"))
	dbUser := v.GetString("db-user")
	dbPassword := v.GetString("db-password")

	// Modify DB options by environment
	dbOptions := map[string]string{"sslmode": "disable"}
	if env == "test" {
		// Leave the test database name hardcoded, since we run tests in the same
		// environment as development, and it's extra confusing to have to swap env
		// variables before running tests.
		dbName = "test_db"
	} else if env == "container" {
		// Require sslmode for containers
		dbOptions["sslmode"] = "require"
	}

	// Construct a safe URL and log it
	s := "postgres://%s:%s@%s:%s/%s?sslmode=%s"
	dbURL := fmt.Sprintf(s, dbUser, "*****", dbHost, dbPort, dbName, dbOptions["sslmode"])
	logger.Debug("Connecting to the database", zap.String("url", dbURL))

	// Configure DB connection details
	dbConnectionDetails := pop.ConnectionDetails{
		Dialect:  "postgres",
		Database: dbName,
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		Options:  dbOptions,
	}
	err := dbConnectionDetails.Finalize()
	if err != nil {
		logger.Error("Failed to finalize DB connection details", zap.Error(err))
		return nil, err
	}

	// Set up the connection
	connection, err := pop.NewConnection(&dbConnectionDetails)
	if err != nil {
		logger.Error("Failed create DB connection", zap.Error(err))
		return nil, err
	}

	// Open the connection
	err = connection.Open()
	if err != nil {
		logger.Error("Failed to open DB connection", zap.Error(err))
		return nil, err
	}

	// Return the open connection
	return connection, nil
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	flag.Parse(os.Args[1:])

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	env := v.GetString("env")

	logger, err := logging.Config(env, v.GetBool("debug-logging"))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// Honeycomb
	useHoneycomb := initHoneycomb(v, logger)

	clientAuthSecretKey := v.GetString("client-auth-secret-key")

	loginGovCallbackProtocol := v.GetString("login-gov-callback-protocol")
	loginGovCallbackPort := v.GetInt("login-gov-callback-port")
	loginGovSecretKey := v.GetString("login-gov-secret-key")
	loginGovHostname := v.GetString("login-gov-hostname")

	// Assert that our secret keys can be parsed into actual private keys
	// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(loginGovSecretKey)); err != nil {
		logger.Fatal("Login.gov private key", zap.Error(err))
	}
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(clientAuthSecretKey)); err != nil {
		logger.Fatal("Client auth private key", zap.Error(err))
	}
	if len(loginGovHostname) == 0 {
		log.Fatal("Must provide the Login.gov hostname parameter, exiting")
	}

	// Create a connection to the DB
	dbConnection, err := initDatabase(v, logger)
	if err != nil {
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	myHostname := v.GetString("http-my-server-name")
	officeHostname := v.GetString("http-office-server-name")
	tspHostname := v.GetString("http-tsp-server-name")

	// Register Login.gov authentication provider for My.(move.mil)
	loginGovProvider := authentication.NewLoginGovProvider(loginGovHostname, loginGovSecretKey, logger)
	err = loginGovProvider.RegisterProvider(
		myHostname,
		v.GetString("login-gov-my-client-id"),
		officeHostname,
		v.GetString("login-gov-office-client-id"),
		tspHostname,
		v.GetString("login-gov-tsp-client-id"),
		loginGovCallbackProtocol,
		loginGovCallbackPort)
	if err != nil {
		logger.Fatal("Registering login provider", zap.Error(err))
	}

	// Session management and authentication middleware
	noSessionTimeout := v.GetBool("no-session-timeout")
	sessionCookieMiddleware := auth.SessionCookieMiddleware(logger, clientAuthSecretKey, noSessionTimeout)
	appDetectionMiddleware := auth.DetectorMiddleware(logger, myHostname, officeHostname, tspHostname)
	userAuthMiddleware := authentication.UserAuthMiddleware(logger)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)
	handlerContext.SetCookieSecret(clientAuthSecretKey)
	if noSessionTimeout {
		handlerContext.SetNoSessionTimeout()
	}

	if v.GetString("email-backend") == "ses" {
		// Setup Amazon SES (email) service
		// TODO: This might be able to be combined with the AWS Session that we're using for S3 down
		// below.
		sesSession, err := awssession.NewSession(&aws.Config{
			Region: aws.String(v.GetString("aws-ses-region")),
		})
		if err != nil {
			logger.Fatal("Failed to create a new AWS client config provider", zap.Error(err))
		}
		sesService := ses.New(sesSession)
		handlerContext.SetNotificationSender(notifications.NewNotificationSender(sesService, logger))
	} else {
		handlerContext.SetNotificationSender(notifications.NewStubNotificationSender(logger))
	}

	build := v.GetString("build")

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(build))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := initRoutePlanner(v, logger)
	handlerContext.SetPlanner(routePlanner)

	// Set SendProductionInvoice for ediinvoice
	handlerContext.SetSendProductionInvoice(v.GetBool("send-prod-invoice"))

	storageBackend := v.GetString("storage-backend")

	var storer storage.FileStorer
	if storageBackend == "s3" {
		zap.L().Info("Using s3 storage backend")
		awsS3Bucket := v.GetString("aws-s3-bucket-name")
		if len(awsS3Bucket) == 0 {
			log.Fatalln(errors.New("must provide aws-s3-bucket-name parameter, exiting"))
		}
		awsS3Region := v.GetString("aws-s3-region")
		if len(awsS3Region) == 0 {
			log.Fatalln(errors.New("Must provide aws-s3-region parameter, exiting"))
		}
		awsS3KeyNamespace := v.GetString("aws-s3-key-namespace")
		if len(awsS3KeyNamespace) == 0 {
			log.Fatalln(errors.New("Must provide aws_s3_key_namespace parameter, exiting"))
		}
		aws := awssession.Must(awssession.NewSession(&aws.Config{
			Region: aws.String(awsS3Region),
		}))
		storer = storage.NewS3(awsS3Bucket, awsS3KeyNamespace, logger, aws)
	} else {
		zap.L().Info("Using filesystem storage backend")
		fsParams := storage.DefaultFilesystemParams(logger)
		storer = storage.NewFilesystem(fsParams)
	}
	handlerContext.SetFileStorer(storer)

	rbs, err := initRealTimeBrokerService(v, logger)
	if err != nil {
		logger.Fatal("Could not instantiate IWS RBS", zap.Error(err))
	}
	handlerContext.SetIWSRealTimeBrokerService(*rbs)

	sddcHostname := v.GetString("http-sddc-server-name")
	dpsAuthSecretKey := v.GetString("dps-auth-secret-key")
	handlerContext.SetDPSAuthParams(
		dpsauth.Params{
			SDDCProtocol:   v.GetString("http-sddc-protocol"),
			SDDCHostname:   sddcHostname,
			SDDCPort:       v.GetString("http-sddc-port"),
			SecretKey:      dpsAuthSecretKey,
			DPSRedirectURL: v.GetString("dps-redirect-url"),
			CookieName:     v.GetString("dps-cookie-name"),
		},
	)

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
	ordersDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, v.GetString("http-orders-server-name"))
	ordersMux.Use(ordersDetectionMiddleware)
	ordersMux.Use(noCacheMiddleware)
	ordersMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("orders-swagger")))
	ordersMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "orders.html")))
	ordersMux.Handle(pat.New("/*"), ordersapi.NewOrdersAPIHandler(handlerContext))
	site.Handle(pat.Get("/orders/v0/*"), ordersMux)

	dpsMux := goji.SubMux()
	dpsDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, v.GetString("http-dps-server-name"))
	dpsMux.Use(dpsDetectionMiddleware)
	dpsMux.Use(noCacheMiddleware)
	dpsMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("dps-swagger")))
	dpsMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "dps.html")))
	dpsMux.Handle(pat.New("/*"), dpsapi.NewDPSAPIHandler(handlerContext))
	site.Handle(pat.New("/dps/v0/*"), dpsMux)

	sddcDPSMux := goji.SubMux()
	sddcDetectionMiddleware := auth.HostnameDetectorMiddleware(logger, sddcHostname)
	sddcDPSMux.Use(sddcDetectionMiddleware)
	sddcDPSMux.Use(noCacheMiddleware)
	site.Handle(pat.New("/dps_auth/*"), sddcDPSMux)
	sddcDPSMux.Handle(pat.Get("/set_cookie"), dpsauth.NewSetCookieHandler(logger, dpsAuthSecretKey, v.GetString("dps-cookie-domain")))

	root := goji.NewMux()
	root.Use(sessionCookieMiddleware)
	root.Use(appDetectionMiddleware) // Comes after the sessionCookieMiddleware as it sets session state
	root.Use(logging.LogRequestMiddleware)
	site.Handle(pat.New("/*"), root)

	apiMux := goji.SubMux()
	root.Handle(pat.New("/api/v1/*"), apiMux)
	apiMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("swagger")))
	apiMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "api.html")))

	externalAPIMux := goji.SubMux()
	apiMux.Handle(pat.New("/*"), externalAPIMux)
	externalAPIMux.Use(noCacheMiddleware)
	externalAPIMux.Use(userAuthMiddleware)
	externalAPIMux.Handle(pat.New("/*"), publicapi.NewPublicAPIHandler(handlerContext))

	internalMux := goji.SubMux()
	root.Handle(pat.New("/internal/*"), internalMux)
	internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(v.GetString("internal-swagger")))
	internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(build, "swagger-ui", "internal.html")))

	// Mux for internal API that enforces auth
	internalAPIMux := goji.SubMux()
	internalMux.Handle(pat.New("/*"), internalAPIMux)
	internalAPIMux.Use(userAuthMiddleware)
	internalAPIMux.Use(noCacheMiddleware)
	internalAPIMux.Handle(pat.New("/*"), internalapi.NewInternalAPIHandler(handlerContext))

	authContext := authentication.NewAuthContext(logger, loginGovProvider, loginGovCallbackProtocol, loginGovCallbackPort)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), authentication.RedirectHandler{Context: authContext})
	authMux.Handle(pat.Get("/login-gov/callback"), authentication.NewCallbackHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))
	authMux.Handle(pat.Get("/logout"), authentication.NewLogoutHandler(authContext, clientAuthSecretKey, noSessionTimeout))

	if env == "development" || env == "test" {
		zap.L().Info("Enabling devlocal auth")
		localAuthMux := goji.SubMux()
		root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
		localAuthMux.Handle(pat.Get("/login"), authentication.NewUserListHandler(authContext, dbConnection))
		localAuthMux.Handle(pat.Post("/login"), authentication.NewAssignUserHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))
		localAuthMux.Handle(pat.Post("/new"), authentication.NewCreateUserHandler(authContext, dbConnection, clientAuthSecretKey, noSessionTimeout))
	}

	if storageBackend == "filesystem" {
		// Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler("tmp")
		root.Handle(pat.Get("/storage/*"), fs)
	}

	// Serve index.html to all requests that haven't matches a previous route,
	root.HandleFunc(pat.Get("/*"), indexHandler(build, v.GetString("new-relic-application-id"), v.GetString("new-relic-license-key"), logger))

	var httpHandler http.Handler
	if useHoneycomb {
		httpHandler = hnynethttp.WrapHandler(site)
	} else {
		httpHandler = site
	}

	errChan := make(chan error)

	moveMilCerts, dodCACertPool, err := initDODCertificates(v, logger)
	if err != nil {
		logger.Fatal("Failed to initialize DOD certificates", zap.Error(err))
	}

	listenInterface := v.GetString("interface")

	go func() {
		noTLSServer := server.Server{
			ListenAddress: listenInterface,
			HTTPHandler:   httpHandler,
			Logger:        logger,
			Port:          v.GetInt("no-tls-port"),
		}
		errChan <- noTLSServer.ListenAndServe()
	}()

	go func() {
		tlsServer := server.Server{
			ClientAuthType: tls.NoClientCert,
			ListenAddress:  listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         logger,
			Port:           v.GetInt("tls-port"),
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
			ListenAddress:  listenInterface,
			HTTPHandler:    httpHandler,
			Logger:         logger,
			Port:           v.GetInt("mutual-tls-port"),
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
