package ghcrateengine

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCRateEngineServiceSuite) Test_fetchTaskOrderFee() {
	testCents := unit.Cents(10000)
	testAvailableToPrimeAt := time.Date(testdatagen.TestYear, time.June, 17, 8, 45, 44, 333, time.UTC)
	suite.setupTaskOrderFeeData(models.ReServiceCodeMS, testCents)

	suite.T().Run("golden path", func(t *testing.T) {
		taskOrderFee, err := fetchTaskOrderFee(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeMS, testAvailableToPrimeAt)
		suite.NoError(err)
		suite.Equal(testCents, taskOrderFee.PriceCents)
	})

	suite.T().Run("no records found", func(t *testing.T) {
		// Look for service code CS that we haven't added
		_, err := fetchTaskOrderFee(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeCS, testAvailableToPrimeAt)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomOtherPrice() {
	testCents := unit.Cents(146)
	servicesSchedule := 1
	isPeakPeriod := true
	suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDPK)

	suite.T().Run("golden path", func(t *testing.T) {
		domOtherPrice, err := fetchDomOtherPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDPK, servicesSchedule, isPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domOtherPrice.PriceCents)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_unpackFetchDomOtherPrice() {
	testCents := unit.Cents(146)
	servicesSchedule := 1
	isPeakPeriod := true
	suite.setUpDomesticPackAndUnpackData(models.ReServiceCodeDUPK)

	suite.T().Run("golden path", func(t *testing.T) {
		domOtherPrice, err := fetchDomOtherPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDUPK, servicesSchedule, isPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domOtherPrice.PriceCents)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchDomServiceAreaPrice() {
	testServiceArea := "123"
	testIsPeakPeriod := true
	testCents := unit.Cents(353)
	suite.setupDomesticServiceAreaPrice(models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod, testCents, 1.125)

	suite.T().Run("golden path", func(t *testing.T) {
		domServiceAreaPrice, err := fetchDomServiceAreaPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDOFSIT, testServiceArea, testIsPeakPeriod)
		suite.NoError(err)
		suite.Equal(testCents, domServiceAreaPrice.PriceCents)
	})

	suite.T().Run("no records found", func(t *testing.T) {
		// Look for service code DDFSIT that we haven't added
		_, err := fetchDomServiceAreaPrice(suite.DB(), testdatagen.DefaultContractCode, models.ReServiceCodeDDFSIT, testServiceArea, testIsPeakPeriod)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) Test_fetchContractYear() {
	testDate := time.Date(testdatagen.TestYear, time.June, 17, 8, 45, 44, 333, time.UTC)
	testEscalationCompounded := 1.0512
	newContractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				EscalationCompounded: testEscalationCompounded,
			},
		})

	suite.T().Run("golden path", func(t *testing.T) {
		contractYear, err := fetchContractYear(suite.DB(), newContractYear.ContractID, testDate)
		suite.NoError(err)
		suite.Equal(testEscalationCompounded, contractYear.EscalationCompounded)
	})

	suite.T().Run("no records found", func(t *testing.T) {
		// Look for a testDate that's a couple of years later.
		_, err := fetchContractYear(suite.DB(), newContractYear.ContractID, testDate.AddDate(2, 0, 0))
		suite.Error(err)
	})
}
