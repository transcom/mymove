package reweigh

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
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
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
		}, nil)
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
		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
		}, nil)
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
	suite.Run("Update a diverted parent shipment reweigh and its child", func() {
		parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: nil,
				},
			},
		}, nil)
		childShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor:     true,
					Diversion:              true,
					DivertedFromShipmentID: &parentShipment.ID,
				},
			},
		}, nil)
		oldParentReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			MTOShipment: parentShipment,
		})
		eTag := etag.GenerateEtag(oldParentReweigh.UpdatedAt)

		// Update Parent
		newReweigh := oldParentReweigh
		newWeight := unit.Pound(200)
		newReweigh.Weight = &newWeight
		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, eTag)
		suite.NoError(err)
		suite.NotNil(updatedReweigh)
		suite.Equal(newWeight, *updatedReweigh.Weight)
		// Check that the child shipment now also has a reweigh as well since the parent just got updated
		reweighFetcher := NewReweighFetcher()
		reweighMap, err := reweighFetcher.ListReweighsByShipmentIDs(suite.AppContextForTest(), []uuid.UUID{childShipment.ID})
		suite.NotNil(reweighMap)
		suite.NoError(err)
		childReweigh := reweighMap[childShipment.ID]
		suite.Equal(childReweigh.ShipmentID, childShipment.ID)
	})
	suite.Run("Existing reweigh chains don't go wonky (Those created prior to 'chaning' logic)", func() {
		parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion: true,
				},
			},
		}, nil)
		childShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					UsesExternalVendor: true,
					Diversion:          true,
				},
			},
		}, nil)
		oldParentReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			MTOShipment: parentShipment,
		})
		eTag := etag.GenerateEtag(oldParentReweigh.UpdatedAt)

		// Update Parent
		newReweigh := oldParentReweigh
		newWeight := unit.Pound(200)
		newReweigh.Weight = &newWeight
		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, eTag)
		suite.NoError(err)
		suite.NotNil(updatedReweigh)
		suite.Equal(newWeight, *updatedReweigh.Weight)
		// Check that the child shipment did not receive a new reweigh
		reweighFetcher := NewReweighFetcher()
		reweighMap, err := reweighFetcher.ListReweighsByShipmentIDs(suite.AppContextForTest(), []uuid.UUID{childShipment.ID})
		suite.NotNil(reweighMap)
		suite.NoError(err)
		_, exists := reweighMap[childShipment.ID]
		suite.False(exists)
	})
	suite.Run("Update a diverted child shipment reweigh and its parent and the parent's grandchild", func() {
		parentShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: nil,
				},
			},
		}, nil)
		childShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: &parentShipment.ID,
				},
			},
		}, nil)
		grandChildShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &currentTime,
				},
			},
			{
				Model: models.MTOShipment{
					Diversion:              true,
					DivertedFromShipmentID: &childShipment.ID,
				},
			},
		}, nil)
		oldChildReweigh := testdatagen.MakeReweigh(suite.DB(), testdatagen.Assertions{
			MTOShipment: parentShipment,
		})
		eTag := etag.GenerateEtag(oldChildReweigh.UpdatedAt)

		// Update Child
		newReweigh := oldChildReweigh
		newWeight := unit.Pound(200)
		newReweigh.Weight = &newWeight
		updatedReweigh, err := reweighUpdater.UpdateReweighCheck(suite.AppContextForTest(), &newReweigh, eTag)
		suite.NoError(err)
		suite.NotNil(updatedReweigh)
		suite.Equal(newWeight, *updatedReweigh.Weight)
		// Check that the parent shipment now has the new lowest reweigh and that the grandchild has it too
		reweighFetcher := NewReweighFetcher()
		reweighMap, err := reweighFetcher.ListReweighsByShipmentIDs(suite.AppContextForTest(), []uuid.UUID{parentShipment.ID, grandChildShipment.ID})
		suite.NotNil(reweighMap)
		suite.NoError(err)
		// Parent
		parentReweigh := reweighMap[parentShipment.ID]
		suite.Equal(parentReweigh.ShipmentID, parentShipment.ID)
		// Grandchild
		grandChildReweigh := reweighMap[grandChildShipment.ID]
		suite.Equal(grandChildReweigh.ShipmentID, grandChildShipment.ID)
	})
}
