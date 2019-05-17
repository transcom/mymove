package cli

func (suite *cliTestSuite) TestConfigSwagger() {
	suite.Setup(InitSwaggerFlags, []string{})
	suite.Nil(CheckSwagger(suite.viper))
}
