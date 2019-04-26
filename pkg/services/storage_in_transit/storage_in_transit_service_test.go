package storageintransit

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type StorageInTransitServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *StorageInTransitServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestStorageInTransitServiceSuite(t *testing.T) {
	hs := &StorageInTransitServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, hs)
}
