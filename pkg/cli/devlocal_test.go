package cli

func (suite *cliTestSuite) TestConfigDevlocal() {
	suite.Setup(InitDevlocalFlags, []string{})
	suite.NoError(CheckDevlocal(suite.viper))
}
