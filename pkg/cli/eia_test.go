package cli

func (suite *cliTestSuite) TestConfigEIA() {
	suite.Setup(InitEIAFlags, []string{})
	suite.NoError(CheckEIA(suite.viper))
}
