package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
)

func (suite *MTOShipmentServiceSuite) TestApproveShipmentDiversion() {
	router := NewShipmentRouter()
	moveRouter := moverouter.NewMoveRouter()
	approver := NewShipmentDiversionApprover(router, moveRouter)

	suite.Run("If the shipment diversion is approved successfully, it should update the shipment status in the DB", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:    models.MTOShipmentStatusSubmitted,
					Diversion: true,
				},
			},
		}, nil)
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		divertedShipment, err := approver.ApproveShipmentDiversion(session, shipment.ID, shipmentEtag)

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
		rejectedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:          models.MTOShipmentStatusRejected,
					RejectionReason: &rejectionReason,
					Diversion:       true,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := approver.ApproveShipmentDiversion(session, rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:    models.MTOShipmentStatusSubmitted,
					Diversion: true,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := approver.ApproveShipmentDiversion(session, staleShipment.ID, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := approver.ApproveShipmentDiversion(session, badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("It calls ApproveDiversion on the ShipmentRouter", func() {
		shipmentRouter := NewShipmentRouter()
		approver := NewShipmentDiversionApprover(shipmentRouter, moveRouter)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:    models.MTOShipmentStatusSubmitted,
					Diversion: true,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := approver.ApproveShipmentDiversion(session, shipment.ID, eTag)
		suite.NoError(err)

		createdShipment := models.MTOShipment{}
		err = suite.DB().Find(&createdShipment, shipment.ID)
		suite.NoError(err)

		suite.FatalNoError(err)
		// if the created shipment has a status of approved, then ApproveDiversion was successful
		suite.Equal(models.MTOShipmentStatusApproved, createdShipment.Status)
	})
}
