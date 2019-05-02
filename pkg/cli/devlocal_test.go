package cli

func (suite *cliTestSuite) TestConfigDevlocal() {
	suite.Setup(InitDevlocalFlags)
	suite.Nil(CheckDevlocal(suite.viper))
}
