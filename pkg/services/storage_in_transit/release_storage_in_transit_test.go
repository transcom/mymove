package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ReleaseStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ReleaseStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestReleaseStorageInTransitSuite(t *testing.T) {

}

func (suite *ReleaseStorageInTransitSuite) TestReleaseStorageInTransit() {

}
