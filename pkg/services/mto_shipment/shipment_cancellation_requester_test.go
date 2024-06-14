package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moveservices "github.com/transcom/mymove/pkg/services/move"
)

func (suite *MTOShipmentServiceSuite) TestRequestShipmentCancellation() {
	router := NewShipmentRouter()
	moveRouter := moveservices.NewMoveRouter()
	requester := NewShipmentCancellationRequester(router, moveRouter)

	suite.Run("If the shipment diversion is requested successfully, it should update the shipment status in the DB", func() {
		// valid pickupdate is anytime after the request to cancel date
		actualPickupDate := time.Now().AddDate(0, 0, 1)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		shipmentToBeCanceled, err := requester.RequestShipmentCancellation(session, shipment.ID, shipmentEtag)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, shipmentToBeCanceled.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusCancellationRequested, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(shipmentToBeCanceled.ID, fetchedShipment.ID)
	})

	suite.Run("When status transition is not allowed, returns a ConflictStatusError", func() {
		rejectionReason := "extraneous shipment"
		actualPickupDate := time.Now().AddDate(0, 0, 1)
		rejectedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusRejected,
					RejectionReason:  &rejectionReason,
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := requester.RequestShipmentCancellation(session, rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := requester.RequestShipmentCancellation(session, staleShipment.ID, staleETag)

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
		_, err := requester.RequestShipmentCancellation(session, badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("It calls RequestCancellation on the ShipmentRouter", func() {
		shipmentRouter := NewShipmentRouter()
		moveRouter := moveservices.NewMoveRouter()
		requester := NewShipmentCancellationRequester(shipmentRouter, moveRouter)
		// valid pickupdate is anytime after the request to cancel date
		actualPickupDate := time.Now().AddDate(0, 0, 1)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		_, err = requester.RequestShipmentCancellation(session, shipment.ID, eTag)

		suite.NoError(err)
		dbShipment := models.MTOShipment{}
		err = suite.DB().Find(&dbShipment, shipment.ID)
		suite.NoError(err)

		suite.FatalNoError(err)
		// if the created shipment has a status of cancellation requested, then RequestCancellation was successful
		suite.Equal(models.MTOShipmentStatusCancellationRequested, dbShipment.Status)
	})

	suite.Run("It calls RequestCancellation on shipment with invalid actualPickupDate", func() {
		shipmentRouter := NewShipmentRouter()
		moveRouter := moveservices.NewMoveRouter()
		requester := NewShipmentCancellationRequester(shipmentRouter, moveRouter)
		actualPickupDate := time.Now()
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:           models.MTOShipmentStatusApproved,
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		createdShipment := models.MTOShipment{
			ActualPickupDate: &actualPickupDate,
		}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		_, err = requester.RequestShipmentCancellation(session, shipment.ID, eTag)

		suite.Equal(err, apperror.NewUpdateError(shipment.ID, "cancellation request date cannot be on or after actual pickup date"))
	})
}
