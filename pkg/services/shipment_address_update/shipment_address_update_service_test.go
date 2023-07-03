package shipmentaddressupdate

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ShipmentAddressUpdateServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestShipmentAddressUpdateServiceSuite(t *testing.T) {
	ts := &ShipmentAddressUpdateServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
