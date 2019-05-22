package cli

func (suite *cliTestSuite) TestConfigDPS() {
	suite.Setup(InitDPSFlags, []string{})
	suite.Nil(CheckDPS(suite.viper))
}
