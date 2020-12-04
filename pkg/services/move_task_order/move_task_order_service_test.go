package movetaskorder

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveTaskOrderServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MoveTaskOrderServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestMoveTaskOrderServiceSuite(t *testing.T) {
	ts := &MoveTaskOrderServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
