package cli

func (suite *cliTestSuite) TestConfigSession() {
	suite.Setup(InitSessionFlags, []string{})
	suite.NoError(CheckSession(suite.viper))
}
