package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/markbates/pop"
	"go.uber.org/zap"
	"goji.io"
	"goji.io/pat"

	"dp3/pkg/api"
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

	entry := flag.String("entry", "client/build/index.html", "the entrypoint to serve.")
	build := flag.String("build", "client/build", "the directory to serve static files from.")
	config := flag.String("config-dir", "server/src/dp3/config", "The location of server config files")
	port := flag.String("port", ":8080", "the `port` to listen on.")
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
	dbConnection, err := pop.Connect("development")
	if err != nil {
		log.Panic(err)
	}

	// initialize api pkg with dbConnection created above
	api.Init(dbConnection)

	// Serves files out of build folder
	fileHandler := http.FileServer(http.Dir(*build))

	// api routes
	api := api.Mux()

	// Base routes
	root := goji.NewMux()
	root.Handle(pat.New("/api/*"), api)
	root.Handle(pat.Get("/static/*"), fileHandler)
	root.Handle(pat.Get("/favicon.ico"), fileHandler)
	root.HandleFunc(pat.Get("/*"), IndexHandler(entry))

	// And request logging
	root.Use(requestLogger)

	zap.L().Info("Starting the server listening", zap.String("port", *port))
	http.ListenAndServe(*port, root)
}

// IndexHandler serves up our index.html
func IndexHandler(entrypoint *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *entrypoint)
	}
}
