package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type CreateStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *CreateStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestCreateStorageInTransitSuite(t *testing.T) {

}

func (suite *CreateStorageInTransitSuite) TestCreateStorageInTransit() {

}
