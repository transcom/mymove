package fakedata

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type FakeDataSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *FakeDataSuite) SetupTest() {
	err := suite.TruncateAll()
	suite.FatalNoError(err)
}

func TestFakeDataSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger, _ := zap.NewDevelopment()

	hs := &FakeDataSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
