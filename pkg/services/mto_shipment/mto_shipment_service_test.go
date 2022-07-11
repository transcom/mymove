package mtoshipment

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MTOShipmentServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestMTOShipmentServiceSuite(t *testing.T) {

	ts := &MTOShipmentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(
			testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction(),
		),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
