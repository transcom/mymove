package ghcrateengine

import (
	"time"
)

func (suite *GHCRateEngineServiceSuite) TestIsPeakPeriod() {
	suite.Run("within peak", func() {
		dateInYear := peakStart.addDate(1, 0)
		date := time.Date(2019, dateInYear.month, 1, 0, 0, 0, 0, time.UTC)
		suite.True(IsPeakPeriod(date))
	})

	suite.Run("on peak start date", func() {
		date := time.Date(2019, peakStart.month, peakStart.day, 0, 0, 0, 0, time.UTC)
		suite.True(IsPeakPeriod(date))
	})

	suite.Run("on peak end date", func() {
		date := time.Date(2019, peakEnd.month, peakEnd.day, 0, 0, 0, 0, time.UTC)
		suite.True(IsPeakPeriod(date))
	})

	suite.Run("just before peak start date", func() {
		dateInYear := peakStart.addDate(0, -1)
		date := time.Date(2019, dateInYear.month, dateInYear.day, 0, 0, 0, 0, time.UTC)
		suite.False(IsPeakPeriod(date))
	})

	suite.Run("just outside peak start date", func() {
		dateInYear := peakEnd.addDate(0, 1)
		date := time.Date(2019, dateInYear.month, dateInYear.day, 0, 0, 0, 0, time.UTC)
		suite.False(IsPeakPeriod(date))
	})

	suite.Run("outside peak", func() {
		dateInYear := peakEnd.addDate(1, 0)
		date := time.Date(2019, dateInYear.month, 1, 0, 0, 0, 0, time.UTC)
		suite.False(IsPeakPeriod(date))
	})
}
