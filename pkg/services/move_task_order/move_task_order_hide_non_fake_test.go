package movetaskorder

import (
	"testing"
)

func (suite *MoveTaskOrderServiceSuite) TestMoveTaskOrderHider_Hide() {
	mtoHider := NewMoveTaskOrderHider(suite.DB())
	suite.T().Run("Clear path", func(t *testing.T) {
		// Create move with service member using fake data
		result, err := mtoHider.Hide()
		suite.NoError(err)
		suite.Equal(0, len(result))
	})
}
