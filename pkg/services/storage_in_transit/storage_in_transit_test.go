package storageintransit

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type StorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *StorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestStorageInTransitSuite(t *testing.T) {
	hs := &StorageInTransitSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, hs)
}
