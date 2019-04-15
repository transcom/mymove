package cli

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type cliServerSuite struct {
	suite.Suite
	viper  *viper.Viper
	logger logger
}
