package main

import (
	"context"
	"log"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/awardqueue"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, configures the database, presenetly.")
	debugLogging := flag.Bool("debug_logging", false, "log messages at the debug level.")
	flag.Parse()

	// Set up logger for the system
	var logger *zap.Logger
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

	// DB connection
	err = pop.AddLookupPaths(*config)
	if err != nil {
		log.Panic(err)
	}
	dbConnection, err := pop.Connect(*env)
	if err != nil {
		log.Panic(err)
	}

	awardQueue := awardqueue.NewAwardQueue(dbConnection, logger)
	err = awardQueue.Run(context.Background())
	if err != nil {
		log.Panic(err)
	}
}
