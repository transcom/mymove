package cli

func (suite *cliTestSuite) TestConfigMigration() {
	suite.Setup(InitMigrationFlags, []string{})
	suite.NoError(CheckMigration(suite.viper))
}
