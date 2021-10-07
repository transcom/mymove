package mtoshipment

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestRequestShipmentReweigh() {
	requester := NewShipmentReweighRequester()

	suite.T().Run("If the shipment reweigh is requested successfully, it creates a reweigh in the DB", func(t *testing.T) {
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		})
		fetchedShipment := models.MTOShipment{}

		reweigh, err := requester.RequestShipmentReweigh(suite.AppContextForTest(), shipment.ID, models.ReweighRequesterTOO)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, reweigh.Shipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(shipment.ID, reweigh.Shipment.ID)
		suite.EqualValues(models.ReweighRequesterTOO, reweigh.RequestedBy)
		suite.WithinDuration(time.Now(), reweigh.RequestedAt, 2*time.Second)
	})

	suite.T().Run("When the shipment is not in a permitted status, returns a ConflictError", func(t *testing.T) {
		rejectionReason := "rejection reason"
		rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:          models.MTOShipmentStatusRejected,
				RejectionReason: &rejectionReason,
			},
		})

		_, err := requester.RequestShipmentReweigh(suite.AppContextForTest(), rejectedShipment.ID, models.ReweighRequesterTOO)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Equal(fmt.Sprintf("id: %s is in a conflicting state Can only reweigh a shipment that is Approved or Diversion Requested. The shipment's current status is %s", rejectedShipment.ID, rejectedShipment.Status), err.Error())
	})

	suite.T().Run("When a reweigh already exists for the shipment, returns ConflictError", func(t *testing.T) {
		reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{})
		existingShipment := reweigh.Shipment

		_, err := requester.RequestShipmentReweigh(suite.AppContextForTest(), existingShipment.ID, models.ReweighRequesterTOO)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Equal(fmt.Sprintf("id: %s is in a conflicting state Cannot request a reweigh on a shipment that already has one.", existingShipment.ID), err.Error())
	})

	suite.T().Run("Passing in a bad shipment id returns a Not Found error", func(t *testing.T) {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := requester.RequestShipmentReweigh(suite.AppContextForTest(), badShipmentID, models.ReweighRequesterTOO)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
