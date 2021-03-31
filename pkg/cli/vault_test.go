package cli

import (
	"os"
)

func (suite *cliTestSuite) TestConfigVault() {
	os.Setenv("AWS_PROFILE", "mock-env")
	suite.Setup(InitVaultFlags, []string{})
	suite.NoError(CheckVault(suite.viper))
}
