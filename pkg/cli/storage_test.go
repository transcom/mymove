package cli

func (suite *cliTestSuite) TestConfigStorage() {
	suite.Setup(InitStorageFlags, []string{})
	suite.NoError(CheckStorage(suite.viper))
}
