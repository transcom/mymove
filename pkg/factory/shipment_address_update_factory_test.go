package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildNonSITAddressUpdate() {

	suite.Run("Successful creation of default shipment address update", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, nil)

		// Validate results, default status is requested
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, addressUpdate.Status)
		suite.Equal(addressUpdate.ContractorRemarks, "Test Contractor Remark")
		suite.Nil(addressUpdate.OfficeRemarks)
		suite.NotNil(addressUpdate.NewAddress)
		suite.NotNil(addressUpdate.NewAddressID)
		suite.NotNil(addressUpdate.OriginalAddress)
		suite.NotNil(addressUpdate.OriginalAddressID)
		suite.NotNil(addressUpdate.ShipmentID)
	})

	suite.Run("Successful creation of shipment address update with requested status trait", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, []Trait{GetTraitNonSITAddressUpdateRequested})

		// Validate shipment address update status is requested
		suite.Equal(models.ShipmentAddressUpdateStatusRequested, addressUpdate.Status)
	})

	suite.Run("Successful creation of shipment address update with approved status trait", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, []Trait{GetTraitNonSITAddressUpdateApproved})

		// Validate shipment address update status is approved
		suite.Equal(models.ShipmentAddressUpdateStatusApproved, addressUpdate.Status)
	})

	suite.Run("Successful creation of shipment address update with rejected status trait", func() {

		addressUpdate := BuildShipmentAddressUpdate(suite.DB(), []Customization{}, []Trait{GetTraitNonSITAddressUpdateRejected})

		// Validate shipment address update status is rejected
		suite.Equal(models.ShipmentAddressUpdateStatusRejected, addressUpdate.Status)
	})

}
