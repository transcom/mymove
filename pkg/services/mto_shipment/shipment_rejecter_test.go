package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestRejectShipment() {
	router := NewShipmentRouter()
	approver := NewShipmentRejecter(router)
	reason := "reason"

	suite.Run("If the shipment rejection is approved successfully, it should update the shipment status in the DB", func() {
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}

		rejectedShipment, err := approver.RejectShipment(suite.AppContextForTest(), shipment.ID, shipmentEtag, &reason)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, rejectedShipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusRejected, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(rejectedShipment.ID, fetchedShipment.ID)
		suite.Equal(&reason, fetchedShipment.RejectionReason)
	})

	suite.Run("When status transition is not allowed, returns a ConflictStatusError", func() {
		rejectionReason := "goods already shipped"
		rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:          models.MTOShipmentStatusRejected,
				RejectionReason: &rejectionReason,
			},
		})
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		_, err := approver.RejectShipment(suite.AppContextForTest(), rejectedShipment.ID, eTag, &reason)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())

		_, err := approver.RejectShipment(suite.AppContextForTest(), staleShipment.ID, staleETag, &reason)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := approver.RejectShipment(suite.AppContextForTest(), badShipmentID, eTag, &reason)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Passing in an empty rejection reason returns an InvalidInputError", func() {
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		emptyReason := ""

		_, err := approver.RejectShipment(suite.AppContextForTest(), shipment.ID, eTag, &emptyReason)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("It calls Reject on the ShipmentRouter", func() {
		shipmentRouter := &mocks.ShipmentRouter{}
		rejecter := NewShipmentRejecter(shipmentRouter)
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		shipmentRouter.On("Reject", mock.AnythingOfType("*appcontext.appContext"), &createdShipment, &reason).Return(nil)

		_, err = rejecter.RejectShipment(suite.AppContextForTest(), shipment.ID, eTag, &reason)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(suite.T(), "Reject", 1)
	})
}
