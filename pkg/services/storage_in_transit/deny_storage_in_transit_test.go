package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DenyStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *DenyStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestDenyStorageInTransitSuite(t *testing.T) {

}

func (suite *DenyStorageInTransitSuite) TestDenyStorageInTransit() {

}
