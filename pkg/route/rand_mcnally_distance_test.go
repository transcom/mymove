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

	suite.T().Run("pickupZip is greater than destinationZip", func(t *testing.T) {
		distance, err := randMcNallyZip3Distance(suite.DB(), "011", "010")
		suite.NoError(err)
		suite.Equal(24, distance)
	})

	suite.T().Run("pickupZip is the same as destinationZip", func(t *testing.T) {
		distance, err := randMcNallyZip3Distance(suite.DB(), "010", "010")
		suite.Equal(0, distance)
		suite.NotNil(err)
		suite.Equal("pickupZip (010) cannot be the same as destinationZip (010)", err.Error())
	})
}
