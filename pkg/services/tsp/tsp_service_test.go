package tsp

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type TSPServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *TSPServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestUserSuite(t *testing.T) {

	hs := &TSPServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}
