package address

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type AddressSuite struct {
	*testingsuite.PopTestSuite
}

func TestAddressServiceSuite(t *testing.T) {
	ts := &AddressSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
