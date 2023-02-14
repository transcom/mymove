package cli

func (suite *cliTestSuite) TestConfigIWS() {
	// create fake data for the iws test
	// use https://pkg.go.dev/testing#B.Setenv for cleanup
	suite.T().Setenv("IWS_RBS_HOST", "iws.rbs.example.com")
	suite.Setup(InitIWSFlags, []string{})
	suite.NoError(CheckIWS(suite.viper))
}
