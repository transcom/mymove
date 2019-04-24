package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PatchStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *PatchStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestPatchStorageInTransitSuite(t *testing.T) {

}

func (suite *PatchStorageInTransitSuite) TestPatchStorageInTransit() {

}
