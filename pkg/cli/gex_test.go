package cli

func (suite *cliTestSuite) TestConfigGEX() {
	suite.Setup(InitGEXFlags, []string{})
	suite.NoError(CheckGEX(suite.viper))
}
