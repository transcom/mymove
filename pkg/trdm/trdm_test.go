package trdm_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type TRDMSuite struct {
	*testingsuite.PopTestSuite
}

func (suite *TRDMSuite) SetupTest() {

}
func TestModelSuite(t *testing.T) {
	hs := &TRDMSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
