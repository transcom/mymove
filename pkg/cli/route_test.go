package cli

func (suite *cliTestSuite) TestConfigRoute() {
	suite.Setup(InitRouteFlags, []string{})
	suite.NoError(CheckRoute(suite.viper))
}
