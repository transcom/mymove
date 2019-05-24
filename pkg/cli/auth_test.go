package cli

func (suite *cliTestSuite) TestConfigAuth() {
	suite.Setup(InitAuthFlags, []string{})
	suite.Nil(CheckAuth(suite.viper))
}
