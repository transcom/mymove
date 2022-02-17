package ppmshipment

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	notificationMocks "github.com/transcom/mymove/pkg/notifications/mocks"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/fetch"
	mockservices "github.com/transcom/mymove/pkg/services/mocks"
	moveservices "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {
	// TODO: organize variables / is there a better way to do this? Any way to make this generic?
	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	planner := &mocks.Planner{}
	planner.On("TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(500, nil)
	moveRouter := moveservices.NewMoveRouter()
	moveWeights := moveservices.NewMoveWeights(mtoshipment.NewShipmentReweighRequester())
	mockShipmentRecalculator := mockservices.PaymentRequestShipmentRecalculator{}
	mockShipmentRecalculator.On("ShipmentRecalculatePaymentRequest",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("uuid.UUID"),
	).Return(&models.PaymentRequests{}, nil)
	mockSender := notificationMocks.NotificationSender{}
	mockSender.On("SendNotification",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.AnythingOfType("*notifications.ReweighRequested"),
	).Return(nil)

	mtoShipmentUpdater := mtoshipment.NewMTOShipmentUpdater(builder, fetcher, planner, moveRouter, moveWeights, &mockSender, &mockShipmentRecalculator)
	ppmShipmentUpdater := NewPPMShipmentUpdater(mtoShipmentUpdater)

	oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

	newPPM := models.PPMShipment{
		ID:          oldPPMShipment.ID,
		SitExpected: models.BoolPointer(true),
	}

	eTag := etag.GenerateEtag(oldPPMShipment.UpdatedAt)
	suite.T().Run("UpdatePPMShipment - Success", func(t *testing.T) {
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, eTag)

		suite.NilOrNoVerrs(err)
		suite.True(*updatedPPMShipment.SitExpected)
		suite.Equal(unit.Pound(1150), *updatedPPMShipment.ProGearWeight)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		ppmForNotFound := models.PPMShipment{
			ID:          uuid.Nil,
			SitExpected: models.BoolPointer(true),
		}
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &ppmForNotFound, eTag)

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Precondition Failed", func(t *testing.T) {
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, "")

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}
