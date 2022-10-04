package ghcimport

import (
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREIntlAccessorialPrices() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREIntlAccessorialPrices(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyIntlAccessorialPrices()
		suite.helperCheckIntlAccessorialPrices()
	})

	suite.Run("run a second time; should fail immediately due to constraint violation", func() {
		setupTestData()
		err := gre.importREIntlAccessorialPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_intl_accessorial_prices_unique_key"))
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyIntlAccessorialPrices() {
	count, err := suite.DB().Count(&models.ReIntlAccessorialPrice{})
	suite.NoError(err)
	suite.Equal(10, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckIntlAccessorialPrices() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = $1", testContractCode).First(&contract)
	suite.NoError(err)

	testCases := []struct {
		serviceCode   models.ReServiceCode
		market        string
		expectedPrice int
		isError       bool
	}{
		{models.ReServiceCodeICRT, "C", 2561, false},
		{models.ReServiceCodeICRTSA, "C", 2561, false},
		{models.ReServiceCodeIUCRT, "C", 654, false},
		{models.ReServiceCodeIDSHUT, "C", 14529, false},
		{models.ReServiceCodeIDSHUT, "O", 15623, false},
		{models.ReServiceCodeIOSHUT, "O", 15623, false},
		{models.ReServiceCodeMS, "O", 0, true},
		{models.ReServiceCodeICRT, "R", 0, true},
	}

	for _, testCase := range testCases {
		// Get service UUID.
		var service models.ReService
		err = suite.DB().Where("code = ?", testCase.serviceCode).First(&service)
		suite.NoError(err)

		var intlAccessorialPrice models.ReIntlAccessorialPrice
		err = suite.DB().
			Where("contract_id = $1", contract.ID).
			Where("service_id = $2", service.ID).
			Where("market = $3", testCase.market).
			First(&intlAccessorialPrice)

		if testCase.isError {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.Equal(unit.Cents(testCase.expectedPrice), intlAccessorialPrice.PerUnitCents, "test case: %+v", testCase)
		}
	}
}
