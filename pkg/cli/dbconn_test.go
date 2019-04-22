package cli

import (
	"os"
)

func (suite *cliTestSuite) TestConfigDatabase() {
	suite.Setup(InitDatabaseFlags)
	suite.Nil(CheckDatabase(suite.viper, suite.logger))
}

func (suite *cliTestSuite) TestInitDatabase() {

	if os.Getenv("TEST_ACC_INIT_DATABASE") != "1" {
		suite.logger.Info("skipping TestInitDatabase")
		return
	}

	suite.Setup(InitDatabaseFlags)
	conn, err := InitDatabase(suite.viper, suite.logger)
	suite.Nil(err)
	suite.NotNil(conn)
}
