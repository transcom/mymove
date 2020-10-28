package ghcimport

import (
	"testing"

	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *GHCRateEngineImportSuite) Test_importREShipmentTypePrices() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.DB())
		suite.NoError(err)

		err = gre.importREShipmentTypePrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyShipmentTypePrices()
		suite.helperCheckShipmentTypePrices()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREShipmentTypePrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_shipment_type_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyShipmentTypePrices()
		suite.helperCheckShipmentTypePrices()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyShipmentTypePrices() {
	count, err := suite.DB().Count(&models.ReShipmentTypePrice{})
	suite.NoError(err)
	suite.Equal(7, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckShipmentTypePrices() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = $1", testContractCode).First(&contract)
	suite.NoError(err)

	// Get service UUID for shipment type
	var service models.ReService
	err = suite.DB().Where("code = 'DMHF'").First(&service)
	suite.NoError(err)

	var shipmentTypePrices models.ReShipmentTypePrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", service.ID).
		First(&shipmentTypePrices)
	suite.NoError(err)
	suite.Equal(1.20, shipmentTypePrices.Factor)
}
