package cli

func (suite *cliTestSuite) TestConfigSwagger() {
	suite.Setup(InitSwaggerFlags, []string{})
	suite.NoError(CheckSwagger(suite.viper))
}
