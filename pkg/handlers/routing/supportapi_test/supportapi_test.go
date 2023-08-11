package supportapi_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers/routing"
)

// SupportAPI tests using full routing
//
// These tests need to be in a package other than handlers/supportapi
// because otherwise import loops occur
// (supportapi -> routing -> supportapi)
type SupportAPISuite struct {
	routing.BaseRoutingSuite
}

func TestSupportAPISuite(t *testing.T) {
	hs := &SupportAPISuite{
		routing.NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
