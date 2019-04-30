package cli

func (suite *cliTestSuite) TestConfigCSRF() {
	suite.Setup(InitCSRFFlags)
	suite.Nil(CheckCSRF(suite.viper))
}
