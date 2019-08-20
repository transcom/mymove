package cli

func (suite *cliTestSuite) TestConfigWebserver() {
	suite.Setup(InitWebserverFlags, []string{})
	suite.NoError(CheckWebserver(suite.viper))
}
