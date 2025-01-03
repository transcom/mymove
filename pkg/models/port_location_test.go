package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestPortLocationValidation() {
	suite.Run("test valid PortLocation", func() {
		validPortLocation := models.PortLocation{
			ID:                   uuid.Must(uuid.NewV4()),
			PortId:               uuid.Must(uuid.NewV4()),
			CitiesId:             uuid.Must(uuid.NewV4()),
			UsPostRegionCitiesId: uuid.Must(uuid.NewV4()),
			CountryId:            uuid.Must(uuid.NewV4()),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validPortLocation, expErrors)
	})

	suite.Run("test missing required fields", func() {
		invalidPortLocation := models.PortLocation{
			ID:        uuid.Must(uuid.NewV4()),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expErrors := map[string][]string{
			"port_id":                  {"PortID can not be blank."},
			"cities_id":                {"CitiesID can not be blank."},
			"us_post_region_cities_id": {"UsPostRegionCitiesID can not be blank."},
			"country_id":               {"CountryID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidPortLocation, expErrors)
	})
}
