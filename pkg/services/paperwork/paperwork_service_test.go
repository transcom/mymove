package paperwork

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaperworkServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// TestAppContext returns the AppContext for the test suite
func (suite *PaperworkServiceSuite) TestAppContext() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger)
}

func TestPaperworkServiceSuite(t *testing.T) {

	ts := &PaperworkServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
