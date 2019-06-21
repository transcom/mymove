package cli

func (suite *cliTestSuite) TestConfigLogging() {
	suite.Setup(InitLoggingFlags, []string{})
	suite.NoError(CheckLogging(suite.viper))
}
