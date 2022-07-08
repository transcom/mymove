package sitextension

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type SitExtensionServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestSitExtensionServiceSuite(t *testing.T) {
	testService := &SitExtensionServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, testService)
	testService.PopTestSuite.TearDown()
}
