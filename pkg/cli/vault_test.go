package cli

import (
	"log"
	"os"

	"go.uber.org/zap"
)

func (suite *cliTestSuite) TestConfigVault() {
	err := os.Setenv("AWS_PROFILE", "mock-env")
	if err != nil {
		log.Fatal("unable to set AWS_PROFILE variable", zap.Error(err))
	}
	suite.Setup(InitVaultFlags, []string{})
	suite.NoError(CheckVault(suite.viper))
}
