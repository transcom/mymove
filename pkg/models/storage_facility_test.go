package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestStorageFacilityValidation() {
	suite.T().Run("test valid StorageFacility", func(t *testing.T) {
		validMTOShipment := models.StorageFacility{}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOShipment, expErrors)
	})

	suite.T().Run("test invalid StorageFacility", func(t *testing.T) {
		facilityName := ""
		lotNumber := ""
		phone := ""
		email := ""
		invalidMTOShipment := models.StorageFacility{
			FacilityName: &facilityName,
			LotNumber:    &lotNumber,
			Phone:        &phone,
			Email:        &email,
		}
		expErrors := map[string][]string{
			"facility_name": {"FacilityName can not be blank."},
			"lot_number":    {"LotNumber can not be blank."},
			"phone":         {"Phone can not be blank."},
			"email":         {"Email can not be blank."},
		}
		suite.verifyValidationErrors(&invalidMTOShipment, expErrors)
	})
}
