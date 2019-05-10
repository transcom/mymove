package cli

func (suite *cliTestSuite) TestConfigAWS() {
	suite.Setup(InitAWSFlags)
	suite.Nil(CheckAWSRegion(suite.viper))
}
