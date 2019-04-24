package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/logging"
)

type cliTestSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger Logger
}

type initFlags func(f *pflag.FlagSet)

// A function to use when there is nothing that needs initializing in our tests
func initNull(flag *pflag.FlagSet) {}

func (suite *cliTestSuite) Setup(fn initFlags) {
	suite.viper = nil

	flag := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	fn(flag)
	flag.Parse([]string{})

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	suite.SetViper(v)
}

func (suite *cliTestSuite) SetViper(v *viper.Viper) {
	suite.viper = v
}

func TestCLISuite(t *testing.T) {

	logger, _ := logging.Config("development", true)
	zap.ReplaceGlobals(logger)

	ss := &cliTestSuite{
		logger: logger,
	}

	suite.Run(t, ss)
}
