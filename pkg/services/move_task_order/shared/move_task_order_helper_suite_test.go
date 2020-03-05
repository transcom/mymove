package movetaskordershared_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveTaskOrderHelperSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MoveTaskOrderHelperSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestMoveTaskOrderServiceSuite(t *testing.T) {
	ts := &MoveTaskOrderHelperSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
