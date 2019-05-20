package cli

func (suite *cliTestSuite) TestConfigPorts() {
	suite.Setup(InitPortFlags, []string{})
	suite.Nil(CheckPorts(suite.viper))
}
