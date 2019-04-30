package cli

func (suite *cliTestSuite) TestConfigIWS() {
	suite.Setup(InitIWSFlags)
	suite.Nil(CheckIWS(suite.viper))
}
