package cli

func (suite *cliTestSuite) TestConfigMigration() {
	suite.Setup(InitMigrationFlags, []string{})
	suite.Nil(CheckMigration(suite.viper))
}
