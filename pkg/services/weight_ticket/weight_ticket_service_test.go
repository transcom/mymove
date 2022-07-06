package weightticket

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type WeightTicketSuite struct {
	testingsuite.PopTestSuite
}

func TestWeightTicketServiceSuite(t *testing.T) {
	ts := &WeightTicketSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
