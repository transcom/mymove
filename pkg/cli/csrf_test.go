package cli

func (suite *cliTestSuite) TestConfigCSRF() {
	// dummy generated with openssl rand -hex 32
	suite.T().Setenv("CSRF_AUTH_KEY", "84ba3ca1ade09785ff4a8f3abf12fbd841a504262baff9fe6d86ab1bdede86d9")
	suite.Setup(InitCSRFFlags, []string{})
	suite.NoError(CheckCSRF(suite.viper))
}
