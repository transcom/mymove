package cli

import "github.com/spf13/pflag"

func InitProtocols(flag *pflag.FlagSet) {
	InitAuthFlags(flag)
	InitDPSFlags(flag)
}

func (suite *cliTestSuite) TestConfigProtocols() {
	suite.Setup(InitProtocols)
	suite.Nil(CheckProtocols(suite.viper))
}
