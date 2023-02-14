package cli

func (suite *cliTestSuite) TestConfigEIA() {
	// create fake data for the eia test
	// use https://pkg.go.dev/testing#B.Setenv for cleanup
	// dummy generated with openssl rand -hex 16
	suite.T().Setenv("EIA_KEY",
		"2bc8b0f2d13dd4f02f654650b6782ef0")
	suite.Setup(InitEIAFlags, []string{})
	suite.NoError(CheckEIA(suite.viper))
}
