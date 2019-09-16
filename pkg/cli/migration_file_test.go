package cli

import (
	"fmt"
)

func (suite *cliTestSuite) TestConfigMigrationFile() {
	flagSet := []string{
		fmt.Sprintf("--%s=%s", MigrationNameFlag, "test_migration"),
	}

	suite.Setup(InitMigrationFileFlags, flagSet)
	suite.NoError(CheckMigrationFile(suite.viper))
}
