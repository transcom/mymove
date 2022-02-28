package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gobuffalo/pop/v5"
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

	// Migration Path Config
	cli.InitMigrationPathFlags(flag)

	// Logging
	cli.InitLoggingFlags(flag)

	// Sort command line flags
	flag.SortFlags = true
}

func checkMigrateConfig(v *viper.Viper, logger *zap.Logger) error {

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

	if err := cli.CheckMigrationPath(v); err != nil {
		return err
	}

	if err := cli.CheckLogging(v); err != nil {
		return err
	}

	return nil
}

func expandPath(in string) string {
	if strings.HasPrefix(in, "s3://") {
		return in
	}
	if strings.HasPrefix(in, "file://") {
		return in
	}
	return "file://" + in
}

func expandPaths(in []string) []string {
	out := make([]string, 0, len(in))
	for _, x := range in {
		// Don't expand empty paths
		if len(x) == 0 {
			continue
		}
		out = append(out, expandPath(x))
	}
	return out
}

func migrateFunction(cmd *cobra.Command, args []string) error {

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
		return errors.Wrap(err, "invalid configuration")
	}

	if v.GetBool(cli.DbDebugFlag) {
		pop.Debug = true
	}

	// Remove any extra quotes around path
	trimmedMigrationPaths := strings.Trim(v.GetString(cli.MigrationPathFlag), "\"")
	migrationPaths := expandPaths(strings.Split(trimmedMigrationPaths, ";"))
	logger.Info(fmt.Sprintf("using migration paths %q", migrationPaths))

	s3Migrations := false
	for _, p := range migrationPaths {
		if strings.HasPrefix(p, "s3://") {
			s3Migrations = true
			break
		}
	}

	var session *awssession.Session
	if v.GetBool(cli.DbIamFlag) || s3Migrations {
		c := &aws.Config{
			Region: aws.String(v.GetString(cli.AWSRegionFlag)),
		}
		s, errorSession := awssession.NewSession(c)
		if errorSession != nil {
			return errors.Wrap(errorSession, "error creating aws session")
		}
		session = s
	}

	var dbCreds *credentials.Credentials
	if v.GetBool(cli.DbIamFlag) {
		if session != nil {
			// We want to get the credentials from the logged in AWS session rather than create directly,
			// because the session conflates the environment, shared, and container metdata config
			// within NewSession.  With stscreds, we use the Secure Token Service,
			// to assume the given role (that has rds db connect permissions).
			dbIamRole := v.GetString(cli.DbIamRoleFlag)
			logger.Info(fmt.Sprintf("assuming AWS role %q for db connection", dbIamRole))
			dbCreds = stscreds.NewCredentials(session, dbIamRole)
			stsService := sts.New(session)
			callerIdentity, callerIdentityErr := stsService.GetCallerIdentity(&sts.GetCallerIdentityInput{})
			if callerIdentityErr != nil {
				logger.Error(errors.Wrap(callerIdentityErr, "error getting aws sts caller identity").Error())
			} else {
				logger.Info(fmt.Sprintf("STS Caller Identity - Account: %s, ARN: %s, UserId: %s", *callerIdentity.Account, *callerIdentity.Arn, *callerIdentity.UserId))
			}
		}
	}

	// Create a connection to the DB with retry logic
	var dbConnection *pop.Connection
	var errDbConn error
	retryCount := 0
	retryMax := v.GetInt(cli.DbRetryMaxFlag)
	retryInterval := v.GetDuration(cli.DbRetryIntervalFlag)

	for retryCount < retryMax {
		dbConnection, errDbConn = cli.InitDatabase(v, dbCreds, logger)
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

	migrationTableName := dbConnection.MigrationTableName()
	logger.Info(fmt.Sprintf("tracking migrations using table %q", migrationTableName))

	migrationManifest := expandPath(v.GetString(cli.MigrationManifestFlag))
	logger.Info(fmt.Sprintf("using migration manifest %q", migrationManifest))

	var s3Client *s3.S3
	if s3Migrations {
		s3Client = s3.New(session)
	}

	migrationFiles := map[string][]string{}

	fileHelper := migrate.NewFileHelper()
	for _, p := range migrationPaths {
		// Don't list files in empty paths
		if len(p) == 0 {
			continue
		}
		filenames, errListFiles := fileHelper.ListFiles(p, s3Client)
		if errListFiles != nil {
			logger.Fatal(fmt.Sprintf("Error listing migrations directory %s", p), zap.Error(errListFiles))
		}
		migrationFiles[p] = filenames
	}

	manifest, err := os.Open(migrationManifest[len("file://"):])
	if err != nil {
		return errors.Wrap(err, "error reading manifest")
	}

	wait := v.GetDuration(cli.MigrationWaitFlag)

	migrator := pop.NewMigrator(dbConnection)
	scanner := bufio.NewScanner(manifest)
	for scanner.Scan() {
		target := scanner.Text()
		if strings.HasPrefix(target, "#") {
			// If line starts with a #, then comment it out.
			continue
		}
		uri := ""
		for dir, filenames := range migrationFiles {
			for _, filename := range filenames {
				if target == filename {
					uri = fmt.Sprintf("%s/%s", dir, filename)
					break
				}
			}
		}
		if len(uri) == 0 {
			return errors.Errorf("Error finding migration for filename %q", target)
		}
		m, err := pop.ParseMigrationFilename(target)
		if err != nil {
			return errors.Wrapf(err, "error parsing migration filename %q", uri)
		}
		if m == nil {
			return errors.Errorf("Error parsing migration filename %q", uri)
		}
		b := &migrate.Builder{Match: m, Path: uri}
		migration, errCompile := b.Compile(s3Client, wait, logger)
		if errCompile != nil {
			return errors.Wrap(errCompile, "Error compiling migration")
		}
		migrator.Migrations[migration.Direction] = append(migrator.Migrations[migration.Direction], *migration)
	}

	errSchemaMigrations := migrator.CreateSchemaMigrations()
	if errSchemaMigrations != nil {
		return errors.Wrap(errSchemaMigrations, "error creating table for tracking migrations")
	}

	errUp := migrator.Up()
	if errUp != nil {
		return errors.Wrap(errUp, "error running migrations")
	}

	return nil
}
