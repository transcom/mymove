package mtoshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestListMTOShipments() {
	appCtx := suite.AppContextForTest()
	mtoShipmentFetcher := NewMTOShipmentFetcher()

	suite.T().Run("Returns not found error when move id doesn't exist", func(t *testing.T) {
		moveID := uuid.Must(uuid.NewV4())
		expectedError := apperror.NewNotFoundError(moveID, "move not found")

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(appCtx, moveID)

		suite.Equalf(err, expectedError, "Expected not found error for non-existent move id")
		suite.Nil(mtoShipments, "Expected shipment slice to be nil")
	})

	suite.T().Run("Returns an empty shipment list when no shipments exist", func(t *testing.T) {
		move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{})

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(appCtx, move.ID)

		suite.NoError(err, "Expected no error for a move without shipments")
		suite.Len(mtoShipments, 0, "Expected a zero length shipment list")
	})

	suite.T().Run("Returns external vendor shipments last", func(t *testing.T) {
		db := appCtx.DB()
		move := testdatagen.MakeMove(db, testdatagen.Assertions{})

		externalVendorShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			Move: move,
			MTOShipment: models.MTOShipment{
				UsesExternalVendor: true,
			},
		})
		firstShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			Move: move,
		})
		secondShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			Move: move,
		})

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(appCtx, move.ID)

		suite.NoError(err, "Expected no error for a move with 3 shipments")
		suite.Len(mtoShipments, 3, "Expected a shipment list of length 3")

		suite.Equal(firstShipment.ID.String(), mtoShipments[0].ID.String())
		suite.Equal(secondShipment.ID.String(), mtoShipments[1].ID.String())
		suite.Equal(externalVendorShipment.ID.String(), mtoShipments[2].ID.String())

	})

	suite.T().Run("Returns multiple shipments for move ordered by created date", func(t *testing.T) {
		db := appCtx.DB()
		move := testdatagen.MakeMove(db, testdatagen.Assertions{})

		firstShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			Move: move,
		})
		secondShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			Move: move,
		})

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(appCtx, move.ID)

		suite.NoError(err, "Expected no error for a move with two shipments")
		suite.Len(mtoShipments, 2, "Expected a shipment list of length 2")

		suite.Equal(firstShipment.ID.String(), mtoShipments[0].ID.String())
		suite.Equal(secondShipment.ID.String(), mtoShipments[1].ID.String())

	})

	suite.T().Run("Loads all shipment associations", func(t *testing.T) {
		db := appCtx.DB()
		move := testdatagen.MakeMove(db, testdatagen.Assertions{})

		storageFacility := testdatagen.MakeStorageFacility(db, testdatagen.Assertions{})

		secondaryPickupAddress := testdatagen.MakeDefaultAddress(db)
		secondaryDeliveryAddress := testdatagen.MakeAddress2(db, testdatagen.Assertions{})

		shipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				StorageFacility:          &storageFacility,
				SecondaryPickupAddress:   &secondaryPickupAddress,
				SecondaryDeliveryAddress: &secondaryDeliveryAddress,
			},
			Move: move,
		})

		serviceItem := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeDCRT,
			},
			MTOShipment: shipment,
			Move:        move,
		})

		agents := testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
			MTOShipment: shipment,
		})

		SITExtension := testdatagen.MakeSITExtension(db, testdatagen.Assertions{
			MTOShipment: shipment,
		})

		reweigh := testdatagen.MakeReweigh(db, testdatagen.Assertions{
			MTOShipment: shipment,
		})

		mtoShipments, err := mtoShipmentFetcher.ListMTOShipments(appCtx, move.ID)

		suite.NoError(err, "Expected no error for a move with shipment associations")
		suite.Len(mtoShipments, 1, "Expected a single shipment with associations")

		actualShipment := mtoShipments[0]

		suite.Equal(serviceItem.ReService.Code, actualShipment.MTOServiceItems[0].ReService.Code)
		suite.Equal(agents.ID.String(), actualShipment.MTOAgents[0].ID.String())
		suite.Equal(shipment.PickupAddress.ID.String(), actualShipment.PickupAddress.ID.String())
		suite.Equal(secondaryPickupAddress.ID.String(), actualShipment.SecondaryPickupAddress.ID.String())
		suite.Equal(shipment.DestinationAddress.ID.String(), actualShipment.DestinationAddress.ID.String())
		suite.Equal(secondaryDeliveryAddress.ID.String(), actualShipment.SecondaryDeliveryAddress.ID.String())
		suite.Len(actualShipment.MTOServiceItems[0].Dimensions, 2)
		suite.Equal(SITExtension.ID.String(), actualShipment.SITExtensions[0].ID.String())
		suite.Equal(storageFacility.Address.ID.String(), actualShipment.StorageFacility.Address.ID.String())
		suite.Equal(reweigh.ID.String(), actualShipment.Reweigh.ID.String())
	})
}
