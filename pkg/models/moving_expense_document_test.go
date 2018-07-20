package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestBasicMovingDocumentInstantiation() {
	expenseDoc := &models.MovingExpenseDocument{}

	expErrors := map[string][]string{
		"move_document_id":    {"MoveDocumentID can not be blank."},
		"reimbursement_id":    {"ReimbursementID can not be blank."},
		"moving_expense_type": {"MovingExpenseType can not be blank."},
	}

	suite.verifyValidationErrors(expenseDoc, expErrors)
}
