package cli

func (suite *cliTestSuite) TestConfigEnvironment() {
	suite.Setup(InitEnvironmentFlags)
	suite.Nil(CheckEnvironment(suite.viper))
}
