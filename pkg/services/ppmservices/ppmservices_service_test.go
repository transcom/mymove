package ppmservices

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMServiceSuite struct {
	testingsuite.PopTestSuite
}

func TestPPMServiceSuite(t *testing.T) {

	hs := &PPMServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
