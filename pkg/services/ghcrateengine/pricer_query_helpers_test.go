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

		taskOrderFee, err := fetchTaskOrderFee(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeMS, testAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(testCents, taskOrderFee.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupTaskOrderFeeData(models.ReServiceCodeMS, testCents)

		// Look for service code CS that we haven't added
		_, err := fetchTaskOrderFee(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeCS, testAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomOtherPrice() {
	testCents := unit.Cents(146)
	servicesSchedule := 1
	isPeakPeriod := true

	suite.Run("golden path", func() {
		suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDPK)

		domOtherPrice, err := fetchDomOtherPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDPK, servicesSchedule, isPeakPeriod)
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

		domOtherPrice, err := fetchDomOtherPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDUPK, servicesSchedule, isPeakPeriod)
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

		domServiceAreaPrice, err := fetchDomServiceAreaPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domServiceAreaPrice.PriceCents)
	})

	suite.Run("no records found", func() {
		suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod, testCents, "Test Contract Year", 1.125)

		// Look for service code DDFSIT that we haven't added
		_, err := fetchDomServiceAreaPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDDFSIT, testServiceArea, testIsPeakPeriod)
		suite.Error(err)
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

		contractYear, err := fetchContractYear(suite.DB(), newContractYear.ContractID, testDate)
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
		_, err := fetchContractYear(suite.DB(), newContractYear.ContractID, testDate.AddDate(2, 0, 0))
		suite.Error(err)
	})
}
