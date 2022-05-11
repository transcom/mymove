package reweigh

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	recalculateTestPickupZip      = "30907"
	recalculateTestDestinationZip = "78234"
	recalculateTestZip3Distance   = 1234
)

func (suite *ReweighSuite) TestReweighUpdater() {

	// Mock out a planner.
	mockPlanner := &routemocks.Planner{}
	mockPlanner.On("Zip3TransitDistance",
		recalculateTestPickupZip,
		recalculateTestDestinationZip,
	).Return(recalculateTestZip3Distance, nil)

	// Get shipment payment request recalculator service
	creator := paymentrequest.NewPaymentRequestCreator(mockPlanner, ghcrateengine.NewServiceItemPricer())
	statusUpdater := paymentrequest.NewPaymentRequestStatusUpdater(query.NewQueryBuilder())
	recalculator := paymentrequest.NewPaymentRequestRecalculator(creator, statusUpdater)
	paymentRequestShipmentRecalculator := paymentrequest.NewPaymentRequestShipmentRecalculator(recalculator)

	reweighUpdater := NewReweighUpdater(movetaskorder.NewMoveTaskOrderChecker(), paymentRequestShipmentRecalculator)
	currentTime := time.Now()
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &currentTime,
		},
	})
	oldReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
		MTOShipment: shipment,
	})
	eTag := etag.GenerateEtag(oldReweigh.UpdatedAt)
	newReweigh := oldReweigh

	// Test Success - Reweigh updated
	suite.T().Run("Updated reweigh - Success", func(t *testing.T) {
		newWeight := unit.Pound(200)
		newReweigh.Weight = &newWeight
		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, eTag)

		suite.NoError(err)
		suite.NotNil(updatedReweigh)
		suite.Equal(newWeight, *updatedReweigh.Weight)
		eTag = etag.GenerateEtag(updatedReweigh.UpdatedAt)
	})
	// Test NotFoundError
	suite.T().Run("Not Found Error", func(t *testing.T) {
		notFoundUUID := "00000000-0000-0000-0000-000000000001"
		notFoundReweigh := newReweigh
		notFoundReweigh.ID = uuid.FromStringOrNil(notFoundUUID)

		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &notFoundReweigh, eTag)

		suite.Nil(updatedReweigh)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundUUID)
	})
	// PreconditionFailedError
	suite.T().Run("Precondition Failed", func(t *testing.T) {
		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, "nada") // base validation

		suite.Nil(updatedReweigh)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}
