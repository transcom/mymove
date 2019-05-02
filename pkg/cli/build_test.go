package cli

func (suite *cliTestSuite) TestConfigBuild() {
	suite.Setup(InitBuildFlags)
	suite.Nil(CheckBuild(suite.viper))
}
