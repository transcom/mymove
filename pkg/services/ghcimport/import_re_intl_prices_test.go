package ghcimport

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREInternationalPrices() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importRERateArea(suite.DB())
		suite.NoError(err)

		err = gre.importREInternationalPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyInternationalPrices()

		// Spot check a linehaul price
		suite.helperCheckInternationalPriceValue()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREInternationalPrices(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "duplicate key value violates unique constraint")
		}

		// Check to see if anything else changed
		suite.helperVerifyInternationalPrices()
		suite.helperCheckInternationalPriceValue()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyInternationalPrices() {
	count, err := suite.DB().Count(&models.ReIntlPrices{})
	suite.NoError(err)
	suite.Equal(132, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckInternationalPriceValue() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Get service UUID.
	var serviceIOOLH models.ReService
	err = suite.DB().Where("code = 'IOOLH'").First(&serviceIOOLH)
	suite.NoError(err)

	var serviceIOOUB models.ReService
	err = suite.DB().Where("code = 'IOOUB'").First(&serviceIOOUB)
	suite.NoError(err)

	// Get origin rate area UUID.
	var origin *models.ReRateArea
	origin, err = models.FetchReRateAreaItem(suite.DB(), "GE")
	suite.NoError(err)

	// Get destination rate area UUID.
	var destination *models.ReRateArea
	destination, err = models.FetchReRateAreaItem(suite.DB(), "US8101000")
	suite.NoError(err)

	var intlPriceIOOLH models.ReIntlPrice
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceIOOLH.ID).
		Where("origin_rate_area_id = ?", origin.ID).
		Where("destination_rate_area_id = ?", destination.ID).
		Where("is_peak_period = false").
		First(&intlPriceIOOLH)
	suite.NoError(err)
	suite.Equal(unit.Cents(1021), intlPriceIOOLH.PerUnitCents)

	var intlPriceIOOLHPeak models.ReIntlPrice
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceIOOLH.ID).
		Where("origin_rate_area_id = ?", origin.ID).
		Where("destination_rate_area_id = ?", destination.ID).
		Where("is_peak_period = true").
		First(&intlPriceIOOLHPeak)
	suite.NoError(err)
	suite.Equal(unit.Cents(1205), intlPriceIOOLHPeak.PerUnitCents)

	var intlPriceIOOUB models.ReIntlPrice
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceIOOUB.ID).
		Where("origin_rate_area_id = ?", origin.ID).
		Where("destination_rate_area_id = ?", destination.ID).
		Where("is_peak_period = false").
		First(&intlPriceIOOUB)
	suite.NoError(err)
	suite.Equal(unit.Cents(1717), intlPriceIOOUB.PerUnitCents)

	var intlPriceIOOUBPeak models.ReIntlPrice
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceIOOUB.ID).
		Where("origin_rate_area_id = ?", origin.ID).
		Where("destination_rate_area_id = ?", destination.ID).
		Where("is_peak_period = true").
		First(&intlPriceIOOUBPeak)
	suite.NoError(err)
	suite.Equal(unit.Cents(2026), intlPriceIOOUBPeak.PerUnitCents)
}
