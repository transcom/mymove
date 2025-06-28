package gunsafeweightticket

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type GunSafeWeightTicketSuite struct {
	*testingsuite.PopTestSuite
}

func TestGunSafeWeightTicketServiceSuite(t *testing.T) {
	ts := &GunSafeWeightTicketSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
