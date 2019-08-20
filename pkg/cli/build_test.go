package cli

func (suite *cliTestSuite) TestConfigBuild() {
	suite.Setup(InitBuildFlags, []string{})
	suite.NoError(CheckBuild(suite.viper))
}
