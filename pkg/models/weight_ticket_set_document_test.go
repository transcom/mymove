package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicWeightTicketSetDocumentInstantiation() {
	expenseDoc := &models.WeightTicketSetDocument{}

	expErrors := map[string][]string{
		"move_document_id":   {"MoveDocumentID can not be blank."},
		"vehicle_nickname":   {"VehicleNickname can not be blank."},
		"vehicle_options":    {"VehicleOptions can not be blank."},
		"empty_weight":       {"0 is not greater than 0."},
		"full_weight":        {"0 is not greater than 0."},
		"weight_ticket_date": {"WeightTicketDate can not be blank."},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}
