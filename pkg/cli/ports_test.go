package cli

func (suite *cliTestSuite) TestConfigPorts() {
	suite.Setup(InitPortFlags)
	suite.Nil(CheckPorts(suite.viper))
}
