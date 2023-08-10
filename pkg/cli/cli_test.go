package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type cliTestSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger *zap.Logger
}

type initFlags func(f *pflag.FlagSet)

func (suite *cliTestSuite) Setup(fn initFlags, flagSet []string) {
	suite.viper = nil

	flag := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	fn(flag)
	suite.NoError(flag.Parse(flagSet))

	v := viper.New()
	err := v.BindPFlags(flag)
	if err != nil {
		suite.logger.Fatal("could not bind flags", zap.Error(err))
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	suite.SetViper(v)
}

func (suite *cliTestSuite) SetViper(v *viper.Viper) {
	suite.viper = v
}

func TestCLISuite(t *testing.T) {

	ss := &cliTestSuite{
		logger: zaptest.NewLogger(t),
	}

	suite.Run(t, ss)
}
