package shipment

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ShipmentServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *ShipmentServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}
func TestShipmentServiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &ShipmentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}
