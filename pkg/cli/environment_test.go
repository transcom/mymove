package cli

func (suite *cliTestSuite) TestConfigEnvironment() {
	suite.Setup(InitEnvironmentFlags, []string{})
	suite.Nil(CheckEnvironment(suite.viper))
}
