package pagination

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaginationServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestPaginationSuite(t *testing.T) {
	ts := &PaginationServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
