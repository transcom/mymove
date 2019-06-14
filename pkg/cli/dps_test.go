package cli

func (suite *cliTestSuite) TestConfigDPS() {
	suite.Setup(InitDPSFlags, []string{})
	suite.NoError(CheckDPS(suite.viper))
}
