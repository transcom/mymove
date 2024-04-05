package mtoshipment

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	shipmentmocks "github.com/transcom/mymove/pkg/services/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

type approveShipmentSubtestData struct {
	appCtx                 appcontext.AppContext
	move                   models.Move
	planner                *mocks.Planner
	shipmentApprover       services.ShipmentApprover
	mockedShipmentApprover services.ShipmentApprover
	mockedShipmentRouter   *shipmentmocks.ShipmentRouter
	reServiceCodes         []models.ReServiceCode
}

// Creates data for the TestApproveShipment function
func (suite *MTOShipmentServiceSuite) createApproveShipmentSubtestData() (subtestData *approveShipmentSubtestData) {
	subtestData = &approveShipmentSubtestData{}

	subtestData.move = factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

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

	// Let's create service codes in the DB
	subtestData.reServiceCodes = []models.ReServiceCode{
		models.ReServiceCodeDLH,
		models.ReServiceCodeFSC,
		models.ReServiceCodeDOP,
		models.ReServiceCodeDDP,
		models.ReServiceCodeDPK,
		models.ReServiceCodeDUPK,
	}

	for _, serviceCode := range subtestData.reServiceCodes {
		factory.BuildReServiceByCode(suite.DB(), serviceCode)
	}

	subtestData.mockedShipmentRouter = &shipmentmocks.ShipmentRouter{}

	router := NewShipmentRouter()

	builder := query.NewQueryBuilder()
	moveRouter := moverouter.NewMoveRouter()
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(planner, builder, moveRouter)
	subtestData.planner = &mocks.Planner{}

	subtestData.shipmentApprover = NewShipmentApprover(router, siCreator, subtestData.planner)
	subtestData.mockedShipmentApprover = NewShipmentApprover(subtestData.mockedShipmentRouter, siCreator, subtestData.planner)
	subtestData.appCtx = suite.AppContextWithSessionForTest(&auth.Session{
		ApplicationName: auth.OfficeApp,
		OfficeUserID:    uuid.Must(uuid.NewV4()),
	})

	return subtestData
}

func (suite *MTOShipmentServiceSuite) TestApproveShipment() {
	suite.Run("If the mtoShipment is approved successfully it should create approved mtoServiceItems", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		shipmentForAutoApprove := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		preApprovalTime := time.Now()
		shipment, approverErr := approver.ApproveShipment(appCtx, shipmentForAutoApprove.ID, shipmentForAutoApproveEtag)

		suite.NoError(approverErr)
		suite.Equal(move.ID, shipment.MoveTaskOrderID)

		err := appCtx.DB().Find(&fetchedShipment, shipmentForAutoApprove.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusApproved, fetchedShipment.Status)
		suite.Equal(shipment.ID, fetchedShipment.ID)

		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipmentForAutoApprove.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)

		suite.Equal(6, len(serviceItems))

		// All ApprovedAt times for service items should be the same, so just get the first one
		// Test that service item was approved within a few seconds of the current time
		suite.Assertions.WithinDuration(preApprovalTime, *serviceItems[0].ApprovedAt, 2*time.Second)

		// If we've gotten the shipment updated and fetched it without error then we can inspect the
		// service items created as a side effect to see if they are approved.
		for i := range serviceItems {
			suite.Equal(models.MTOServiceItemStatusApproved, serviceItems[i].Status)
			suite.Equal(subtestData.reServiceCodes[i], serviceItems[i].ReService.Code)
		}
	})

	suite.Run("approves shipment of type PPM and loads PPMShipment association", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		shipmentForAutoApprove := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model:    subtestData.move,
				LinkOnly: true,
			},
		}, nil)
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.Shipment.UpdatedAt)

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		shipment, approverErr := approver.ApproveShipment(appCtx, shipmentForAutoApprove.Shipment.ID, shipmentForAutoApproveEtag)

		suite.NoError(approverErr)
		suite.Equal(move.ID, shipment.MoveTaskOrderID)

		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
		suite.Equal(shipment.ID, shipmentForAutoApprove.Shipment.ID)

		suite.Equal(shipmentForAutoApprove.ID, shipment.PPMShipment.ID)
		suite.Equal(models.PPMShipmentStatusSubmitted, shipment.PPMShipment.Status)
	})

	suite.Run("If we act on a shipment with a weight that has a 0 upper weight it should still work", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(11000)
		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress2})
		pickupAddress := factory.BuildAddress(suite.DB(), nil, nil)

		shipmentHeavy := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipmentHeavy.ID)
		suite.FatalNoError(err)
		err = suite.DB().Load(&createdShipment)
		suite.FatalNoError(err)

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			createdShipment.PickupAddress.PostalCode,
			createdShipment.DestinationAddress.PostalCode,
		).Return(500, nil)

		shipmentHeavyEtag := etag.GenerateEtag(shipmentHeavy.UpdatedAt)
		_, err = approver.ApproveShipment(appCtx, shipmentHeavy.ID, shipmentHeavyEtag)
		suite.NoError(err)

		fetchedShipment := models.MTOShipment{}
		err = suite.DB().Find(&fetchedShipment, shipmentHeavy.ID)
		suite.NoError(err)
		// We also should have a required delivery date
		suite.NotNil(fetchedShipment.RequiredDeliveryDate)
	})

	suite.Run("When status transition is not allowed, returns a ConflictStatusError", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover

		rejectionReason := "a reason"
		rejectedShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:          models.MTOShipmentStatusRejected,
					RejectionReason: &rejectionReason,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(rejectedShipment.UpdatedAt)

		_, err := approver.ApproveShipment(appCtx, rejectedShipment.ID, eTag)

		suite.Error(err)
		suite.IsType(ConflictStatusError{}, err)
	})

	suite.Run("Passing in a stale identifier returns a PreconditionFailedError", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover

		staleETag := etag.GenerateEtag(time.Now())
		staleShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)

		_, err := approver.ApproveShipment(appCtx, staleShipment.ID, staleETag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})

	suite.Run("Passing in a bad shipment id returns a Not Found error", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		approver := subtestData.shipmentApprover

		eTag := etag.GenerateEtag(time.Now())
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := approver.ApproveShipment(appCtx, badShipmentID, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.Run("It calls Approve on the ShipmentRouter", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.mockedShipmentApprover
		shipmentRouter := subtestData.mockedShipmentRouter

		shipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status: models.MTOShipmentStatusSubmitted,
				},
			},
		}, nil)
		eTag := etag.GenerateEtag(shipment.UpdatedAt)

		createdShipment := models.MTOShipment{}
		err := suite.DB().Find(&createdShipment, shipment.ID)
		suite.FatalNoError(err)
		err = suite.DB().Load(&createdShipment, "MoveTaskOrder", "PickupAddress", "DestinationAddress")
		suite.FatalNoError(err)

		shipmentRouter.On("Approve", mock.AnythingOfType("*appcontext.appContext"), &createdShipment).Return(nil)

		_, err = approver.ApproveShipment(appCtx, shipment.ID, eTag)

		suite.NoError(err)
		shipmentRouter.AssertNumberOfCalls(suite.T(), "Approve", 1)
	})

	suite.Run("If the mtoShipment uses external vendor not allowed to approve shipment", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		shipmentForAutoApprove := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					Status:             models.MTOShipmentStatusSubmitted,
					UsesExternalVendor: true,
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				},
			},
		}, nil)
		shipmentForAutoApproveEtag := etag.GenerateEtag(shipmentForAutoApprove.UpdatedAt)
		fetchedShipment := models.MTOShipment{}
		serviceItems := models.MTOServiceItems{}

		// Verify that required delivery date is not calculated when it does not need to be
		planner.AssertNumberOfCalls(suite.T(), "TransitDistance", 0)

		shipment, approverErr := approver.ApproveShipment(appCtx, shipmentForAutoApprove.ID, shipmentForAutoApproveEtag)

		suite.Contains(approverErr.Error(), "shipment uses external vendor, cannot be approved for GHC Prime")
		suite.Equal(uuid.UUID{}, shipment.ID)

		err := appCtx.DB().Find(&fetchedShipment, shipmentForAutoApprove.ID)
		suite.NoError(err)

		suite.Equal(models.MTOShipmentStatusSubmitted, fetchedShipment.Status)
		suite.Nil(shipment.ApprovedDate)
		suite.Nil(fetchedShipment.ApprovedDate)

		err = appCtx.DB().EagerPreload("ReService").Where("mto_shipment_id = ?", shipmentForAutoApprove.ID).All(&serviceItems)
		suite.NoError(err)

		suite.Equal(0, len(serviceItems))
	})

	suite.Run("Test that correct addresses are being used to calculate required delivery date", func() {
		subtestData := suite.createApproveShipmentSubtestData()
		appCtx := subtestData.appCtx
		move := subtestData.move
		approver := subtestData.shipmentApprover
		planner := subtestData.planner

		expectedReServiceCodes := []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
			models.ReServiceCodeDNPK,
		}

		for _, serviceCode := range expectedReServiceCodes {
			factory.FetchOrBuildReServiceByCode(appCtx.DB(), serviceCode)
		}

		// This is testing that the Required Delivery Date is calculated correctly.
		// In order for the Required Delivery Date to be calculated, the following conditions must be true:
		// 1. The shipment is moving to the APPROVED status
		// 2. The shipment must already have the following fields present:
		// MTOShipmentTypeHHG: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, DestinationAddress
		// MTOShipmentTypeHHGIntoNTSDom: ScheduledPickupDate, PrimeEstimatedWeight, PickupAddress, StorageFacility
		// MTOShipmentTypeHHGOutOfNTSDom: ScheduledPickupDate, NTSRecordedWeight, StorageFacility, DestinationAddress
		// 3. The shipment must not already have a Required Delivery Date
		// Note that MakeMTOShipment will automatically add a Required Delivery Date if the ScheduledPickupDate
		// is present, therefore we need to use MakeMTOShipmentMinimal and add the Pickup and Destination addresses
		estimatedWeight := unit.Pound(1400)

		destinationAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress4})
		pickupAddress := factory.BuildAddress(suite.DB(), nil, []factory.Trait{factory.GetTraitAddress3})
		storageFacility := factory.BuildStorageFacility(suite.DB(), nil, nil)

		hhgShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		ntsShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
					ScheduledPickupDate:  &testdatagen.DateInsidePeakRateCycle,
					PrimeEstimatedWeight: &estimatedWeight,
					Status:               models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    pickupAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		ntsrShipment := factory.BuildMTOShipmentMinimal(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:        models.MTOShipmentTypeHHGOutOfNTSDom,
					ScheduledPickupDate: &testdatagen.DateInsidePeakRateCycle,
					NTSRecordedWeight:   &estimatedWeight,
					Status:              models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
			{
				Model:    destinationAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)

		var TransitDistancePickupArg string
		var TransitDistanceDestinationArg string

		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("string"),
		).Return(500, nil).Run(func(args mock.Arguments) {
			TransitDistancePickupArg = args.Get(1).(string)
			TransitDistanceDestinationArg = args.Get(2).(string)
		})

		testCases := []struct {
			shipment            models.MTOShipment
			pickupLocation      *models.Address
			destinationLocation *models.Address
		}{
			{hhgShipment, hhgShipment.PickupAddress, hhgShipment.DestinationAddress},
			{ntsShipment, ntsShipment.PickupAddress, &ntsShipment.StorageFacility.Address},
			{ntsrShipment, &ntsrShipment.StorageFacility.Address, ntsrShipment.DestinationAddress},
		}

		for _, testCase := range testCases {
			shipmentEtag := etag.GenerateEtag(testCase.shipment.UpdatedAt)
			_, err := approver.ApproveShipment(appCtx, testCase.shipment.ID, shipmentEtag)
			suite.NoError(err)

			fetchedShipment := models.MTOShipment{}
			err = suite.DB().Find(&fetchedShipment, testCase.shipment.ID)
			suite.NoError(err)
			// We also should have a required delivery date
			suite.NotNil(fetchedShipment.RequiredDeliveryDate)
			// Check that TransitDistance was called with the correct addresses
			suite.Equal(testCase.pickupLocation.PostalCode, TransitDistancePickupArg)
			suite.Equal(testCase.destinationLocation.PostalCode, TransitDistanceDestinationArg)
		}
	})
}
