package cli

func (suite *cliTestSuite) TestConfigDebug() {
	suite.Setup(InitDebugFlags, []string{})
	suite.NoError(CheckDebugFlags(suite.viper))
}
