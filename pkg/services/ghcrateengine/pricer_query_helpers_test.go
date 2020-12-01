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
