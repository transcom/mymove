package cli

import (
	"fmt"
)

func (suite *cliTestSuite) TestConfigMigrationGenPath() {
	flagSet := []string{
		fmt.Sprintf("--%s=%s", MigrationGenPathFlag, "../../migrations/app/schema"),
	}

	suite.Setup(InitMigrationGenPathFlags, flagSet)
	suite.NoError(CheckMigrationGenPath(suite.viper))
}
