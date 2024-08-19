package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *FactorySuite) TestBuildSITAddressUpdate() {
	suite.Run("Successful creation of default SITAddressUpdate", func() {
		// Under test:      BuildSITAddressUpdate
		// Mocked:          None
		// Set up:          Create an SITAddressUpdate with no customizations or traits
		// Expected outcome:SITAddressUpdate should be created with default values

		// SETUP
		defaultOldAddress := BuildAddress(suite.DB(), nil, nil)
		defaultNewAddress := BuildAddress(suite.DB(), nil, []Trait{GetTraitAddress2})
		// Create a default SITAddressUpdate to compare values
		defaultSIT := models.SITAddressUpdate{
			Distance: 40,
			Status:   models.SITAddressUpdateStatusRequested,
		}

		// FUNCTION UNDER TEST
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.Equal(defaultSIT.Distance, sitAddressUpdate.Distance)
		suite.Equal(defaultSIT.Status, sitAddressUpdate.Status)
		suite.Nil(sitAddressUpdate.OfficeRemarks)
		suite.NotNil(sitAddressUpdate.MTOServiceItem)
		suite.False(sitAddressUpdate.MTOServiceItem.ID.IsNil())
		suite.False(sitAddressUpdate.MTOServiceItemID.IsNil())
		suite.Equal(sitAddressUpdate.OldAddressID, *sitAddressUpdate.MTOServiceItem.SITDestinationFinalAddressID)
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.False(sitAddressUpdate.OldAddress.ID.IsNil())
		suite.False(sitAddressUpdate.OldAddressID.IsNil())
		suite.Equal(defaultOldAddress.PostalCode, sitAddressUpdate.OldAddress.PostalCode)
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.False(sitAddressUpdate.NewAddress.ID.IsNil())
		suite.False(sitAddressUpdate.NewAddressID.IsNil())
		suite.Equal(defaultNewAddress.PostalCode, sitAddressUpdate.NewAddress.PostalCode)
	})

	suite.Run("Successful creation of customized SITAddressUpdate", func() {
		// Under test:      BuildSITAddressUpdate
		// Mocked:          None
		// Set up:          Create SITAddressUpdate with customization
		// Expected outcome:SITAddressUpdate should be created with customized values

		// SETUP
		customUpdate := models.SITAddressUpdate{
			ID:                uuid.Must(uuid.NewV4()),
			ContractorRemarks: models.StringPointer("custom contractor remarks"),
			OfficeRemarks:     models.StringPointer("office remarks"),
			Distance:          40,
			Status:            models.SITAddressUpdateStatusRejected,
		}

		customServiceItem := models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusRejected,
		}

		customOldAddress := models.Address{
			ID:         uuid.Must(uuid.NewV4()),
			PostalCode: "77083",
		}

		customNewAddress := models.Address{
			ID:         uuid.Must(uuid.NewV4()),
			PostalCode: "90210",
		}

		// FUNCTION UNDER TEST
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), []Customization{
			{Model: customUpdate},
			{Model: customServiceItem},
			{
				Model: customOldAddress,
				Type:  &Addresses.SITAddressUpdateOldAddress,
			},
			{
				Model: customNewAddress,
				Type:  &Addresses.SITAddressUpdateNewAddress,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(customUpdate.ID, sitAddressUpdate.ID)
		suite.Equal(*customUpdate.ContractorRemarks, *sitAddressUpdate.ContractorRemarks)
		suite.Equal(*customUpdate.OfficeRemarks, *sitAddressUpdate.OfficeRemarks)
		suite.Equal(customUpdate.Distance, sitAddressUpdate.Distance)
		suite.Equal(customUpdate.Status, sitAddressUpdate.Status)
		suite.NotNil(sitAddressUpdate.MTOServiceItem)
		suite.False(sitAddressUpdate.MTOServiceItem.ID.IsNil())
		suite.False(sitAddressUpdate.MTOServiceItemID.IsNil())
		suite.Equal(customServiceItem.ID, sitAddressUpdate.MTOServiceItem.ID)
		suite.Equal(customServiceItem.Status, sitAddressUpdate.MTOServiceItem.Status)
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.False(sitAddressUpdate.OldAddress.ID.IsNil())
		suite.False(sitAddressUpdate.OldAddressID.IsNil())
		suite.Equal(customOldAddress.ID, sitAddressUpdate.OldAddress.ID)
		suite.Equal(customOldAddress.PostalCode, sitAddressUpdate.OldAddress.PostalCode)
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.False(sitAddressUpdate.NewAddress.ID.IsNil())
		suite.False(sitAddressUpdate.NewAddressID.IsNil())
		suite.Equal(customNewAddress.ID, sitAddressUpdate.NewAddress.ID)
		suite.Equal(customNewAddress.PostalCode, sitAddressUpdate.NewAddress.PostalCode)
	})

	suite.Run("Successful creation of customized SITAddressUpdate using GetTraitSITAddressUpdateWithMoveSetUp", func() {
		// Under test:      BuildSITAddressUpdate with GetTraitSITAddressUpdateWithMoveSetUp
		// Mocked:          None
		// Set up:          Create SITAddressUpdate with customization from trait
		// Expected outcome:SITAddressUpdate should be created with customized values

		// FUNCTION UNDER TEST
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), nil, []Trait{GetTraitSITAddressUpdateWithMoveSetUp})

		// VALIDATE RESULTS
		originalPostalCode := "90210"
		suite.Equal(originalPostalCode, sitAddressUpdate.OldAddress.PostalCode)
		suite.Equal("92114", sitAddressUpdate.NewAddress.PostalCode)
		suite.Equal(140, sitAddressUpdate.Distance)
		suite.Equal(models.SITAddressUpdateStatusRequested, sitAddressUpdate.Status)

		dependentsAuthorized := sitAddressUpdate.MTOServiceItem.MoveTaskOrder.Orders.Entitlement.DependentsAuthorized
		suite.Equal(true, *dependentsAuthorized)

		entitlement := sitAddressUpdate.MTOServiceItem.MoveTaskOrder.Orders.Entitlement
		sitDaysAllowance := 200
		suite.Equal(sitDaysAllowance, *entitlement.StorageInTransit)

		suite.Equal(models.MoveStatusAPPROVED, sitAddressUpdate.MTOServiceItem.MoveTaskOrder.Status)
		suite.NotNil(sitAddressUpdate.MTOServiceItem.MoveTaskOrder.AvailableToPrimeAt)
		suite.NotNil(sitAddressUpdate.MTOServiceItem.MoveTaskOrder.ApprovedAt)

		shipment := sitAddressUpdate.MTOServiceItem.MTOShipment
		suite.Equal(unit.Pound(1400), *shipment.PrimeEstimatedWeight)
		suite.Equal(unit.Pound(2000), *shipment.PrimeActualWeight)
		suite.Equal(models.MTOShipmentTypeHHG, shipment.ShipmentType)
		suite.Equal(models.MTOShipmentStatusApproved, shipment.Status)
		suite.NotNil(shipment.RequestedPickupDate)
		suite.NotNil(shipment.RequestedDeliveryDate)
		suite.Equal(sitDaysAllowance, *shipment.SITDaysAllowance)

		reserviceCode := models.ReServiceCodeDDDSIT

		suite.Equal(reserviceCode, sitAddressUpdate.MTOServiceItem.ReService.Code)
		suite.Equal(models.MTOServiceItemStatusApproved, sitAddressUpdate.MTOServiceItem.Status)
		suite.NotNil(sitAddressUpdate.MTOServiceItem.SITEntryDate)
		suite.Equal(originalPostalCode, *sitAddressUpdate.MTOServiceItem.SITPostalCode)
		suite.Equal("peak season all trucks in use", *sitAddressUpdate.MTOServiceItem.Reason)
	})

	suite.Run("Successful return of linkOnly SITAddressUpdate", func() {
		// Under test:       BuildSITAddressUpdate
		// Set up:           Pass in a linkOnly SITAddressUpdate
		// Expected outcome: No new SITAddressUpdate should be created.

		// Check num SITAddressUpdates
		precount, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)

		id := uuid.Must(uuid.NewV4())
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), []Customization{
			{
				Model: models.SITAddressUpdate{
					ID: id,
				},
				LinkOnly: true,
			},
		}, nil)

		// VALIDATE RESULTS
		suite.Equal(id, sitAddressUpdate.ID)

		// Count how many notification are in the DB, no new
		// SITAddressUpdate should have been created
		count, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of stubbed SITAddressUpdate", func() {
		// Under test:      BuildSITAddressUpdate
		// Set up:          Create a stubbed SITAddressUpdate
		// Expected outcome:No new SITAddressUpdate should be created

		// Check num SITAddressUpdates
		precount, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)

		sitAddressUpdate := BuildSITAddressUpdate(nil, nil, nil)

		// VALIDATE RESULTS
		suite.NotNil(sitAddressUpdate.MTOServiceItem)
		suite.True(sitAddressUpdate.MTOServiceItem.ID.IsNil())
		suite.True(sitAddressUpdate.MTOServiceItemID.IsNil())
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.True(sitAddressUpdate.OldAddress.ID.IsNil())
		suite.True(sitAddressUpdate.OldAddressID.IsNil())
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.True(sitAddressUpdate.NewAddress.ID.IsNil())
		suite.True(sitAddressUpdate.NewAddressID.IsNil())

		// Count how many notification are in the DB, no new
		// SITAddressUpdate should have been created
		count, err := suite.DB().Count(&models.SITAddressUpdate{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("Successful creation of SITAddressUpdate with trait over 50 Miles", func() {
		// Under test:      BuildSITAddressUpdate
		// Mocked:          None
		// Set up:          Create a SITAddressUpdate but pass in a trait that sets
		//                  old & new address that is over 50 miles
		// Expected outcome:SITAddressUpdate should have the old address, new address,
		//    distance and status filled out

		fiftyMiles := 50
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), nil, []Trait{
			GetTraitSITAddressUpdateOver50Miles,
		})
		suite.NotNil(sitAddressUpdate.NewAddressID)
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.NotNil(sitAddressUpdate.OldAddressID)
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.Greater(sitAddressUpdate.Distance, fiftyMiles)
		suite.Equal(models.SITAddressUpdateStatusRequested, sitAddressUpdate.Status)
	})

	suite.Run("Successful creation of SITAddressUpdate with trait under 50 Miles", func() {
		// Under test:      BuildSITAddressUpdate
		// Mocked:          None
		// Set up:          Create a SITAddressUpdate but pass in a trait that sets
		//                  old & new address that is under 50 miles
		// Expected outcome:SITAddressUpdate should have the old address, new address,
		//    distance and status filled out

		fiftyMiles := 50
		sitAddressUpdate := BuildSITAddressUpdate(suite.DB(), nil, []Trait{
			GetTraitSITAddressUpdateUnder50Miles,
		})
		suite.NotNil(sitAddressUpdate.NewAddressID)
		suite.NotNil(sitAddressUpdate.NewAddress)
		suite.NotNil(sitAddressUpdate.OldAddressID)
		suite.NotNil(sitAddressUpdate.OldAddress)
		suite.LessOrEqual(sitAddressUpdate.Distance, fiftyMiles)
		suite.Equal(models.SITAddressUpdateStatusApproved, sitAddressUpdate.Status)
	})
}
