package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReZip3Validations() {
	suite.Run("test valid ReZip3", func() {
		validReZip3 := models.ReZip3{
			ContractID:            uuid.Must(uuid.NewV4()),
			Zip3:                  "606",
			BasePointCity:         "New York",
			State:                 "NY",
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReZip3, expErrors)
	})

	suite.Run("test invalid ReZip3", func() {
		emptyReZip3 := models.ReZip3{}
		expErrors := map[string][]string{
			"contract_id":              {"ContractID can not be blank."},
			"zip3":                     {"Zip3 not in range(3, 3)"},
			"base_point_city":          {"BasePointCity can not be blank."},
			"state":                    {"State can not be blank."},
			"domestic_service_area_id": {"DomesticServiceAreaID can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReZip3, expErrors)
	})

	suite.Run("test when zip3 is not a length of 3", func() {
		invalidReZip3 := models.ReZip3{
			ContractID:            uuid.Must(uuid.NewV4()),
			Zip3:                  "60",
			BasePointCity:         "New York",
			State:                 "NY",
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{
			"zip3": {"Zip3 not in range(3, 3)"},
		}
		suite.verifyValidationErrors(&invalidReZip3, expErrors)
	})
}
