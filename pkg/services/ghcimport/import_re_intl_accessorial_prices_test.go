package ghcimport

import (
	"testing"

	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREIntlAccessorialPrices() {
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

		err = gre.importREIntlAccessorialPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyIntlAccessorialPrices()
		suite.helperCheckIntlAccessorialPrices()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREIntlAccessorialPrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_intl_accessorial_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyIntlAccessorialPrices()
		suite.helperCheckIntlAccessorialPrices()
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
		serviceCode   string
		market        string
		expectedPrice int
		isError       bool
	}{
		{"ICRT", "C", 2561, false},
		{"ICRTSA", "C", 2561, false},
		{"IUCRT", "C", 654, false},
		{"IDSHUT", "C", 14529, false},
		{"IDSHUT", "O", 15623, false},
		{"IOSHUT", "O", 15623, false},
		{"MS", "O", 0, true},
		{"ICRT", "R", 0, true},
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
