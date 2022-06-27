package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReRateAreaValidation() {
	suite.Run("test valid ReRateArea", func() {
		validReRateArea := models.ReRateArea{
			ContractID: uuid.Must(uuid.NewV4()),
			IsOconus:   true,
			Code:       "123abc",
			Name:       "California",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReRateArea, expErrors)
	})

	suite.Run("test empty ReRateArea", func() {
		emptyReRateArea := models.ReRateArea{}
		expErrors := map[string][]string{
			"contract_id": {"ContractID can not be blank."},
			"code":        {"Code can not be blank."},
			"name":        {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReRateArea, expErrors)
	})
}
