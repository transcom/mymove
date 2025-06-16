package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildShipmentAddressUpdate() {

	suite.Run("Successful creation of default shipment address update", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, nil)

		// Validate results, default status is requested
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, addressUpdate.Status)
		suite.Equal(addressUpdate.ContractorRemarks, "Customer reached out to me this week & let me know they want to move closer to a sick dependent who needs care.")
		suite.Nil(addressUpdate.OfficeRemarks)
		suite.NotNil(addressUpdate.NewAddress)
		suite.NotNil(addressUpdate.NewAddressID)
		suite.NotNil(addressUpdate.OriginalAddress)
		suite.NotNil(addressUpdate.OriginalAddressID)
		suite.NotNil(addressUpdate.ShipmentID)
	})

	suite.Run("Successful creation of shipment address update with requested status trait", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, []Trait{GetTraitShipmentAddressUpdateRequested})

		suite.Equal(models.ShipmentAddressUpdateStatusRequested, addressUpdate.Status)
		suite.Equal(models.MTOShipmentStatusApprovalsRequested, addressUpdate.Shipment.Status)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, addressUpdate.Shipment.MoveTaskOrder.Status)
		suite.Equal("CRQST1", addressUpdate.Shipment.MoveTaskOrder.Locator)
		suite.NotNil(addressUpdate.Shipment.MoveTaskOrder.AvailableToPrimeAt)
		suite.NotNil(addressUpdate.Shipment.MoveTaskOrder.ApprovedAt)

	})

	suite.Run("Successful creation of shipment address update with approved status trait", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, []Trait{GetTraitShipmentAddressUpdateApproved})

		suite.Equal(models.ShipmentAddressUpdateStatusApproved, addressUpdate.Status)
		suite.Equal(models.MTOShipmentStatusApproved, addressUpdate.Shipment.Status)
		suite.Equal(models.MoveStatusAPPROVED, addressUpdate.Shipment.MoveTaskOrder.Status)
		suite.Equal("CRQST2", addressUpdate.Shipment.MoveTaskOrder.Locator)
		suite.NotNil(addressUpdate.Shipment.MoveTaskOrder.AvailableToPrimeAt)
		suite.NotNil(addressUpdate.Shipment.MoveTaskOrder.ApprovedAt)
	})

	suite.Run("Successful creation of shipment address update with rejected status trait", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, []Trait{GetTraitShipmentAddressUpdateRejected})

		suite.Equal(models.ShipmentAddressUpdateStatusRejected, addressUpdate.Status)
		suite.Equal(models.MTOShipmentStatusApproved, addressUpdate.Shipment.Status)
		suite.Equal(models.MoveStatusAPPROVED, addressUpdate.Shipment.MoveTaskOrder.Status)
		suite.Equal("CRQST3", addressUpdate.Shipment.MoveTaskOrder.Locator)
		suite.NotNil(addressUpdate.Shipment.MoveTaskOrder.AvailableToPrimeAt)
		suite.NotNil(addressUpdate.Shipment.MoveTaskOrder.ApprovedAt)
	})

}
