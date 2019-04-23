package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PlaceIntoSITStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *PlaceIntoSITStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestPlaceIntoSITStorageInTransitSuite(t *testing.T) {

}

func (suite *PlaceIntoSITStorageInTransitSuite) TestPlaceIntoSITStorageInTransit() {

}
