package mtoshipment

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentAddress() {
	mtoShipmentAddressUpdater := NewMTOShipmentAddressUpdater()

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
					ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
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
	suite.Run("Destination address update for HHG should fail", func() {
		availableToPrimeMove := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		address := factory.BuildAddress(suite.DB(), nil, nil)
		externalShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    availableToPrimeMove,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ShipmentType: models.MTOShipmentTypeHHG,
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
			suite.IsType(apperror.ConflictError{}, err)
			suite.Contains(err.Error(), "This endpoint cannot be used to update HHG shipment destination addresses")
		}
		// With mustBeAvailableToPrime = false, we should get the same error
		_, err = mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(suite.AppContextForTest(), &updatedAddress, externalShipment.ID, eTag, false)
		if suite.Error(err) {
			suite.IsType(apperror.ConflictError{}, err)
			suite.Contains(err.Error(), "This endpoint cannot be used to update HHG shipment destination addresses")
		}
	})
}
