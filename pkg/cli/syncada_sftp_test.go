package cli

func (suite *cliTestSuite) TestConfigSyncadaSFTP() {
	suite.Setup(InitSyncadaSFTPFlags, []string{})
	suite.NoError(CheckSyncadaSFTP(suite.viper))
}
