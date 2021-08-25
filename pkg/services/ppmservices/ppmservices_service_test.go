package ppmservices

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMServiceSuite struct {
	testingsuite.PopTestSuite
	testingsuite.AppContextTestHelper
	logger *zap.Logger
}

func (suite *PPMServiceSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestPPMServiceSuite(t *testing.T) {

	hs := &PPMServiceSuite{
		PopTestSuite:         testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		AppContextTestHelper: testingsuite.NewAppContextTestHelper(),
		logger:               zap.NewNop(), // Use a no-op logger during testing, // Use a no-op logger during testing
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
