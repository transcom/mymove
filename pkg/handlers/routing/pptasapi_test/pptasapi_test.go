package pptasapi_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/handlers/routing"
)

// Pptas tests using full routing
//
// These tests need to be in a package other than handlers/pptas
// because otherwise import loops occur
// (pptas -> routing -> pptas)
type PPTASAPISuite struct {
	routing.BaseRoutingSuite
}

func TestPptasSuite(t *testing.T) {
	hs := &PPTASAPISuite{
		routing.NewBaseRoutingSuite(),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
