package cli

func (suite *cliTestSuite) TestConfigGEX() {
	suite.Setup(InitGEXFlags, []string{})
	suite.Nil(CheckGEX(suite.viper))
}
