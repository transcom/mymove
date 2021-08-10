package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_fetchTaskOrderFee() {
	testCents := unit.Cents(10000)
	testAvailableToPrimeAt := time.Date(testdatagen.TestYear, time.June, 17, 8, 45, 44, 333, time.UTC)
	suite.Run("golden path", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, testCents)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		taskOrderFee, err := fetchTaskOrderFee(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeMS, testAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(testCents, taskOrderFee.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, testCents)

		// Look for service code CS that we haven't added
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := fetchTaskOrderFee(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeCS, testAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomOtherPrice() {
	testCents := unit.Cents(146)
	servicesSchedule := 1
	isPeakPeriod := true

	suite.Run("golden path", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDPK)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		domOtherPrice, err := fetchDomOtherPrice(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeDPK, servicesSchedule, isPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domOtherPrice.PriceCents)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_unpackFetchDomOtherPrice() {
	testCents := unit.Cents(146)
	servicesSchedule := 1
	isPeakPeriod := true

	suite.Run("golden path", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		domOtherPrice, err := fetchDomOtherPrice(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeDUPK, servicesSchedule, isPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domOtherPrice.PriceCents)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomServiceAreaPrice() {
	testServiceArea := "123"
	testIsPeakPeriod := true
	testCents := unit.Cents(353)

	suite.Run("golden path", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod, testCents, "Test Contract Year", 1.125)

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		domServiceAreaPrice, err := fetchDomServiceAreaPrice(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domServiceAreaPrice.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod, testCents, "Test Contract Year", 1.125)

		// Look for service code DDFSIT that we haven't added
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := fetchDomServiceAreaPrice(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeDDFSIT, testServiceArea, testIsPeakPeriod)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchAccessorialPrice() {
	suite.Run("golden path", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		domAccessorialPrice, err := fetchAccessorialPrice(appCfg, testdatagen.DefaultContractCode, models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule)

		suite.NoError(err)
		suite.Equal(ddshutTestBasePriceCents, domAccessorialPrice.PerUnitCents)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchContractYear() {
	testDate := time.Date(testdatagen.TestYear, time.June, 17, 8, 45, 44, 333, time.UTC)
	testEscalationCompounded := 1.0512

	suite.Run("golden path", func() {
		newContractYear := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: testEscalationCompounded,
				},
			})

		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		contractYear, err := fetchContractYear(appCfg, newContractYear.ContractID, testDate)
		suite.NoError(err)
		suite.Equal(testEscalationCompounded, contractYear.EscalationCompounded)
	})

	suite.Run("no records found", func() {
		newContractYear := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					EscalationCompounded: testEscalationCompounded,
				},
			})

		// Look for a testDate that's a couple of years later.
		appCfg := appconfig.NewAppConfig(suite.DB(), suite.logger)
		_, err := fetchContractYear(appCfg, newContractYear.ContractID, testDate.AddDate(2, 0, 0))
		suite.Error(err)
	})
}
