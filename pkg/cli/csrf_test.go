package cli

func (suite *cliTestSuite) TestConfigCSRF() {
	suite.Setup(InitCSRFFlags, []string{})
	suite.NoError(CheckCSRF(suite.viper))
}
