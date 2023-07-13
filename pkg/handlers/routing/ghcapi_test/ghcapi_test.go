package ghcapi_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers/routing"
)

// GhcAPI tests using full routing
//
// These tests need to be in a package other than handlers/ghcapi
// because otherwise import loops occur
// (ghcapi -> routing -> ghcapi)
type GhcAPISuite struct {
	routing.BaseRoutingSuite
}

func TestGhcAPISuite(t *testing.T) {
	hs := &GhcAPISuite{
		routing.NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
