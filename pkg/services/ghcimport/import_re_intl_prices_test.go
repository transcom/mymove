package ghcimport

import (
	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREInternationalPrices() {
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

		err = gre.importREInternationalPrices(suite.AppContextForTest())
		suite.NoError(err)
	}

	suite.Run("import success", func() {
		setupTestData()
		suite.helperVerifyInternationalPrices()

		// Spot check prices
		suite.helperCheckInternationalPriceValues()
	})

	suite.Run("run a second time; should fail immediately due to constraint violation", func() {
		setupTestData()
		err := gre.importREInternationalPrices(suite.AppContextForTest())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_intl_prices_unique_key"))
		}
	})
}
func (suite *GHCRateEngineImportSuite) Test_getRateAreaIDForKind() {
	gre := &GHCRateEngineImporter{
		ContractCode: testContractCode,
	}

	setupTestData := func() {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.AppContextForTest())
		suite.NoError(err)

		err = gre.importRERateArea(suite.AppContextForTest())
		suite.NoError(err)
	}
	// Doing this here instead of a separate test function so we don't have to reload prerequisite tables
	suite.Run("tests for getRateAreaIDForKind", func() {
		setupTestData()
		testCases := []struct {
			name        string
			rateArea    string
			kind        string
			shouldError bool
		}{
			{"good NSRA", "NSRA2", "NSRA", false},
			{"good OCONUS", "US8101000", "OCONUS", false},
			{"good CONUS", "US47", "CONUS", false},
			{"bad NSRA", "XYZ", "NSRA", true},
			{"bad OCONUS", "US47", "OCONUS", true},
			{"bad CONUS", "NSRA13", "CONUS", true},
			{"bad kind", "NSRA2", "NNNN", true},
		}

		var contract models.ReContract
		err := suite.DB().Where("code = ?", testContractCode).First(&contract)
		suite.NoError(err)

		for _, testCase := range testCases {
			suite.Run(testCase.name, func() {
				id, err := gre.getRateAreaIDForKind(testCase.rateArea, testCase.kind)
				if testCase.shouldError {
					suite.Error(err)
					suite.Equal(uuid.Nil, id)
				} else {
					suite.NoError(err)

					// Fetch the UUID from the database and see if it matches
					origin, err := models.FetchReRateAreaItem(suite.DB(), contract.ID, testCase.rateArea)
					suite.NoError(err)
					suite.Equal(origin.ID, id)
				}
			})
		}
	})
}

func (suite *GHCRateEngineImportSuite) helperVerifyInternationalPrices() {
	count, err := suite.DB().Count(&models.ReIntlPrice{})
	suite.NoError(err)
	suite.Equal(276, count)
}

func (suite *GHCRateEngineImportSuite) helperCheckInternationalPriceValues() {
	// Get contract UUID.
	var contract models.ReContract
	err := suite.DB().Where("code = ?", testContractCode).First(&contract)
	suite.NoError(err)

	// Spot check one non-peak/peak record of each type
	testCases := []struct {
		serviceCode         models.ReServiceCode
		originRateArea      string
		destinationRateArea string
		isPeakPeriod        bool
		expectedPrice       int
	}{
		// 3a: OCONUS to OCONUS
		{models.ReServiceCodeIOOLH, "GE", "US8101000", false, 1021},
		{models.ReServiceCodeIOOUB, "GE", "US8101000", false, 1717},
		{models.ReServiceCodeIOOLH, "GE", "US8101000", true, 1205},
		{models.ReServiceCodeIOOUB, "GE", "US8101000", true, 2026},
		// 3b: CONUS to OCONUS
		{models.ReServiceCodeICOLH, "US47", "AS11", false, 3090},
		{models.ReServiceCodeICOUB, "US47", "AS11", false, 3398},
		{models.ReServiceCodeICOLH, "US47", "AS11", true, 3646},
		{models.ReServiceCodeICOUB, "US47", "AS11", true, 4010},
		// 3c: OCONUS to CONUS
		{models.ReServiceCodeIOCLH, "US8101000", "US68", false, 1757},
		{models.ReServiceCodeIOCUB, "US8101000", "US68", false, 3445},
		{models.ReServiceCodeIOCLH, "US8101000", "US68", true, 2073},
		{models.ReServiceCodeIOCUB, "US8101000", "US68", true, 4065},
		// 3e: NSRA to NSRA
		{models.ReServiceCodeNSTH, "NSRA2", "NSRA13", false, 4849},
		{models.ReServiceCodeNSTUB, "NSRA2", "NSRA13", false, 4793},
		{models.ReServiceCodeNSTH, "NSRA2", "NSRA13", true, 5722},
		{models.ReServiceCodeNSTUB, "NSRA2", "NSRA13", true, 5656},
		// 3e: NSRA to OCONUS
		{models.ReServiceCodeNSTH, "NSRA13", "AS11", false, 5172},
		{models.ReServiceCodeNSTUB, "NSRA13", "AS11", false, 1175},
		{models.ReServiceCodeNSTH, "NSRA13", "AS11", true, 6103},
		{models.ReServiceCodeNSTUB, "NSRA13", "AS11", true, 1386},
		// 3e: OCONUS to NSRA
		{models.ReServiceCodeNSTH, "GE", "NSRA2", false, 4872},
		{models.ReServiceCodeNSTUB, "GE", "NSRA2", false, 1050},
		{models.ReServiceCodeNSTH, "GE", "NSRA2", true, 5749},
		{models.ReServiceCodeNSTUB, "GE", "NSRA2", true, 1239},
		// 3e: NSRA to CONUS
		{models.ReServiceCodeNSTH, "NSRA2", "US4965500", false, 931},
		{models.ReServiceCodeNSTUB, "NSRA2", "US4965500", false, 1717},
		{models.ReServiceCodeNSTH, "NSRA2", "US4965500", true, 1099},
		{models.ReServiceCodeNSTUB, "NSRA2", "US4965500", true, 2026},
		// 3e: CONUS to NSRA
		{models.ReServiceCodeNSTH, "US68", "NSRA13", false, 1065},
		{models.ReServiceCodeNSTUB, "US68", "NSRA13", false, 1689},
		{models.ReServiceCodeNSTH, "US68", "NSRA13", true, 1257},
		{models.ReServiceCodeNSTUB, "US68", "NSRA13", true, 1993},
	}

	for _, testCase := range testCases {
		var service models.ReService
		err = suite.DB().Where("code = ?", testCase.serviceCode).First(&service)
		suite.NoError(err)

		// Get origin rate area UUID.
		origin, err := models.FetchReRateAreaItem(suite.DB(), contract.ID, testCase.originRateArea)
		suite.NoError(err)

		// Get destination rate area UUID.
		destination, err := models.FetchReRateAreaItem(suite.DB(), contract.ID, testCase.destinationRateArea)
		suite.NoError(err)

		var intlPrice models.ReIntlPrice
		err = suite.DB().
			Where("contract_id = ?", contract.ID).
			Where("service_id = ?", service.ID).
			Where("origin_rate_area_id = ?", origin.ID).
			Where("destination_rate_area_id = ?", destination.ID).
			Where("is_peak_period = ?", testCase.isPeakPeriod).
			First(&intlPrice)
		suite.NoError(err)
		suite.Equal(unit.Cents(testCase.expectedPrice), intlPrice.PerUnitCents, "test case: %+v", testCase)
	}
}
