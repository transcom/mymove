package pwsviolation

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PWSViolationsSuite struct {
	*testingsuite.PopTestSuite
}

func TestPWSViolationsServiceSuite(t *testing.T) {
	ts := &PWSViolationsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
