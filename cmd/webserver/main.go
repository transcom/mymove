package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/auth/authentication"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/storage"
)

var logger *zap.Logger

// max request body size is 20 mb
const maxBodySize int64 = 200 * 1000 * 1000

// max request headers size is 1 mb
const maxHeaderSize int = 1 * 1000 * 1000

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

func main() {

	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	listenInterface := flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	myHostname := flag.String("http_my_server_name", "localhost", "Hostname according to environment.")
	officeHostname := flag.String("http_office_server_name", "officelocal", "Hostname according to environment.")
	port := flag.String("port", "8080", "the HTTP `port` to listen on.")
	internalSwagger := flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	apiSwagger := flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
	clientAuthSecretKey := flag.String("client_auth_secret_key", "", "Client auth secret JWT key.")
	noSessionTimeout := flag.Bool("no_session_timeout", false, "whether user sessions should timeout.")

	httpsPort := flag.String("https_port", "8443", "the `port` to listen on.")
	httpsCert := flag.String("https_cert", "", "TLS certificate.")
	httpsKey := flag.String("https_key", "", "TLS private key.")

	loginGovCallbackProtocol := flag.String("login_gov_callback_protocol", "https://", "Protocol for non local environments.")
	loginGovCallbackPort := flag.String("login_gov_callback_port", "443", "The port for callback urls.")
	loginGovSecretKey := flag.String("login_gov_secret_key", "", "Login.gov auth secret JWT key.")
	loginGovMyClientID := flag.String("login_gov_my_client_id", "", "Client ID registered with login gov.")
	loginGovOfficeClientID := flag.String("login_gov_office_client_id", "", "Client ID registered with login gov.")
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
	s3Bucket := flag.String("aws_s3_bucket_name", "", "S3 bucket used for file storage")
	s3Region := flag.String("aws_s3_region", "", "AWS region used for S3 file storage")
	s3KeyNamespace := flag.String("aws_s3_key_namespace", "", "Key prefix for all objects written to S3")
	awsSesRegion := flag.String("aws_ses_region", "", "AWS region used for SES")

	flag.Parse()

	logger, err := logging.Config(*env, *debugLogging)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

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
	err = loginGovProvider.RegisterProvider(*myHostname, *loginGovMyClientID, *officeHostname, *loginGovOfficeClientID, *loginGovCallbackProtocol, *loginGovCallbackPort)
	if err != nil {
		logger.Fatal("Registering login provider", zap.Error(err))
	}

	// Session management and authentication middleware
	sessionCookieMiddleware := auth.SessionCookieMiddleware(logger, *clientAuthSecretKey, *noSessionTimeout)
	appDetectionMiddleware := auth.DetectorMiddleware(logger, *myHostname, *officeHostname)
	userAuthMiddleware := authentication.UserAuthMiddleware(logger)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)
	handlerContext.SetCookieSecret(*clientAuthSecretKey)
	if *noSessionTimeout {
		handlerContext.SetNoSessionTimeout()
	}

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
	handlerContext.SetSesService(sesService)

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(*build))

	// Get route planner for handlers to calculate transit distances
	// routePlanner := route.NewBingPlanner(logger, bingMapsEndpoint, bingMapsKey)
	routePlanner := route.NewHEREPlanner(logger, hereGeoEndpoint, hereRouteEndpoint, hereAppID, hereAppCode)
	handlerContext.SetPlanner(routePlanner)

	var storer handlers.FileStorer
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
		absTmpPath, err := filepath.Abs("tmp")
		if err != nil {
			log.Fatalln(errors.New("could not get absolute path for tmp"))
		}
		storagePath := path.Join(absTmpPath, "storage")
		webRoot := "/" + "storage"
		storer = storage.NewFilesystem(storagePath, webRoot, logger)
	}
	handlerContext.SetFileStorer(storer)

	// Base routes
	site := goji.NewMux()
	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	site.Use(httpsComplianceMiddleware)
	site.Use(limitBodySizeMiddleware)

	// Stub health check
	site.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, r *http.Request) {})

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
	externalAPIMux.Handle(pat.New("/*"), handlers.NewPublicAPIHandler(handlerContext))

	internalMux := goji.SubMux()
	root.Handle(pat.New("/internal/*"), internalMux)
	internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*internalSwagger))
	internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "internal.html")))

	// Mux for internal API that enforces auth
	internalAPIMux := goji.SubMux()
	internalMux.Handle(pat.New("/*"), internalAPIMux)
	internalAPIMux.Use(userAuthMiddleware)
	internalAPIMux.Use(noCacheMiddleware)
	internalAPIMux.Handle(pat.New("/*"), handlers.NewInternalAPIHandler(handlerContext))

	authContext := authentication.NewAuthContext(logger, loginGovProvider, *loginGovCallbackProtocol, *loginGovCallbackPort)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), authentication.RedirectHandler{Context: authContext})
	authMux.Handle(pat.Get("/login-gov/callback"), authentication.NewCallbackHandler(authContext, dbConnection, *clientAuthSecretKey, *noSessionTimeout))
	authMux.Handle(pat.Get("/logout"), authentication.NewLogoutHandler(authContext, *clientAuthSecretKey, *noSessionTimeout))

	if *env == "development" {
		zap.L().Info("Enabling devlocal auth")
		localAuthMux := goji.SubMux()
		root.Handle(pat.New("/devlocal-auth/*"), localAuthMux)
		localAuthMux.Handle(pat.Get("/login"), authentication.NewUserListHandler(authContext, dbConnection))
		localAuthMux.Handle(pat.Post("/login"), authentication.NewAssignUserHandler(authContext, dbConnection, *clientAuthSecretKey, *noSessionTimeout))
	}

	if *storageBackend == "filesystem" {
		// Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler("tmp")
		root.Handle(pat.Get("/storage/*"), fs)
	}

	root.Handle(pat.Get("/static/*"), clientHandler)
	root.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	root.Handle(pat.Get("/downloads/*"), clientHandler)
	root.Handle(pat.Get("/favicon.ico"), clientHandler)
	root.HandleFunc(pat.Get("/*"), fileHandler(path.Join(*build, "index.html")))

	// Start http/https listener(s)
	errChan := make(chan error)
	go func() { // start http listener
		addr := fmt.Sprintf("%s:%s", *listenInterface, *port)
		zap.L().Info("Starting http server listening", zap.String("address", addr))
		s := &http.Server{
			Addr:           addr,
			Handler:        site,
			MaxHeaderBytes: maxHeaderSize,
		}
		errChan <- s.ListenAndServe()
	}()
	go func() { // start https listener
		addr := fmt.Sprintf("%s:%s", *listenInterface, *httpsPort)
		zap.L().Info("Starting https server listening", zap.String("address", addr))
		errChan <- listenAndServeTLS(addr, []byte(*httpsCert), []byte(*httpsKey), site)
	}()
	log.Fatal(<-errChan)
}

// fileHandler serves up a single file
func fileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}

func listenAndServeTLS(addr string, certPEMBlock, keyPEMBlock []byte, handler http.Handler) error {
	// Configure TLS
	cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2"}, // enable HTTP/2
	}

	// Create listener
	ln, err := tls.Listen("tcp", addr, config)
	if err != nil {
		return err
	}
	defer ln.Close()

	// Start server
	srv := &http.Server{Addr: addr, Handler: handler, MaxHeaderBytes: maxHeaderSize}
	return srv.Serve(ln)
}
