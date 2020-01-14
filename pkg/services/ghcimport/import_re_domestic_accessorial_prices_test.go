package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticAccessorialPrices() {
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

		err = gre.importREDomesticAccessorialPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyDomesticAccessorialPrices()
		suite.helperCheckDomesticAccessorialPrices()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREDomesticAccessorialPrices(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "duplicate key value violates unique constraint")
		}

		// Check to see if anything else changed
		suite.helperVerifyDomesticAccessorialPrices()
		suite.helperCheckDomesticAccessorialPrices()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticAccessorialPrices() {
	count, err := suite.DB().Count(&models.ReDomesticAccessorialPrice{})
	suite.NoError(err)
	suite.Equal(9, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckDomesticAccessorialPrices() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = $1", testContractCode).First(&contract)
	suite.NoError(err)

	// Get service UUID.
	var serviceDCRT models.ReService
	err = suite.DB().Where("code = 'DCRT'").First(&serviceDCRT)
	suite.NoError(err)

	var serviceDUCRT models.ReService
	err = suite.DB().Where("code = 'DUCRT'").First(&serviceDUCRT)
	suite.NoError(err)

	var serviceDDSHUT models.ReService
	err = suite.DB().Where("code = 'DDSHUT'").First(&serviceDDSHUT)
	suite.NoError(err)

	var serviceNotValid models.ReService
	err = suite.DB().Where("code = 'MS'").First(&serviceNotValid)
	suite.NoError(err)

	var domesticAccessorialPriceDCRT models.ReDomesticAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceDCRT.ID).
		Where("services_schedule = $3", 1).
		First(&domesticAccessorialPriceDCRT)
	suite.NoError(err)
	suite.Equal(unit.Cents(2369), domesticAccessorialPriceDCRT.PerUnitCents)

	var domesticAccessorialPriceDUCRT models.ReDomesticAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceDUCRT.ID).
		Where("services_schedule = $3", 1).
		First(&domesticAccessorialPriceDUCRT)
	suite.NoError(err)
	suite.Equal(unit.Cents(595), domesticAccessorialPriceDUCRT.PerUnitCents)

	var domesticAccessorialPriceDDSHUT models.ReDomesticAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceDDSHUT.ID).
		Where("services_schedule = $3", 3).
		First(&domesticAccessorialPriceDDSHUT)
	suite.NoError(err)
	suite.Equal(unit.Cents(576), domesticAccessorialPriceDDSHUT.PerUnitCents)

	var domesticAccessorialPriceServiceNotValid models.ReDomesticAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceNotValid.ID).
		Where("services_schedule = $3", 3).
		First(&domesticAccessorialPriceServiceNotValid)
	suite.Error(err)
}
