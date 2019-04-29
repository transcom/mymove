package cli

func (suite *cliTestSuite) TestConfigEmail() {
	suite.Setup(InitEmailFlags)
	suite.Nil(CheckEmail(suite.viper))
}
