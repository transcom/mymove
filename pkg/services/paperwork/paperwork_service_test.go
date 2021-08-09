package paperwork

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaperworkServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func TestPaperworkServiceSuite(t *testing.T) {

	ts := &PaperworkServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
