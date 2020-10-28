package ghcimport

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgerrcode"

	"github.com/transcom/mymove/pkg/db/dberr"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineImportSuite) Test_importREInternationalOtherPrices() {
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

		err = gre.importREInternationalOtherPrices(suite.DB())
		suite.NoError(err)
		suite.helperVerifyInternationalOtherPrices()

		// Spot check a staging row's prices
		suite.helperCheckInternationalOtherPriceRecords()
	})

	suite.T().Run("run a second time; should fail immediately due to constraint violation", func(t *testing.T) {
		err := gre.importREInternationalOtherPrices(suite.DB())
		if suite.Error(err) {
			suite.True(dberr.IsDBErrorForConstraint(err, pgerrcode.UniqueViolation, "re_intl_other_prices_unique_key"))
		}

		// Check to see if anything else changed
		suite.helperVerifyInternationalOtherPrices()
		suite.helperCheckInternationalOtherPriceRecords()
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
		service       string
		expectedPrice int
	}{
		{"IHPK", 8186},
		{"IHUPK", 915},
		{"IUBPK", 8482},
		{"IUBUPK", 847},
		{"IOFSIT", 507},
		{"IDFSIT", 507},
		{"IOASIT", 14},
		{"IDASIT", 14},
		{"IOPSIT", 17001},
		{"IDDSIT", 30186},
	}

	for _, test := range testServices {
		suite.helperCheckOneOtherInternationalPriceRecord(test.expectedPrice, contract.ID, test.service, rateArea.ID)
	}
}

func (suite *GHCRateEngineImportSuite) helperCheckOneOtherInternationalPriceRecord(expected int, contractID uuid.UUID, serviceCode string, rateAreaID uuid.UUID) {
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
