package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

const (
	hhgFlag string = "hhg"
	ppmFlag string = "ppm"
)

func initFlags(flag *pflag.FlagSet) {
	flag.Int(hhgFlag, 0, "The number of access codes for HHG moves to create")
	flag.Int(ppmFlag, 0, "The number of access codes for PPM moves to create")

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Logging Levels
	cli.InitLoggingFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if verrs.HasAny() {
		log.Fatalf("validation Errors %v", verrs)
	}
	if err != nil {
		log.Fatalf("Failed to save %v", err)
	}
}

func checkConfig(v *viper.Viper, logger logger) error {

	logger.Debug("checking config")

	err := cli.CheckDatabase(v, logger)
	if err != nil {
		return err
	}
	err = cli.CheckLogging(v)
	if err != nil {
		return err
	}

	return nil
}

func main() {

	flag := pflag.CommandLine
	initFlags(flag)
	err := flag.Parse(os.Args[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)
	logger, err := logging.Config(
		logging.WithEnvironment(dbEnv),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	err = checkConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, nil, logger)
	if err != nil {
		// No connection object means that the configuraton failed to validate and we should not startup
		// A valid connection object that still has an error indicates that the DB is not up and we should not startup
		logger.Fatal("Connecting to DB", zap.Error(err))
	}

	hhg := v.GetInt(hhgFlag)
	ppm := v.GetInt(ppmFlag)

	if hhg == 0 && ppm == 0 {
		log.Fatal("Usage: generate_access_codes -ppm 1000 -hhg 4000")
	}
	// go run cmd/generate_access_codes/main.go -ppm 2000 -hhg 500
	for i := 0; i < hhg; i++ {
		accessCode := models.AccessCode{
			Code:     models.GenerateLocator(),
			MoveType: models.SelectedMoveTypeHHG,
		}

		mustSave(dbConnection, &accessCode)
	}

	for i := 0; i < ppm; i++ {
		accessCode := models.AccessCode{
			Code:     models.GenerateLocator(),
			MoveType: models.SelectedMoveTypePPM,
		}

		mustSave(dbConnection, &accessCode)
	}

	fmt.Println("Completed generate_access_codes")
}
