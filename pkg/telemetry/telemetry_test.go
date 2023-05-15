package telemetry

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type TelemetrySuite struct {
	*testingsuite.PopTestSuite
}

func TestTelemetrySuite(t *testing.T) {
	hs := &TelemetrySuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
