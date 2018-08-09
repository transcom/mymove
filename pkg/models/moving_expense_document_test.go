package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicMovingDocumentInstantiation() {
	expenseDoc := &models.MovingExpenseDocument{}

	expErrors := map[string][]string{
		"move_document_id":       {"MoveDocumentID can not be blank."},
		"payment_method":         {"PaymentMethod can not be blank."},
		"moving_expense_type":    {"MovingExpenseType can not be blank."},
		"requested_amount_cents": {"0 is not greater than 0."},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}
