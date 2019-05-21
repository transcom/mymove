package cli

func (suite *cliTestSuite) TestConfigRoute() {
	suite.Setup(InitRouteFlags, []string{})
	suite.Nil(CheckRoute(suite.viper))
}
