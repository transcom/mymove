package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREIntlAccessorialPrices() {
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

		err = gre.importREIntlAccessorialPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyIntlAccessorialPrices()
		suite.helperCheckIntlAccessorialPrices()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREIntlAccessorialPrices(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "duplicate key value violates unique constraint")
		}

		// Check to see if anything else changed
		suite.helperVerifyIntlAccessorialPrices()
		suite.helperCheckIntlAccessorialPrices()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyIntlAccessorialPrices() {
	count, err := suite.DB().Count(&models.ReIntlAccessorialPrice{})
	suite.NoError(err)
	suite.Equal(6, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckIntlAccessorialPrices() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = $1", testContractCode).First(&contract)
	suite.NoError(err)

	// Get service UUID.
	var serviceICRT models.ReService
	err = suite.DB().Where("code = 'ICRT'").First(&serviceICRT)
	suite.NoError(err)

	var serviceIUCRT models.ReService
	err = suite.DB().Where("code = 'IUCRT'").First(&serviceIUCRT)
	suite.NoError(err)

	var serviceIDSHUT models.ReService
	err = suite.DB().Where("code = 'IDSHUT'").First(&serviceIDSHUT)
	suite.NoError(err)

	var intlAccessorialPriceICRT models.ReIntlAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceICRT.ID).
		Where("market = $3", "C").
		First(&intlAccessorialPriceICRT)
	suite.NoError(err)
	suite.Equal(unit.Cents(2561), intlAccessorialPriceICRT.PerUnitCents)

	var intlAccessorialPriceIUCRT models.ReIntlAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceIUCRT.ID).
		Where("market = $3", "C").
		First(&intlAccessorialPriceIUCRT)
	suite.NoError(err)
	suite.Equal(unit.Cents(654), intlAccessorialPriceIUCRT.PerUnitCents)

	var intlAccessorialPriceIDSHUT models.ReIntlAccessorialPrice
	err = suite.DB().
		Where("contract_id = $1", contract.ID).
		Where("service_id = $2", serviceIDSHUT.ID).
		Where("market = $3", "O").
		First(&intlAccessorialPriceIDSHUT)
	suite.NoError(err)
	suite.Equal(unit.Cents(15623), intlAccessorialPriceIDSHUT.PerUnitCents)
}
