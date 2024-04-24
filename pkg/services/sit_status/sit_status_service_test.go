package sitstatus

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SITStatusServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestSITStatusServiceSuite(t *testing.T) {

	ts := &SITStatusServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
