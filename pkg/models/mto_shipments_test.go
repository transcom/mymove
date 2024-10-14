package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestMTOShipmentValidation() {
	suite.Run("test valid MTOShipment", func() {
		// mock weights
		estimatedWeight := unit.Pound(1000)
		actualWeight := unit.Pound(980)
		sitDaysAllowance := 90
		tacType := models.LOATypeHHG
		sacType := models.LOATypeHHG
		marketCode := models.MarketCodeDomestic
		validMTOShipment := models.MTOShipment{
			MoveTaskOrderID:      uuid.Must(uuid.NewV4()),
			Status:               models.MTOShipmentStatusApproved,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			SITDaysAllowance:     &sitDaysAllowance,
			TACType:              &tacType,
			SACType:              &sacType,
			MarketCode:           marketCode,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validMTOShipment, expErrors)
	})

	suite.Run("test empty MTOShipment", func() {
		emptyMTOShipment := models.MTOShipment{}
		expErrors := map[string][]string{
			"move_task_order_id": {"MoveTaskOrderID can not be blank."},
			"status":             {"Status is not in the list [APPROVED, REJECTED, SUBMITTED, DRAFT, CANCELLATION_REQUESTED, CANCELED, DIVERSION_REQUESTED]."},
		}
		suite.verifyValidationErrors(&emptyMTOShipment, expErrors)
	})

	suite.Run("test rejected MTOShipment", func() {
		rejectionReason := "bad shipment"
		marketCode := models.MarketCodeDomestic
		rejectedMTOShipment := models.MTOShipment{
			MoveTaskOrderID: uuid.Must(uuid.NewV4()),
			Status:          models.MTOShipmentStatusRejected,
			MarketCode:      marketCode,
			RejectionReason: &rejectionReason,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&rejectedMTOShipment, expErrors)
	})

	suite.Run("test validation failures", func() {
		// mock weights
		estimatedWeight := unit.Pound(-1000)
		actualWeight := unit.Pound(-980)
		billableWeightCap := unit.Pound(-1)
		billableWeightJustification := ""
		sitDaysAllowance := -1
		serviceOrderNumber := ""
		tacType := models.LOAType("FAKE")
		marketCode := models.MarketCode("x")
		invalidMTOShipment := models.MTOShipment{
			MoveTaskOrderID:             uuid.Must(uuid.NewV4()),
			Status:                      models.MTOShipmentStatusRejected,
			PrimeEstimatedWeight:        &estimatedWeight,
			PrimeActualWeight:           &actualWeight,
			BillableWeightCap:           &billableWeightCap,
			BillableWeightJustification: &billableWeightJustification,
			SITDaysAllowance:            &sitDaysAllowance,
			ServiceOrderNumber:          &serviceOrderNumber,
			StorageFacilityID:           &uuid.Nil,
			TACType:                     &tacType,
			SACType:                     &tacType,
			MarketCode:                  marketCode,
		}
		expErrors := map[string][]string{
			"prime_estimated_weight":        {"-1000 is not greater than 0."},
			"prime_actual_weight":           {"-980 is not greater than 0."},
			"rejection_reason":              {"RejectionReason can not be blank."},
			"billable_weight_cap":           {"-1 is less than zero."},
			"billable_weight_justification": {"BillableWeightJustification can not be blank."},
			"sitdays_allowance":             {"-1 is not greater than -1."},
			"service_order_number":          {"ServiceOrderNumber can not be blank."},
			"storage_facility_id":           {"StorageFacilityID can not be blank."},
			"tactype":                       {"TACType is not in the list [HHG, NTS]."},
			"sactype":                       {"SACType is not in the list [HHG, NTS]."},
			"market_code":                   {"MarketCode is not in the list [d, i]."},
		}
		suite.verifyValidationErrors(&invalidMTOShipment, expErrors)
	})
}
