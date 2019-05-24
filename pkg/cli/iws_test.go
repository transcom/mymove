package cli

func (suite *cliTestSuite) TestConfigIWS() {
	suite.Setup(InitIWSFlags, []string{})
	suite.Nil(CheckIWS(suite.viper))
}
