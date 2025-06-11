package calendar

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/utils"
)

func (suite *CalendarSuite) TestDateSelectionCheckerFailure() {
	service := NewDateSelectionChecker()

	suite.Run("IsDateWeekendHoliday fails for nil app context", func() {
		var appCtx appcontext.AppContext = nil
		info, err := service.IsDateWeekendHoliday(appCtx, "US", time.Now())
		suite.Nil(info, "info should be nil when app context is nil")
		suite.Error(err, "Expected an error when app context is nil")
		suite.Contains(err.Error(), "app context is nil", "Error should mention nil app context")
	})

	countryCodefailureTestCases := []struct {
		name        string
		countryCode string
		expectError string
	}{
		{
			name:        "empty country code",
			countryCode: "",
		},
		{
			name:        "single character country code",
			countryCode: "a",
		},
		{
			name:        "3 character country code",
			countryCode: "123",
		},
	}

	for _, tc := range countryCodefailureTestCases {
		suite.Run("IsDateWeekendHoliday check fails for "+tc.name, func() {
			info, err := service.IsDateWeekendHoliday(suite.AppContextForTest(), tc.countryCode, time.Now())
			suite.Nil(info, "info should be nil when countryCode is invalid")
			suite.Error(err, "Expected an error when countryCode is invalid")
			suite.Contains(err.Error(), "countryCode should be precisely two characters", "Error should mention invalid countryCode")
		})
	}

	suite.Run("IsDateWeekendHoliday fails for zero date", func() {
		var date = time.Time{}
		info, err := service.IsDateWeekendHoliday(suite.AppContextForTest(), "US", date)
		suite.Nil(info, "info should be nil when date is zero")
		suite.Error(err, "Expected an error when date is zero")
		suite.Contains(err.Error(), "date value is zero", "Error should mention date value is zero")
	})

	suite.Run("IsDateWeekendHoliday fails for invalid country code", func() {
		info, err := service.IsDateWeekendHoliday(suite.AppContextForTest(), "XX", time.Now())
		suite.Nil(info, "info should be nil when date is zero")
		suite.Error(err, "Expected an error when country code is not valid")
		suite.Contains(err.Error(), "the country code provided was not found", "Error should mention the country code not found")
	})

}

func (suite *CalendarSuite) TestDateSelectionChecker() {
	service := NewDateSelectionChecker()
	suite.Run("date is both holiday and weekend", func() {

		syria := factory.FetchOrBuildCountry(suite.DB(), []factory.Customization{
			{
				Model: models.Country{
					Country: "SY",
				},
			},
		}, nil)
		// The service returns the country name in title case
		expectedCountryName := utils.ToTitleCase(syria.CountryName)
		expectedDate := time.Date(2025, 3, 8, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), syria.Country, expectedDate)
		suite.Equal(syria.Country, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.True(info.IsHoliday)
		suite.True(info.IsWeekend)
	})

	suite.Run("date is only a weekend", func() {

		unitedStates := factory.FetchOrBuildCountry(suite.DB(), []factory.Customization{
			{
				Model: models.Country{
					Country: "US",
				},
			},
		}, nil)

		expectedCountryName := utils.ToTitleCase(unitedStates.CountryName)
		expectedDate := time.Date(2025, 6, 14, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), unitedStates.Country, expectedDate)
		suite.Equal(unitedStates.Country, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.False(info.IsHoliday)
		suite.True(info.IsWeekend)
	})

	suite.Run("date is only a holiday", func() {

		france := factory.FetchOrBuildCountry(suite.DB(), []factory.Customization{
			{
				Model: models.Country{
					Country: "FR",
				},
			},
		}, nil)

		expectedCountryName := utils.ToTitleCase(france.CountryName)
		expectedDate := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)
		info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), france.Country, expectedDate)
		suite.Equal(france.Country, info.CountryCode)
		suite.Equal(expectedCountryName, info.CountryName)
		suite.Equal(expectedDate, info.Date)
		suite.True(info.IsHoliday)
		suite.False(info.IsWeekend)
	})

	workdayTestCases := []struct {
		name        string
		countryCode string
		date        time.Time
	}{
		{
			name:        "United States",
			countryCode: "US",
			date:        time.Date(2025, 6, 13, 0, 0, 0, 0, time.UTC), // Friday is a workday in the US.
		},
		{
			name:        "Bahrain",
			countryCode: "BH",
			date:        time.Date(2025, 6, 14, 0, 0, 0, 0, time.UTC), // Saturday is a workday in Bahrain.
		},
		{
			name:        "Kuwait",
			countryCode: "KW",
			date:        time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC), // Sunday is a workday in Kuwait.
		},
	}

	for _, tc := range workdayTestCases {
		suite.Run("date is a workday in "+tc.name, func() {
			country := factory.FetchOrBuildCountry(suite.DB(), []factory.Customization{
				{
					Model: models.Country{
						Country: tc.countryCode,
					},
				},
			}, nil)

			expectedCountryName := utils.ToTitleCase(country.CountryName)
			expectedDate := tc.date
			info, _ := service.IsDateWeekendHoliday(suite.AppContextForTest(), country.Country, tc.date)
			suite.Equal(country.Country, info.CountryCode)
			suite.Equal(expectedCountryName, info.CountryName)
			suite.Equal(expectedDate, info.Date)
			suite.False(info.IsHoliday)
			suite.False(info.IsWeekend)
		})
	}
}
