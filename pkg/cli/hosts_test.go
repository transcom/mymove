package cli

func (suite *cliTestSuite) TestConfigHosts() {
	suite.Setup(InitHostFlags)
	suite.Nil(CheckHosts(suite.viper))
}
