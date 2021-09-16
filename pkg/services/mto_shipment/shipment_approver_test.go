package mtoshipment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	shipmentmocks "github.com/transcom/mymove/pkg/services/mocks"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestApproveShipment() {
	move := testdatagen.MakeAvailableMove(suite.DB())

	ghcDomesticTransitTime := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     0,
		WeightLbsUpper:     10000,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	verrs, err := suite.DB().ValidateAndCreate(&ghcDomesticTransitTime)
	suite.False(verrs.HasAny())
	suite.FatalNoError(err)

	// Let's also create a transit time object with a zero upper bound for weight (this can happen in the table).
	ghcDomesticTransitTime0LbsUpper := models.GHCDomesticTransitTime{
		MaxDaysTransitTime: 12,
		WeightLbsLower:     10001,
		WeightLbsUpper:     0,
		DistanceMilesLower: 0,
		DistanceMilesUpper: 10000,
	}
	verrs, err = suite.DB().ValidateAndCreate(&ghcDomesticTransitTime0LbsUpper)
	suite.False(verrs.HasAny())
	suite.FatalNoError(err)

	router := NewShipmentRouter()
	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(builder, moveRouter)
	planner := &mocks.Planner{}
	approver := NewShipmentApprover(router, siCreator, planner)

	suite.T().Run("If the mtoShipment is approved successfully it should create approved mtoServiceItems", func(t *testing.T) {
		shipmentForAutoApprove := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}
		var expectedReServiceCodes []models.ReServiceCode
		expectedReServiceCodes = append(expectedReServiceCodes,
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		)

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		shipment, approverErr := approver.ApproveShipment(suite.TestAppContext(), shipmentForAutoApprove.ID, shipmentForAutoApproveEtag)

		suite.NoError(approverErr)
		suite.Equal(move.ID, shipment.MoveTaskOrderID)

		err = suite.DB().Find(&fetchedShipment, shipmentForAutoApprove.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusApproved, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)

		err = suite.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipmentForAutoApprove.ID).All(&serviceItems)
		suite.NoError(err)

		suite.Equal(6, len(serviceItems))

		// All ApprovedAt times for service items should be the same, so just get the first one
		actualApprovedAt := serviceItems[0].ApprovedAt
		// If we've gotten the shipment updated and fetched it without error then we can inspect the
		// service items created as a side effect to see if they are approved.
		for i := range serviceItems {
			suite.Equal(models.MTOServiceItemStatusApproved, serviceItems[i].Status)
			suite.Equal(expectedReServiceCodes[i], serviceItems[i].ReService.Code)
			// Test that service item was approved within a few seconds of the current time
			suite.Assertions.WithinDuration(time.Now(), *actualApprovedAt, 2*time.Second)
		}
	})

	suite.T().Run("If we act on a shipment with a weight that has a 0 upper weight it should still work", func(t *testing.T) {
		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(11000)
		destinationAddress := testdatagen.MakeAddress2(suite.DB(), testdatagen.Assertions{})
		pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})

		shipmentHeavy := testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
				ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
				PrimeEstimatedWeight: &estimatedWeight,
				Status:               models.MTOShipmentStatusSubmitted,
				DestinationAddress:   &destinationAddress,
				DestinationAddressID: &destinationAddress.ID,
				PickupAddress:        &pickupAddress,
				PickupAddressID:      &pickupAddress.ID,
			},
		})

		createdShipment := models.MTOShipment{}
		err = suite.DB().Find(&createdShipment, shipmentHeavy.ID)
		suite.FatalNoError(err)
		err = suite.DB().Load(&createdShipment)
		suite.FatalNoError(err)

		planner.On("TransitDistance",
			createdShipment.PickupAddress,
			createdShipment.DestinationAddress,
		).Return(500, nil)

		shipmentHeavyEtag := etag.GenerateEtag(shipmentHeavy.UpdatedAt)
		_, err = approver.ApproveShipment(suite.TestAppContext(), shipmentHeavy.ID, shipmentHeavyEtag)
		suite.NoError(err)

		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipmentHeavy.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.T().Run("When status transition is not allowed, returns a ConflictStatusError", func(t *testing.T) {
		rejectionReason := "a reason"
		rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status:          models.MTOShipmentStatusRejected,
				RejectionReason: &rejectionReason,
			},
		})
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		_, err = approver.ApproveShipment(suite.TestAppContext(), rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.T().Run("Passing in a stale identifier returns a PreconditionFailedError", func(t *testing.T) {
		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})

		_, err = approver.ApproveShipment(suite.TestAppContext(), staleShipment.ID, staleETag)

		suite.Error(err)
		suite.IsType(services.PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id returns a Not Found error", func(t *testing.T) {
		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err = approver.ApproveShipment(suite.TestAppContext(), badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(services.NotFoundError{}, err)
	})

	suite.T().Run("It calls Approve on the ShipmentRouter", func(t *testing.T) {
		shipmentRouter := &shipmentmocks.ShipmentRouter{}
		approver := NewShipmentApprover(shipmentRouter, siCreator, planner)
		shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				Status: models.MTOShipmentStatusSubmitted,
			},
		})
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err = suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)
		err = suite.DB().Load(&createdShipment, "MoveTaskOrder", "PickupAddress", "DestinationAddress")
		suite.FatalNoError(err)

		shipmentRouter.On("Approve", mock.AnythingOfType("*appcontext.appContext"), &createdShipment).Return(nil)

		_, err := approver.ApproveShipment(suite.TestAppContext(), shipment.ID, eTag)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(t, "Approve", 1)
	})
}
