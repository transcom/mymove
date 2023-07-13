package primeapi_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers/routing"
)

// PrimeAPI tests using full routing
//
// These tests need to be in a package other than handlers/primeapi
// because otherwise import loops occur
// (primeapi -> routing -> primeapi)
type PrimeAPISuite struct {
	routing.BaseRoutingSuite
}

func TestPrimeAPISuite(t *testing.T) {
	hs := &PrimeAPISuite{
		routing.NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
