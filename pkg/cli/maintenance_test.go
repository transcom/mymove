package cli

func (suite *cliTestSuite) TestMaintenance() {
	suite.Setup(InitMaintenanceFlags, []string{})
	suite.NoError(CheckMaintenance(suite.viper))
}
