package ghcimport

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREInternationalOtherPrices() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importRERateArea(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREInternationalOtherPrices(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyInternationalOtherPrices()

		// Spot check a staging row's prices
		suite.helperCheckInternationalOtherPriceRecords()
	})

	suite.Run("run twice; should immediately fail the second time due to constraint violation", func() {
		setupTestData()
		err := gre.importREInternationalOtherPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_intl_other_prices_unique_key"))
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyInternationalOtherPrices() {
	count, err := suite.DB().Count(&models.ReIntlOtherPrice{})
	suite.NoError(err)
	suite.Equal(180, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckInternationalOtherPriceRecords() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Get rate area UUID.
	var rateArea *models.ReRateArea
	rateArea, err = models.FetchReRateAreaItem(suite.DB(), contract.ID, "US68")
	suite.NoError(err)

	// Get service UUID.
	testServices := []struct {
		service       models.ReServiceCode
		expectedPrice int
	}{
		{models.ReServiceCodeIHPK, 8186},
		{models.ReServiceCodeIHUPK, 915},
		{models.ReServiceCodeIUBPK, 8482},
		{models.ReServiceCodeIUBUPK, 847},
		{models.ReServiceCodeIOFSIT, 507},
		{models.ReServiceCodeIDFSIT, 507},
		{models.ReServiceCodeIOASIT, 14},
		{models.ReServiceCodeIDASIT, 14},
		{models.ReServiceCodeIOPSIT, 17001},
		{models.ReServiceCodeIDDSIT, 30186},
	}

	for _, test := range testServices {
		suite.helperCheckOneOtherInternationalPriceRecord(test.expectedPrice, contract.ID, test.service, rateArea.ID)
	}
}

func (suite *GHCRateEngineImportSuite) helperCheckOneOtherInternationalPriceRecord(expected int, contractID uuid.UUID, serviceCode models.ReServiceCode, rateAreaID uuid.UUID) {
	var service models.ReService
	err := suite.DB().Where("code = ?", serviceCode).First(&service)
	suite.NoError(err)

	var intlOtherPrice models.ReIntlOtherPrice
	err = suite.DB().
		Where("contract_id = ?", contractID).
		Where("service_id = ?", service.ID).
		Where("is_peak_period = true").
		Where("rate_area_id = ?", rateAreaID).
		First(&intlOtherPrice)
	suite.NoError(err)
	suite.Equal(unit.Cents(expected), intlOtherPrice.PerUnitCents)
}
