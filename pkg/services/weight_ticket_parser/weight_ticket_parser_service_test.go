package weightticketparser

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type WeightTicketParserServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestShipmentSummaryWorksheetServiceSuite(t *testing.T) {
	ts := &WeightTicketParserServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
