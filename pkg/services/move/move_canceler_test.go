package move

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MoveServiceSuite) TestMoveCanceler() {
	moveCanceler := NewMoveCanceler()

	suite.Run("successfully cancels a move", func() {
		move := factory.BuildMove(suite.DB(), nil, nil)

		suite.NotEqual(move.Status, models.MoveStatusCANCELED)

		m, err := moveCanceler.CancelMove(suite.AppContextForTest(), move.ID)
		suite.NoError(err)
		suite.Equal(m.Status, models.MoveStatusCANCELED)
	})
}
