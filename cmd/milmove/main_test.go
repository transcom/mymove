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
	flag.Parse([]string{})

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	logger, err := logging.Config(v.GetString(cli.DbEnvFlag), v.GetBool(cli.VerboseFlag))
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

	if testEnv := os.Getenv("TEST_ACC_ENV"); len(testEnv) > 0 {
		filename := fmt.Sprintf("%s/config/env/%s.env", os.Getenv("TEST_ACC_CWD"), testEnv)
		logger.Info(fmt.Sprintf("Loading environment variables from file %s", filename))
		ss.applyContext(ss.patchContext(ss.loadContext(filename)))
	}

	suite.Run(t, ss)
}

// TestCheckConfig is the acceptance test for the milmove webserver
// This will run all checks against the local environment and fail if something isn't configured
func (suite *webServerSuite) TestCheckConfig() {
	suite.Nil(checkServeConfig(suite.viper, suite.logger))
}

func (suite *webServerSuite) loadContext(variablesFile string) map[string]string {
	ctx := map[string]string{}
	if len(variablesFile) > 0 {
		// Read contents of variables file into vars
		vars, err := ioutil.ReadFile(variablesFile)
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

func (suite *webServerSuite) patchContext(ctx map[string]string) map[string]string {
	for k, v := range ctx {
		if strings.HasPrefix(v, "/bin/") {
			ctx[k] = filepath.Join(os.Getenv("TEST_ACC_CWD"), v[1:])
		}
	}
	return ctx
}

func (suite *webServerSuite) applyContext(ctx map[string]string) {
	for k, v := range ctx {
		suite.logger.Info("overriding " + k)
		suite.viper.Set(strings.Replace(strings.ToLower(k), "_", "-", -1), v)
	}
}
