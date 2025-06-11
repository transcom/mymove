package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCountryWeekendCreateFails() {

	country := factory.FetchOrBuildCountry(suite.DB(), []factory.Customization{
		{
			Model: models.Country{
				Country: "AQ",
			},
		},
	}, nil)

	newCountryWeekend := models.CountryWeekend{
		ID:                 uuid.Must(uuid.NewV4()),
		CountryId:          country.ID,
		IsMondayWeekend:    false,
		IsTuesdayWeekend:   false,
		IsWednesdayWeekend: false,
		IsThursdayWeekend:  false,
		IsFridayWeekend:    false,
		IsSaturdayWeekend:  true,
		IsSundayWeekend:    true,
	}

	err := suite.DB().Create(&newCountryWeekend)
	suite.Error(err)
	suite.Contains(err.Error(), "violates not-null constraint", "All model fields are readonly which should trigger a not-null constraint violation")
}

func (suite *ModelSuite) TestCountryWeekendUpdateFails() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)

	countryWeekend := models.CountryWeekend{}
	err := suite.DB().Where("country_id = ?", country.ID).First(&countryWeekend)
	suite.NoError(err)
	suite.NotNil(countryWeekend, "the default country is missing weekend data")

	countryWeekend.IsMondayWeekend = !countryWeekend.IsMondayWeekend
	err = suite.DB().Save(&countryWeekend)

	suite.Error(err)
	suite.Contains(err.Error(), "syntax error at or near \"WHERE\"", "All model fields are readonly which should trigger a syntax error on the update statement")
}

func (suite *ModelSuite) TestIsWeekend() {

	weekends := models.CountryWeekend{
		IsSaturdayWeekend: true,
		IsSundayWeekend:   true,
	}

	testCases := []struct {
		date     time.Time
		expected bool
	}{
		{time.Date(2025, 6, 9, 0, 0, 0, 0, time.UTC), false},  // Monday
		{time.Date(2025, 6, 10, 0, 0, 0, 0, time.UTC), false}, // Tuesday
		{time.Date(2025, 6, 11, 0, 0, 0, 0, time.UTC), false}, // Wednesday
		{time.Date(2025, 6, 12, 0, 0, 0, 0, time.UTC), false}, // Thursday
		{time.Date(2025, 6, 13, 0, 0, 0, 0, time.UTC), false}, // Friday
		{time.Date(2025, 6, 14, 0, 0, 0, 0, time.UTC), true},  // Saturday
		{time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC), true},  // Sunday
	}

	for _, tc := range testCases {
		result := weekends.IsWeekend(tc.date)
		suite.Equal(tc.expected, result, "IsWeekend should correctly identify a weekend day")
	}
}
