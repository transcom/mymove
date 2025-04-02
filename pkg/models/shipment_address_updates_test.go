package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestShipmentAddressUpdateValidation() {
	suite.Run("test valid ShipmentAddressUpdate", func() {
		originalAddress := factory.BuildDefaultAddress(suite.DB())
		newAddress := factory.BuildDefaultAddress(suite.DB())
		shipmentID := uuid.Must(uuid.NewV4())
		validAddressChange := models.ShipmentAddressUpdate{
			ContractorRemarks: "prev tenant of house at original address filled it w balloons and floated away",
			Status:            models.ShipmentAddressUpdateStatusRequested,
			ShipmentID:        shipmentID,
			OriginalAddressID: originalAddress.ID,
			NewAddressID:      newAddress.ID,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validAddressChange, expErrors, nil)
	})
	suite.Run("test empty ShipmentAddressUpdate", func() {
		emptyAddressChange := models.ShipmentAddressUpdate{}
		expErrors := map[string][]string{
			"shipment_id":         {"ShipmentID can not be blank."},
			"original_address_id": {"OriginalAddressID can not be blank."},
			"new_address_id":      {"NewAddressID can not be blank."},
			"status":              {"Status is not in the list [REQUESTED, REJECTED, APPROVED]."},
			"contractor_remarks":  {"ContractorRemarks can not be blank."},
		}

		suite.verifyValidationErrors(&emptyAddressChange, expErrors, nil)
	})
	suite.Run("test invalid ShipmentAddressUpdate", func() {
		validAddressChange := models.ShipmentAddressUpdate{
			ContractorRemarks: "",
			Status:            "not a real status",
			ShipmentID:        uuid.Nil,
			OriginalAddressID: uuid.Nil,
			NewAddressID:      uuid.Nil,
		}
		expErrors := map[string][]string{
			"shipment_id":         {"ShipmentID can not be blank."},
			"original_address_id": {"OriginalAddressID can not be blank."},
			"new_address_id":      {"NewAddressID can not be blank."},
			"status":              {"Status is not in the list [REQUESTED, REJECTED, APPROVED]."},
			"contractor_remarks":  {"ContractorRemarks can not be blank."},
		}

		suite.verifyValidationErrors(&validAddressChange, expErrors, nil)
	})
}
