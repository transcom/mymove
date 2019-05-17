package cli

func (suite *cliTestSuite) TestConfigEIA() {
	suite.Setup(InitEIAFlags, []string{})
	suite.Nil(CheckEIA(suite.viper))
}
