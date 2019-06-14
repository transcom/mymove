package cli

func (suite *cliTestSuite) TestConfigMiddleware() {
	suite.Setup(InitMiddlewareFlags, []string{})
	suite.NoError(CheckMiddleWare(suite.viper))
}
