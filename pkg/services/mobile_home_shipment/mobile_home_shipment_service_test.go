package mobilehomeshipment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MobileHomeShipmentSuite struct {
	*testingsuite.PopTestSuite
}

func TestMobileHomeShipmentServiceSuite(t *testing.T) {
	ts := &MobileHomeShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}