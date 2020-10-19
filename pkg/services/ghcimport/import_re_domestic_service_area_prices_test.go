package ghcimport

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
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

		err = gre.loadServiceMap(suite.DB())
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
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_domestic_service_area_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyDomesticServiceAreaPrices()
		suite.helperCheckDomesticServiceAreaPriceValue()
	})
}

func (suite *GHCRateEngineImportSuite) Test_importREDomesticServiceAreaPricesFailures() {
	suite.T().Run("stage_domestic_service_area_prices table missing", func(t *testing.T) {
		renameQuery := fmt.Sprintf("ALTER TABLE stage_domestic_service_area_prices RENAME TO missing_stage_domestic_service_area_prices")
		renameErr := suite.DB().RawQuery(renameQuery).Exec()
		suite.NoError(renameErr)

		gre := &GHCRateEngineImporter{
			Logger:       suite.logger,
			ContractCode: testContractCode,
		}

		err := gre.importREContract(suite.DB())
		suite.NoError(err)
		suite.NotNil(gre.ContractID)

		err = gre.importREDomesticServiceAreaPrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBError(err, pgerrcode.UndefinedTable))
		}

		renameQuery = fmt.Sprintf("ALTER TABLE missing_stage_domestic_service_area_prices RENAME TO stage_domestic_service_area_prices")
		renameErr = suite.DB().RawQuery(renameQuery).Exec()
		suite.NoError(renameErr)
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
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(16), contract.ID, "DSH", serviceArea.ID)

	// Get domestic service area price DOP
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(581), contract.ID, "DOP", serviceArea.ID)

	// Get domestic service area price DDP
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(581), contract.ID, "DDP", serviceArea.ID)

	// Get domestic service area price DOFSIT
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(1597), contract.ID, "DOFSIT", serviceArea.ID)

	// Get domestic service area price DDFSIT
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(1597), contract.ID, "DDFSIT", serviceArea.ID)

	// Get domestic service area price DOASIT
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(62), contract.ID, "DOASIT", serviceArea.ID)

	// Get domestic service area price DDASIT
	suite.verifyDomesticSerivceAreaPrice(unit.Cents(62), contract.ID, "DDASIT", serviceArea.ID)
}

func (suite *GHCRateEngineImportSuite) verifyDomesticSerivceAreaPrice(expected unit.Cents, contractID uuid.UUID, serviceCode string, serviceAreaID uuid.UUID) {
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
