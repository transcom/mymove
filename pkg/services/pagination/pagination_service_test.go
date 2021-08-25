package pagination

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaginationServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger Logger
}

func TestPaginationSuite(t *testing.T) {
	ts := &PaginationServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing,
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
