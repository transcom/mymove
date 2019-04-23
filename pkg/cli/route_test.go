package cli

func (suite *cliTestSuite) TestConfigRoute() {
	suite.Setup(InitRouteFlags)
	suite.Nil(CheckRoute(suite.viper))
}
