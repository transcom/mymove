package cli

import (
	"os"
)

func (suite *cliTestSuite) TestConfigVault() {
	err := os.Setenv("AWS_PROFILE", "mock-env")
	suite.NoError(err)
	suite.Setup(InitVaultFlags, []string{})
	suite.NoError(CheckVault(suite.viper))
}
