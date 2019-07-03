package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicWeightTicketSetDocumentInstantiation() {
	expenseDoc := &models.WeightTicketSetDocument{}

	expErrors := map[string][]string{
		"move_document_id": {"MoveDocumentID can not be blank."},
		"vehicle_nickname": {"VehicleNickname can not be blank."},
		"vehicle_options":  {"VehicleOptions can not be blank."},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}
