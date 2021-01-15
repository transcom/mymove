package certs

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging"
)

type certTestSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger Logger
}

type initFlags func(f *pflag.FlagSet)

func (suite *certTestSuite) Setup(fn initFlags, flagSet []string) {
	suite.viper = nil

	flag := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	fn(flag)
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return values are used to set/unset environment variables needed for session creation in the unit test's local database
	//RA: Setting/unsetting of environment variables does not present any risks and are solely used for unit testing purposes
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	flag.Parse(flagSet) // nolint:errcheck

	v := viper.New()
	err := v.BindPFlags(flag)
	if err != nil {
		suite.logger.Fatal("could not bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	suite.SetViper(v)
}

func (suite *certTestSuite) SetViper(v *viper.Viper) {
	suite.viper = v
}

func TestCertSuite(t *testing.T) {

	logger, err := logging.Config(logging.WithEnvironment("development"), logging.WithLoggingLevel("debug"))
	if err != nil {
		log.Fatalf("Failed to initialize Zap logging due to %v", err)
	}
	zap.ReplaceGlobals(logger)

	ss := &certTestSuite{
		logger: logger,
	}

	suite.Run(t, ss)
}

func (suite *certTestSuite) TestDODCertificates() {

	if os.Getenv("TEST_ACC_DOD_CERTIFICATES") != "1" {
		suite.logger.Info("skipping TestDODCertificates")
		return
	}

	suite.Setup(cli.InitCertFlags, []string{})
	_, _, err := InitDoDCertificates(suite.viper, suite.logger)
	suite.NoError(err)
}
