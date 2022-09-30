package ghcimport

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPrices() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREDomesticServiceArea(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREDomesticServiceAreaPrices(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyDomesticServiceAreaPrices()

		// Spot check domestic service area prices for one row
		suite.helperCheckDomesticServiceAreaPriceValue()
	})

	suite.Run("run a second time; should fail immediately due to constraint violation", func() {
		setupTestData()
		err := gre.importREDomesticServiceAreaPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_domestic_service_area_prices_unique_key"))
		}
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPricesFailures() {
	suite.Run("stage_domestic_service_area_prices table missing", func() {
		renameQuery := "ALTER TABLE stage_domestic_service_area_prices RENAME TO missing_stage_domestic_service_area_prices"
		renameErr := suite.DB().RawQuery(renameQuery).Exec()
		suite.NoError(renameErr)

		gre := &GHCRateEngineImporter{
			ContractCode: testContractCode,
		}

		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)
		suite.NotNil(gre.ContractID)

		err = gre.importREDomesticServiceAreaPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBError(err, pgerrcode.UndefinedTable))
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticServiceAreaPrices() {
	count, err := suite.DB().Count(&models.ReDomesticServiceAreaPrice{})
	suite.NoError(err)
	suite.Equal(70, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckDomesticServiceAreaPriceValue() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Get domestic service area UUID.
	var serviceArea models.ReDomesticServiceArea
	err = suite.DB().
		Where("contract_id = ?", contract.ID).
		Where("service_area = '592'").
		First(&serviceArea)
	suite.NoError(err)

	// Get domestic service area price DSH
	suite.verifyDomesticServiceAreaPrice(unit.Cents(16), contract.ID, models.ReServiceCodeDSH, serviceArea.ID)

	// Get domestic service area price DOP
	suite.verifyDomesticServiceAreaPrice(unit.Cents(581), contract.ID, models.ReServiceCodeDOP, serviceArea.ID)

	// Get domestic service area price DDP
	suite.verifyDomesticServiceAreaPrice(unit.Cents(581), contract.ID, models.ReServiceCodeDDP, serviceArea.ID)

	// Get domestic service area price DOFSIT
	suite.verifyDomesticServiceAreaPrice(unit.Cents(1597), contract.ID, models.ReServiceCodeDOFSIT, serviceArea.ID)

	// Get domestic service area price DDFSIT
	suite.verifyDomesticServiceAreaPrice(unit.Cents(1597), contract.ID, models.ReServiceCodeDDFSIT, serviceArea.ID)

	// Get domestic service area price DOASIT
	suite.verifyDomesticServiceAreaPrice(unit.Cents(62), contract.ID, models.ReServiceCodeDOASIT, serviceArea.ID)

	// Get domestic service area price DDASIT
	suite.verifyDomesticServiceAreaPrice(unit.Cents(62), contract.ID, models.ReServiceCodeDDASIT, serviceArea.ID)
}

func (suite *GHCRateEngineImportSuite) verifyDomesticServiceAreaPrice(expected unit.Cents, contractID uuid.UUID, serviceCode models.ReServiceCode, serviceAreaID uuid.UUID) {
	var service models.ReService
	err := suite.DB().Where("code = ?", serviceCode).First(&service)
	suite.NoError(err)

	var domesticServiceAreaPrice models.ReDomesticServiceAreaPrice
	err = suite.DB().
		Where("contract_id = ?", contractID).
		Where("service_id = ?", service.ID).
		Where("domestic_service_area_id = ?", serviceAreaID).
		Where("is_peak_period = false").First(&domesticServiceAreaPrice)
	suite.NoError(err)
	suite.Equal(expected, domesticServiceAreaPrice.PriceCents)
}
