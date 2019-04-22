package cli

import "github.com/spf13/pflag"

func InitHosts(flag *pflag.FlagSet) {
	InitAuthFlags(flag)
	InitDatabaseFlags(flag)
	InitDPSFlags(flag)
	InitHostFlags(flag)
	InitIWSFlags(flag)
}
func (suite *cliTestSuite) TestConfigHosts() {
	suite.Setup(InitHosts)
	suite.Nil(CheckHosts(suite.viper))
}
