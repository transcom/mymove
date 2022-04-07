package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

func dbloadFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(args)
	if err != nil {
		return errors.Wrap(err, "could not parse flags")
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	loggingEnv := v.GetString(cli.LoggingEnvFlag)

	logger, _, errLogging := logging.Config(
		logging.WithEnvironment(loggingEnv),
		logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)),
		logging.WithStacktraceLength(v.GetInt(cli.StacktraceLengthFlag)),
	)
	if errLogging != nil {
		return errors.Wrapf(errLogging, "failed to initialize zap logging")
	}

	zap.ReplaceGlobals(logger)

	logger.Info("dbloader starting up")

	if v.GetBool(cli.DbDebugFlag) {
		pop.Debug = true
	}

	// Create a connection to the DB with retry logic
	var dbConnection *pop.Connection
	var errDbConn error
	retryCount := 0
	retryMax := v.GetInt(cli.DbRetryMaxFlag)
	retryInterval := v.GetDuration(cli.DbRetryIntervalFlag)

	for retryCount < retryMax {
		dbConnection, errDbConn = cli.InitDatabase(v, nil, logger)
		if errDbConn != nil {
			if dbConnection == nil {
				// No connection object means that the configuraton failed to validate and we should kill server startup
				logger.Fatal("Invalid DB Configuration", zap.Error(errDbConn))
			} else {
				// A valid connection object that still has an error indicates that the DB is not up and
				// thus is not ready for migrations. Attempt to retry connecting.
				logger.Error(fmt.Sprintf("DB is not ready for connections, sleeping for %q", retryInterval), zap.Error(errDbConn))
				time.Sleep(retryInterval)
			}
		} else {
			break
		}

		// Retry logic should break after max retries
		retryCount++
		if retryCount >= retryMax {
			logger.Fatal(fmt.Sprintf("DB was not ready for connections after %d retries", retryMax), zap.Error(errDbConn))
		}
	}

	schemaPath := v.GetString(cli.MigrationSchemaPathFlag)

	roles := filepath.Join(schemaPath, "roles.sql")
	f, err := os.Open(roles)
	if err != nil {
		return errors.Wrap(err, "Cannot open roles file: "+roles)
	}

	err = dbConnection.Dialect.LoadSchema(f)
	if err != nil {
		return errors.Wrap(err, "Cannot load db roles")
	}
	f.Close()
	logger.Info(fmt.Sprintf("Successfully loaded roles from %s", roles))

	schema := filepath.Join(schemaPath, "schema.sql")
	f, err = os.Open(schema)
	if err != nil {
		return errors.Wrap(err, "Cannot open schema file: "+schema)
	}

	err = dbConnection.Dialect.LoadSchema(f)
	if err != nil {
		return errors.Wrap(err, "Cannot load db schema")
	}
	logger.Info(fmt.Sprintf("Successfully loaded schema from %s", schema))
	f.Close()

	return nil
}
