package ghcimport

import (
	"fmt"
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPrices() {
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

		err = gre.importREDomesticServiceAreaPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyDomesticServiceAreaPrices()

		// Spot check domestic service area prices for one row
		suite.helperCheckDomesticServiceAreaPriceValue()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREDomesticServiceAreaPrices(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "duplicate key value violates unique constraint")
		}

		// Check to see if anything else changed
		suite.helperVerifyDomesticServiceAreaPrices()
		suite.helperCheckDomesticServiceAreaPriceValue()
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPricesFailures() {
	suite.T().Run("stage_domestic_service_area_prices table missing", func(t *testing.T) {
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "stage_domestic_service_area_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)

		gre := &GHCRateEngineImporter{
			Logger:       suite.logger,
			ContractCode: testContractCode,
		}

		err := gre.importREContract(suite.DB())
		suite.NoError(err)
		suite.NotNil(gre.contractID)

		err = gre.importREDomesticServiceAreaPrices(suite.DB())
		if suite.Error(err) {
			suite.Equal("Error looking up StageDomesticServiceAreaPrice data: unable to fetch records: pq: relation \"stage_domestic_service_area_prices\" does not exist", err.Error())
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticServiceAreaPrices() {
	count, err := suite.DB().Count(&models.ReDomesticServiceAreaPrices{})
	suite.NoError(err)
	suite.Equal(32, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckDomesticServiceAreaPriceValue() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Get domestic service area UUID.
	var serviceArea models.ReDomesticServiceArea
	err = suite.DB().Where("service_area = '592'").First(&serviceArea)
	suite.NoError(err)

	var serviceDSH models.ReService
	err = suite.DB().Where("code = 'DSH'").First(&serviceDSH)
	suite.NoError(err)

	// Get domestic service area price DSH
	var domesticServiceAreaPrice models.ReDomesticServiceAreaPrice
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceDSH.ID).
		Where("domestic_service_area_id = ?", serviceArea.ID).
		Where("is_peak_period = false").First(&domesticServiceAreaPrice)
	suite.NoError(err)
	suite.Equal(unit.Cents(16), domesticServiceAreaPrice.PriceCents)

	var serviceDODP models.ReService
	err = suite.DB().Where("code = 'DODP'").First(&serviceDODP)
	suite.NoError(err)

	// Get domestic service area price DODP
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceDODP.ID).
		Where("domestic_service_area_id = ?", serviceArea.ID).
		Where("is_peak_period = false").First(&domesticServiceAreaPrice)
	suite.NoError(err)
	suite.Equal(unit.Cents(581), domesticServiceAreaPrice.PriceCents)

	var serviceDFSIT models.ReService
	err = suite.DB().Where("code = 'DFSIT'").First(&serviceDFSIT)
	suite.NoError(err)

	// Get domestic service area price DFSIT
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceDFSIT.ID).
		Where("domestic_service_area_id = ?", serviceArea.ID).
		Where("is_peak_period = false").First(&domesticServiceAreaPrice)
	suite.NoError(err)
	suite.Equal(unit.Cents(1597), domesticServiceAreaPrice.PriceCents)

	var serviceDASIT models.ReService
	err = suite.DB().Where("code = 'DASIT'").First(&serviceDASIT)
	suite.NoError(err)

	// Get domestic service area price DASIT
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_id = ?", serviceDASIT.ID).
		Where("domestic_service_area_id = ?", serviceArea.ID).
		Where("is_peak_period = false").First(&domesticServiceAreaPrice)
	suite.NoError(err)
	suite.Equal(unit.Cents(62), domesticServiceAreaPrice.PriceCents)
}
