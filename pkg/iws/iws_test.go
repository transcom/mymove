package iws

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type iwsSuite struct {
	testingsuite.BaseTestSuite
}

func (suite *iwsSuite) SetupSuite() {
}

func TestIwsSuite(t *testing.T) {
	suite.Run(t, new(iwsSuite))
}
