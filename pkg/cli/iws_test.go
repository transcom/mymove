package cli

func (suite *cliTestSuite) TestConfigIWS() {
	suite.Setup(InitIWSFlags, []string{})
	suite.NoError(CheckIWS(suite.viper))
}
