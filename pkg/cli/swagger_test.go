package cli

func (suite *cliTestSuite) TestConfigSwagger() {
	suite.Setup(InitSwaggerFlags)
	suite.Nil(CheckSwagger(suite.viper))
}
