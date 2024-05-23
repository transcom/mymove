package ppmcloseout

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PPMCloseoutSuite struct {
	*testingsuite.PopTestSuite
}

func TestPPMCloseoutServiceSuite(t *testing.T) {
	ts := &PPMCloseoutSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
