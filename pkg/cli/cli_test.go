package cli

import (
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
	logger logger
}

type initFlags func(f *pflag.FlagSet)

// A function to use when there is nothing that needs initializing in our tests
func initNull(flag *pflag.FlagSet) {}

func (suite *cliTestSuite) Setup(fn initFlags) {
	flag := pflag.CommandLine
	fn(flag)
	flag.Parse([]string{})

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

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
