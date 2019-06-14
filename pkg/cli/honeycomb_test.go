package cli

import (
	"os"
)

func (suite *cliTestSuite) TestHoneycomb() {

	if os.Getenv("TEST_ACC_HONEYCOMB") != "1" {
		suite.logger.Info("skipping TestHoneycomb")
		return
	}

	suite.Setup(InitHoneycombFlags, []string{})
	suite.NoError(CheckHoneycomb(suite.viper))
	enabled := InitHoneycomb(suite.viper, suite.logger)
	suite.True(enabled)
}
