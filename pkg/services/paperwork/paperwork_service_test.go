package paperwork

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaperworkServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestPaperworkServiceSuite(t *testing.T) {

	ts := &PaperworkServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
