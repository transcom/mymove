package cli

func (suite *cliTestSuite) TestConfigCSRF() {
	suite.Setup(InitCSRFFlags, []string{})
	suite.Nil(CheckCSRF(suite.viper))
}
