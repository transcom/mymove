package mtoshipment

import (
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/address"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentAddress() {
	planner := &mocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	addressCreator := address.NewAddressCreator()
	addressUpdater := address.NewAddressUpdater()
	mtoShipmentAddressUpdater := NewMTOShipmentAddressUpdater(planner, addressCreator, addressUpdater)

	// TESTCASE SCENARIO
	// Under test: UpdateMTOShipmentAddress
	// Mocked:     None
	// Set up:     We request an address update on an external shipment with the mustBeAvailableToPrime flag = true
	//             And again with mustBeAvailableToPrime flag = false
	// Expected outcome:
	//             With mustBeAvailableToPrime = true, we should receive an error
	//             With mustBeAvailableToPrime = false, there should be no error
	suite.Run("Using external vendor shipment", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		address := factory.BuildAddress(suite.DB(), nil, nil)
		externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
					UsesExternalVendor: true,
				},
			},
			{
				Model:    address,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
		}, nil)
		eTag := etag.GenerateEtag(address.UpdatedAt)

		updatedAddress := address
		updatedAddress.StreetAddress1 = "123 Somewhere Ln"

		//  With mustBeAvailableToPrime = true, we should receive an error
		_, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, externalShipment.ID, eTag, true)
		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "looking for mtoShipment")
		}
		// With mustBeAvailableToPrime = false, there should be no error
		returnAddress, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, externalShipment.ID, eTag, false)
		suite.NoError(err)
		suite.Equal(updatedAddress.StreetAddress1, returnAddress.StreetAddress1)
	})

	suite.Run("Test updating service item destination address on shipment address change", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)
		pickUpAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1234 Some Street",
					City:           "COLUMBIA",
					State:          "SC",
					PostalCode:     "29229",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
					UsesExternalVendor: true,
					Status:             models.MTOShipmentStatusSubmitted,
				},
			},
			{
				Model:    deliveryAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    pickUpAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		threeMonthsAgo := time.Now().AddDate(0, -3, 0)
		twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
		sitServiceItems := factory.BuildOriginSITServiceItems(suite.DB(), availableToPrimeMove, externalShipment, &threeMonthsAgo, &twoMonthsAgo)
		sitServiceItems = append(sitServiceItems, factory.BuildDestSITServiceItems(suite.DB(), availableToPrimeMove, externalShipment, &twoMonthsAgo, nil)...)
		suite.Equal(8, len(sitServiceItems))

		eTag := etag.GenerateEtag(deliveryAddress.UpdatedAt)

		updatedAddress := deliveryAddress
		updatedAddress.StreetAddress1 = "123 Somewhere Ln"

		//  With mustBeAvailableToPrime = true, we should receive an error
		_, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, externalShipment.ID, eTag, true)
		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "looking for mtoShipment")
		}
		// With mustBeAvailableToPrime = false, there should be no error
		returnAddress, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, externalShipment.ID, eTag, false)
		suite.NoError(err)
		suite.Equal(updatedAddress.StreetAddress1, returnAddress.StreetAddress1)

		mtoServiceItems, _ := UpdateSITServiceItemDestinationAddressToMTOShipmentAddress(&sitServiceItems, &updatedAddress, suite.AppContextForTest())
		suite.Equal(4, len(*mtoServiceItems))
		for _, mtoServiceItem := range *mtoServiceItems {
			suite.Equal(externalShipment.DestinationAddressID, mtoServiceItem.SITDestinationFinalAddressID)
		}
	})

	suite.Run("Updating address validators", func() {
		testCases := map[string]struct {
			status    models.MTOShipmentStatus
			happyPath bool
		}{
			"Terminated shipment is a bad path": {
				models.MTOShipmentStatusTerminatedForCause,
				false,
			},
			"Submitted shipment is a happy path": {
				models.MTOShipmentStatusSubmitted,
				true,
			},
		}
		for _, tc := range testCases {
			availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
			address := factory.BuildAddress(suite.DB(), nil, nil)
			externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
				{
					Model:    availableToPrimeMove,
					LinkOnly: true,
				},
				{
					Model: models.MTOShipment{
						ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
						UsesExternalVendor: true,
						Status:             tc.status,
					},
				},
				{
					Model:    address,
					Type:     &factory.Addresses.DeliveryAddress,
					LinkOnly: true,
				},
			}, nil)
			eTag := etag.GenerateEtag(address.UpdatedAt)

			updatedAddress := address
			updatedAddress.StreetAddress1 = "123 Somewhere Ln"

			returnAddress, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, externalShipment.ID, eTag, false)
			// If an error occurred when one isn't expected
			if tc.happyPath {
				suite.FatalNoError(err, "Happy path scenario failed, the validators should have been satisfied and no error returned")
				suite.Equal(updatedAddress.StreetAddress1, returnAddress.StreetAddress1)
			}
			// If an error didn't occur when it is expected
			if !tc.happyPath {
				suite.Error(err, "No error occurred when the validator should have returned an error for the test case")
			}
		}
	})

	suite.Run("Test updating origin SITDeliveryMiles on shipment pickup address change", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)
		pickUpAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1234 Some Street",
					City:           "COLUMBIA",
					State:          "SC",
					PostalCode:     "29229",
					IsOconus:       models.BoolPointer(false),
				},
			},
		}, nil)

		externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTS,
					UsesExternalVendor: true,
					Status:             models.MTOShipmentStatusApproved,
				},
			},
			{
				Model:    deliveryAddress,
				Type:     &factory.Addresses.DeliveryAddress,
				LinkOnly: true,
			},
			{
				Model:    pickUpAddress,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		threeMonthsAgo := time.Now().AddDate(0, -3, 0)
		twoMonthsAgo := threeMonthsAgo.AddDate(0, 1, 0)
		sitServiceItems := factory.BuildOriginSITServiceItems(suite.DB(), availableToPrimeMove, externalShipment, &threeMonthsAgo, &twoMonthsAgo)
		sitServiceItems = append(sitServiceItems, factory.BuildDestSITServiceItems(suite.DB(), availableToPrimeMove, externalShipment, &twoMonthsAgo, nil)...)
		suite.Equal(8, len(sitServiceItems))

		eTag := etag.GenerateEtag(deliveryAddress.UpdatedAt)

		oldAddress := deliveryAddress
		oldAddress.PostalCode = "75116"
		newAddress := deliveryAddress
		newAddress.PostalCode = "67492"

		//  With mustBeAvailableToPrime = true, we should receive an error
		_, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &newAddress, externalShipment.ID, eTag, true)
		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "looking for mtoShipment")
		}

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(465, nil)

		addressCreator := address.NewAddressCreator()

		mtoServiceItems, _ := UpdateOriginSITServiceItemSITDeliveryMiles(planner, addressCreator, &externalShipment, &newAddress, &oldAddress, suite.AppContextForTest())
		suite.Equal(2, len(*mtoServiceItems))
		for _, mtoServiceItem := range *mtoServiceItems {
			if mtoServiceItem.ReService.Code == "DOSFSC" || mtoServiceItem.ReService.Code == "DOPSIT" {
				suite.Equal(*mtoServiceItem.SITDeliveryMiles, 465)
			}
		}
	})

	suite.Run("UB shipment without any OCONUS address should error", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		conusAddress := factory.BuildAddress(suite.DB(), nil, nil)

		// default factory is OCONUS dest and CONUS pickup
		ubShipment := factory.BuildUBShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
		}, nil)

		suite.True(*ubShipment.DestinationAddress.IsOconus)
		suite.False(*ubShipment.PickupAddress.IsOconus)

		updatedAddress := conusAddress
		updatedAddress.ID = *ubShipment.DestinationAddressID
		eTag := etag.GenerateEtag(ubShipment.DestinationAddress.UpdatedAt)

		_, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, ubShipment.ID, eTag, false)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Contains(err.Error(), "At least one address for a UB shipment must be OCONUS")
	})

	suite.Run("Successful - UpdateMTOShipmentAddress - Test updating international origin SITDeliveryMiles on shipment pickup address change", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		address := factory.BuildAddress(suite.DB(), nil, nil)
		actualAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "177 Q st",
					City:           "Solomons",
					State:          "MD",
					PostalCode:     "20688",
				},
			},
		}, nil)

		primeActualWeight := unit.Pound(1234)
		primeEstimatedWeight := unit.Pound(1234)

		externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					UsesExternalVendor:   true,
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: &primeActualWeight,
					PrimeActualWeight:    &primeEstimatedWeight,
				},
			},
			{
				Model:    address,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model:    externalShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOSFSC,
				},
			},
			{
				Model:    actualAddress,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    actualAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(address.UpdatedAt)

		newAddress := address
		newAddress.PostalCode = "67492"
		newAddress.City = "WOODBINE"

		var serviceItems []models.MTOServiceItem

		// verify pre-update mto service items for both origin FSC SIT have not been set
		err := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", externalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		// expecting only IOSFSC and IDSFSC created for tests
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.Nil(serviceItems[i].PricingEstimate)
			suite.True(serviceItems[i].SITDeliveryMiles == (*int)(nil))
		}

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			"67492",
			"20688",
		).Return(5, nil)

		mtoShipmentAddressIntlSITUpdater := NewMTOShipmentAddressUpdater(planner, addressCreator, addressUpdater)

		_, err = mtoShipmentAddressIntlSITUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &newAddress, externalShipment.ID, eTag, false)
		suite.Nil(err)

		err = suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", externalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC)
			suite.NotNil(serviceItems[i].PricingEstimate)
			suite.Equal(*serviceItems[i].SITDeliveryMiles, 5)
		}
	})

	suite.Run("Successful - UpdateMTOShipmentAddress - OCONUS Original Pickup - should not calculate pricing", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		address := factory.BuildAddress(suite.DB(), nil, nil)
		actualAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "177 Q st",
					City:           "FAIRBANKS",
					State:          "AK",
					PostalCode:     "99708",
					IsOconus:       models.BoolPointer(true), //OCONUS - prevent pricing
				},
			},
		}, nil)

		primeActualWeight := unit.Pound(1234)
		primeEstimatedWeight := unit.Pound(1234)

		externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType:         models.MTOShipmentTypeHHG,
					UsesExternalVendor:   true,
					Status:               models.MTOShipmentStatusApproved,
					PrimeEstimatedWeight: &primeActualWeight,
					PrimeActualWeight:    &primeEstimatedWeight,
					MarketCode:           models.MarketCodeInternational,
				},
			},
			{
				Model:    address,
				Type:     &factory.Addresses.PickupAddress,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model:    externalShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOSFSC,
				},
			},
			{
				Model:    actualAddress,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
				LinkOnly: true,
			},
			{
				Model:    actualAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
			{
				Model: models.MTOServiceItem{
					Status:          models.MTOServiceItemStatusApproved,
					PricingEstimate: nil,
				},
			},
		}, nil)

		eTag := etag.GenerateEtag(address.UpdatedAt)

		newAddress := address
		newAddress.PostalCode = "67492"
		newAddress.City = "WOODBINE"

		var serviceItems []models.MTOServiceItem

		// verify pre-update mto service items for both origin FSC SIT have not been set
		err := suite.AppContextForTest().DB().EagerPreload("ReService").Where("mto_shipment_id = ?", externalShipment.ID).Order("created_at asc").All(&serviceItems)
		suite.NoError(err)
		// expecting only IOSFSC and IDSFSC created for tests
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.Nil(serviceItems[i].PricingEstimate)
			suite.True(serviceItems[i].SITDeliveryMiles == (*int)(nil))
		}

		planner := &mocks.Planner{}
		planner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(5, nil).Times(0)

		mtoShipmentAddressIntlSITUpdater := NewMTOShipmentAddressUpdater(planner, addressCreator, addressUpdater)

		_, err = mtoShipmentAddressIntlSITUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &newAddress, externalShipment.ID, eTag, false)
		suite.Nil(err)

		err = suite.AppContextForTest().DB().Eager("SITOriginHHGOriginalAddress", "SITOriginHHGActualAddress",
			"ReService").Where("mto_shipment_id = ?", externalShipment.ID).Order("created_at asc").All(&serviceItems)

		suite.NoError(err)
		suite.Equal(1, len(serviceItems))
		for i := 0; i < len(serviceItems); i++ {
			suite.True(serviceItems[i].ReService.Code == models.ReServiceCodeIOSFSC)
			// verify mileage and pricing were not calcuated because of OCONUS origin pickup - SITOriginHHGOriginalAddress
			suite.Equal(*serviceItems[i].PricingEstimate, unit.Cents(0))
			suite.NotNil(serviceItems[i].SITDeliveryMiles)
			suite.NotNil(serviceItems[i].SITOriginHHGActualAddressID)
			// verify SITOriginHHGActualAddress was not changed
			suite.Equal(serviceItems[i].SITOriginHHGActualAddress.PostalCode, actualAddress.PostalCode)
			// verify SITOriginHHGOriginalAddress was not changed
			suite.Equal(serviceItems[i].SITOriginHHGOriginalAddress.PostalCode, actualAddress.PostalCode)
		}
	})
}
