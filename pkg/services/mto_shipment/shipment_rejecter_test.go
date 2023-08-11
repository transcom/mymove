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

func (suite *MTOShipmentServiceSuite) TestRejectShipment() {
	router := NewShipmentRouter()
	approver := NewShipmentRejecter(router)
	reason := "reason"

	suite.Run("If the shipment rejection is approved successfully, it should update the shipment status in the DB", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		rejectedShipment, err := approver.RejectShipment(session, shipment.ID, shipmentEtag, &reason)

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

		_, err := approver.RejectShipment(session, rejectedShipment.ID, eTag, &reason)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := approver.RejectShipment(session, staleShipment.ID, staleETag, &reason)

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

		_, err := approver.RejectShipment(session, badShipmentID, eTag, &reason)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("Passing in an empty rejection reason returns an InvalidInputError", func() {
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		emptyReason := ""
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})

		_, err := approver.RejectShipment(session, shipment.ID, eTag, &emptyReason)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
	})

	suite.Run("It calls Reject on the ShipmentRouter", func() {
		shipmentRouter := &mocks.ShipmentRouter{}
		rejecter := NewShipmentRejecter(shipmentRouter)
		shipment := factory.BuildMTOShipmentMinimal(suite.DB(), nil, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		session := suite.AppContextWithSessionForTest(&auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    uuid.Must(uuid.NewV4()),
		})
		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		shipmentRouter.On("Reject", mock.AnythingOfType("*appcontext.appContext"), &createdShipment, &reason).Return(nil)

		_, err = rejecter.RejectShipment(session, shipment.ID, eTag, &reason)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(suite.T(), "Reject", 1)
	})
}
