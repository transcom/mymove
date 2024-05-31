package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/mocks"
)

func (suite *MTOShipmentServiceSuite) TestRequestShipmentDiversion() {
	router := NewShipmentRouter()
	requester := NewShipmentDiversionRequester(router)
	diversionReason := "Test Reason"

	suite.Run("If the shipment diversion is requested successfully, it should update the shipment status in the DB", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
				},
			},
		}, nil)
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		divertedShipment, err := requester.RequestShipmentDiversion(session, shipment.ID, shipmentEtag, &diversionReason)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, divertedShipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusDiversionRequested, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(divertedShipment.ID, fetchedShipment.ID)
		suite.Equal(divertedShipment.DiversionReason, fetchedShipment.DiversionReason)
	})

	suite.Run("When status transition is not allowed, returns a ConflictStatusError", func() {
		rejectionReason := "a rejection reason"
		rejectedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:          models.MTOShipmentStatusRejected,
					RejectionReason: &rejectionReason,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := requester.RequestShipmentDiversion(session, rejectedShipment.ID, eTag, &diversionReason)

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

		_, err := requester.RequestShipmentDiversion(session, staleShipment.ID, staleETag, &diversionReason)

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
		_, err := requester.RequestShipmentDiversion(session, badShipmentID, eTag, &diversionReason)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("It calls RequestDiversion on the ShipmentRouter", func() {
		shipmentRouter := &mocks.ShipmentRouter{}
		requester := NewShipmentDiversionRequester(shipmentRouter)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusApproved,
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

		shipmentRouter.On("RequestDiversion", mock.AnythingOfType("*appcontext.appContext"), &createdShipment, &diversionReason).Return(nil)

		_, err = requester.RequestShipmentDiversion(session, shipment.ID, eTag, &diversionReason)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(suite.T(), "RequestDiversion", 1)
	})
}
