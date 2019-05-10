package cli

func (suite *cliTestSuite) TestConfigVault() {
	suite.Setup(InitVaultFlags, []string{})
	suite.Nil(CheckVault(suite.viper))
}
