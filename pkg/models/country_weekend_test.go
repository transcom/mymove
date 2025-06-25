package models_test

import (
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
