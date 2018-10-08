package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/namsral/flag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/internal/pkg/dutystationsloader"
)

func main() {
	config := flag.String("config-dir", "config", "The location of server config files")
	verbose := flag.Bool("verbose", false, "Sets debug logging level")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	validate := flag.Bool("validate", false, "Only run file validations")
	output := flag.String("output", "", "Where to output the migration file")
	stationsPath := flag.String("stations", "", "Input file for duty stations")
	officesPath := flag.String("offices", "", "Input file for transportation offices")
	flag.Parse()

	zapConfig := zap.NewDevelopmentConfig()
	logger, _ := zapConfig.Build()

	zapConfig.Level.SetLevel(zap.InfoLevel)
	if *verbose {
		zapConfig.Level.SetLevel(zap.DebugLevel)
	}

	//DB connection
	err := pop.AddLookupPaths(*config)
	if err != nil {
		logger.Panic("Error initializing db connection", zap.Error(err))
	}
	db, err := pop.Connect(*env)
	if err != nil {
		logger.Panic("Error initializing db connection", zap.Error(err))
	}

	// If we just want to validate files we can exit
	if validate != nil && *validate {
		os.Exit(0)
	}

	builder := dutystationsloader.NewMigrationBuilder(db, logger)
	insertions, err := builder.Build(*stationsPath, *officesPath)
	if err != nil {
		logger.Panic("Error while building migration", zap.Error(err))
	}

	var migration strings.Builder
	migration.WriteString("-- Migration generated using cmd/load_duty_stations\n")
	migration.WriteString(fmt.Sprintf("-- Duty stations file: %v\n", *stationsPath))
	migration.WriteString(fmt.Sprintf("-- Transportation offices file: %v\n", *officesPath))
	migration.WriteString("\n")
	migration.WriteString(insertions)

	f, err := os.OpenFile(*output, os.O_TRUNC|os.O_WRONLY, os.ModeAppend)
	defer f.Close()
	if err != nil {
		log.Panic(err)
	}
	_, err = f.WriteString(migration.String())
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Complete! Migration written to %v\n", *output)
}
