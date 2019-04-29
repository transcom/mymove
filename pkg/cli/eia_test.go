package cli

func (suite *cliTestSuite) TestConfigEIA() {
	suite.Setup(InitEIAFlags)
	suite.Nil(CheckEIA(suite.viper))
}
