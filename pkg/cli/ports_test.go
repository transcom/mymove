package cli

func (suite *cliTestSuite) TestConfigPorts() {
	suite.Setup(InitPortFlags, []string{})
	suite.NoError(CheckPorts(suite.viper))
}
