package movetaskorder_test

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveTaskOrderServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func TestMoveTaskOrderServiceSuite(t *testing.T) {
	ts := &MoveTaskOrderServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
