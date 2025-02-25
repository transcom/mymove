package ghcimport

import (
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticAccessorialPrices() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importREDomesticAccessorialPrices(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyDomesticAccessorialPrices()
		suite.helperCheckDomesticAccessorialPrices()
	})

	suite.Run("run a second time; should fail immediately due to constraint violation", func() {
		setupTestData()
		err := gre.importREDomesticAccessorialPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_domestic_accessorial_prices_unique_key"))
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyDomesticAccessorialPrices() {
	count, err := suite.DB().Count(&models.ReDomesticAccessorialPrice{})
	suite.NoError(err)
	suite.Equal(15, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckDomesticAccessorialPrices() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = $1", testContractCode).First(&contract)
	suite.NoError(err)

	testCases := []struct {
		serviceCode   models.ReServiceCode
		schedule      int
		expectedPrice int
		isError       bool
	}{
		{models.ReServiceCodeDCRT, 1, 2369, false},
		{models.ReServiceCodeDCRTSA, 1, 2369, false},
		{models.ReServiceCodeDUCRT, 1, 595, false},
		{models.ReServiceCodeDDSHUT, 1, 505, false},
		{models.ReServiceCodeDDSHUT, 3, 576, false},
		{models.ReServiceCodeDOSHUT, 1, 505, false},
		{models.ReServiceCodeMS, 3, 0, true},
		{models.ReServiceCodeDCRT, 5, 0, true},
	}

	for _, testCase := range testCases {
		// Get service UUID.
		var service models.ReService
		err = suite.DB().Where("code = ?", testCase.serviceCode).First(&service)
		suite.NoError(err)

		var domesticAccessorialPrice models.ReDomesticAccessorialPrice
		err = suite.DB().
			Where("contract_id = $1", contract.ID).
			Where("service_id = $2", service.ID).
			Where("services_schedule = $3", testCase.schedule).
			First(&domesticAccessorialPrice)

		if testCase.isError {
			suite.Error(err)
		} else {
			suite.NoError(err)
			suite.Equal(unit.Cents(testCase.expectedPrice), domesticAccessorialPrice.PerUnitCents, "test case: %+v", testCase)
		}
	}
}
