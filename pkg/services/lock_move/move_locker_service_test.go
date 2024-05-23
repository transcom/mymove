package lockmove

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveLockerServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestMoveLockerServiceSuite(t *testing.T) {

	hs := &MoveLockerServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
