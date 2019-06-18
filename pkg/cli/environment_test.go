package cli

func (suite *cliTestSuite) TestConfigEnvironment() {
	suite.Setup(InitEnvironmentFlags, []string{})
	suite.NoError(CheckEnvironment(suite.viper))
}
