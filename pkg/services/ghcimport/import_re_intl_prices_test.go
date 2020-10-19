package ghcimport

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREInternationalPrices() {
	gre := &GHCRateEngineImporter{
		Logger:       suite.logger,
		ContractCode: testContractCode,
	}

	suite.T().Run("import success", func(t *testing.T) {
		// Prerequisite tables must be loaded.
		err := gre.importREContract(suite.DB())
		suite.NoError(err)

		err = gre.importRERateArea(suite.DB())
		suite.NoError(err)

		err = gre.loadServiceMap(suite.DB())
		suite.NoError(err)

		err = gre.importREInternationalPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyInternationalPrices()

		// Spot check prices
		suite.helperCheckInternationalPriceValues()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREInternationalPrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_intl_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyInternationalPrices()
		suite.helperCheckInternationalPriceValues()
	})

	// Doing this here instead of a separate test function so we don't have to reload prerequisite tables
	suite.T().Run("tests for getRateAreaIDForKind", func(t *testing.T) {
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
			suite.T().Run(testCase.name, func(t *testing.T) {
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
		serviceCode         string
		originRateArea      string
		destinationRateArea string
		isPeakPeriod        bool
		expectedPrice       int
	}{
		// 3a: OCONUS to OCONUS
		{"IOOLH", "GE", "US8101000", false, 1021},
		{"IOOUB", "GE", "US8101000", false, 1717},
		{"IOOLH", "GE", "US8101000", true, 1205},
		{"IOOUB", "GE", "US8101000", true, 2026},
		// 3b: CONUS to OCONUS
		{"ICOLH", "US47", "AS11", false, 3090},
		{"ICOUB", "US47", "AS11", false, 3398},
		{"ICOLH", "US47", "AS11", true, 3646},
		{"ICOUB", "US47", "AS11", true, 4010},
		// 3c: OCONUS to CONUS
		{"IOCLH", "US8101000", "US68", false, 1757},
		{"IOCUB", "US8101000", "US68", false, 3445},
		{"IOCLH", "US8101000", "US68", true, 2073},
		{"IOCUB", "US8101000", "US68", true, 4065},
		// 3e: NSRA to NSRA
		{"NSTH", "NSRA2", "NSRA13", false, 4849},
		{"NSTUB", "NSRA2", "NSRA13", false, 4793},
		{"NSTH", "NSRA2", "NSRA13", true, 5722},
		{"NSTUB", "NSRA2", "NSRA13", true, 5656},
		// 3e: NSRA to OCONUS
		{"NSTH", "NSRA13", "AS11", false, 5172},
		{"NSTUB", "NSRA13", "AS11", false, 1175},
		{"NSTH", "NSRA13", "AS11", true, 6103},
		{"NSTUB", "NSRA13", "AS11", true, 1386},
		// 3e: OCONUS to NSRA
		{"NSTH", "GE", "NSRA2", false, 4872},
		{"NSTUB", "GE", "NSRA2", false, 1050},
		{"NSTH", "GE", "NSRA2", true, 5749},
		{"NSTUB", "GE", "NSRA2", true, 1239},
		// 3e: NSRA to CONUS
		{"NSTH", "NSRA2", "US4965500", false, 931},
		{"NSTUB", "NSRA2", "US4965500", false, 1717},
		{"NSTH", "NSRA2", "US4965500", true, 1099},
		{"NSTUB", "NSRA2", "US4965500", true, 2026},
		// 3e: CONUS to NSRA
		{"NSTH", "US68", "NSRA13", false, 1065},
		{"NSTUB", "US68", "NSRA13", false, 1689},
		{"NSTH", "US68", "NSRA13", true, 1257},
		{"NSTUB", "US68", "NSRA13", true, 1993},
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
