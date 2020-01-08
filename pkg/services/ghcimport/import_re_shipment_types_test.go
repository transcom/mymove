package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREShipmentTypes() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importREShipmentTypes(suite.DB())
		suite.NoError(err)
		suite.helperVerifyShipmentTypes()
		suite.helperCheckShipmentTypes()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREShipmentTypes(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "re_shipment_types_code_key")
		}

		// Check to see if anything else changed
		suite.helperVerifyShipmentTypes()
		suite.helperCheckShipmentTypes()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyShipmentTypes() {
	count, err := suite.DB().Count(&models.ReShipmentTypes{})
	suite.NoError(err)
	suite.Equal(7, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckShipmentTypes() {
	var shipmentTypes models.ReShipmentType
	err := suite.DB().
		Where("code = 'DMHF'").
		First(&shipmentTypes)
	suite.NoError(err)
	suite.Equal("Mobile Homes", shipmentTypes.Name)
}
