package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DeleteStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *DeleteStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestDeleteStorageInTransitSuite(t *testing.T) {

}

func (suite *ReleaseStorageInTransitSuite) TestDeleteStorageInTransit() {

}
