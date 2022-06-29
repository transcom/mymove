package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestStorageFacilityValidation() {
	suite.Run("test valid StorageFacility", func() {
		validMTOShipment := models.StorageFacility{
			FacilityName: "Test Storage Facility",
			AddressID:    uuid.Must(uuid.NewV4()),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOShipment, expErrors)
	})

	suite.Run("test invalid StorageFacility", func() {
		lotNumber := ""
		phone := ""
		email := ""
		invalidMTOShipment := models.StorageFacility{
			LotNumber: &lotNumber,
			Phone:     &phone,
			Email:     &email,
		}
		expErrors := map[string][]string{
			"address_id":    {"AddressID can not be blank."},
			"facility_name": {"FacilityName can not be blank."},
			"lot_number":    {"LotNumber can not be blank."},
			"phone":         {"Phone can not be blank."},
			"email":         {"Email can not be blank."},
		}
		suite.verifyValidationErrors(&invalidMTOShipment, expErrors)
	})
}
