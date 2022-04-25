package mtoshipment

import (
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
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
		availableToPrimeMove := testdatagen.MakeAvailableMove(suite.DB())
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
		externalShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			Move: availableToPrimeMove,
			MTOShipment: models.MTOShipment{
				ShipmentType:       models.MTOShipmentTypeHHGOutOfNTSDom,
				UsesExternalVendor: true,
				DestinationAddress: &address,
			},
		})
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
}
