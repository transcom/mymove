package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/markbates/pop"
	"go.uber.org/zap"
	// "goji.io"
	// "goji.io/pat"
	"github.com/go-openapi/loads"

	"github.com/transcom/mymove/pkg/gen/genserver"
	"github.com/transcom/mymove/pkg/gen/genserver/operations"
	"github.com/transcom/mymove/pkg/gen/genserver/operations/issues"
	"github.com/transcom/mymove/pkg/handlers"
)

var logger *zap.Logger

// TODO(nick - 12/21/17) - this is a simple logger for debugging testing
// It needs replacing with something we can use in production
func requestLogger(h http.Handler) http.Handler {
	zap.L().Info("Request logger installed")
	wrapper := func(w http.ResponseWriter, r *http.Request) {
		zap.L().Info("Request", zap.String("url", r.URL.String()))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(wrapper)
}

func main() {

	// entry := flag.String("entry", "build/index.html", "the entrypoint to serve.")
	// build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presenetly.")
	port := flag.String("port", "8080", "the `port` to listen on.")
	// swagger := flag.String("swagger", "swagger.yaml", "The location of the swagger API definition")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
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

	//DB connection
	pop.AddLookupPaths(*config)
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	// initialize api pkg with dbConnection created above
	handlers.Init(dbConnection)

	swaggerSpec, err := loads.Analyzed(genserver.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewMymoveAPI(swaggerSpec)
	api.Logger = log.Printf

	api.IssuesCreateIssueHandler = issues.CreateIssueHandlerFunc(handlers.CreateIssueHandler)

	server := genserver.NewServer(api)
	server.Port, err = strconv.Atoi(*port)
	if err != nil {
		log.Fatal(err)
	}
	defer server.Shutdown()

	zap.L().Info("Starting the GEN server listening", zap.String("port", *port))
	server.Serve()

	// // Serves files out of build folder
	// fileHandler := http.FileServer(http.Dir(*build))

	// // api routes
	// api := api.Mux()

	// // Base routes
	// root := goji.NewMux()
	// root.Handle(pat.New("/api/*"), api)
	// root.Handle(pat.Get("/static/*"), fileHandler)
	// root.Handle(pat.Get("/favicon.ico"), fileHandler)
	// root.HandleFunc(pat.Get("/*"), IndexHandler(entry))

	// // And request logging
	// root.Use(requestLogger)

	// zap.L().Info("Starting the server listening", zap.String("port", *port))
	// http.ListenAndServe(*port, root)
}

// IndexHandler serves up our index.html
func IndexHandler(entrypoint *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *entrypoint)
	}
}
