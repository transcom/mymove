package cli

func (suite *cliTestSuite) TestConfigVault() {
	suite.Setup(InitVaultFlags, []string{})
	suite.NoError(CheckVault(suite.viper))
}
