package cli

func (suite *cliTestSuite) TestConfigDPS() {
	suite.Setup(InitDPSFlags)
	suite.Nil(CheckDPS(suite.viper))
}
