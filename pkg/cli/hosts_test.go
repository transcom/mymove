package cli

func (suite *cliTestSuite) TestConfigHosts() {
	suite.Setup(InitHostFlags, []string{})
	suite.Nil(CheckHosts(suite.viper))
}
