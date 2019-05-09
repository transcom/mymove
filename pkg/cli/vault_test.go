package cli

func (suite *cliTestSuite) TestConfigVault() {
	suite.Setup(InitVaultFlags)
	suite.Nil(CheckVault(suite.viper))
}
