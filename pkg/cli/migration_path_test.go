package cli

func (suite *cliTestSuite) TestConfigMigrationPath() {
	suite.Setup(InitMigrationPathFlags, []string{})
	suite.NoError(CheckMigrationPath(suite.viper))
}
