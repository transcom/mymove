package factory

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildMTOShipment() {
	defaultShipmentType := models.MTOShipmentTypeHHG
	defaultStatus := models.MTOShipmentStatusDraft
	defaultBaseStatus := models.MTOShipmentStatusSubmitted

	suite.Run("Successful creation of basic MTOShipment", func() {
		// Under test:      BuildBaseMTOShipment
		// Set up:          Create a basic mtoShipment
		// Expected outcome: Create a move, and mtoShipment

		// SETUP
		mtoShipment := BuildBaseMTOShipment(suite.DB(), nil, nil)

		suite.Equal(defaultShipmentType, mtoShipment.ShipmentType)
		suite.Equal(defaultBaseStatus, mtoShipment.Status)
	})

	suite.Run("Successful creation of stubbed basic MTOShipment", func() {
		// Under test:      BuildBaseMTOShipment
		// Set up:          Create a basic mtoShipment, but don't pass in a db
		// Expected outcome: Create a move, and mtoShipment but should not be in database

		// SETUP
		precount, err := suite.DB().Count(&models.MTOShipment{})
		suite.NoError(err)

		mtoShipment := BuildBaseMTOShipment(nil, nil, nil)

		suite.Equal(defaultShipmentType, mtoShipment.ShipmentType)
		suite.Equal(defaultBaseStatus, mtoShipment.Status)

		count, err := suite.DB().Count(&models.Move{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of custom MTO Shipment", func() {
		// Under test:      BuildBaseMTOShipment
		// Set up:          Create a custom mtoShipment and custom move
		// Expected outcome: Create a move, and mtoShipment

		// SETUP
		locator := "ABC123"
		ppmType := "FULL"

		move := models.Move{
			PPMType: &ppmType,
			Locator: locator,
		}

		mtoShipment := BuildBaseMTOShipment(suite.DB(), []Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeBoatHaulAway,
					Status:       models.MTOShipmentStatusApproved,
				},
			},
			{
				Model: move,
			},
		}, nil)

		suite.Equal(mtoShipment.ShipmentType, models.MTOShipmentTypeBoatHaulAway)
		suite.Equal(mtoShipment.Status, models.MTOShipmentStatusApproved)

		suite.Equal(mtoShipment.MoveTaskOrder.PPMType, &ppmType)
		suite.Equal(mtoShipment.MoveTaskOrder.Locator, locator)
	})

	suite.Run("Successful creation of custom MTO Shipment with UB shipment type", func() {
		// Under test:      BuildBaseMTOShipment
		// Set up:          Create a custom mtoShipment with UB shipment type and custom move
		// Expected outcome: Create a move, and mtoShipment

		// SETUP
		locator := "ABC123"

		move := models.Move{
			Locator: locator,
		}

		mtoShipment := BuildBaseMTOShipment(suite.DB(), []Customization{
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeUnaccompaniedBaggage,
					Status:       models.MTOShipmentStatusApproved,
				},
			},
			{
				Model: move,
			},
		}, nil)

		suite.Equal(mtoShipment.ShipmentType, models.MTOShipmentTypeUnaccompaniedBaggage)
		suite.Equal(mtoShipment.Status, models.MTOShipmentStatusApproved)
		suite.Equal(mtoShipment.MoveTaskOrder.Locator, locator)
	})

	suite.Run("Successful creation of default MTOShipment with other associated set relationships", func() {
		defaultCustomerRemarks := models.StringPointer("Please treat gently")
		// Under test:      BuildMTOShipment
		// Set up:          Create a default mtoShipment
		// Expected outcome: Create a move, pickupAddress, DeliveryAddress and mtoShipment

		// SETUP
		defaultPickupAddress := models.Address{
			StreetAddress1: "123 Any Street",
			StreetAddress2: models.StringPointer("P.O. Box 12345"),
			StreetAddress3: models.StringPointer("c/o Some Person"),
		}

		defaultDeliveryAddress := models.Address{
			StreetAddress1: "987 Any Avenue",
			StreetAddress2: models.StringPointer("P.O. Box 9876"),
			StreetAddress3: models.StringPointer("c/o Some Person"),
		}

		partialType := "PARTIAL"
		defaultMove := models.Move{
			PPMType: &partialType,
		}

		defaultRequestedPickupDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
		defaultScheduledPickupDate := time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)
		defaultActualPickupDate := time.Date(GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)
		defaultRequestedDeliveryDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
		defaultScheduledDeliveryDate := time.Date(GHCTestYear, time.March, 17, 0, 0, 0, 0, time.UTC)

		mtoShipment := BuildMTOShipment(suite.DB(), nil, nil)

		suite.Equal(defaultShipmentType, mtoShipment.ShipmentType)
		suite.Equal(defaultStatus, mtoShipment.Status)
		suite.Equal(defaultCustomerRemarks, mtoShipment.CustomerRemarks)

		// Check Pickup Address
		suite.Equal(defaultPickupAddress.StreetAddress1, mtoShipment.PickupAddress.StreetAddress1)
		suite.Equal(defaultPickupAddress.StreetAddress2, mtoShipment.PickupAddress.StreetAddress2)
		suite.Equal(defaultPickupAddress.StreetAddress3, mtoShipment.PickupAddress.StreetAddress3)

		// Check Delivery Address
		suite.Equal(defaultDeliveryAddress.StreetAddress1, mtoShipment.DestinationAddress.StreetAddress1)
		suite.Equal(defaultDeliveryAddress.StreetAddress2, mtoShipment.DestinationAddress.StreetAddress2)
		suite.Equal(defaultDeliveryAddress.StreetAddress3, mtoShipment.DestinationAddress.StreetAddress3)

		// Check move
		suite.Equal(defaultMove.PPMType, mtoShipment.MoveTaskOrder.PPMType)

		suite.Nil(mtoShipment.StorageFacility)

		// Check dates
		suite.NotNil(mtoShipment.RequestedPickupDate)
		suite.Equal(defaultRequestedPickupDate, *mtoShipment.RequestedPickupDate)
		suite.NotNil(mtoShipment.ScheduledPickupDate)
		suite.Equal(defaultScheduledPickupDate, *mtoShipment.ScheduledPickupDate)
		suite.NotNil(mtoShipment.ActualPickupDate)
		suite.Equal(defaultActualPickupDate, *mtoShipment.ActualPickupDate)
		suite.NotNil(mtoShipment.RequestedDeliveryDate)
		suite.Equal(defaultRequestedDeliveryDate, *mtoShipment.RequestedDeliveryDate)
		suite.NotNil(mtoShipment.ScheduledDeliveryDate)
		suite.Equal(defaultScheduledDeliveryDate, *mtoShipment.ScheduledDeliveryDate)
	})

	suite.Run("Successful creation of custom MTOShipment with pickup details and other associated set relationships", func() {
		// Under test:      BuildMTOShipment
		// Set up:          Create a custom mtoShipment
		// Expected outcome: Create a move, storageFacility, pickupAddress, DeliveryAddress and mtoShipment

		// SETUP
		var estimatedWeight = unit.Pound(1400)
		var actualWeight = unit.Pound(2000)

		customMTOShipment := models.MTOShipment{
			ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			ApprovedDate:         models.TimePointer(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
		}

		customMove := models.Move{
			ID:                 uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
			Status:             models.MoveStatusAPPROVALSREQUESTED,
			AvailableToPrimeAt: models.TimePointer(time.Now()),
			ApprovedAt:         models.TimePointer(time.Now()),
		}

		customPickupAddress := models.Address{
			StreetAddress1: "101 This is Awesome Street",
		}

		customSecondaryPickupAddress := models.Address{
			StreetAddress1: "201 Other Street",
		}

		customTertiaryPickupAddress := models.Address{
			StreetAddress1: "301 Other Street",
		}

		customDeliveryAddress := models.Address{
			StreetAddress1: "301 Another Good Street",
		}

		customSecondaryDeliveryAddress := models.Address{
			StreetAddress1: "401 Big MTO Street",
		}

		customTertiaryDeliveryAddress := models.Address{
			StreetAddress1: "301 Big MTO Street",
		}

		customStorageFacility := models.StorageFacility{
			Email: models.StringPointer("old@email.com"),
		}

		mtoShipment := BuildMTOShipment(suite.DB(), []Customization{
			{
				Model: customMTOShipment,
			},
			{
				Model: customMove,
			},
			{
				Model: customStorageFacility,
			},
			{
				Model: customPickupAddress,
				Type:  &Addresses.PickupAddress,
			},
			{
				Model: customDeliveryAddress,
				Type:  &Addresses.DeliveryAddress,
			},
			{
				Model: customSecondaryPickupAddress,
				Type:  &Addresses.SecondaryPickupAddress,
			},
			{
				Model: customSecondaryDeliveryAddress,
				Type:  &Addresses.SecondaryDeliveryAddress,
			},
			{
				Model: customTertiaryPickupAddress,
				Type:  &Addresses.TertiaryPickupAddress,
			},
			{
				Model: customTertiaryDeliveryAddress,
				Type:  &Addresses.TertiaryDeliveryAddress,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customMTOShipment.PrimeEstimatedWeight, mtoShipment.PrimeEstimatedWeight)
		suite.Equal(customMTOShipment.PrimeActualWeight, mtoShipment.PrimeActualWeight)
		suite.Equal(customMTOShipment.ShipmentType, mtoShipment.ShipmentType)
		suite.Equal(customMTOShipment.ApprovedDate, mtoShipment.ApprovedDate)
		suite.Equal(customMTOShipment.Status, mtoShipment.Status)

		// Check Pickup Address
		suite.Equal(customPickupAddress.StreetAddress1, mtoShipment.PickupAddress.StreetAddress1)
		// Check Secondary PickupAddress
		suite.Equal(customSecondaryPickupAddress.StreetAddress1, mtoShipment.SecondaryPickupAddress.StreetAddress1)
		suite.Equal(models.BoolPointer(true), mtoShipment.HasSecondaryPickupAddress)
		// Check Tertiary PickupAddress
		suite.Equal(customTertiaryPickupAddress.StreetAddress1, mtoShipment.TertiaryPickupAddress.StreetAddress1)
		suite.Equal(models.BoolPointer(true), mtoShipment.HasTertiaryPickupAddress)

		// Check Storage Facility
		suite.Equal(customStorageFacility.Email, mtoShipment.StorageFacility.Email)

		// Check move
		suite.Equal(customMove.Status, mtoShipment.MoveTaskOrder.Status)
		suite.Equal(customMove.AvailableToPrimeAt, mtoShipment.MoveTaskOrder.AvailableToPrimeAt)
		suite.Equal(customMove.ApprovedAt, mtoShipment.MoveTaskOrder.ApprovedAt)
	})

	suite.Run("Successful creation of custom MTOShipment with delivery details and other associated set relationships", func() {
		// Under test:      BuildMTOShipment
		// Set up:          Create a custom mtoShipment
		// Expected outcome: Create a move, storageFacility, pickupAddress, DeliveryAddress and mtoShipment

		// SETUP
		var estimatedWeight = unit.Pound(1400)
		var actualWeight = unit.Pound(2000)

		customMTOShipment := models.MTOShipment{
			ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         models.TimePointer(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
		}

		customMove := models.Move{
			ID:                 uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
			Status:             models.MoveStatusAPPROVALSREQUESTED,
			AvailableToPrimeAt: models.TimePointer(time.Now()),
			ApprovedAt:         models.TimePointer(time.Now()),
		}

		customDeliveryAddress := models.Address{
			StreetAddress1: "301 Another Good Street",
		}

		customSecondaryDeliveryAddress := models.Address{
			StreetAddress1: "401 Big MTO Street",
		}

		customTertiaryDeliveryAddress := models.Address{
			StreetAddress1: "401 Big MTO Street",
		}

		customStorageFacility := models.StorageFacility{
			Email: models.StringPointer("old@email.com"),
		}

		mtoShipment := BuildMTOShipment(suite.DB(), []Customization{
			{
				Model: customMTOShipment,
			},
			{
				Model: customMove,
			},
			{
				Model: customStorageFacility,
			},
			{
				Model: customDeliveryAddress,
				Type:  &Addresses.DeliveryAddress,
			},
			{
				Model: customSecondaryDeliveryAddress,
				Type:  &Addresses.SecondaryDeliveryAddress,
			},
			{
				Model: customSecondaryDeliveryAddress,
				Type:  &Addresses.TertiaryDeliveryAddress,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customMTOShipment.PrimeEstimatedWeight, mtoShipment.PrimeEstimatedWeight)
		suite.Equal(customMTOShipment.PrimeActualWeight, mtoShipment.PrimeActualWeight)
		suite.Equal(customMTOShipment.ShipmentType, mtoShipment.ShipmentType)
		suite.Equal(customMTOShipment.ApprovedDate, mtoShipment.ApprovedDate)
		suite.Equal(customMTOShipment.Status, mtoShipment.Status)

		// Check Delivery Address
		suite.Equal(customDeliveryAddress.StreetAddress1, mtoShipment.DestinationAddress.StreetAddress1)
		// Check Secondary DeliveryAddress
		suite.Equal(customSecondaryDeliveryAddress.StreetAddress1, mtoShipment.SecondaryDeliveryAddress.StreetAddress1)
		suite.Equal(models.BoolPointer(true), mtoShipment.HasSecondaryDeliveryAddress)

		// Check Tertiary DeliveryAddress
		suite.Equal(customTertiaryDeliveryAddress.StreetAddress1, mtoShipment.TertiaryDeliveryAddress.StreetAddress1)
		suite.Equal(models.BoolPointer(true), mtoShipment.HasTertiaryDeliveryAddress)

		// Check Storage Facility
		suite.Equal(customStorageFacility.Email, mtoShipment.StorageFacility.Email)

		// Check move
		suite.Equal(customMove.Status, mtoShipment.MoveTaskOrder.Status)
		suite.Equal(customMove.AvailableToPrimeAt, mtoShipment.MoveTaskOrder.AvailableToPrimeAt)
		suite.Equal(customMove.ApprovedAt, mtoShipment.MoveTaskOrder.ApprovedAt)
	})

	suite.Run("Successful return of linkOnly mtoShipment", func() {
		// Under test:      BuildMTOShipment
		// Set up:          Create a mtoShipment and pass in a linkOnly flag
		// Expected outcome: Create a mtoShipment

		// SETUP
		precount, err := suite.DB().Count(&models.MTOShipment{})
		suite.NoError(err)

		customMTOShipment := models.MTOShipment{
			ID:     uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
			Status: models.MTOShipmentStatusApproved,
		}

		mtoShipment := BuildMTOShipment(suite.DB(), []Customization{
			{
				Model:    customMTOShipment,
				LinkOnly: true,
			},
		}, nil)

		count, err := suite.DB().Count(&models.MTOShipment{})
		suite.Equal(precount, count)
		suite.NoError(err)
		suite.Equal(customMTOShipment.Status, mtoShipment.Status)
	})

	suite.Run("Successful creation of stubbed MTOShipment", func() {
		move := BuildMove(suite.DB(), nil, nil)
		mtoShipment := BuildMTOShipment(nil, []Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		// the stubbed shipment still needs non nil ids for the pickup
		// and delivery addresses
		suite.NotNil(mtoShipment.PickupAddressID)
		suite.NotNil(mtoShipment.DestinationAddressID)
	})

	suite.Run("Successful creation of NTSShipment", func() {
		ntsShipment := BuildNTSShipment(suite.DB(), nil, nil)

		suite.Equal(models.MTOShipmentTypeHHGIntoNTSDom, ntsShipment.ShipmentType)
		suite.False(ntsShipment.MoveTaskOrderID.IsNil())
		suite.False(ntsShipment.MoveTaskOrder.ID.IsNil())

		suite.NotNil(ntsShipment.PickupAddressID)
		suite.NotNil(ntsShipment.PickupAddress)
		suite.False(ntsShipment.PickupAddressID.IsNil())
		suite.False(ntsShipment.PickupAddress.ID.IsNil())

		suite.NotNil(ntsShipment.SecondaryPickupAddressID)
		suite.NotNil(ntsShipment.SecondaryPickupAddress)
		suite.False(ntsShipment.SecondaryPickupAddressID.IsNil())
		suite.False(ntsShipment.SecondaryPickupAddress.ID.IsNil())
		suite.NotNil(ntsShipment.HasSecondaryPickupAddress)
		suite.True(*ntsShipment.HasSecondaryPickupAddress)

		suite.NotNil(ntsShipment.TertiaryPickupAddressID)
		suite.NotNil(ntsShipment.TertiaryPickupAddress)
		suite.False(ntsShipment.TertiaryPickupAddressID.IsNil())
		suite.False(ntsShipment.TertiaryPickupAddress.ID.IsNil())
		suite.NotNil(ntsShipment.HasTertiaryPickupAddress)
		suite.True(*ntsShipment.HasTertiaryPickupAddress)

		suite.NotNil(ntsShipment.CustomerRemarks)
		suite.Equal("Please treat gently", *ntsShipment.CustomerRemarks)
		suite.Equal(models.MTOShipmentStatusDraft, ntsShipment.Status)
		suite.NotNil(ntsShipment.PrimeActualWeight)
		suite.Nil(ntsShipment.StorageFacility)
		suite.NotNil(ntsShipment.ScheduledPickupDate)
		suite.Nil(ntsShipment.RequestedDeliveryDate)
		suite.Nil(ntsShipment.ActualDeliveryDate)
		suite.Nil(ntsShipment.ScheduledDeliveryDate)
	})

	suite.Run("Successful creation of NTSShipment with storage facility", func() {
		storageFacility := BuildStorageFacility(suite.DB(), nil, nil)
		ntsShipment := BuildNTSShipment(suite.DB(), []Customization{
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		suite.NotNil(ntsShipment.StorageFacilityID)
		suite.Equal(storageFacility.ID, *ntsShipment.StorageFacilityID)
	})

	suite.Run("Successful creation of NTSRShipment", func() {
		ntsrShipment := BuildNTSRShipment(suite.DB(), nil, nil)

		suite.Equal(models.MTOShipmentTypeHHGOutOfNTSDom, ntsrShipment.ShipmentType)
		suite.False(ntsrShipment.MoveTaskOrderID.IsNil())
		suite.False(ntsrShipment.MoveTaskOrder.ID.IsNil())
		suite.NotNil(ntsrShipment.DestinationAddressID)
		suite.NotNil(ntsrShipment.DestinationAddress)
		suite.False(ntsrShipment.DestinationAddressID.IsNil())
		suite.False(ntsrShipment.DestinationAddress.ID.IsNil())

		suite.NotNil(ntsrShipment.SecondaryDeliveryAddressID)
		suite.NotNil(ntsrShipment.SecondaryDeliveryAddress)
		suite.False(ntsrShipment.SecondaryDeliveryAddressID.IsNil())
		suite.False(ntsrShipment.SecondaryDeliveryAddress.ID.IsNil())
		suite.NotNil(ntsrShipment.HasSecondaryDeliveryAddress)
		suite.True(*ntsrShipment.HasSecondaryDeliveryAddress)

		suite.NotNil(ntsrShipment.TertiaryDeliveryAddressID)
		suite.NotNil(ntsrShipment.TertiaryDeliveryAddress)
		suite.False(ntsrShipment.TertiaryDeliveryAddressID.IsNil())
		suite.False(ntsrShipment.TertiaryDeliveryAddress.ID.IsNil())
		suite.NotNil(ntsrShipment.HasTertiaryDeliveryAddress)
		suite.True(*ntsrShipment.HasTertiaryDeliveryAddress)

		suite.NotNil(ntsrShipment.CustomerRemarks)
		suite.Equal("Please treat gently", *ntsrShipment.CustomerRemarks)
		suite.Equal(models.MTOShipmentStatusDraft, ntsrShipment.Status)
		suite.Nil(ntsrShipment.PrimeActualWeight)
		suite.Nil(ntsrShipment.StorageFacility)
		suite.Nil(ntsrShipment.ScheduledPickupDate)
		suite.NotNil(ntsrShipment.RequestedDeliveryDate)
		suite.Nil(ntsrShipment.ActualDeliveryDate)
	})

	suite.Run("Successful creation of NTSRShipment with storage facility", func() {
		storageFacility := BuildStorageFacility(suite.DB(), nil, nil)
		ntsrShipment := BuildNTSRShipment(suite.DB(), []Customization{
			{
				Model:    storageFacility,
				LinkOnly: true,
			},
		}, nil)
		suite.NotNil(ntsrShipment.StorageFacilityID)
		suite.Equal(storageFacility.ID, *ntsrShipment.StorageFacilityID)
	})
	suite.Run("Successful creation of NTSRShipment with pickup address", func() {
		address := BuildAddress(suite.DB(), nil, nil)
		ntsrShipment := BuildNTSRShipment(suite.DB(), []Customization{
			{
				Model:    address,
				LinkOnly: true,
				Type:     &Addresses.PickupAddress,
			},
		}, nil)
		suite.NotNil(ntsrShipment.PickupAddress)
		suite.NotNil(ntsrShipment.PickupAddressID)
		suite.Equal(address.ID, *ntsrShipment.PickupAddressID)
	})

}
