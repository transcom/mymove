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

	suite.Run("fails to cancel move with approved hhg shipment", func() {
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)

		_, err := moveCanceler.CancelMove(suite.AppContextForTest(), move.ID)
		suite.Error(err)
	})
}
