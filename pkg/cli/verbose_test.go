package cli

func (suite *cliTestSuite) TestConfigVerbose() {
	suite.Setup(InitVerboseFlags, []string{})
	suite.NoError(CheckVerbose(suite.viper))
}
