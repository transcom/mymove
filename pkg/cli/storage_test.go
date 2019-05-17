package cli

func (suite *cliTestSuite) TestConfigStorage() {
	suite.Setup(InitStorageFlags, []string{})
	suite.Nil(CheckStorage(suite.viper))
}
