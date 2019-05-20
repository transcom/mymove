package cli

func (suite *cliTestSuite) TestConfigVerbose() {
	suite.Setup(InitVerboseFlags, []string{})
	suite.Nil(CheckVerbose(suite.viper))
}
