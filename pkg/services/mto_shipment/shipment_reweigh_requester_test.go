package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestRequestShipmentReweigh() {
	requester := NewShipmentReweighRequester()

	suite.Run("If the shipment reweigh is requested successfully, it creates a reweigh in the DB", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		fetchedShipment := models.MTOShipment{}
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		reweigh, err := requester.RequestShipmentReweigh(session, shipment.ID, models.ReweighRequesterTOO)

		suite.NoError(err)

		var reweighShipment models.MTOShipment
		err = suite.DB().Where("id = ?", reweigh.ShipmentID).First(&reweighShipment)
		suite.NoError(err, "Get shipment associated with reweigh")

		suite.Equal(shipment.MoveTaskOrderID, reweighShipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.EqualValues(models.ReweighRequesterTOO, reweigh.RequestedBy)
		suite.WithinDuration(time.Now(), reweigh.RequestedAt, 2*time.Second)
	})

	suite.Run("When the shipment is not in a permitted status, returns a ConflictError", func() {
		rejectionReason := "rejection reason"
		rejectedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:          models.MTOShipmentStatusRejected,
					RejectionReason: &rejectionReason,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := requester.RequestShipmentReweigh(session, rejectedShipment.ID, models.ReweighRequesterTOO)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s is in a conflicting state Can only reweigh a shipment that is Approved or Diversion Requested. The shipment's current status is %s", rejectedShipment.ID, rejectedShipment.Status), err.Error())
	})

	suite.Run("When a reweigh already exists for the shipment, returns ConflictError", func() {
		reweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{})
		existingShipment := reweigh.Shipment
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := requester.RequestShipmentReweigh(session, existingShipment.ID, models.ReweighRequesterTOO)

		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s is in a conflicting state Cannot request a reweigh on a shipment that already has one.", existingShipment.ID), err.Error())
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		_, err := requester.RequestShipmentReweigh(session, badShipmentID, models.ReweighRequesterTOO)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
