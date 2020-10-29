package ghcimport

import (
	"testing"

	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREDomesticAccessorialPrices() {
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

		err = gre.importREDomesticAccessorialPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyDomesticAccessorialPrices()
		suite.helperCheckDomesticAccessorialPrices()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREDomesticAccessorialPrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_domestic_accessorial_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyDomesticAccessorialPrices()
		suite.helperCheckDomesticAccessorialPrices()
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
		serviceCode   string
		schedule      int
		expectedPrice int
		isError       bool
	}{
		{"DCRT", 1, 2369, false},
		{"DCRTSA", 1, 2369, false},
		{"DUCRT", 1, 595, false},
		{"DDSHUT", 1, 505, false},
		{"DDSHUT", 3, 576, false},
		{"DOSHUT", 1, 505, false},
		{"MS", 3, 0, true},
		{"DCRT", 5, 0, true},
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
