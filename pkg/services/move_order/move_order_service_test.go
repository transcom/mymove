package moveorder

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveOrderServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *MoveOrderServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestMoveOrderServiceSuite(t *testing.T) {
	ts := &MoveOrderServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
