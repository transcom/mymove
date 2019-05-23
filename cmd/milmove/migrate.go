package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

// initMigrateFlags - Order matters!
func initMigrateFlags(flag *pflag.FlagSet) {
	// Migration Config
	cli.InitMigrationFlags(flag)

	// DB Config
	cli.InitDatabaseFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Don't sort flags
	flag.SortFlags = false
}

func checkMigrateConfig(v *viper.Viper, logger logger) error {

	logger.Info("checking migration config")

	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckDatabase(v, logger); err != nil {
		return err
	}

	if err := cli.CheckVerbose(v); err != nil {
		return err
	}

	return nil
}

func migrateFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(os.Args[1:])
	if err != nil {
		return errors.Wrap(err, "Could not parse flags")
	}

	flag := cmd.Flags()

	v := viper.New()
	err = v.BindPFlags(flag)
	if err != nil {
		return errors.Wrap(err, "Could not bind flags")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	dbEnv := v.GetString(cli.DbEnvFlag)

	logger, err := logging.Config(dbEnv, v.GetBool(cli.VerboseFlag))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}

	fields := make([]zap.Field, 0)
	if len(gitBranch) > 0 {
		fields = append(fields, zap.String("git_branch", gitBranch))
	}
	if len(gitCommit) > 0 {
		fields = append(fields, zap.String("git_commit", gitCommit))
	}
	logger = logger.With(fields...)
	zap.ReplaceGlobals(logger)

	logger.Info("migrator starting up")

	err = checkMigrateConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

	// Create a connection to the DB
	dbConnection, err := cli.InitDatabase(v, logger)
	if err != nil {
		if dbConnection == nil {
			// No connection object means that the configuraton failed to validate and we should kill server startup
			logger.Fatal("Invalid DB Configuration", zap.Error(err))
		} else {
			// A valid connection object that still has an error indicates that the DB is not up and
			// thus is not ready for migrations
			logger.Fatal("DB is not ready for connections", zap.Error(err))
		}
	}

	migrationPath := v.GetString(cli.MigrationPathFlag)
	logger.Info(fmt.Sprintf("using migration path %q", migrationPath))
	mig, err := pop.NewFileMigrator(migrationPath, dbConnection)
	if err != nil {
		return err
	}
	return mig.Up()
}
