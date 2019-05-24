package accesscode

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type AccessCodeServiceSuite struct {
	testingsuite.PopTestSuite
}

func (suite *AccessCodeServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestAccessCodeServiceSuite(t *testing.T) {
	ts := &AccessCodeServiceSuite{
		testingsuite.NewPopTestSuite(),
	}
	suite.Run(t, ts)
}
