package cli

func (suite *cliTestSuite) TestConfigMiddleware() {
	suite.Setup(InitMiddlewareFlags)
	suite.Nil(CheckMiddleWare(suite.viper))
}
