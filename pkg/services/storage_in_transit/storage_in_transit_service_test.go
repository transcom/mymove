package storageintransit

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type StorageInTransitServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *StorageInTransitServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestStorageInTransitSuite(t *testing.T) {
	hs := &StorageInTransitServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix(("storage_in_transit"))),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, hs)
}

func TestStorageInTransitServiceSuite(t *testing.T) {
	hs := &StorageInTransitServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("storage_in_transit_service")),
	}
	suite.Run(t, hs)
}
