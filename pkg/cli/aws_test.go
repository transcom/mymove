package cli

func (suite *cliTestSuite) TestConfigAWS() {
	suite.Setup(InitAWSFlags)
	suite.Nil(CheckAWS(suite.viper))
}
