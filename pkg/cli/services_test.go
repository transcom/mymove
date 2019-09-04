package cli

func (suite *cliTestSuite) TestConfigServices() {
	suite.Setup(InitServiceFlags, []string{})
	suite.NoError(CheckServices(suite.viper))
}
