package shipmentsummaryworksheet

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ShipmentSummaryWorksheetServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestShipmentSummaryWorksheetServiceSuite(t *testing.T) {
	ts := &ShipmentSummaryWorksheetServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
