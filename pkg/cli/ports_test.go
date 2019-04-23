package cli

import "github.com/spf13/pflag"

func InitPorts(flag *pflag.FlagSet) {
	InitAuthFlags(flag)
	InitDPSFlags(flag)
	InitDatabaseFlags(flag)
	InitPortFlags(flag)
}

func (suite *cliTestSuite) TestConfigPorts() {
	suite.Setup(InitPorts)
	suite.Nil(CheckPorts(suite.viper))
}
