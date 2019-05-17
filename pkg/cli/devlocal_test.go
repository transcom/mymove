package cli

func (suite *cliTestSuite) TestConfigDevlocal() {
	suite.Setup(InitDevlocalFlags, []string{})
	suite.Nil(CheckDevlocal(suite.viper))
}
