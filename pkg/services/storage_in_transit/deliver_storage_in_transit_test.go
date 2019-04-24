package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type DeliverStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *DeliverStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestDeliverStorageInTransitSuite(t *testing.T) {

}

func (suite *DeliverStorageInTransitSuite) TestDeliverStorageInTransit() {

}
