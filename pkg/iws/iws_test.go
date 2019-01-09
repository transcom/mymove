package iws

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type iwsSuite struct {
	testingsuite.BaseTestSuite
	logger *zap.Logger
}

func (suite *iwsSuite) SetupSuite() {
	var err error
	suite.logger, err = zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
}

func TestIwsSuite(t *testing.T) {
	suite.Run(t, new(iwsSuite))
}
