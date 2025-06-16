package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCountryHolidayCreateFails() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)

	newCountryHoliday := models.CountryHoliday{
		ID:              uuid.Must(uuid.NewV4()),
		CountryId:       country.ID,
		HolidayName:     "Fake Holiday",
		ObservationDate: time.Now(),
	}

	err := suite.DB().Create(&newCountryHoliday)
	suite.Error(err)
	suite.Contains(err.Error(), "violates not-null constraint", "All model fields are readonly which should trigger a not-null constraint violation")
}

func (suite *ModelSuite) TestCountryHolidayUpdateFails() {
	country := factory.FetchOrBuildCountry(suite.DB(), nil, nil)

	holidays := models.CountryHolidays{}
	err := suite.DB().Where("country_id = ?", country.ID).All(&holidays)
	suite.NoError(err)
	suite.NotEmpty(holidays, "the default country is missing holidays")

	holiday := holidays[0]
	holiday.HolidayName = "My New Holiday Name"
	err = suite.DB().Save(&holiday)

	suite.Error(err)
	suite.Contains(err.Error(), "syntax error at or near \"WHERE\"", "All model fields are readonly which should trigger a syntax error on the update statement")
}
