package mtoshipment

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestApproveShipmentDiversion() {
	router := NewShipmentRouter(suite.DB())
	approver := NewShipmentDiversionApprover(suite.DB(), router)

	suite.T().Run("If the shipment diversion is approved successfully, it should update the shipment status in the DB", func(t *testing.T) {
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusDiversionRequested,
			},
		})
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}

		divertedShipment, err := approver.ApproveShipmentDiversion(shipment.ID, shipmentEtag)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, divertedShipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusApproved, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(divertedShipment.ID, fetchedShipment.ID)
	})

	suite.T().Run("When status transition is not allowed, returns a ConflictStatusError", func(t *testing.T) {
		rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusRejected,
			},
		})
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		_, err := approver.ApproveShipmentDiversion(rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.T().Run("Passing in a stale identifier returns a PreconditionFailedError", func(t *testing.T) {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusDiversionRequested,
			},
		})

		_, err := approver.ApproveShipmentDiversion(staleShipment.ID, staleETag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id returns a Not Found error", func(t *testing.T) {
		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := approver.ApproveShipmentDiversion(badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("It calls ApproveDiversion on the ShipmentRouter", func(t *testing.T) {
		shipmentRouter := &mocks.ShipmentRouter{}
		approver := NewShipmentDiversionApprover(suite.DB(), shipmentRouter)
		shipment := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusDiversionRequested,
			},
		})
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		shipmentRouter.On("ApproveDiversion", &createdShipment).Return(nil)

		_, err = approver.ApproveShipmentDiversion(shipment.ID, eTag)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(t, "ApproveDiversion", 1)
	})
}
