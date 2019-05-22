package cli

func (suite *cliTestSuite) TestConfigEmail() {
	suite.Setup(InitEmailFlags, []string{})
	suite.Nil(CheckEmail(suite.viper))
}
