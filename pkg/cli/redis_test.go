package cli

func (suite *cliTestSuite) TestConfigRedis() {
	suite.Setup(InitRedisFlags, []string{})
	suite.NoError(CheckRedis(suite.viper))
}
