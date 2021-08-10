package upload

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type UploadsServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func TestUploadsServiceSuite(t *testing.T) {

	hs := &UploadsServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
