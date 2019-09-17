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

func (suite *PaginationServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestPaginationSuite(t *testing.T) {
	ts := &PaginationServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(),
	}

	suite.Run(t, ts)
}
