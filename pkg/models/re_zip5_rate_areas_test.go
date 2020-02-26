package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReZip5RateAreaValidations() {
	suite.T().Run("test valid ReZip5RateArea", func(t *testing.T) {
		validReZip5RateArea := models.ReZip5RateArea{
			RateAreaID: uuid.Must(uuid.NewV4()),
			Zip5:       "60610",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReZip5RateArea, expErrors)
	})

	suite.T().Run("test invalid ReZip5RateArea", func(t *testing.T) {
		emptyReZip5RateArea := models.ReZip5RateArea{}
		expErrors := map[string][]string{
			"rate_area_id": {"RateAreaID can not be blank."},
			"zip5":         {"Zip5 not in range(5, 5)"},
		}
		suite.verifyValidationErrors(&emptyReZip5RateArea, expErrors)
	})

	suite.T().Run("test when zip5 is not a length of 5", func(t *testing.T) {
		invalidReZip5RateArea := models.ReZip5RateArea{
			RateAreaID: uuid.Must(uuid.NewV4()),
			Zip5:       "6034",
		}
		expErrors := map[string][]string{
			"zip5": {"Zip5 not in range(5, 5)"},
		}
		suite.verifyValidationErrors(&invalidReZip5RateArea, expErrors)
	})
}
