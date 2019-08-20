package shipmentlineitem

import (
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"

	"go.uber.org/zap"

	"testing"
)

type ShipmentLineItemServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *ShipmentLineItemServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestShipmentLineItemSuite(t *testing.T) {
	hs := &ShipmentLineItemServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(),
	}
	suite.Run(t, hs)
}
