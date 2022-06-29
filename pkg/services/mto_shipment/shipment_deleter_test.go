package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestShipmentDeleter() {
	suite.Run("Returns an error when shipment is not found", func() {
		shipmentDeleter := NewShipmentDeleter()
		uuid := uuid.Must(uuid.NewV4())

		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), uuid)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Returns an error when the Move is neither in Draft nor in NeedsServiceCounseling status", func() {
		shipmentDeleter := NewShipmentDeleter()
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		move := shipment.MoveTaskOrder
		move.Status = models.MoveStatusServiceCounselingCompleted
		suite.MustSave(&move)

		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)

		suite.Error(err)
		suite.IsType(apperror.ForbiddenError{}, err)
	})

	suite.Run("Soft deletes the shipment when it is found", func() {
		shipmentDeleter := NewShipmentDeleter()
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())

		validStatuses := []struct {
			desc   string
			status models.MoveStatus
		}{
			{"Draft", models.MoveStatusDRAFT},
			{"Needs Service Counseling", models.MoveStatusNeedsServiceCounseling},
		}
		for _, validStatus := range validStatuses {
			move := shipment.MoveTaskOrder
			move.Status = validStatus.status
			suite.MustSave(&move)

			moveID, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)
			suite.NoError(err)
			// Verify that the shipment's Move ID is returned because the
			// handler needs it to generate the TriggerEvent.
			suite.Equal(shipment.MoveTaskOrderID, moveID)

			// Verify the shipment still exists in the DB
			var shipmentInDB models.MTOShipment
			err = suite.DB().Find(&shipmentInDB, shipment.ID)
			suite.NoError(err)

			actualDeletedAt := shipmentInDB.DeletedAt
			suite.WithinDuration(time.Now(), *actualDeletedAt, 2*time.Second)

			// Reset the deleted_at field to nil to allow the shipment to be
			// deleted a second time when testing the other move status (a
			// shipment can only be deleted once)
			shipmentInDB.DeletedAt = nil
			suite.MustSave(&shipment)
		}
	})

	suite.Run("Returns not found error when the shipment is already deleted", func() {
		shipmentDeleter := NewShipmentDeleter()
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)

		suite.NoError(err)

		// Try to delete the shipment a second time
		_, err = shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Soft deletes the associated PPM shipment", func() {
		shipmentDeleter := NewShipmentDeleter()
		ppmShipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
		})

		moveID, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), ppmShipment.ShipmentID)
		suite.NoError(err)
		// Verify that the shipment's Move ID is returned because the
		// handler needs it to generate the TriggerEvent.
		suite.Equal(ppmShipment.Shipment.MoveTaskOrderID, moveID)

		// Verify the shipment still exists in the DB
		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("PPMShipment").Find(&shipmentInDB, ppmShipment.ShipmentID)
		suite.NoError(err)

		actualDeletedAt := shipmentInDB.DeletedAt
		suite.WithinDuration(time.Now(), *actualDeletedAt, 2*time.Second)

		actualDeletedAt = shipmentInDB.PPMShipment.DeletedAt
		suite.WithinDuration(time.Now(), *actualDeletedAt, 2*time.Second)
	})
}

func (suite *MTOShipmentServiceSuite) TestPrimeShipmentDeleter() {
	suite.Run("Doesn't return an error when allowed to delete a shipment", func() {
		shipmentDeleter := NewPrimeShipmentDeleter()
		now := time.Now()
		shipment := testdatagen.MakePPMShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
			PPMShipment: models.PPMShipment{
				Status: models.PPMShipmentStatusSubmitted,
			},
		})

		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)
		suite.Error(err)
	})

	suite.Run("Returns an error when a shipment is not available to prime", func() {
		shipmentDeleter := NewPrimeShipmentDeleter()

		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: nil,
			},
		})

		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)
		suite.Error(err)
	})

	suite.Run("Returns an error when a shipment is not a PPM", func() {
		shipmentDeleter := NewPrimeShipmentDeleter()
		now := time.Now()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
		})

		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)
		suite.Error(err)
	})

	suite.Run("Returns an error when PPM status is WAITING_ON_CUSTOMER", func() {
		shipmentDeleter := NewPrimeShipmentDeleter()
		now := time.Now()
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &now,
			},
			MTOShipment: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHG,
			},
		})

		_, err := shipmentDeleter.DeleteShipment(suite.AppContextForTest(), shipment.ID)
		suite.Error(err)
	})
}
