package calendar

import (
	"time"
)

func (suite *CalendarSuite) TestDateSelectionChecker() {
	expectedCountryCode := "US"
	expectedCountryName := "United States"
	service := NewDateSelectionChecker()
	suite.Run("date is both holiday and weekend - US", func() {
		expectedDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), expectedCountryCode, expectedDate)
		suite.Equal(expectedCountryCode, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.True(info.IsHoliday)
		suite.True(info.IsWeekend)
	})

	suite.Run("date is only weekend - US", func() {
		expectedDate := time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), expectedCountryCode, expectedDate)
		suite.Equal(expectedCountryCode, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.False(info.IsHoliday)
		suite.True(info.IsWeekend)
	})

	suite.Run("date is only holiday - US", func() {
		expectedDate := time.Date(2024, 7, 4, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), expectedCountryCode, expectedDate)
		suite.Equal(expectedCountryCode, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.True(info.IsHoliday)
		suite.False(info.IsWeekend)
	})

	suite.Run("regular work day - US", func() {
		expectedDate := time.Date(2024, 8, 6, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), expectedCountryCode, expectedDate)
		suite.Equal(expectedCountryCode, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.False(info.IsHoliday)
		suite.False(info.IsWeekend)
	})
}
