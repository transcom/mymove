package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type IndexStorageInTransitsSuite struct {
	testingsuite.PopTestSuite
}

func (suite *IndexStorageInTransitsSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestIndexStorageInTransitsSuite(t *testing.T) {

}

func (suite *CreateStorageInTransitSuite) TestIndexStorageInTransits() {

}
