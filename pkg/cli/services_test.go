package cli

func (suite *cliTestSuite) TestConfigServices() {
	// create fake data for the eia test
	// use https://pkg.go.dev/testing#B.Setenv for cleanup
	suite.T().Setenv("SERVE_API_INTERNAL", "true")
	suite.Setup(InitServiceFlags, []string{})
	suite.NoError(CheckServices(suite.viper))
}
