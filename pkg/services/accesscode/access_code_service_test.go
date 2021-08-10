package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type AccessCodeServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *AccessCodeServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestAccessCodeServiceSuite(t *testing.T) {
	ts := &AccessCodeServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
