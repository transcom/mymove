package main

import (
	"fmt"
	"log"
	"net/http"
	"path"

	"github.com/go-openapi/loads"
	"github.com/markbates/pop"
	"github.com/namsral/flag" // This flag package accepts ENV vars as well as cmd line flags
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"github.com/transcom/mymove/pkg/gen/restapi"
	"github.com/transcom/mymove/pkg/gen/restapi/operations"
	form1299op "github.com/transcom/mymove/pkg/gen/restapi/operations/form1299s"
	issueop "github.com/transcom/mymove/pkg/gen/restapi/operations/issues"
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

	build := flag.String("build", "build", "the directory to serve static files from.")
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presenetly.")
	listenInterface := flag.String("interface", "", "The interface spec to listen for connections on. Default is all.")
	port := flag.String("port", "8080", "the `port` to listen on.")
	swagger := flag.String("swagger", "swagger.yaml", "The location of the swagger API definition")
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

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewMymoveAPI(swaggerSpec)

	api.IssuesCreateIssueHandler = issueop.CreateIssueHandlerFunc(handlers.CreateIssueHandler)
	api.IssuesIndexIssuesHandler = issueop.IndexIssuesHandlerFunc(handlers.IndexIssuesHandler)
	api.Form1299sCreateForm1299Handler = form1299op.CreateForm1299HandlerFunc(handlers.CreateForm1299Handler)
	api.Form1299sIndexForm1299sHandler = form1299op.IndexForm1299sHandlerFunc(handlers.IndexForm1299sHandler)

	// Serves files out of build folder
	clientHandler := http.FileServer(http.Dir(*build))

	// Base routes
	root := goji.NewMux()
	root.Handle(pat.Get("/api/v1/swagger.yaml"), fileHandler(*swagger))
	root.Handle(pat.Get("/api/v1/docs"), fileHandler(path.Join(*build, "swagger-ui", "index.html")))
	root.Handle(pat.New("/api/*"), api.Serve(nil)) // Serve(nil) returns an http.Handler for the swagger api
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
