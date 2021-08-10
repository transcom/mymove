package mtoshipment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/mocks"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestRejectShipment() {
	router := NewShipmentRouter()
	approver := NewShipmentRejecter(router)
	reason := "reason"

	suite.T().Run("If the shipment rejection is approved successfully, it should update the shipment status in the DB", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		shipmentEtag := etag.GenerateEtag(shipment.UpdatedAt)
		fetchedShipment := models.MTOShipment{}

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		rejectedShipment, err := approver.RejectShipment(appCfg, shipment.ID, shipmentEtag, &reason)

		suite.NoError(err)
		suite.Equal(shipment.MoveTaskOrderID, rejectedShipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipment.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusRejected, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)
		suite.Equal(rejectedShipment.ID, fetchedShipment.ID)
		suite.Equal(&reason, fetchedShipment.RejectionReason)
	})

	suite.T().Run("When status transition is not allowed, returns a ConflictStatusError", func(t *testing.T) {
		rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusRejected,
			},
		})
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := approver.RejectShipment(appCfg, rejectedShipment.ID, eTag, &reason)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.T().Run("Passing in a stale identifier returns a PreconditionFailedError", func(t *testing.T) {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := approver.RejectShipment(appCfg, staleShipment.ID, staleETag, &reason)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id returns a Not Found error", func(t *testing.T) {
		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := approver.RejectShipment(appCfg, badShipmentID, eTag, &reason)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("Passing in an empty rejection reason returns an InvalidInputError", func(t *testing.T) {
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)
		emptyReason := ""

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := approver.RejectShipment(appCfg, shipment.ID, eTag, &emptyReason)

		suite.Error(err)
		suite.IsType(services.InvalidInputError{}, err)
	})

	suite.T().Run("It calls Reject on the ShipmentRouter", func(t *testing.T) {
		shipmentRouter := &mocks.ShipmentRouter{}
		rejecter := NewShipmentRejecter(shipmentRouter)
		shipment := testdatagen.MakeDefaultMTOShipmentMinimal(suite.DB())
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)

		shipmentRouter.On("Reject", mock.AnythingOfType("*appconfig.appConfig"), &createdShipment, &reason).Return(nil)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err = rejecter.RejectShipment(appCfg, shipment.ID, eTag, &reason)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(t, "Reject", 1)
	})
}
