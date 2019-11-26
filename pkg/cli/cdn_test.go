package cli

func (suite *cliTestSuite) TestConfigCDN() {
	suite.Setup(InitCDNFlags, []string{})
	suite.NoError(CheckCDNValues(suite.viper))
}
