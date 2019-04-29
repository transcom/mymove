package cli

func (suite *cliTestSuite) TestConfigStorage() {
	suite.Setup(InitStorageFlags)
	suite.Nil(CheckStorage(suite.viper))
}
