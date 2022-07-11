package movehistory

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveHistoryServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestMoveServiceSuite(t *testing.T) {

	hs := &MoveHistoryServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
