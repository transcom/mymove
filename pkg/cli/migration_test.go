package cli

func (suite *cliTestSuite) TestConfigMigration() {
	suite.Setup(InitMigrationFlags)
	suite.Nil(CheckMigration(suite.viper))
}
