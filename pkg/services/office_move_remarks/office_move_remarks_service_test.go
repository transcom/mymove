package officemoveremarks

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type OfficeMoveRemarksSuite struct {
	testingsuite.PopTestSuite
}

func TestOfficeMoveRemarksServiceSuite(t *testing.T) {
	ts := OfficeMoveRemarksSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, &ts)
	ts.PopTestSuite.TearDown()
}
