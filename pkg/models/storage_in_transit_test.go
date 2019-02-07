package models_test

import (
	"testing"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestStorageInTransitValidations() {
	suite.T().Run("test valid storage in transit", func(t *testing.T) {
		validStorageInTransit := testdatagen.MakeDefaultStorageInTransit(suite.DB())
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validStorageInTransit, expErrors)
	})

	suite.T().Run("test invalid/empty storage in transit", func(t *testing.T) {
		invalidStorageInTransit := &models.StorageInTransit{}
		expErrors := map[string][]string{
			"shipment_id":          {"ShipmentID can not be blank."},
			"status":               {"Status is not in the list [REQUESTED, APPROVED, DENIED, IN_SIT, RELEASED, DELIVERED]."},
			"location":             {"Location is not in the list [ORIGIN, DESTINATION]."},
			"estimated_start_date": {"EstimatedStartDate can not be blank."},
			"warehouse_id":         {"WarehouseID can not be blank."},
			"warehouse_name":       {"WarehouseName can not be blank."},
			"warehouse_address_id": {"WarehouseAddressID can not be blank."},
		}
		suite.verifyValidationErrors(invalidStorageInTransit, expErrors)
	})
}
