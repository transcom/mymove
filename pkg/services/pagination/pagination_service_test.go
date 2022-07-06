package pagination

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaginationServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestPaginationSuite(t *testing.T) {
	ts := &PaginationServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
