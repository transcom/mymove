package pricing

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type FetchServiceItemPriceTestSuite struct {
	*testingsuite.PopTestSuite
}

func TestMoveServiceSuite(t *testing.T) {

	hs := &FetchServiceItemPriceTestSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
