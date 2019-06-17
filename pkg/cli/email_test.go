package cli

func (suite *cliTestSuite) TestConfigEmail() {
	suite.Setup(InitEmailFlags, []string{})
	suite.NoError(CheckEmail(suite.viper))
}
