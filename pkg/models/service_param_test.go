package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestServiceParamValidation() {
	suite.T().Run("test valid ServiceParam", func(t *testing.T) {
		validServiceParam := models.ServiceParam{
			ServiceID:             uuid.Must(uuid.NewV4()),
			ServiceItemParamKeyID: uuid.Must(uuid.NewV4()),
			IsOptional:            false,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceParam, expErrors)
	})

	suite.T().Run("test empty ServiceParam", func(t *testing.T) {
		invalidServiceParam := models.ServiceParam{}

		expErrors := map[string][]string{
			"service_id":                {"ServiceID can not be blank."},
			"service_item_param_key_id": {"ServiceItemParamKeyID can not be blank."},
		}

		suite.verifyValidationErrors(&invalidServiceParam, expErrors)
	})
}
