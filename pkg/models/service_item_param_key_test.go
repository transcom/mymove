package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestServiceItemParamKeyValidation() {
	suite.T().Run("test valid ServiceItemParamKey", func(t *testing.T) {
		validServiceItemParamKey := models.ServiceItemParamKey{
			Key:         "Key",
			Description: "Description",
			Type:        "STRING",
			Origin:      "PRIME",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceItemParamKey, expErrors)
	})

	suite.T().Run("test empty ServiceItemParamKey", func(t *testing.T) {
		invalidServiceItemParamKey := models.ServiceItemParamKey{}

		expErrors := map[string][]string{
			"key":         {"Key can not be blank."},
			"description": {"Description can not be blank."},
			"type":        {"Type can not be blank.", "Type is not in the list [STRING, DATE, INTEGER, DECIMAL]."},
			"origin":      {"Origin can not be blank.", "Origin is not in the list [PRIME, SYSTEM]."},
		}

		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})

	suite.T().Run("test invalid type for ServiceItemParamKey", func(t *testing.T) {
		invalidServiceItemParamKey := models.ServiceItemParamKey{
			Key:         "Key",
			Description: "Description",
			Type:        "TIME",
			Origin:      "PRIME",
		}
		expErrors := map[string][]string{
			"type": {"Type is not in the list [STRING, DATE, INTEGER, DECIMAL]."},
		}
		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})

	suite.T().Run("test invalid origin for ServiceItemParamKey", func(t *testing.T) {
		invalidServiceItemParamKey := models.ServiceItemParamKey{
			Key:         "Key",
			Description: "Description",
			Type:        "DATE",
			Origin:      "OPTIMUS",
		}
		expErrors := map[string][]string{
			"origin": {"Origin is not in the list [PRIME, SYSTEM]."},
		}
		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})
}
