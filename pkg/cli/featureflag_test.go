package cli

func (suite *cliTestSuite) TestCheckFeatureFlag() {
	suite.Setup(InitFeatureFlags, []string{})
	suite.NoError(CheckFeatureFlag(suite.viper))
}
