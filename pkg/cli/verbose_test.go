package cli

func (suite *cliTestSuite) TestConfigVerbose() {
	suite.Setup(InitVerboseFlags)
	suite.Nil(CheckVerbose(suite.viper))
}
