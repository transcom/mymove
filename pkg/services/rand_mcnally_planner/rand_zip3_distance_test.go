package randmcnally

import (
	"testing"
)

func (suite *RandMcNallyPlannerServiceSuite) TestRandMcNallyZip3Distance() {
	suite.T().Run("test basic distance check", func(t *testing.T) {
		randMcnally := NewRandMcNallyZip3Distance(suite.logger)
		distance, err := randMcnally.RandMcNallyZip3Distance("003", "030")
		suite.NoError(err)
		suite.Equal(-1, distance)
	})
}
