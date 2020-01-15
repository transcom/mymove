package ghcimport

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticOtherPrices() {
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

		err = gre.importREDomesticOtherPrices(suite.DB())
		suite.NoError(err)

		suite.helperVerifyDomesticOtherPrices()
		suite.helperCheckDomesticOtherPriceValue()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREDomesticOtherPrices(suite.DB())
		if suite.Error(err) {
			suite.Contains(err.Error(), "duplicate key value violates unique constraint")
		}

		// Check to see if anything else changed
		suite.helperVerifyDomesticOtherPrices()
		suite.helperCheckDomesticOtherPriceValue()
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREDomesticOtherPricesFailures() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	err := gre.importREContract(suite.DB())
	suite.NoError(err)
	suite.NotNil(gre.contractID)

	err = gre.loadServiceMap(suite.DB())
	suite.NoError(err)

	suite.T().Run("stage_domestic_other_sit_prices table missing", func(t *testing.T) {
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "stage_domestic_other_sit_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)

		err = gre.importREDomesticOtherPrices(suite.DB())
		if suite.Error(err) {
			suite.Equal("error looking up StageDomesticOtherSitPrice data: unable to fetch records: pq: relation \"stage_domestic_other_sit_prices\" does not exist", err.Error())
		}
	})
	suite.T().Run("stage_domestic_other_pack_prices table missing", func(t *testing.T) {
		// drop a staging table that we are depending on to do import
		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s;", "stage_domestic_other_pack_prices")
		dropErr := suite.DB().RawQuery(dropQuery).Exec()
		suite.NoError(dropErr)

		err = gre.importREDomesticOtherPrices(suite.DB())
		if suite.Error(err) {
			suite.Equal("error looking up StageDomesticOtherPackPrice data: unable to fetch records: pq: relation \"stage_domestic_other_pack_prices\" does not exist", err.Error())
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

	suite.verifyDomesticOtherPrice(unit.Cents(7395), contract.ID, false, "DPK", 3)
	suite.verifyDomesticOtherPrice(unit.Cents(8000), contract.ID, true, "DPK", 3)
	suite.verifyDomesticOtherPrice(unit.Cents(597), contract.ID, false, "DUPK", 2)
	suite.verifyDomesticOtherPrice(unit.Cents(650), contract.ID, true, "DUPK", 2)
	suite.verifyDomesticOtherPrice(unit.Cents(23440), contract.ID, false, "DOPSIT", 2)
	suite.verifyDomesticOtherPrice(unit.Cents(24122), contract.ID, true, "DOPSIT", 2)
	suite.verifyDomesticOtherPrice(unit.Cents(24625), contract.ID, false, "DDDSIT", 3)
	suite.verifyDomesticOtherPrice(unit.Cents(25030), contract.ID, true, "DDDSIT", 3)
}

func (suite *GHCRateEngineImportSuite) verifyDomesticOtherPrice(expected unit.Cents, contractID uuid.UUID, isPeakPeriod bool, serviceCode string, schedule int) {
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
