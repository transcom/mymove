package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReDomesticServiceAreaValidation() {
	suite.T().Run("test valid ReDomesticServiceArea", func(t *testing.T) {
		validReDomesticServiceArea := models.ReDomesticServiceArea{
			BasePointCity:   "New York",
			State:           "NY",
			ServiceArea:     9,
			ServiceSchedule: 2,
			SITPDSchedule:   2,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticServiceArea, expErrors)
	})

	suite.T().Run("test invalid ReDomesticServiceArea", func(t *testing.T) {
		invalidReDomesticServiceArea := models.ReDomesticServiceArea{}
		expErrors := map[string][]string{
			"base_point_city":    {"BasePointCity can not be blank."},
			"state":              {"State can not be blank."},
			"service_area":       {"ServiceArea can not be blank."},
			"service_schedule":   {"0 is not greater than 0."},
			"s_i_t_p_d_schedule": {"0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReDomesticServiceArea, expErrors)
	})

	suite.T().Run("test schedules over 3 for ReDomesticServiceArea", func(t *testing.T) {
		invalidReDomesticServiceArea := models.ReDomesticServiceArea{
			BasePointCity:   "New York",
			State:           "NY",
			ServiceArea:     9,
			ServiceSchedule: 4,
			SITPDSchedule:   5,
		}
		expErrors := map[string][]string{
			"service_schedule":   {"4 is not less than 4."},
			"s_i_t_p_d_schedule": {"5 is not less than 4."},
		}
		suite.verifyValidationErrors(&invalidReDomesticServiceArea, expErrors)
	})

	suite.T().Run("test schedules less than 1 for ReDomesticServiceArea", func(t *testing.T) {
		invalidReDomesticServiceArea := models.ReDomesticServiceArea{
			BasePointCity:   "New York",
			State:           "NY",
			ServiceArea:     9,
			ServiceSchedule: -3,
			SITPDSchedule:   -1,
		}
		expErrors := map[string][]string{
			"service_schedule":   {"-3 is not greater than 0."},
			"s_i_t_p_d_schedule": {"-1 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReDomesticServiceArea, expErrors)
	})
}
