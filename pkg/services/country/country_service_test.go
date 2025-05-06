package country

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type CountrySuite struct {
	*testingsuite.PopTestSuite
}

func TestCountryServiceSuite(t *testing.T) {
	ts := &CountrySuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)

	ts.PopTestSuite.TearDown()
}
