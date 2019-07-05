package cli

func (suite *cliTestSuite) TestCheckFeatureFlag() {
	suite.Setup(InitDPSFlags, []string{})
	suite.NoError(CheckFeatureFlag(suite.viper))
}