package ghcimport

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticOtherPrices() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREDomesticOtherPrices(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyDomesticOtherPrices()
		suite.helperCheckDomesticOtherPriceValue()
	})

	suite.Run("run a second time; should fail immediately due to constraint violation", func() {
		setupTestData()
		err := gre.importREDomesticOtherPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_domestic_other_prices_unique_key"))
		}
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREDomesticOtherPricesFailures() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}
	setupTestData := func() {
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)
		suite.NotNil(gre.ContractID)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("stage_domestic_other_sit_prices table missing", func() {
		setupTestData()
		renameQuery := "ALTER TABLE stage_domestic_other_sit_prices RENAME TO missing_stage_domestic_other_sit_prices"
		renameErr := suite.DB().RawQuery(renameQuery).Exec()
		suite.NoError(renameErr)

		err := gre.importREDomesticOtherPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBError(err, pgerrcode.UndefinedTable))
		}
	})

	suite.Run("stage_domestic_other_pack_prices table missing", func() {
		setupTestData()
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "stage_domestic_other_pack_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)

		err := gre.importREDomesticOtherPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBError(err, pgerrcode.UndefinedTable))
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticOtherPrices() {
	count, err := suite.DB().Count(&models.ReDomesticOtherPrice{})
	suite.NoError(err)
	suite.Equal(24, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckDomesticOtherPriceValue() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	suite.verifyDomesticOtherPrice(unit.Cents(7395), contract.ID, false, models.ReServiceCodeDPK, 3)
	suite.verifyDomesticOtherPrice(unit.Cents(8000), contract.ID, true, models.ReServiceCodeDPK, 3)
	suite.verifyDomesticOtherPrice(unit.Cents(597), contract.ID, false, models.ReServiceCodeDUPK, 2)
	suite.verifyDomesticOtherPrice(unit.Cents(650), contract.ID, true, models.ReServiceCodeDUPK, 2)
	suite.verifyDomesticOtherPrice(unit.Cents(23440), contract.ID, false, models.ReServiceCodeDOPSIT, 2)
	suite.verifyDomesticOtherPrice(unit.Cents(24122), contract.ID, true, models.ReServiceCodeDOPSIT, 2)
	suite.verifyDomesticOtherPrice(unit.Cents(24625), contract.ID, false, models.ReServiceCodeDDDSIT, 3)
	suite.verifyDomesticOtherPrice(unit.Cents(25030), contract.ID, true, models.ReServiceCodeDDDSIT, 3)
}

func (suite *GHCRateEngineImportSuite) verifyDomesticOtherPrice(expected unit.Cents, contractID uuid.UUID, isPeakPeriod bool, serviceCode models.ReServiceCode, schedule int) {
	var service models.ReService
	err := suite.DB().Where("code = ?", serviceCode).First(&service)
	suite.NoError(err)

	var domesticOtherPrice models.ReDomesticOtherPrice
	err = suite.DB().
		Where("contract_id = ?", contractID).
		Where("service_id = ?", service.ID).
		Where("is_peak_period = ?", isPeakPeriod).
		Where("schedule = ?", schedule).First(&domesticOtherPrice)
	suite.NoError(err)
	suite.Equal(expected, domesticOtherPrice.PriceCents)
}
