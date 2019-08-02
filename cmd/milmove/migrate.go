package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/ecs"
	"github.com/transcom/mymove/pkg/logging"
	"github.com/transcom/mymove/pkg/migrate"
)

// initMigrateFlags - Order matters!
func initMigrateFlags(flag *pflag.FlagSet) {

	// DB Config
	cli.InitDatabaseFlags(flag)

	// DB Retry Config
	cli.InitDatabaseRetryFlags(flag)

	// Migration Config
	cli.InitMigrationFlags(flag)

	// aws-vault Config
	cli.InitVaultFlags(flag)

	// Logging
	cli.InitLoggingFlags(flag)

	// Verbose
	cli.InitVerboseFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func checkMigrateConfig(v *viper.Viper, logger logger) error {

	logger.Info("checking migration config")

	if err := cli.CheckDatabase(v, logger); err != nil {
		return err
	}

	if err := cli.CheckDatabaseRetry(v); err != nil {
		return err
	}

	if err := cli.CheckMigration(v); err != nil {
		return err
	}

	if err := cli.CheckVault(v); err != nil {
		return err
	}

	if err := cli.CheckLogging(v); err != nil {
		return err
	}

	if err := cli.CheckVerbose(v); err != nil {
		return err
	}

	return nil
}

func migrateFunction(cmd *cobra.Command, args []string) error {

	err := cmd.ParseFlags(args)
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

	loggingEnv := v.GetString(cli.LoggingEnvFlag)

	logger, err := logging.Config(loggingEnv, v.GetBool(cli.VerboseFlag))
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

	if v.GetBool(cli.LogTaskMetadataFlag) {
		resp, httpGetErr := http.Get("http://169.254.170.2/v2/metadata")
		if httpGetErr != nil {
			logger.Error(errors.Wrap(httpGetErr, "could not fetch task metadata").Error())
		} else {
			body, readAllErr := ioutil.ReadAll(resp.Body)
			if readAllErr != nil {
				logger.Error(errors.Wrap(readAllErr, "could not read task metadata").Error())
			} else {
				taskMetadata := &ecs.TaskMetadata{}
				unmarshallErr := json.Unmarshal(body, taskMetadata)
				if unmarshallErr != nil {
					logger.Error(errors.Wrap(unmarshallErr, "could not parse task metadata").Error())
				} else {
					logger = logger.With(
						zap.String("ecs_cluster", taskMetadata.Cluster),
						zap.String("ecs_task_def_family", taskMetadata.Family),
						zap.String("ecs_task_def_revision", taskMetadata.Revision),
					)
				}
			}
			err = resp.Body.Close()
			if err != nil {
				logger.Error(errors.Wrap(err, "could not close task metadata response").Error())
			}
		}
	}

	zap.ReplaceGlobals(logger)

	logger.Info("migrator starting up")

	err = checkMigrateConfig(v, logger)
	if err != nil {
		logger.Fatal("invalid configuration", zap.Error(err))
	}

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
		dbConnection, errDbConn = cli.InitDatabase(v, logger)
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

	migrationPath := v.GetString(cli.MigrationPathFlag)
	logger.Info(fmt.Sprintf("using migration path %q", migrationPath))

	migrationManifest := v.GetString(cli.MigrationManifestFlag)
	logger.Info(fmt.Sprintf("using migration manifest at %q", migrationManifest))

	manifest, err := os.Open(migrationManifest)
	if err != nil {
		return errors.Wrap(err, "error reading manifest")
	}
	migrations := map[string]struct{}{}
	scanner := bufio.NewScanner(manifest)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			// skip blank lines
			continue
		}
		if strings.HasPrefix(line, "#") {
			// If line starts with a #, then comment it out.
			continue
		}
		migrations[line] = struct{}{}
	}

	fm := &pop.FileMigrator{
		Migrator: pop.NewMigrator(dbConnection),
		Path:     migrationPath,
	}
	fm.SchemaPath = migrationPath

	runner := func(mf pop.Migration, tx *pop.Connection) error {

		f, err := os.Open(mf.Path)
		if err != nil {
			return err
		}

		// if a secure migration, this step will execute.
		// See https://github.com/gobuffalo/fizz/pull/54
		content, err := pop.MigrationContent(mf, tx, f, false)
		if err != nil {
			return errors.Wrapf(err, "error processing %s", mf.Path)
		}

		if content == "" {
			return nil
		}

		err = tx.RawQuery(content).Exec()
		if err != nil {
			return errors.Wrapf(err, "error executing %s, sql: %s", mf.Path, content)
		}

		return nil

	}

	errFindMigrations := migrate.FindMigrations(fm, migrations, runner, logger)
	if errFindMigrations != nil {
		return errFindMigrations
	}

	return fm.Up()
}
