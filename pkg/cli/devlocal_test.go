package cli

func (suite *cliTestSuite) TestConfigDevlocal() {
	var viper = suite.viper
	viper.Set("devlocal-auth", true)
	suite.Setup(InitDevlocalFlags, []string{})
	suite.NoError(CheckDevlocal(viper))
}
