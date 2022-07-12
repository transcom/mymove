package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestServiceParamValidation() {
	suite.Run("test valid ServiceParam", func() {
		validServiceParam := models.ServiceParam{
			ServiceID:             uuid.Must(uuid.NewV4()),
			ServiceItemParamKeyID: uuid.Must(uuid.NewV4()),
			IsOptional:            false,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceParam, expErrors)
	})

	suite.Run("test empty ServiceParam", func() {
		invalidServiceParam := models.ServiceParam{}

		expErrors := map[string][]string{
			"service_id":                {"ServiceID can not be blank."},
			"service_item_param_key_id": {"ServiceItemParamKeyID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidServiceParam, expErrors)
	})
}
