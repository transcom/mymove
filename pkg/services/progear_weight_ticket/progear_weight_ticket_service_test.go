package progearweightticket

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type ProgearWeightTicketSuite struct {
	*testingsuite.PopTestSuite
}

func TestProgearWeightTicketServiceSuite(t *testing.T) {
	ts := &ProgearWeightTicketSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
