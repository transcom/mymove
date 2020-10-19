package ghcimport

import (
	"testing"

	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticLinehaulPrices() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importREDomesticServiceArea(suite.DB())
		suite.NoError(err)

		err = gre.importREDomesticLinehaulPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyDomesticLinehaulCount()

		// Spot check a linehaul price
		suite.helperCheckDomesticLinehaulValue()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREDomesticLinehaulPrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_domestic_linehaul_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyDomesticLinehaulCount()
		suite.helperCheckDomesticLinehaulValue()
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticLinehaulCount() {
	count, err := suite.DB().Count(&models.ReDomesticLinehaulPrice{})
	suite.NoError(err)
	suite.Equal(240, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckDomesticLinehaulValue() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Get domestic service area UUID.
	var serviceArea models.ReDomesticServiceArea
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_area = '452'").
		First(&serviceArea)
	suite.NoError(err)

	// Get linehaul price.
	var linehaul models.ReDomesticLinehaulPrice
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("weight_lower = 5000").
		Where("weight_upper = 9999").
		Where("miles_lower = 2501").
		Where("miles_upper = 3000").
		Where("is_peak_period = false").
		Where("domestic_service_area_id = ?", serviceArea.ID).
		First(&linehaul)
	suite.NoError(err)

	suite.Equal(unit.Millicents(745600), linehaul.PriceMillicents)
}
