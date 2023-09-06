package cli

func (suite *cliTestSuite) TestConfigAuth() {
	suite.T().Setenv("okta-customer-client-id", "someValue")
	suite.T().Setenv("okta-office-client-id", "someValue")
	suite.T().Setenv("okta-admin-client-id", "someValue")
	suite.Setup(InitAuthFlags, []string{})
	suite.NoError(CheckAuth(suite.viper))
}
