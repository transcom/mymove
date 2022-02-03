package mtoshipment

import (
	"testing"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentAddress() {
	appCtx := suite.AppContextForTest()
	mtoShipmentAddressUpdater := NewMTOShipmentAddressUpdater()

	suite.T().Run("Using external vendor shipment", func(t *testing.T) {
		availableToPrimeMove := testdatagen.MakeAvailableMove(appCtx.DB())
		address := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{})
		externalShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
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

		_, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(appCtx, &updatedAddress, externalShipment.ID, eTag, true)
		if suite.Error(err) {
			suite.IsType(apperror.NotFoundError{}, err)
			suite.Contains(err.Error(), "looking for mtoShipment")
		}

		returnAddress, err := mtoShipmentAddressUpdater.UpdateMTOShipmentAddress(appCtx, &updatedAddress, externalShipment.ID, eTag, false)
		suite.NoError(err)
		suite.Equal(updatedAddress.StreetAddress1, returnAddress.StreetAddress1)
	})
}
