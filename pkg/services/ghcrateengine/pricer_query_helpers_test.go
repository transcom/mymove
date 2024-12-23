package ghcrateengine

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_fetchTaskOrderFee() {
	testCents := unit.Cents(10000)
	testAvailableToPrimeAt := time.Date(testdatagen.TestYear, time.June, 17, 8, 45, 44, 333, time.UTC)
	suite.Run("golden path", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, testCents)

		taskOrderFee, err := fetchTaskOrderFee(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeMS, testAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(testCents, taskOrderFee.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, testCents)

		// Look for service code CS that we haven't added
		_, err := fetchTaskOrderFee(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeCS, testAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomOtherPrice() {
	suite.Run("golden path", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)
		domOtherPrice, err := fetchDomOtherPrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod)

		suite.NoError(err)
		suite.Equal(dpkTestBasePriceCents, domOtherPrice.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupDomesticOtherPrice(models.ReServiceCodeDPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod, dpkTestBasePriceCents, dpkTestContractYearName, dpkTestEscalationCompounded)

		// Look for service code IHPK that we haven't added
		_, err := fetchDomOtherPrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeIHPK, dpkTestServicesScheduleOrigin, dpkTestIsPeakPeriod)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomServiceAreaPrice() {
	testServiceArea := "123"
	testIsPeakPeriod := true
	testCents := unit.Cents(353)

	suite.Run("golden path", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod, testCents, "Base Period Year 1", 1.125)

		domServiceAreaPrice, err := fetchDomServiceAreaPrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domServiceAreaPrice.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod, testCents, "Base Period Year 1", 1.125)

		// Look for service code DDFSIT that we haven't added
		_, err := fetchDomServiceAreaPrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeDDFSIT, testServiceArea, testIsPeakPeriod)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchAccessorialPrice() {
	suite.Run("golden path", func() {
		suite.setupDomesticAccessorialPrice(models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule, ddshutTestBasePriceCents, testdatagen.DefaultContractCode, ddshutTestEscalationCompounded)
		domAccessorialPrice, err := fetchAccessorialPrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeDDSHUT, ddshutTestServiceSchedule)

		suite.NoError(err)
		suite.Equal(ddshutTestBasePriceCents, domAccessorialPrice.PerUnitCents)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchInternationalAccessorialPrice() {
	suite.Run("golden path", func() {
		suite.setupInternationalAccessorialPrice(models.ReServiceCodeIDSHUT, idshutTestServiceSchedule, idshutTestBasePriceCents, testdatagen.DefaultContractCode, idshutTestEscalationCompounded)
		internationalAccessorialPrice, err := fetchAccessorialPrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeDDSHUT, idshutTestServiceSchedule)

		suite.NoError(err)
		suite.Equal(idshutTestBasePriceCents, internationalAccessorialPrice.PerUnitCents)
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

		contractYear, err := fetchContractYear(suite.AppContextForTest(), newContractYear.ContractID, testDate)
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
		_, err := fetchContractYear(suite.AppContextForTest(), newContractYear.ContractID, testDate.AddDate(2, 0, 0))
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchShipmentTypePrice() {
	suite.Run("golden path", func() {
		suite.setupShipmentTypePrice(models.ReServiceCodeDNPK, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)
		shipmentTypePrice, err := fetchShipmentTypePrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeDNPK, models.MarketConus)

		suite.NoError(err)
		suite.Equal(dnpkTestFactor, shipmentTypePrice.Factor)
	})

	suite.Run("no records found", func() {
		suite.setupShipmentTypePrice(models.ReServiceCodeDNPK, models.MarketConus, dnpkTestFactor, dnpkTestContractYearName, dnpkTestEscalationCompounded)

		// Look for service code INPK that we haven't added
		_, err := fetchShipmentTypePrice(suite.AppContextForTest(), testdatagen.DefaultContractCode, models.ReServiceCodeINPK, models.MarketOconus)
		suite.Error(err)
	})
}
