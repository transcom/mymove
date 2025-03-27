package acknowledgemovesshipments

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *AcknowledgeMovesAndShipmentsServiceSuite) TestUpdateMoveAcknowledgement() {
	moveAndShipmentAcknowledgementUpdater := NewMoveAndShipmentAcknowledgementUpdater()

	suite.Run("Move and Shipment acknowledgement dates are updated successfully", func() {

		var threeDaysAgo = time.Now().AddDate(0, 0, -3)
		var twoDaysAgo = time.Now().AddDate(0, 0, -2)
		var yesterday = time.Now().AddDate(0, 0, -1)
		move1 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		move2 := factory.BuildMoveWithShipment(suite.DB(), nil, nil)

		// Add an additional shipment to move 2
		factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move2,
				LinkOnly: true,
			},
		}, nil)
		move1.PrimeAcknowledgedAt = &threeDaysAgo
		move1.MTOShipments[0].PrimeAcknowledgedAt = &twoDaysAgo

		// grab move 2 from the db to pull in the 2nd shipment
		move2 = models.Move{
			ID: move2.ID,
		}
		err := suite.DB().EagerPreload("MTOShipments").Find(&move2, move2.ID)
		suite.NoError(err)

		move2.PrimeAcknowledgedAt = &threeDaysAgo
		move2.MTOShipments[0].PrimeAcknowledgedAt = &twoDaysAgo
		move2.MTOShipments[1].PrimeAcknowledgedAt = &yesterday
		moves := models.Moves{
			move1,
			move2,
		}
		err = moveAndShipmentAcknowledgementUpdater.AcknowledgeMovesAndShipments(suite.AppContextForTest(), &moves)
		suite.NoError(err)

		dbMove1 := models.Move{}

		// Validate move 1
		err = suite.DB().EagerPreload("MTOShipments").Find(&dbMove1, move1.ID)
		suite.NoError(err)
		suite.Equal(move1.ID, dbMove1.ID)
		suite.Equal(move1.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), dbMove1.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))
		// Move 1 shipment 1
		suite.Equal(move1.MTOShipments[0].ID, dbMove1.MTOShipments[0].ID)
		suite.Equal(move1.MTOShipments[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), dbMove1.MTOShipments[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))

		dbMove2 := models.Move{}
		// Validate move 2
		err = suite.DB().EagerPreload("MTOShipments").Find(&dbMove2, move2.ID)
		suite.NoError(err)
		suite.Equal(move2.ID, dbMove2.ID)
		suite.Equal(move2.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), dbMove2.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))

		suite.Equal(len(move2.MTOShipments), len(dbMove2.MTOShipments))

		// Verify that both the shipments match
		matchingShipments := 0
		for _, shipment := range move2.MTOShipments {
			for _, dbShipment := range dbMove2.MTOShipments {
				if shipment.ID == dbShipment.ID {
					suite.Equal(shipment.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), dbShipment.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond))
					matchingShipments++
					break
				}
			}
		}
		suite.Equal(matchingShipments, len(move2.MTOShipments))
	})

	suite.Run("Move and Shipment acknowledgement date are not updated when they are not provided", func() {

		move := factory.BuildMoveWithShipment(suite.DB(), nil, nil)
		moves := models.Moves{
			move,
		}
		err := moveAndShipmentAcknowledgementUpdater.AcknowledgeMovesAndShipments(suite.AppContextForTest(), &moves)
		suite.NoError(err)

		dbMove := models.Move{}
		err = suite.DB().EagerPreload("MTOShipments").Find(&dbMove, move.ID)
		suite.NoError(err)
		suite.Equal(move.ID, dbMove.ID)
		suite.Nil(dbMove.PrimeAcknowledgedAt)
		suite.Equal(move.MTOShipments[0].ID, dbMove.MTOShipments[0].ID)
		suite.Nil(dbMove.MTOShipments[0].PrimeAcknowledgedAt)
	})

	suite.Run("Move and Shipment acknowledgement dates are NOT updated if they are already populated in the DB", func() {

		var fourDaysAgo = time.Now().AddDate(0, 0, -4)
		var threeDaysAgo = time.Now().AddDate(0, 0, -3)
		var twoDaysAgo = time.Now().AddDate(0, 0, -2)
		move := factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					PrimeAcknowledgedAt: &fourDaysAgo,
				},
			},
			{
				Model: models.MTOShipment{
					PrimeAcknowledgedAt: &fourDaysAgo,
				},
			},
		}, nil)

		//Attempting to update prime acknowledged at dates that are already popualted in the DB.
		move.PrimeAcknowledgedAt = &threeDaysAgo
		move.MTOShipments[0].PrimeAcknowledgedAt = &twoDaysAgo
		moves := models.Moves{
			move,
		}
		err := moveAndShipmentAcknowledgementUpdater.AcknowledgeMovesAndShipments(suite.AppContextForTest(), &moves)
		suite.NoError(err)

		dbMove := models.Move{}
		err = suite.DB().EagerPreload("MTOShipments").Find(&dbMove, move.ID)
		suite.NoError(err)
		suite.Equal(move.ID, dbMove.ID)
		suite.NotEqual(move.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), dbMove.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), "Prime acknowledged at date was not updated")
		suite.Equal(fourDaysAgo.UTC().Truncate(time.Millisecond), dbMove.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), "Prime acknowledged at date remained 4 days ago")
		suite.Equal(move.MTOShipments[0].ID, dbMove.MTOShipments[0].ID)
		suite.NotEqual(move.MTOShipments[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), dbMove.MTOShipments[0].PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), "Prime acknowledged at date was not updated")
		suite.Equal(fourDaysAgo.UTC().Truncate(time.Millisecond), dbMove.PrimeAcknowledgedAt.UTC().Truncate(time.Millisecond), "Prime acknowledged at date remained 4 days ago")
	})
}
