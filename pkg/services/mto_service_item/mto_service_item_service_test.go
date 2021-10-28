package mtoserviceitem

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOServiceItemServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// AppContextForTest returns the AppContext for the test suite
func (suite *MTOServiceItemServiceSuite) AppContextForTest() appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, nil)
}

func TestMTOServiceItemServiceSuite(t *testing.T) {
	ts := &MTOServiceItemServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
		zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
