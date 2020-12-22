package randmcnally

import (
	"testing"

	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *RandMcNallyPlannerServiceSuite) TestRandMcNallyZip3Distance() {
	testdatagen.MakeDefaultZip3Distance(suite.DB())

	suite.T().Run("test basic distance check", func(t *testing.T) {
		randMcnally := NewRandMcNallyZip3Distance(suite.DB(), suite.logger)
		distance, err := randMcnally.RandMcNallyZip3Distance("010", "011")
		suite.NoError(err)
		suite.Equal(24, distance)
	})

	suite.T().Run("pickupZip is greater than destinationZip", func(t *testing.T) {
		randMcnally := NewRandMcNallyZip3Distance(suite.DB(), suite.logger)
		distance, err := randMcnally.RandMcNallyZip3Distance("011", "010")
		suite.NoError(err)
		suite.Equal(24, distance)
	})

	suite.T().Run("pickupZip is the same as destinationZip", func(t *testing.T) {
		randMcnally := NewRandMcNallyZip3Distance(suite.DB(), suite.logger)
		distance, err := randMcnally.RandMcNallyZip3Distance("010", "010")
		suite.Equal(0, distance)
		suite.NotNil(err)
		suite.Equal("Data received from requester is bad: BAD_DATA: pickupZip cannot be the same as destinationZip", err.Error())
	})
}
