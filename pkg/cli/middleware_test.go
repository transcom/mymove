package cli

func (suite *cliTestSuite) TestConfigMiddleware() {
	suite.Setup(InitMiddlewareFlags, []string{})
	suite.Nil(CheckMiddleWare(suite.viper))
}
