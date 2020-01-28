package cli

func (suite *cliTestSuite) TestCAC() {
	suite.Setup(InitCACFlags, []string{})
	suite.NoError(CheckCAC(suite.viper))
}
