package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReZip3Validations() {
	suite.T().Run("test valid ReZip3", func(t *testing.T) {
		validReZip3 := models.ReZip3{
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
			Zip3:                  "606",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReZip3, expErrors)
	})

	suite.T().Run("test invalid ReZip3", func(t *testing.T) {
		invalidReZip3 := &models.ReZip3{}
		expErrors := map[string][]string{
			"domestic_service_area_id": {"DomesticServiceAreaID can not be blank."},
			"zip3":                     {"Zip3 not in range(3, 3)"},
		}
		suite.verifyValidationErrors(invalidReZip3, expErrors)
	})

	suite.T().Run("test when zip3 is not a length of 3", func(t *testing.T) {
		invalidReZip3 := &models.ReZip3{
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
			Zip3:                  "60",
		}
		expErrors := map[string][]string{
			"zip3": {"Zip3 not in range(3, 3)"},
		}
		suite.verifyValidationErrors(invalidReZip3, expErrors)
	})
}
