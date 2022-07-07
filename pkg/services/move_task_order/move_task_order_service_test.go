package movetaskorder_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveTaskOrderServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestMoveTaskOrderServiceSuite(t *testing.T) {
	ts := &MoveTaskOrderServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
