package cli

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type cliTestSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger logger
}
