package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"

	aws "github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/storage"
)

var logger *zap.Logger

// TODO(nick - 12/21/17) - this is a simple logger for debugging testing
// It needs replacing with something we can use in production
func requestLoggerMiddleware(inner http.Handler) http.Handler {
	zap.L().Info("requestLoggerMiddleware installed")
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("x-forwarded-for", r.Header.Get("x-forwarded-for")),
			zap.String("x-forwarded-host", r.Header.Get("x-forwarded-host")),
			zap.String("x-forwarded-proto", r.Header.Get("x-forwarded-proto")),
		)
		inner.ServeHTTP(w, r)
	}
	return http.HandlerFunc(wrapper)
}

// max request body size is 20 mb
const maxBodySize int64 = 200 * 1000 * 1000

// max request headers size is 1 mb
const maxHeaderSize int = 1 * 1000 * 1000

func limitBodySizeMiddleware(inner http.Handler) http.Handler {
	zap.L().Info("limitBodySizeMiddleware installed")
	mw := func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		inner.ServeHTTP(w, r)
		return
	}
	return http.HandlerFunc(mw)
}

func main() {

	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	listenInterface := flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	protocol := flag.String("protocol", "https://", "Protocol for non local environments.")
	hostname := flag.String("http_server_name", "localhost", "Hostname according to environment.")
	port := flag.String("port", "8080", "the HTTP `port` to listen on.")
	callbackPort := flag.String("callback_port", "443", "The port for callback urls.")
	internalSwagger := flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	apiSwagger := flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
	clientAuthSecretKey := flag.String("client_auth_secret_key", "", "Client auth secret JWT key.")
	noSessionTimeout := flag.Bool("no_session_timeout", false, "whether user sessions should timeout.")

	httpsPort := flag.String("https_port", "8443", "the `port` to listen on.")
	httpsCert := flag.String("https_cert", "", "TLS certificate.")
	httpsKey := flag.String("https_key", "", "TLS private key.")

	loginGovSecretKey := flag.String("login_gov_secret_key", "", "Login.gov auth secret JWT key.")
	loginGovClientID := flag.String("login_gov_client_id", "", "Client ID registered with login gov.")
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

	flag.Parse()

	// Set up logger for the system
	var err error
	if *debugLogging {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	// Assert that our secret keys can be parsed into actual private keys
	// TODO: Store the parsed key in handlers/AppContext instead of parsing every time
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*loginGovSecretKey)); err != nil {
		log.Fatalln(err)
	}
	if _, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(*clientAuthSecretKey)); err != nil {
		log.Fatalln(err)
	}
	if *loginGovHostname == "" {
		log.Fatalln(errors.New("Must provide the Login.gov hostname parameter, exiting"))
	}

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(*build))

	// Register Login.gov authentication provider
	if *env == "development" {
		*protocol = "http://"
		*callbackPort = "3000"
	}
	fullHostname := fmt.Sprintf("%s%s:%s", *protocol, *hostname, *callbackPort)
	loginGovProvider := auth.NewLoginGovProvider(*loginGovHostname, *loginGovSecretKey, *loginGovClientID, logger)
	loginGovProvider.RegisterProvider(fullHostname)

	// Populates user info using cookie and renews token
	tokenMiddleware := auth.TokenParsingMiddleware(logger, *clientAuthSecretKey, *noSessionTimeout)
	userAuthMiddleware := auth.UserAuthMiddleware(dbConnection)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger)

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
		aws := awssession.Must(awssession.NewSession(&aws.Config{
			Region: s3Region,
		}))

		storer = storage.NewS3(*s3Bucket, logger, aws)
	} else {
		zap.L().Info("Using filesystem storage backend")
		absTmpPath, err := filepath.Abs("tmp")
		if err != nil {
			log.Fatalln(errors.New("could not get absolute path for tmp"))
		}
		storagePath := path.Join(absTmpPath, "storage")
		webRoot := fullHostname + "/" + "storage"
		storer = storage.NewFilesystem(storagePath, webRoot, logger)
	}
	handlerContext.SetFileStorer(storer)

	// Base routes
	root := goji.NewMux()

	// Add middleware: they are evaluated in the reverse order in which they
	// are added, but the resulting http.Handlers execute in "normal" order
	// (i.e., the http.Handler returned by the first Middleware added gets
	// called first).
	root.Use(requestLoggerMiddleware)
	root.Use(limitBodySizeMiddleware)
	root.Use(tokenMiddleware)

	// Stub health check
	root.HandleFunc(pat.Get("/health"), func(w http.ResponseWriter, r *http.Request) {})

	apiMux := goji.SubMux()
	root.Handle(pat.New("/api/v1/*"), apiMux)
	apiMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*apiSwagger))
	apiMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "api.html")))
	apiMux.Handle(pat.New("/*"), handlers.NewPublicAPIHandler(handlerContext))

	internalMux := goji.SubMux()
	root.Handle(pat.New("/internal/*"), internalMux)
	internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*internalSwagger))
	internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "internal.html")))

	// Mux for internal API that enforces auth
	internalAPIMux := goji.SubMux()
	internalAPIMux.Use(userAuthMiddleware)
	internalMux.Handle(pat.New("/*"), internalAPIMux)
	internalAPIMux.Handle(pat.New("/*"), handlers.NewInternalAPIHandler(handlerContext))

	authContext := auth.NewAuthContext(fullHostname, logger, loginGovProvider)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), auth.AuthorizationRedirectHandler(authContext))
	authMux.Handle(pat.Get("/login-gov/callback"), auth.NewAuthorizationCallbackHandler(dbConnection, *clientAuthSecretKey, *noSessionTimeout, fullHostname, logger, loginGovProvider))
	authMux.Handle(pat.Get("/logout"), auth.AuthorizationLogoutHandler(authContext))

	if *storageBackend == "filesystem" {
		// Add a file handler to provide access to files uploaded in development
		fs := storage.NewFilesystemHandler("tmp")
		root.Handle(pat.Get("/storage/*"), fs)
	}

	root.Handle(pat.Get("/static/*"), clientHandler)
	root.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	root.Handle(pat.Get("/favicon.ico"), clientHandler)
	root.HandleFunc(pat.Get("/*"), fileHandler(path.Join(*build, "index.html")))

	// Start http/https listener(s)
	errChan := make(chan error)
	go func() { // start http listener
		addr := fmt.Sprintf("%s:%s", *listenInterface, *port)
		zap.L().Info("Starting http server listening", zap.String("address", addr))
		s := &http.Server{
			Addr:           addr,
			Handler:        root,
			MaxHeaderBytes: maxHeaderSize,
		}
		errChan <- s.ListenAndServe()
	}()
	go func() { // start https listener
		addr := fmt.Sprintf("%s:%s", *listenInterface, *httpsPort)
		zap.L().Info("Starting https server listening", zap.String("address", addr))
		errChan <- listenAndServeTLS(addr, []byte(*httpsCert), []byte(*httpsKey), root)
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
