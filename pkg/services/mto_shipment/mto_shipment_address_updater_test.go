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
					City:           "Some City",
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

	suite.Run("Test updating origin SITDeliveryMiles on shipment pickup address change", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		deliveryAddress := factory.BuildAddress(suite.DB(), nil, nil)
		pickUpAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "1234 Some Street",
					City:           "Some City",
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
		mtoServiceItems, _ := UpdateOriginSITServiceItemSITDeliveryMiles(planner, &externalShipment, &newAddress, &oldAddress, suite.AppContextForTest())
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

}
