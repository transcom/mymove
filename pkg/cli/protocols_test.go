package cli

func (suite *cliTestSuite) TestConfigProtocols() {
	suite.Setup(initNull)
	suite.Nil(CheckProtocols(suite.viper))
}
