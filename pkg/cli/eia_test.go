package cli

func (suite *cliTestSuite) TestConfigEIA() {
	suite.Nil(CheckEIA(suite.viper))
}
