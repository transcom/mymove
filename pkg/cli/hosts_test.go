package cli

func (suite *cliTestSuite) TestConfigHosts() {
	suite.Setup(InitHostFlags, []string{})
	suite.NoError(CheckHosts(suite.viper))
}
