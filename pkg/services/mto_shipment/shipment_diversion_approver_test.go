package mtoshipment

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApproveShipmentDiversion() {
	router := NewShipmentRouter()
	approver := NewShipmentDiversionApprover(router)

	suite.Run("If the shipment diversion is approved successfully, it should update the shipment status in the DB", func() {
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				Diversion: true,
			},
		})
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}

		divertedShipment, err := approver.ApproveShipmentDiversion(suite.AppContextForTest(), shipment.ID, shipmentEtag)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, divertedShipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusApproved, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(divertedShipment.ID, fetchedShipment.ID)
	})

	suite.Run("When status transition is not allowed, returns a ConflictStatusError", func() {
		rejectionReason := "we lost it"
		rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:          models.MTOShipmentStatusRejected,
				RejectionReason: &rejectionReason,
				Diversion:       true,
			},
		})
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		_, err := approver.ApproveShipmentDiversion(suite.AppContextForTest(), rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				Diversion: true,
			},
		})

		_, err := approver.ApproveShipmentDiversion(suite.AppContextForTest(), staleShipment.ID, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := approver.ApproveShipmentDiversion(suite.AppContextForTest(), badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("It calls ApproveDiversion on the ShipmentRouter", func() {
		shipmentRouter := &mocks.ShipmentRouter{}
		approver := NewShipmentDiversionApprover(shipmentRouter)
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status:    models.MTOShipmentStatusSubmitted,
				Diversion: true,
			},
		})
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		shipmentRouter.On("ApproveDiversion", mock.AnythingOfType("*appcontext.appContext"), &createdShipment).Return(nil)

		_, err = approver.ApproveShipmentDiversion(suite.AppContextForTest(), shipment.ID, eTag)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(suite.T(), "ApproveDiversion", 1)
	})
}
