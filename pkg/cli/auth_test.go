package cli

func (suite *cliTestSuite) TestConfigAuth() {
	suite.Setup(InitAuthFlags, []string{})
	suite.NoError(CheckAuth(suite.viper))
}
