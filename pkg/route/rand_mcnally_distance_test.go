package route

import (
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *GHCTestSuite) TestRandMcNallyZip3Distance() {
	testdatagen.MakeDefaultZip3Distance(suite.DB())

	suite.T().Run("test basic distance check", func(t *testing.T) {
		distance, err := randMcNallyZip3Distance(suite.DB(), "010", "011")
		suite.NoError(err)
		suite.Equal(24, distance)
	})

	suite.T().Run("fromZip3 is greater than toZip3", func(t *testing.T) {
		distance, err := randMcNallyZip3Distance(suite.DB(), "011", "010")
		suite.NoError(err)
		suite.Equal(24, distance)
	})

	suite.T().Run("fromZip3 is the same as toZip3", func(t *testing.T) {
		distance, err := randMcNallyZip3Distance(suite.DB(), "010", "010")
		suite.Equal(0, distance)
		suite.NotNil(err)
		suite.Equal("fromZip3 (010) cannot be the same as toZip3 (010)", err.Error())
	})
}
