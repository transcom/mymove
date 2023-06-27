package internalapi_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers/routing"
)

// InternalAPI tests using full routing
//
// These tests need to be in a package other than handlers/internalapi
// because otherwise import loops occur
// (internalapi -> routing -> internalapi)
type InternalAPISuite struct {
	routing.BaseRoutingSuite
}

func TestInternalAPISuite(t *testing.T) {
	hs := &InternalAPISuite{
		routing.NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
