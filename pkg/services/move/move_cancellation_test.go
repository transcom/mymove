package move

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MoveServiceSuite) TestCancelMove() {
	moveCancellation := NewMoveCancellation()

	suite.Run("successfully cancels a move", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		suite.NotEqual(move.Status, models.MoveStatusCANCELED)

		m, err := moveCancellation.CancelMove(suite.AppContextForTest(), move.ID)
		suite.NoError(err)
		suite.Equal(m.Status, models.MoveStatusCANCELED)
	})
}
