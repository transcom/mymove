package cli

func (suite *cliTestSuite) TestConfigGEX() {
	suite.Setup(InitGEXFlags)
	suite.Nil(CheckGEX(suite.viper))
}
