package auth

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type authSuite struct {
	testingsuite.BaseTestSuite
	logger *zap.Logger
}

func TestAuthSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &authSuite{logger: logger}
	suite.Run(t, hs)
}
