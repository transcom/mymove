//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used set up environment variables
//RA: Given the functions causing the lint errors are used to set environment variables for testing purposes, it does not present a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

type webServerSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger logger
}

func TestWebServerSuite(t *testing.T) {

	flag := pflag.CommandLine
	initServeFlags(flag)
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values in the file are used set up environment variables
	//RA: Given the functions causing the lint errors are used to set environment variables for testing purposes, it does not present a risk
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Modified Severity: N/A
	flag.Parse([]string{}) // nolint:errcheck

	v := viper.New()
	bindErr := v.BindPFlags(flag)
	if bindErr != nil {
		log.Fatal("failed to bind flags", zap.Error(bindErr))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config(logging.WithEnvironment(v.GetString(cli.DbEnvFlag)), logging.WithLoggingLevel(v.GetString(cli.LoggingLevelFlag)))
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

	ss := &webServerSuite{
		viper:  v,
		logger: logger,
	}

	suite.Run(t, ss)
}

// TestCheckServeConfigApp is the acceptance test for the milmove webserver
// This will run all checks against the local environment and fail if something isn't configured
func (suite *webServerSuite) TestCheckServeConfigApp() {
	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filenameApp := fmt.Sprintf("%s/config/env/%s.app.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		suite.logger.Info(fmt.Sprintf("Loading environment variables from file %s", filenameApp))
		suite.applyContext(suite.patchContext(suite.loadContext(filenameApp)))
	}

	suite.Nil(checkServeConfig(suite.viper, suite.logger))
}

// TestCheckServeConfigAppClientTLS is the acceptance test for the milmove webserver
// This will run all checks against the local environment and fail if something isn't configured
func (suite *webServerSuite) TestCheckServeConfigAppClientTLS() {
	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filenameApp := fmt.Sprintf("%s/config/env/%s.app-client-tls.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		suite.logger.Info(fmt.Sprintf("Loading environment variables from file %s", filenameApp))
		suite.applyContext(suite.patchContext(suite.loadContext(filenameApp)))
	}

	suite.Nil(checkServeConfig(suite.viper, suite.logger))
}

// TestCheckServeConfigMigrate is the acceptance test for the milmove migration command
// This will run all checks against the local environment and fail if something isn't configured
func (suite *webServerSuite) TestCheckServeConfigMigrate() {
	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filenameApp := fmt.Sprintf("%s/config/env/%s.migrations.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		suite.logger.Info(fmt.Sprintf("Loading environment variables from file %s", filenameApp))
		suite.applyContext(suite.patchContext(suite.loadContext(filenameApp)))
	}

	suite.Nil(checkMigrateConfig(suite.viper, suite.logger))
}

func (suite *webServerSuite) loadContext(variablesFile string) map[string]string {
	ctx := map[string]string{}
	if len(variablesFile) > 0 {
		if _, variablesFileStatErr := os.Stat(variablesFile); os.IsNotExist(variablesFileStatErr) {
			suite.logger.Fatal(fmt.Sprintf("File %q does not exist", variablesFile))
		}
		// Read contents of variables file into vars
		vars, err := ioutil.ReadFile(filepath.Clean(variablesFile))
		if err != nil {
			suite.logger.Fatal(fmt.Sprintf("error reading variables from file %s", variablesFile))
		}

		// Adds variables from file into context
		for _, x := range strings.Split(string(vars), "\n") {
			// If a line is empty or starts with #, then skip.
			if len(x) > 0 && x[0] != '#' {
				// Split each line on the first equals sign into []string{name, value}
				pair := strings.SplitAfterN(x, "=", 2)
				ctx[pair[0][0:len(pair[0])-1]] = pair[1]
			}
		}
	}
	return ctx
}

// patchContext updates specific variables based on value
func (suite *webServerSuite) patchContext(ctx map[string]string) map[string]string {
	for k, v := range ctx {
		if strings.HasPrefix(v, "/bin/") {
			newValue := filepath.Join(os.Getenv("TEST_ACC_CWD"), v[1:])
			ctx[k] = newValue
		}
	}

	// Always set the root cert to something on the local system.
	newValue := filepath.Join(os.Getenv("TEST_ACC_CWD"), "bin/rds-ca-us-gov-west-1-2017-root.pem")
	ctx["DB_SSL_ROOT_CERT"] = newValue

	// Always set the migration path to something on the local system.
	appSecure := filepath.Join(os.Getenv("TEST_ACC_CWD"), "migrations/app/secure")
	appSchema := filepath.Join(os.Getenv("TEST_ACC_CWD"), "migrations/app/schema")
	ctx["MIGRATION_PATH"] = fmt.Sprintf("file://%v;file://%v", appSchema, appSecure)

	return ctx
}

func (suite *webServerSuite) applyContext(ctx map[string]string) {
	for k, v := range ctx {
		suite.logger.Info("overriding " + k)
		suite.viper.Set(strings.Replace(strings.ToLower(k), "_", "-", -1), v)
	}
}
