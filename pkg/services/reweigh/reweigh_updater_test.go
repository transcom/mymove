package reweigh

import (
	"time"

	"github.com/transcom/mymove/pkg/apperror"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
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
	mockPlanner.On("ZipTransitDistance",
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

	// Test Success - Reweigh updated
	suite.Run("Updated reweigh - Success", func() {
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
		newWeight := unit.Pound(200)
		newReweigh.Weight = &newWeight
		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, eTag)

		suite.NoError(err)
		suite.NotNil(updatedReweigh)
		suite.Equal(newWeight, *updatedReweigh.Weight)
	})
	// Test NotFoundError
	suite.Run("Not Found Error", func() {
		notFoundReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			Stub: true,
			Reweigh: models.Reweigh{
				ID: uuid.Must(uuid.NewV4()),
			},
		})
		eTag := etag.GenerateEtag(time.Now())

		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &notFoundReweigh, eTag)

		suite.Nil(updatedReweigh)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), notFoundReweigh.ID.String())
	})
	// PreconditionFailedError
	suite.Run("Precondition Failed", func() {
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: models.Move{
				AvailableToPrimeAt: &currentTime,
			},
		})
		oldReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			MTOShipment: shipment,
		})
		// bad etag value
		eTag := etag.GenerateEtag(time.Now())
		newReweigh := oldReweigh

		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, eTag) // base validation

		suite.Nil(updatedReweigh)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}
