package sitaddressupdate

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SITAddressUpdateServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestSITAddressUpdateServiceSuite(t *testing.T) {
	ts := &SITAddressUpdateServiceSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
