package boatshipment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type BoatShipmentSuite struct {
	*testingsuite.PopTestSuite
}

func TestBoatShipmentServiceSuite(t *testing.T) {
	ts := &BoatShipmentSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
