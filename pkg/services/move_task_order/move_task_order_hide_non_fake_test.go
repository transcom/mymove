package movetaskorder_test

import (
	"testing"

	. "github.com/transcom/mymove/pkg/services/move_task_order"
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
