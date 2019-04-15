package cli

func (suite *cliServerSuite) TestConfigEIA() {
	suite.Nil(CheckEIA(suite.viper))
}
