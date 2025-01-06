package shipmentsummaryworksheet

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ShipmentSummaryWorksheetServiceSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *ShipmentSummaryWorksheetServiceSuite) SetupSuite() {
	suite.PreloadData(func() {
		factory.SetupDefaultAllotments(suite.DB())
	})
}

func TestShipmentSummaryWorksheetServiceSuite(t *testing.T) {
	ts := &ShipmentSummaryWorksheetServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
