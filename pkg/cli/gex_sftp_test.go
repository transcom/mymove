package cli

func (suite *cliTestSuite) TestConfigGEXSFTP() {
	suite.Setup(InitGEXSFTPFlags, []string{})
	suite.NoError(CheckGEXSFTP(suite.viper))
}
