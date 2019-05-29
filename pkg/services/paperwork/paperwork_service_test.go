package paperwork

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type PaperworkServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *PaperworkServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestPaperworkServiceSuite(t *testing.T) {

	hs := &PaperworkServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, hs)
}
