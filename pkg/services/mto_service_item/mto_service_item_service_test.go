package mtoserviceitem

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOServiceItemServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *MTOServiceItemServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestMTOServiceItemServiceSuite(t *testing.T) {
	ts := &MTOServiceItemServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
