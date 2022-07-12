package models_test

import (
	"fmt"
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestServiceItemParamKeyValidation() {
	validServiceItemParamNames := strings.Join(models.ValidServiceItemParamNameStrings, ", ")
	validServiceItemParamTypes := strings.Join(models.ValidServiceItemParamTypes, ", ")
	validServiceItemParamOrigins := strings.Join(models.ValidServiceItemParamOrigins, ", ")

	suite.Run("test valid ServiceItemParamKey", func() {
		validServiceItemParamKey := models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameZipPickupAddress,
			Description: "Description",
			Type:        "STRING",
			Origin:      "PRIME",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validServiceItemParamKey, expErrors)
	})

	suite.Run("test empty ServiceItemParamKey", func() {
		invalidServiceItemParamKey := models.ServiceItemParamKey{}

		expErrors := map[string][]string{
			"key":         {"Key can not be blank.", fmt.Sprintf("Key is not in the list [%s].", validServiceItemParamNames)},
			"description": {"Description can not be blank."},
			"type":        {"Type can not be blank.", fmt.Sprintf("Type is not in the list [%s].", validServiceItemParamTypes)},
			"origin":      {"Origin can not be blank.", fmt.Sprintf("Origin is not in the list [%s].", validServiceItemParamOrigins)},
		}

		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})

	suite.Run("test invalid key name for ServiceItemParamKey", func() {
		invalidServiceItemParamKey := models.ServiceItemParamKey{
			Key:         "foo",
			Description: "Description",
			Type:        "STRING",
			Origin:      "PRIME",
		}
		expErrors := map[string][]string{
			"key": {fmt.Sprintf("Key is not in the list [%s].", validServiceItemParamNames)},
		}
		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})

	suite.Run("test invalid type for ServiceItemParamKey", func() {
		invalidServiceItemParamKey := models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameZipPickupAddress,
			Description: "Description",
			Type:        "TIME",
			Origin:      "PRIME",
		}
		expErrors := map[string][]string{
			"type": {fmt.Sprintf("Type is not in the list [%s].", validServiceItemParamTypes)},
		}
		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})

	suite.Run("test invalid origin for ServiceItemParamKey", func() {
		invalidServiceItemParamKey := models.ServiceItemParamKey{
			Key:         models.ServiceItemParamNameZipPickupAddress,
			Description: "Description",
			Type:        "DATE",
			Origin:      "OPTIMUS",
		}
		expErrors := map[string][]string{
			"origin": {fmt.Sprintf("Origin is not in the list [%s].", validServiceItemParamOrigins)},
		}
		suite.verifyValidationErrors(&invalidServiceItemParamKey, expErrors)
	})
}
