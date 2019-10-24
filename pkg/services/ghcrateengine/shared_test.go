package ghcrateengine

import (
	"testing"
	"time"
)

func (suite *GHCRateEngineServiceSuite) TestIsPeakPeriod() {
	suite.T().Run("within peak", func(t *testing.T) {
		dateInYear := peakStart.addDate(1, 0)
		date := time.Date(2019, dateInYear.month, 1, 0, 0, 0, 0, time.UTC)
		suite.True(IsPeakPeriod(date))
	})

	suite.T().Run("on peak start date", func(t *testing.T) {
		date := time.Date(2019, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)
		suite.True(IsPeakPeriod(date))
	})

	suite.T().Run("on peak end date", func(t *testing.T) {
		date := time.Date(2019, peakEnd.month, peakEnd.day, 0, 0, 0, 0, time.UTC)
		suite.True(IsPeakPeriod(date))
	})

	suite.T().Run("just before peak start date", func(t *testing.T) {
		dateInYear := peakStart.addDate(0, -1)
		date := time.Date(2019, dateInYear.month, dateInYear.day, 0, 0, 0, 0, time.UTC)
		suite.False(IsPeakPeriod(date))
	})

	suite.T().Run("just outside peak start date", func(t *testing.T) {
		dateInYear := peakEnd.addDate(0, 1)
		date := time.Date(2019, dateInYear.month, dateInYear.day, 0, 0, 0, 0, time.UTC)
		suite.False(IsPeakPeriod(date))
	})

	suite.T().Run("outside peak", func(t *testing.T) {
		dateInYear := peakEnd.addDate(1, 0)
		date := time.Date(2019, dateInYear.month, 1, 0, 0, 0, 0, time.UTC)
		suite.False(IsPeakPeriod(date))
	})
}
