package adminapi_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers/routing"
)

// AdminAPI tests using full routing
//
// These tests need to be in a package other than handlers/adminapi
// because otherwise import loops occur
// (adminapi -> routing -> adminapi)
type AdminAPISuite struct {
	routing.BaseRoutingSuite
}

func TestAdminAPISuite(t *testing.T) {
	hs := &AdminAPISuite{
		routing.NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
