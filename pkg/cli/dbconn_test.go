package cli

import (
	"os"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type cliServerSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger logger
}

func (suite *cliServerSuite) TestConfigDatabase() {
	suite.Nil(CheckDatabase(suite.viper, suite.logger))
}

func (suite *cliServerSuite) TestInitDatabase() {

	if os.Getenv("TEST_ACC_INIT_DATABASE") != "1" {
		suite.logger.Info("skipping TestInitDatabase")
		return
	}

	conn, err := InitDatabase(suite.viper, suite.logger)
	suite.Nil(err)
	suite.NotNil(conn)
}
