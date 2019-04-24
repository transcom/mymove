package storageintransit

import (
	"testing"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ApproveStorageInTransitSuite struct {
	testingsuite.PopTestSuite
}

func (suite *ApproveStorageInTransitSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestApproveStorageInTransitSuite(t *testing.T) {

}

func (suite *ApproveStorageInTransitSuite) TestApproveStorageInTransit() {

}
