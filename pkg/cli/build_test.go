package cli

func (suite *cliTestSuite) TestConfigBuild() {
	suite.Setup(InitBuildFlags, []string{})
	suite.Nil(CheckBuild(suite.viper))
}
