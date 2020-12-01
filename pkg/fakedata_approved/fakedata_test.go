package fakedata_approved

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type FakeDataSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *FakeDataSuite) SetupTest() {
	suite.DB().TruncateAll()
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