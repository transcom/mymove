package cli

func (suite *cliTestSuite) TestConfigReceiver() {
	suite.Setup(InitReceiverFlags, []string{})
	suite.NoError(CheckReceiver(suite.viper))
}
