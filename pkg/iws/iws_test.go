package iws

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
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
