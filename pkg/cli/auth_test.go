package cli

func (suite *cliTestSuite) TestConfigAuth() {
	suite.Setup(InitAuthFlags)
	suite.Nil(CheckAuth(suite.viper))
}
