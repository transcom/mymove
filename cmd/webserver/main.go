package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/pop"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/storage"
)

var logger *zap.Logger

// TODO(nick - 12/21/17) - this is a simple logger for debugging testing
// It needs replacing with something we can use in production
func requestLogger(h http.Handler) http.Handler {
	zap.L().Info("Request logger installed")
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request", zap.String("method", r.Method), zap.String("url", r.URL.String()))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(wrapper)
}

func main() {

	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presently.")
	listenInterface := flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	protocol := flag.String("protocol", "https://", "Protocol for non local environments.")
	hostname := flag.String("http_server_name", "localhost", "Hostname according to environment.")
	port := flag.String("port", "8080", "the `port` to listen on.")
	callbackPort := flag.String("callback_port", "443", "The port for callback urls.")
	internalSwagger := flag.String("internal-swagger", "swagger/internal.yaml", "The location of the internal API swagger definition")
	apiSwagger := flag.String("swagger", "swagger/api.yaml", "The location of the public API swagger definition")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
	clientAuthSecretKey := flag.String("client_auth_secret_key", "", "Client auth secret JWT key.")
	noSessionTimeout := flag.Bool("no_session_timeout", false, "whether user sessions should timeout.")

	loginGovSecretKey := flag.String("login_gov_secret_key", "", "Login.gov auth secret JWT key.")
	loginGovClientID := flag.String("login_gov_client_id", "", "Client ID registered with login gov.")
	loginGovHostname := flag.String("login_gov_hostname", "", "Hostname for communicating with login gov.")

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

	bucket := os.Getenv("AWS_S3_BUCKET_NAME")
	if len(bucket) == 0 {
		log.Fatalln("AWS_S3_BUCKET_NAME not configured")
	}

	aws := awssession.Must(awssession.NewSession())
	storer := storage.NewS3(bucket, logger, aws)

	handlerContext := handlers.NewHandlerContext(dbConnection, logger, storer)

	// Base routes
	root := goji.NewMux()
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
	internalMux.Use(auth.RequireAuthMiddleware)
	internalMux.Handle(pat.Get("/swagger.yaml"), fileHandler(*internalSwagger))
	internalMux.Handle(pat.Get("/docs"), fileHandler(path.Join(*build, "swagger-ui", "internal.html")))
	internalMux.Handle(pat.New("/*"), handlers.NewInternalAPIHandler(handlerContext))

	authContext := auth.NewAuthContext(fullHostname, logger, loginGovProvider)
	authMux := goji.SubMux()
	root.Handle(pat.New("/auth/*"), authMux)
	authMux.Handle(pat.Get("/login-gov"), auth.AuthorizationRedirectHandler(authContext))
	authMux.Handle(pat.Get("/login-gov/callback"), auth.NewAuthorizationCallbackHandler(dbConnection, *clientAuthSecretKey, *noSessionTimeout, fullHostname, logger, loginGovProvider))
	authMux.Handle(pat.Get("/logout"), auth.AuthorizationLogoutHandler(authContext))

	root.Handle(pat.Get("/static/*"), clientHandler)
	root.Handle(pat.Get("/swagger-ui/*"), clientHandler)
	root.Handle(pat.Get("/favicon.ico"), clientHandler)
	root.HandleFunc(pat.Get("/*"), fileHandler(path.Join(*build, "index.html")))

	// And request logging
	root.Use(requestLogger)

	address := fmt.Sprintf("%s:%s", *listenInterface, *port)
	zap.L().Info("Starting the server listening", zap.String("address", address))
	log.Fatal(http.ListenAndServe(address, root))
}

// fileHandler serves up a single file
func fileHandler(entrypoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
}
